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
	"github.com/AspieSoft/goutil/fs/v2"
	"github.com/AspieSoft/goutil/v7"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tystuyfzand/less-go"
)


var config = map[string]string{
	"theme_name": "Smart Theme",
	"theme_version": "v0.0.1",
	"theme_uri": "github.com/AspieSoft/smart-theme",
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
	port := ""
	if len(os.Args) > 1 {
		if p, err := strconv.Atoi(os.Args[1]); err == nil && p >= 3000 && p <= 65535 /* golang only accepts 16 bit port numbers */ {
			port = strconv.Itoa(p)
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

	handleCompileTheme(port)

	if port == "" {
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := strings.Trim(r.URL.Path, "/")
		if url == "" {
			http.ServeFile(w, r, "./test/html/index.html")
			return
		}else{
			cUrl := string(regClean.ReplaceAll(regExt.ReplaceAll([]byte(url), []byte{}), []byte{}))

			for _, ext := range validExtList {
				if strings.HasSuffix(url, "."+ext) {
					if strings.HasPrefix(cUrl, "theme/") {
						if stat, err := os.Stat("./dist/"+strings.Replace(cUrl, "theme/", "", 1)+"."+ext); err == nil && !stat.IsDir() {
							http.ServeFile(w, r, "./dist/"+strings.Replace(cUrl, "theme/", "", 1)+"."+ext)
							return
						}
					}

					if stat, err := os.Stat("./test/"+cUrl+"."+ext); err == nil && !stat.IsDir() {
						http.ServeFile(w, r, "./test/"+cUrl+"."+ext)
						return
					}

					w.WriteHeader(404)
					w.Write([]byte("error 404"))
					return
				}
			}

			if stat, err := os.Stat("./test/html/"+cUrl[2:]+".html"); err == nil && !stat.IsDir() {
				http.ServeFile(w, r, "./test/html/"+cUrl[2:]+".html")
				return
			}

			for _, ext := range validExtList {
				if stat, err := os.Stat("./test/"+cUrl+"."+ext); err == nil && !stat.IsDir() {
					http.ServeFile(w, r, "./test/"+cUrl+"."+ext)
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

func handleCompileTheme(port string){
	if themeDir, err := filepath.Abs("./src"); err == nil {
		themeDist, err := filepath.Abs("./dist")
		if err != nil {
			return
		}

		//todo: replace less with sass/scss
		// compileSass(themeDir, themeDist)
		compileCSS(themeDir, themeDist)
		compileJS(themeDir, themeDist)

		if port == "" {
			return
		}

		watcher := fs.Watcher()
		watcher.OnAny = func(path, op string) {
			path = string(regex.Comp(`^%1[\\/]?`, themeDir).RepStrLit([]byte(path), []byte{}))
			if (strings.HasSuffix(path, ".scss") && !strings.HasSuffix(path, ".min.scss")) || (strings.HasSuffix(path, ".sass") && !strings.HasSuffix(path, ".min.sass")) {
				// compileSass(themeDir, themeDist)
			}else if strings.HasSuffix(path, ".less") && !strings.HasSuffix(path, ".min.less") {
				// compileLess(themeDir, themeDist)
			}else if strings.HasSuffix(path, ".css") && !strings.HasSuffix(path, ".min.css") {
				compileCSS(themeDir, themeDist)
			}else if strings.HasSuffix(path, ".js") && !strings.HasSuffix(path, ".min.js") {
				compileJS(themeDir, themeDist)
			}
		}
		watcher.WatchDir(themeDir)
	}
}


func compileSass(themeDir string, themeDist string){
	
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

func compileLess(themeDir string, themeDist string){
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

		if path, err := fs.JoinPath(themeDist, "style.min.css"); err == nil {
			os.WriteFile(path, res, 0755)
		}

		if path, err := fs.JoinPath(themeDir, "normalize.min.css"); err == nil {
			if buf, err := os.ReadFile(path); err == nil {
				res = append(append(buf, '\n'), res...)
			}
		}

		if path, err := fs.JoinPath(themeDist, "style.norm.min.css"); err == nil {
			os.WriteFile(path, res, 0755)
		}
	}
}


func compileCSS(themeDir string, themeDist string){
	res := []byte{}

	if files, err := os.ReadDir(themeDir); err == nil {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".less") && !strings.Contains(file.Name(), "config") {
				if filePath, err := fs.JoinPath(themeDir, string(regex.Comp(`\.less$`).RepStrLit([]byte(file.Name()), []byte(".css")))); err == nil {
					if buf, err := os.ReadFile(filePath); err == nil {
						res = append(res, buf...)
					}
				}
			}else if strings.Contains(file.Name(), "config") || strings.Contains(file.Name(), "fonts") {
				if srcPath, err := fs.JoinPath(themeDir, file.Name()); err == nil {
					if distPath, err := fs.JoinPath(themeDist, file.Name()); err == nil {
						fs.Copy(srcPath, distPath)
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

	if path, err := fs.JoinPath(themeDist, "style.css"); err == nil {
		os.WriteFile(path, res, 0755)
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	if res, err := m.Bytes("text/css", res); err == nil {
		res = append([]byte("/*! "+config["theme_name"]+" "+config["theme_version"]+" | "+config["theme_license"]+" | "+config["theme_uri"]+" */\n"), res...)

		if path, err := fs.JoinPath(themeDist, "style.min.css"); err == nil {
			os.WriteFile(path, res, 0755)
		}

		if path, err := fs.JoinPath(themeDir, "normalize.min.css"); err == nil {
			if buf, err := os.ReadFile(path); err == nil {
				res = append(append(buf, '\n'), res...)
			}
		}

		if path, err := fs.JoinPath(themeDist, "style.norm.min.css"); err == nil {
			os.WriteFile(path, res, 0755)
		}
	}
}


func compileJS(themeDir string, themeDist string){
	res := []byte{}

	if files, err := os.ReadDir(themeDir); err == nil {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".js") && !strings.HasSuffix(file.Name(), ".min.js") {
				if filePath, err := fs.JoinPath(themeDir, file.Name()); err == nil {
					if buf, err := os.ReadFile(filePath); err == nil {
						res = append(res, buf...)
						res = append(res, '\n')
					}
				}
			}
		}
	}

	/* if path, err := fs.JoinPath(themeDir, "script.js"); err == nil {
		if buf, err := os.ReadFile(path); err == nil {
			res = append(res, buf...)
			res = append(res, '\n')
		}
	} */

	res = append([]byte{';'}, res...)
	res = append(res, ';')

	if path, err := fs.JoinPath(themeDist, "script.js"); err == nil {
		os.WriteFile(path, res, 0755)
	}

	//todo: fix js minify returning error
	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)
	if res, err := m.Bytes("text/javascript", res); err == nil {
		if path, err := fs.JoinPath(themeDist, "script.min.js"); err == nil {
			os.WriteFile(path, res, 0755)
		}
	}else if path, err := fs.JoinPath(themeDist, "script.min.js"); err == nil {
		os.WriteFile(path, res, 0755)
	}
}
