package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/AspieSoft/go-regex-re2/v2"
	"github.com/AspieSoft/goutil/fs/v2"
	"github.com/AspieSoft/goutil/v7"
	"github.com/lucasb-eyer/go-colorful"
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

		compileLess(themeDir, themeDist)
		// compileCSS(themeDir, themeDist)
		compileJS(themeDir, themeDist)

		if port == "" {
			return
		}

		watcher := fs.Watcher()
		watcher.OnAny = func(path, op string) {
			path = string(regex.Comp(`^%1[\\/]?`, themeDir).RepStrLit([]byte(path), []byte{}))
			if strings.HasSuffix(path, ".less") && !strings.HasSuffix(path, ".min.less") {
				compileLess(themeDir, themeDist)
			}else if strings.HasSuffix(path, ".css") && !strings.HasSuffix(path, ".min.css") {
				// compileCSS(themeDir, themeDist)
			}else if strings.HasSuffix(path, ".js") && !strings.HasSuffix(path, ".min.js") {
				compileJS(themeDir, themeDist)
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

func compileLess(themeDir string, themeDist string){
	res := []byte{}
	lessConfig := []byte{}

	if path, err := fs.JoinPath(themeDir, "config.less"); err == nil {
		if file, err := os.ReadFile(path); err == nil {
			lessConfig = compileGoLessConfig(file)
			lessConfig = append(lessConfig, '\n')
		}
	}

	if path, err := fs.JoinPath(themeDir, "style.less"); err == nil {
		if file, err := os.ReadFile(path); err == nil {
			importList := []string{"style.less", "config.less"}
			file = regex.Comp(`@import\s*(["'\'])([^"'\']+\.less)(["'\']);`).RepFunc(file, func(data func(int) []byte) []byte {
				if filePath := string(data(2)); !goutil.Contains(importList, filePath) {
					importList = append(importList, filePath)
					return importLessFile(themeDir, filePath, &importList)
				}
				return []byte{}
			})

			res = append(res, file...)

			/* out, err := less.Render(string(append(lessConfig, file...)), map[string]interface{}{"compress": true})
			if err != nil {
				fmt.Println(err)
			}else{
				res = append(res, []byte(out)...)
			} */
		}
	}

	out, err := less.Render(string(append(lessConfig, res...)), map[string]interface{}{"compress": false})
	if err != nil {
		fmt.Println(err)
	}else{
		res = []byte(out)
	}

	res = compileGoLess(res)

	out, err = less.Render(string(append(lessConfig, res...)), map[string]interface{}{"compress": true})
	if err != nil {
		fmt.Println(err)
	}else{
		res = []byte(out)
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	if res, err := m.Bytes("text/css", res); err == nil {
		res = append([]byte("/*! "+config["theme_name"]+" "+config["theme_version"]+" | "+config["theme_license"]+" | "+config["theme_uri"]+" */\n"), res...)

		if path, err := fs.JoinPath(themeDist, "style.css"); err == nil {
			os.WriteFile(path, res, 0755)
		}

		if path, err := fs.JoinPath(themeDist, "style.min.css"); err == nil {
			os.WriteFile(path, res, 0755)
		}

		if path, err := fs.JoinPath(themeDir, "normalize.min.css"); err == nil {
			if buf, err := os.ReadFile(path); err == nil {
				res = append(buf, res...)
			}
		}

		if path, err := fs.JoinPath(themeDist, "style.norm.min.css"); err == nil {
			os.WriteFile(path, res, 0755)
		}
	}
}

// Compile special functions from go.less
func compileGoLess(buf []byte) []byte {
	buf = regex.Comp(`---COMP_auto_contrast:\s*hsl\(([0-9]+)(?:deg|),\s*([0-9]+(?:%|\.[0-9]+|)),\s*([0-9]+(?:%|\.[0-9]+|))\)(,\s*[\w_\- ]+);?`).RepFunc(buf, func(data func(int) []byte) []byte {
		hsl := [3]float64{}

		// get hsl value
		{
			if h, err := strconv.ParseFloat(string(data(1)), 32); err == nil {
				for h < 0 {
					h += 360
				}
				for h > 360 {
					h -= 360
				}
				if h < 0 {
					h = 0
				}

				hsl[0] = h
			}
	
			if bytes.HasSuffix(data(2), []byte{'%'}) {
				if s, err := strconv.ParseFloat(string(bytes.TrimRight(data(2), "%")), 32); err == nil {
					hsl[1] = s / 100
				}
			}else if s, err := strconv.ParseFloat(string(data(2)), 32); err == nil {
				hsl[1] = s
			}
	
			if bytes.HasSuffix(data(3), []byte{'%'}) {
				if l, err := strconv.ParseFloat(string(bytes.TrimRight(data(3), "%")), 32); err == nil {
					hsl[2] = l / 100
				}
			}else if l, err := strconv.ParseFloat(string(data(3)), 32); err == nil {
				hsl[2] = l
			}
		}

		var elm string
		if len(data(4)) != 0 {
			elm = string(bytes.TrimLeft(data(4), ", "))
		}

		res := regex.JoinBytes(
			`--bg-h: `, hslStringPart(hsl[0], "deg"), `;`,
			`--bg-s: `, hslStringPart(hsl[1], "%"), `;`,
			`--bg-l: `, hslStringPart(hsl[2], "%"), `;`,
		)

		elmList := strings.Split(elm, " ")
		for len(elmList) < 4 {
			elmList = append(elmList, "")
		}

		if contrastRatio([3]float64{0, 1, 1}, hsl) > contrastRatio([3]float64{0, 0, 0}, hsl) {
			if elmList[0] == "text" {
				res = append(res, regex.JoinBytes(
					`--fg-h: var(--color-h);`,
					`--fg-s: @text-light-s;`,
					`--fg-l: @text-light-l;`,
					`--fg-text: var(--text-dark);`,
				)...)
			}else{
				res = append(res, regex.JoinBytes(
					`--fg-h: var(--color-h);`,
					`--fg-s: var(--color-s);`,
					`--fg-l: var(--color-l);`,
					`--fg-text: var(--color-text);`,
				)...)
			}

			res = append(res, regex.JoinBytes(`--text: var(--text-light);`)...)

			if elmList[1] == "text" {
				res = append(res, regex.JoinBytes(`--link: var(--text-light);`)...)
			}else{
				res = append(res, regex.JoinBytes(`--link: var(--link-light);`)...)
			}

			if elmList[2] == "text" {
				res = append(res, regex.JoinBytes(`--heading: var(--text-light);`)...)
			}else{
				res = append(res, regex.JoinBytes(`--heading: var(--heading-light, var(--text-light));`)...)
			}

			if elmList[3] == "text" {
				res = append(res, regex.JoinBytes(`--strongheading: linear-gradient(45deg, hsl(calc(var(--fg-h) + 5), calc(var(--fg-s) - 10%), 55%), hsl(calc(var(--fg-h) - 5), calc(var(--fg-s) + 10%), 70%));`)...)
			}else{
				res = append(res, regex.JoinBytes(`--strongheading: var(--strongheading-light, linear-gradient(45deg, hsl(calc(var(--fg-h) + 5), calc(var(--fg-s) - 10%), 55%), hsl(calc(var(--fg-h) - 5), calc(var(--fg-s) + 10%), 70%)));`)...)
			}

			res = append(res, regex.JoinBytes(
				`&>*{--shadow: var(--shadow-light);}`,
				`--textshadow: var(--textshadow-light);`,
			)...)
		}else{
			if elmList[0] == "text" {
				res = append(res, regex.JoinBytes(
					`--fg-h: var(--color-h);`,
					`--fg-s: @text-dark-s;`,
					`--fg-l: @text-dark-l;`,
					`--fg-text: var(--text-light);`,
				)...)
			}else{
				res = append(res, regex.JoinBytes(
					`--fg-h: var(--color-h);`,
					`--fg-s: var(--color-s);`,
					`--fg-l: var(--color-l);`,
					`--fg-text: var(--color-text);`,
				)...)
			}

			res = append(res, regex.JoinBytes(`--text: var(--text-dark);`)...)

			if elmList[1] == "text" {
				res = append(res, regex.JoinBytes(`--link: var(--text-dark);`)...)
			}else{
				res = append(res, regex.JoinBytes(`--link: var(--link-dark);`)...)
			}

			if elmList[2] == "text" {
				res = append(res, regex.JoinBytes(`--heading: var(--text-dark);`)...)
			}else{
				res = append(res, regex.JoinBytes(`--heading: var(--heading-dark, var(--text-dark));`)...)
			}

			if elmList[3] == "text" {
				res = append(res, regex.JoinBytes(`--strongheading: linear-gradient(45deg, hsl(calc(var(--fg-h) + 5), calc(var(--fg-s) - 10%), 30%), hsl(calc(var(--fg-h) - 5), calc(var(--fg-s) + 10%), 40%));`)...)
			}else{
				res = append(res, regex.JoinBytes(`--strongheading: var(--strongheading-dark, linear-gradient(45deg, hsl(calc(var(--fg-h) + 5), calc(var(--fg-s) - 10%), 30%), hsl(calc(var(--fg-h) - 5), calc(var(--fg-s) + 10%), 40%)));`)...)
			}

			res = append(res, regex.JoinBytes(
				`&>*{--shadow: var(--shadow-dark);}`,
				`--textshadow: var(--textshadow-dark);`,
			)...)
		}

		return res
	})

	//todo: may be able to replace this method with compile config method and go.less file
	// setting text contrast in the config

	return buf
}

// Auto compile config file contrasts
func compileGoLessConfig(buf []byte) []byte {
	// remove deg from numbered vars
	buf = regex.Comp(`:([0-9]+)deg;`).RepStr(buf, []byte(":$1;"))

	primary := [3]float64{}
	accent := [3]float64{}
	warn := [3]float64{}

	// fix color hue distances
	{
		var colorPrimary colorful.Color
		buf = regex.Comp(`@primary-h:\s*([0-9]+);`).RepFunc(buf, func(data func(int) []byte) []byte {
			if f, err := strconv.ParseFloat(string(data(1)), 32); err == nil {
				primary[0] = f
			}else{
				return data(0)
			}
	
			for primary[0] < 0 {
				primary[0] += 360
			}
			for primary[0] > 360 {
				primary[0] -= 360
			}
			if primary[0] < 0 {
				primary[0] = 0
			}
	
			colorPrimary = colorful.Hsl(float64(primary[0]), 100, 100)
	
			return regex.JoinBytes(`@primary-h: `, hslStringPart(primary[0], "deg"), ';')
		})
	
		var colorAccent colorful.Color
		buf = regex.Comp(`@accent-h:\s*([0-9]+)(?:deg|);`).RepFunc(buf, func(data func(int) []byte) []byte {
			if f, err := strconv.ParseFloat(string(data(1)), 32); err == nil {
				accent[0] = f
			}else{
				return data(0)
			}
	
			for accent[0] < 0 {
				accent[0] += 360
			}
			for accent[0] > 360 {
				accent[0] -= 360
			}
			if accent[0] < 0 {
				accent[0] = 0
			}
	
			colorAccent = colorful.Hsl(float64(accent[0]), 100, 100)
	
			if !goutil.IsZeroOfUnderlyingType(colorPrimary) {
				if colorPrimary.AlmostEqualRgb(colorAccent) {
					accent[0] = primary[0]
					colorAccent = colorPrimary
				}
			}
	
			return regex.JoinBytes(`@accent-h: `, hslStringPart(accent[0], "deg"), ';')
		})
	
		buf = regex.Comp(`@warn-h:\s*([0-9]+)(?:deg|);`).RepFunc(buf, func(data func(int) []byte) []byte {
			if f, err := strconv.ParseFloat(string(data(1)), 32); err == nil {
				warn[0] = f
			}else{
				return data(0)
			}
	
			for warn[0] < 0 {
				warn[0] += 360
			}
			for warn[0] > 360 {
				warn[0] -= 360
			}
			if warn[0] < 0 {
				warn[0] = 0
			}
	
			if !goutil.IsZeroOfUnderlyingType(colorPrimary) || !goutil.IsZeroOfUnderlyingType(colorAccent) {
				if goutil.IsZeroOfUnderlyingType(colorPrimary) {
					colorPrimary = colorAccent
				}else if goutil.IsZeroOfUnderlyingType(colorAccent) {
					colorAccent = colorPrimary
				}
	
				// fix warn distance from primary
				warnFix := float64(0)
				for colorful.Hsl(float64(warn[0]), 100, 100).AlmostEqualRgb(colorPrimary) {
					if warnFix == 0 {
						if accent[0] < primary[0] {
							warnFix = 20
						}else if accent[0] > primary[0] {
							warnFix = -20
						}else if warn[0] > primary[0] {
							warnFix = 20
						}else{
							warnFix = -20
						}
					}
					warn[0] += warnFix
				}
	
				// ensure warn color is between 0-360
				for warn[0] > 360 {
					warn[0] -= 360
				}
				for warn[0] < 0 {
					warn[0] += 360
				}
				if warn[0] < 0 {
					warn[0] = 0
				}else if warn[0] > 360 {
					warn[0] = 360
				}
	
				if accent[0] != primary[0] {
					for colorful.Hsl(float64(warn[0]), 100, 100).AlmostEqualRgb(colorAccent) {
						if warnFix == 0 {
							if primary[0] < accent[0] {
								warnFix = 20
							}else if primary[0] > accent[0] {
								warnFix = -20
							}else if warn[0] > accent[0] {
								warnFix = 20
							}else{
								warnFix = -20
							}
						}
						warn[0] += warnFix
					}
	
					// ensure warn color is between 0-360
					for warn[0] > 360 {
						warn[0] -= 360
					}
					for warn[0] < 0 {
						warn[0] += 360
					}
					if warn[0] < 0 {
						warn[0] = 0
					}else if warn[0] > 360 {
						warn[0] = 360
					}
				}
	
				// ensure warn is between 0-64 or 264-360 (red)
				if warn[0] < 180 && warn[0] > 64 {
					warn[0] = 64
				}else if warn[0] >= 180 && warn[0] < 264 {
					warn[0] = 264
				}
	
				tryColors := []float64{
					0, // red
					270, // purple
					32, // orange
					330, // pink
					54, // yellow
					300, // light purple
					0, // red
				}
				tryColorsIndex := 0
				for colorful.Hsl(float64(warn[0]), 100, 100).AlmostEqualRgb(colorPrimary) || colorful.Hsl(float64(warn[0]), 100, 100).AlmostEqualRgb(colorAccent) {
					warn[0] = tryColors[tryColorsIndex]
					tryColorsIndex++
					if tryColorsIndex >= len(tryColors) {
						break
					}
				}
			}
	
			return regex.JoinBytes(`@warn-h: `, hslStringPart(warn[0], "deg"), ';')
		})
	}

	// get other color hsl values
	{
		buf = regex.Comp(`@primary-([sl]):\s*([0-9]+(?:%|\.[0-9]+|));`).RepFunc(buf, func(data func(int) []byte) []byte {
			pf := 's'
			i := 1
			if len(data(1)) != 0 && data(1)[0] == 'l' {
				pf = 'l'
				i = 2
			}
	
			if bytes.HasSuffix(data(2), []byte{'%'}) {
				if f, err := strconv.ParseFloat(string(bytes.TrimRight(data(2), "%")), 32); err == nil {
					primary[i] = f / 100
				}
			}else if f, err := strconv.ParseFloat(string(data(2)), 32); err == nil {
				primary[i] = f
			}
	
			if pf == 's' {
				if primary[1] < 0.15 {
					primary[1] = 0.15
				}
			}
	
			return regex.JoinBytes(`@primary-`, pf, `: `, hslStringPart(primary[i], "%"), ';')
		})
	
		buf = regex.Comp(`@accent-([sl]):\s*([0-9]+(?:%|\.[0-9]+|));`).RepFunc(buf, func(data func(int) []byte) []byte {
			pf := 's'
			i := 1
			if len(data(1)) != 0 && data(1)[0] == 'l' {
				pf = 'l'
				i = 2
			}
	
			if bytes.HasSuffix(data(2), []byte{'%'}) {
				if f, err := strconv.ParseFloat(string(bytes.TrimRight(data(2), "%")), 32); err == nil {
					accent[i] = f / 100
				}
			}else if f, err := strconv.ParseFloat(string(data(2)), 32); err == nil {
				accent[i] = f
			}
	
			if pf == 's' {
				if accent[1] < 0.15 {
					accent[1] = 0.15
				}
			}
	
			return regex.JoinBytes(`@accent-`, pf, `: `, hslStringPart(accent[i], "%"), ';')
		})
	
		buf = regex.Comp(`@warn-([sl]):\s*([0-9]+(?:%|\.[0-9]+|));`).RepFunc(buf, func(data func(int) []byte) []byte {
			pf := 's'
			i := 1
			if len(data(1)) != 0 && data(1)[0] == 'l' {
				pf = 'l'
				i = 2
			}
	
			if bytes.HasSuffix(data(2), []byte{'%'}) {
				if f, err := strconv.ParseFloat(string(bytes.TrimRight(data(2), "%")), 32); err == nil {
					warn[i] = f / 100
				}
			}else if f, err := strconv.ParseFloat(string(data(2)), 32); err == nil {
				warn[i] = f
			}
	
			if pf == 's' {
				if warn[1] < 0.15 {
					warn[1] = 0.15
				}
			}
	
			return regex.JoinBytes(`@warn-`, pf, `: `, hslStringPart(warn[i], "%"), ';')
		})
	}

	// get best color text contrast
	{
		buf = regex.Comp(`@primary-text:.*?;`).RepFunc(buf, func(data func(int) []byte) []byte {
			if contrastRatio([3]float64{0, 1, 1}, primary) > contrastRatio([3]float64{0, 0, 0}, primary) {
				return []byte(`@primary-text: var(--text-light);`)
			}
			return []byte(`@primary-text: var(--text-dark);`)
		})
	
		buf = regex.Comp(`@accent-text:.*?;`).RepFunc(buf, func(data func(int) []byte) []byte {
			if contrastRatio([3]float64{0, 1, 1}, accent) > contrastRatio([3]float64{0, 0, 0}, accent) {
				return []byte(`@accent-text: var(--text-light);`)
			}
			return []byte(`@accent-text: var(--text-dark);`)
		})
	
		buf = regex.Comp(`@warn-text:.*?;`).RepFunc(buf, func(data func(int) []byte) []byte {
			if contrastRatio([3]float64{0, 1, 1}, warn) > contrastRatio([3]float64{0, 0, 0}, warn) {
				return []byte(`@warn-text: var(--text-light);`)
			}
			return []byte(`@warn-text: var(--text-dark);`)
		})
	}

	// fix dark and light colors
	{
		var darkL float64
		var lightL float64
		var darkLdark float64
		var lightLdark float64
		buf = regex.Comp(`@(dark|light)-([sl])(-dark|):\s*([0-9]+(?:%|\.[0-9]+|));`).RepFunc(buf, func(data func(int) []byte) []byte {
			pf := 's'
			if len(data(2)) != 0 && data(2)[0] == 'l' {
				pf = 'l'
			}
			
			sl := float64(-1)
			if bytes.HasSuffix(data(4), []byte{'%'}) {
				if f, err := strconv.ParseFloat(string(bytes.TrimRight(data(4), "%")), 32); err == nil {
					sl = f / 100
				}
			}else if f, err := strconv.ParseFloat(string(data(4)), 32); err == nil {
				sl = f
			}
	
			if sl == -1 {
				return data(0)
			}
	
			if pf == 's' && sl > 0.15 {
				sl = 0.15
			}else if bytes.Equal(data(1), []byte("dark")) && sl > 0.49 {
				sl = 0.49
				if len(data(3)) != 0 {
					if sl > 0.3 {
						sl = 0.3
					}
					darkLdark = sl
				}else{
					darkL = sl
				}
			}else if bytes.Equal(data(1), []byte("light")) && sl < 0.5 {
				sl = 0.5
				if len(data(3)) != 0 {
					if sl < 0.15 {
						sl = 0.15
					}else if sl > 0.45 {
						sl = 0.45
					}
					lightLdark = sl
				}else{
					lightL = sl
				}
			}
	
			return regex.JoinBytes(`@`, data(1), `-`, pf, data(3), `: `, hslStringPart(sl, "%"), ';')
		})
	
		if lightL - darkL < 0.15 {
			buf = regex.Comp(`@light-l:\s*([0-9]+(?:%|\.[0-9]+|));`).RepFunc(buf, func(data func(int) []byte) []byte {
				l := float64(-1)
				if bytes.HasSuffix(data(1), []byte{'%'}) {
					if f, err := strconv.ParseFloat(string(bytes.TrimRight(data(1), "%")), 32); err == nil {
						l = f / 100
					}
				}else if f, err := strconv.ParseFloat(string(data(1)), 32); err == nil {
					l = f
				}
	
				if l == -1 {
					return data(0)
				}
	
				return regex.JoinBytes(`@light-l: `, hslStringPart(l + (lightL - darkL), "%"), ';')
			})
		}
	
		if lightLdark - darkLdark < 0.15 {
			buf = regex.Comp(`@light-l-dark:\s*([0-9]+(?:%|\.[0-9]+|));`).RepFunc(buf, func(data func(int) []byte) []byte {
				l := float64(-1)
				if bytes.HasSuffix(data(1), []byte{'%'}) {
					if f, err := strconv.ParseFloat(string(bytes.TrimRight(data(1), "%")), 32); err == nil {
						l = f / 100
					}
				}else if f, err := strconv.ParseFloat(string(data(1)), 32); err == nil {
					l = f
				}
	
				if l == -1 {
					return data(0)
				}
	
				return regex.JoinBytes(`@light-l-dark: `, hslStringPart(l + (lightLdark - darkLdark), "%"), ';')
			})
		}
	}

	return buf
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
				res = append(append(buf), res...)
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
						res = append(res)
					}
				}
			}
		}
	}

	/* if path, err := fs.JoinPath(themeDir, "script.js"); err == nil {
		if buf, err := os.ReadFile(path); err == nil {
			res = append(res, buf...)
			res = append(res)
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


func contrastRatio(fg [3]float64, bg [3]float64) int16 {
	fgColor := colorful.Hsl(fg[0], fg[1] / 100, fg[2] / 100)
	bgColor := colorful.Hsl(bg[0], bg[1] / 100, bg[2] / 100)
	return int16(math.Max(fgColor.DistanceRgb(bgColor), fgColor.DistanceLuv(bgColor)) * 100)
}

func hslStringPart(f float64, t string) string {
	if t == "deg" {
		return strconv.FormatFloat(f, 'f', 0, 32) + "deg"
	}else if t == "%" {
		return strconv.FormatFloat(f*100, 'f', 0, 32) + "%"
	}

	return strconv.FormatFloat(f, 'f', 2, 32)
}
