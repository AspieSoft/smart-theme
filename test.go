package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/AspieSoft/go-regex-re2/v2"
	"github.com/AspieSoft/goutil/fs"
	"github.com/AspieSoft/goutil/v7"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tystuyfzand/less-go"
)


var config = map[string]string{
	"theme_name": "theme",
	"theme_version": "v1.0.0",
	"theme_uri": "github.com/example/theme",
	"theme_license": "MIT License",
}

var validExtList []string = []string{
	"css",
	"js",
	"less",
	"md",
	"txt",
	"pdf",
	"png",
	"jpg",
	"gif",
	"mp3",
	"mp4",
	"mov",
	"ogg",
	"wav",
	"webp",
	"webm",
	"weba",
	"ttf",
	"woff",
}


func main(){
	port := "3000"
	if len(os.Args) > 1 {
		if p, err := strconv.Atoi(os.Args[1]); err == nil && p >= 3000 && p <= 65535 /* golang only accepts 16 bit port numbers */ {
			port = strconv.Itoa(p)
		}
	}

	if buf, err := os.ReadFile("./config.json"); err == nil {
		if conf, err := goutil.JSON.Parse(buf); err == nil {
			for key, val := range conf {
				if str := goutil.Conv.ToString(val); str != "" {
					config[key] = str
				}
			}
		}
	}

	regExt, err := regexp.Compile(`\.[\w_-]+$`)
	if err != nil {
		log.Fatal(err)
		return
	}
	regClean, err := regexp.Compile(`[^\w_\-\/\.]+`)
	if err != nil {
		log.Fatal(err)
		return
	}

	handleCssMinify()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := strings.Trim(r.URL.Path, "/")
		if url == "" {
			http.ServeFile(w, r, "./html/index.html")
			return
		}else{
			cUrl := "./"+string(regClean.ReplaceAll(regExt.ReplaceAll([]byte(url), []byte{}), []byte{}))

			for _, ext := range validExtList {
				if strings.HasSuffix(url, "."+ext) {
					if stat, err := os.Stat(cUrl+"."+ext); err == nil && !stat.IsDir() {
						http.ServeFile(w, r, cUrl+"."+ext)
						return
					}

					w.WriteHeader(404)
					w.Write([]byte("error 404"))
					return
				}
			}

			if stat, err := os.Stat("./html/"+cUrl[2:]+".html"); err == nil && !stat.IsDir() {
				http.ServeFile(w, r, "./html/"+cUrl[2:]+".html")
				return
			}

			for _, ext := range validExtList {
				if stat, err := os.Stat(cUrl+"."+ext); err == nil && !stat.IsDir() {
					http.ServeFile(w, r, cUrl+"."+ext)
					return
				}
			}
		}
		
		w.WriteHeader(404)
		w.Write([]byte("error 404"))
	})

	fmt.Println("Running Server On Port "+port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleCssMinify(){
	if themeDir, err := filepath.Abs("./theme"); err == nil {
		compileCSS(themeDir)
		compileJS(themeDir)

		watcher := fs.Watcher()
		watcher.OnAny = func(path, op string) {
			path = string(regex.Comp(`^%1[\\/]?`, themeDir).RepStrLit([]byte(path), []byte{}))
			if strings.HasSuffix(path, ".css") && !strings.HasSuffix(path, ".min.css") {
				compileCSS(themeDir)
			}else if path == "script.js" {
				compileJS(themeDir)
			}
		}
		watcher.WatchDir(themeDir)
	}
}

func importLessFile(themeDir string, path string, importList *[]string) []byte {
	if path, err := fs.JoinPath(themeDir, path); err == nil {
		if file, err := os.ReadFile(path); err == nil {
			return regex.Comp(`@import\s*(["'\'])([^"'\']+\.less)(["'\']);`).RepFunc(file, func(data func(int) []byte) []byte {
				if filePath := string(data(2)); !goutil.Contains(*importList, filePath) {
					*importList = append(*importList, filePath)
					return importLessFile(themeDir, filePath, importList)
				}
				return []byte{}
			})
		}
	}
	return []byte{}
}

func compileLess(themeDir string){
	res := []byte{}

	if path, err := fs.JoinPath(themeDir, "style.less"); err == nil {
		if file, err := os.ReadFile(path); err == nil {
			importList := []string{"style.less"}
			file = regex.Comp(`@import\s*(["'\'])([^"'\']+\.less)(["'\']);`).RepFunc(file, func(data func(int) []byte) []byte {
				if filePath := string(data(2)); !goutil.Contains(importList, filePath) {
					importList = append(importList, filePath)
					return importLessFile(themeDir, filePath, &importList)
				}
				return []byte{}
			})

			out, err := less.Render(string(file), map[string]interface{}{"compress": true})
			if err != nil {
				fmt.Println(err)
			}else{
				res = append(res, []byte(out)...)
			}
		}
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	if res, err := m.Bytes("text/css", res); err == nil {
		res = append([]byte("/*! "+config["theme_name"]+" "+config["theme_version"]+" | "+config["theme_license"]+" | "+config["theme_uri"]+" */\n"), res...)

		if path, err := fs.JoinPath(themeDir, "style.min.css"); err == nil {
			os.WriteFile(path, res, 0755)
		}

		if path, err := fs.JoinPath(themeDir, "normalize.min.css"); err == nil {
			if buf, err := os.ReadFile(path); err == nil {
				res = append(append(buf, '\n'), res...)
			}
		}

		if path, err := fs.JoinPath(themeDir, "style.norm.min.css"); err == nil {
			os.WriteFile(path, res, 0755)
		}
	}
}

func compileCSS(themeDir string){
	res := []byte{}

	if files, err := os.ReadDir(themeDir); err == nil {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".less") {
				if filePath, err := fs.JoinPath(themeDir, string(regex.Comp(`\.less$`).RepStrLit([]byte(file.Name()), []byte(".css")))); err == nil {
					if buf, err := os.ReadFile(filePath); err == nil {
						res = append(res, buf...)
					}
				}
			}
		}
	}

	/* if path, err := fs.JoinPath(themeDir, "style.css"); err == nil {
		if buf, err := os.ReadFile(path); err == nil {
			res = append(res, buf...)
		}
	} */

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	if res, err := m.Bytes("text/css", res); err == nil {
		res = append([]byte("/*! "+config["theme_name"]+" "+config["theme_version"]+" | "+config["theme_license"]+" | "+config["theme_uri"]+" */\n"), res...)

		if path, err := fs.JoinPath(themeDir, "style.min.css"); err == nil {
			os.WriteFile(path, res, 0755)
		}

		if path, err := fs.JoinPath(themeDir, "normalize.min.css"); err == nil {
			if buf, err := os.ReadFile(path); err == nil {
				res = append(append(buf, '\n'), res...)
			}
		}

		if path, err := fs.JoinPath(themeDir, "style.norm.min.css"); err == nil {
			os.WriteFile(path, res, 0755)
		}
	}
}

func compileJS(themeDir string){
	res := []byte{}

	if path, err := fs.JoinPath(themeDir, "script.js"); err == nil {
		if buf, err := os.ReadFile(path); err == nil {
			res = append(res, buf...)
			res = append(res, '\n')
		}
	}

	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)
	if res, err := m.Bytes("text/javascript", res); err == nil {
		res = append([]byte{';'}, res...)
		res = append(res, ';')

		if path, err := fs.JoinPath(themeDir, "script.min.js"); err == nil {
			os.WriteFile(path, res, 0755)
		}
	}
}
