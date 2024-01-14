package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/AspieSoft/go-regex-re2/v2"
	"github.com/AspieSoft/goutil/fs/v2"
	"github.com/AspieSoft/goutil/v7"
	"github.com/lucasb-eyer/go-colorful"
	"gopkg.in/yaml.v2"
)

type cssConfig struct {
	Colors cssConfigColors
	Text cssConfigText
	Link cssConfigText
	Heading cssConfigText
	StrongHeading cssConfigTextGradient
	StrongHeadingAccent cssConfigTextGradient

	Primary cssConfigColorType
	Accent cssConfigColorType
	Warn cssConfigColorType
	Dark cssConfigColorType
	Light cssConfigColorType

	Base cssConfigColorType
	Card cssConfigColorType

	Header cssConfigColorType
	HeaderImg cssConfigColorType
	Footer cssConfigColorType
}

type cssConfigColors struct {
	Primary [4]int16
	Accent [4]int16
	Warn [4]int16
	Dark [4]int16
	Light [4]int16
}

type cssConfigText struct {
	Dark [4]int16
	Light [4]int16
}

type cssConfigTextGradient struct {
	Dark string
	Light string
	Fallback string
}

type cssConfigColorType struct {
	Color cssConfigText
	Gradient cssConfigTextGradient
	Text string
	Link string
	Heading string
	StrongHeading string
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
	args := goutil.MapArgs()

	port := ""
	if args["port"] != "" {
		if p, err := strconv.Atoi(args["port"]); err == nil && p >= 3000 && p <= 65535 /* golang only accepts 16 bit port numbers */ {
			port = strconv.Itoa(p)
		}
	}else if args["0"] != "" {
		if args["0"] == "test" {
			port = "3000"
		}else if p, err := strconv.Atoi(args["0"]); err == nil && p >= 3000 && p <= 65535 /* golang only accepts 16 bit port numbers */ {
			port = strconv.Itoa(p)
		}
	}else if args["test"] != "" || args["t"] == "true" {
		port = "3000"
	}

	src := "../src"
	if args["src"] != "" {
		src = args["src"]
	}

	dist := "../dist"
	if args["dist"] != "" {
		dist = args["dist"]
	}

	compileTheme(src, dist)

	if port != "" {
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
							if path, err := fs.JoinPath("./dist", strings.Replace(cUrl, "theme/", "", 1)+"."+ext); err == nil {
								if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
									http.ServeFile(w, r, path)
									return
								}
							}
						}

						if path, err := fs.JoinPath("./test", cUrl+"."+ext); err == nil {
							if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
								http.ServeFile(w, r, path)
								return
							}
						}
	
						w.WriteHeader(404)
						w.Write([]byte("error 404"))
						return
					}
				}

				if path, err := fs.JoinPath("./test/html", cUrl[2:]+".html"); err == nil {
					if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
						http.ServeFile(w, r, path)
						return
					}
				}
			}

			w.WriteHeader(404)
			w.Write([]byte("error 404"))
		})
	}
}

func compileTheme(src string, dist string){
	os.RemoveAll(dist)
	os.MkdirAll(dist, 0755)

	compileThemeConfig(src, dist);
}

func compileThemeConfig(src string, dist string){
	config := cssConfig{}

	if path, err := fs.JoinPath(src, "config.yml"); err == nil {
		if file, err := os.ReadFile(path); err == nil {
			if err := yaml.Unmarshal(file, &config); err != nil {
				panic(err)
			}
		}
	}

	// fix color distances
	{
		// fix accent distance from primary
		colorPrimary := colorful.Hsl(float64(config.Colors.Primary[0]), 100, 100)
		if colorPrimary.AlmostEqualRgb(colorful.Hsl(float64(config.Colors.Accent[0]), 100, 100)) {
			config.Colors.Accent[0] = config.Colors.Primary[0]
		}
	
		// ensure warn is between 0-64 or 264-360 (near red)
		if config.Colors.Warn[0] < 180 && config.Colors.Warn[0] > 64 {
			config.Colors.Warn[0] = 64
		}else if config.Colors.Warn[0] >= 180 && config.Colors.Warn[0] < 264 {
			config.Colors.Warn[0] = 264
		}
	
		// fix warn distance from primary
		warnFix := int16(0)
		for colorful.Hsl(float64(config.Colors.Warn[0]), 100, 100).AlmostEqualRgb(colorPrimary) {
			if warnFix == 0 {
				if config.Colors.Accent[0] < config.Colors.Primary[0] {
					warnFix = 20
				}else if config.Colors.Accent[0] > config.Colors.Primary[0] {
					warnFix = -20
				}else if config.Colors.Warn[0] > config.Colors.Primary[0] {
					warnFix = 20
				}else{
					warnFix = -20
				}
			}
			config.Colors.Warn[0] += warnFix
		}
	
		// ensure warn color is between 0-360
		for config.Colors.Warn[0] > 360 {
			config.Colors.Warn[0] -= 360
		}
		for config.Colors.Warn[0] < 0 {
			config.Colors.Warn[0] += 360
		}
		if config.Colors.Warn[0] < 0 {
			config.Colors.Warn[0] = 0
		}else if config.Colors.Warn[0] > 360 {
			config.Colors.Warn[0] = 360
		}
	
		colorAccent := colorful.Hsl(float64(config.Colors.Accent[0]), 100, 100)
		if config.Colors.Accent[0] != config.Colors.Primary[0] {
			for colorful.Hsl(float64(config.Colors.Warn[0]), 100, 100).AlmostEqualRgb(colorAccent) {
				if warnFix == 0 {
					if config.Colors.Primary[0] < config.Colors.Accent[0] {
						warnFix = 20
					}else if config.Colors.Primary[0] > config.Colors.Accent[0] {
						warnFix = -20
					}else if config.Colors.Warn[0] > config.Colors.Accent[0] {
						warnFix = 20
					}else{
						warnFix = -20
					}
				}
				config.Colors.Warn[0] += warnFix
			}
	
			// ensure warn color is between 0-360
			for config.Colors.Warn[0] > 360 {
				config.Colors.Warn[0] -= 360
			}
			for config.Colors.Warn[0] < 0 {
				config.Colors.Warn[0] += 360
			}
			if config.Colors.Warn[0] < 0 {
				config.Colors.Warn[0] = 0
			}else if config.Colors.Warn[0] > 360 {
				config.Colors.Warn[0] = 360
			}
		}
	
		// ensure warn is between 0-64 or 264-360 (red)
		if config.Colors.Warn[0] < 180 && config.Colors.Warn[0] > 64 {
			config.Colors.Warn[0] = 64
		}else if config.Colors.Warn[0] >= 180 && config.Colors.Warn[0] < 264 {
			config.Colors.Warn[0] = 264
		}
	
		tryColors := []int16{
			0, // red
			270, // purple
			32, // orange
			330, // pink
			54, // yellow
			300, // light purple
			0, // red
		}
		tryColorsIndex := 0
		for colorful.Hsl(float64(config.Colors.Warn[0]), 100, 100).AlmostEqualRgb(colorPrimary) || colorful.Hsl(float64(config.Colors.Warn[0]), 100, 100).AlmostEqualRgb(colorAccent) {
			config.Colors.Warn[0] = tryColors[tryColorsIndex]
			tryColorsIndex++
			if tryColorsIndex >= len(tryColors) {
				break
			}
		}

		// fix dark and light hue relative to primary
		// config.Colors.Dark[0] = config.Colors.Primary[0] + clamp(config.Colors.Dark[0], -20, 20)
		config.Colors.Dark[1] = clamp(config.Colors.Dark[1], 0, 15)
		config.Colors.Dark[2] = clamp(config.Colors.Dark[2], 0, 49)

		// config.Colors.Light[0] = config.Colors.Primary[0] + clamp(config.Colors.Light[0], -20, 20)
		config.Colors.Light[1] = clamp(config.Colors.Light[1], 0, 15)
		config.Colors.Light[2] = clamp(config.Colors.Light[2], 50, 100)

		// ensure colors are between 0-360
		for config.Colors.Dark[0] > 360 {
			config.Colors.Dark[0] -= 360
		}
		for config.Colors.Dark[0] < 0 {
			config.Colors.Dark[0] += 360
		}
		if config.Colors.Dark[0] < 0 {
			config.Colors.Dark[0] = 0
		}else if config.Colors.Dark[0] > 360 {
			config.Colors.Dark[0] = 360
		}

		for config.Colors.Light[0] > 360 {
			config.Colors.Light[0] -= 360
		}
		for config.Colors.Light[0] < 0 {
			config.Colors.Light[0] += 360
		}
		if config.Colors.Light[0] < 0 {
			config.Colors.Light[0] = 0
		}else if config.Colors.Light[0] > 360 {
			config.Colors.Light[0] = 360
		}
	}


	// set text colors
	{
		// fix text colors
		// config.Link.Dark[0] = config.Colors.Primary[0] + clamp(config.Text.Dark[0], -20, 20)
		config.Text.Dark[1] = clamp(config.Text.Dark[1], 0, 10)
		config.Text.Dark[2] = clamp(config.Text.Dark[2], 0, 49)

		// config.Link.Light[0] = config.Colors.Primary[0] + clamp(config.Text.Light[0], -20, 20)
		config.Text.Light[1] = clamp(config.Text.Light[1], 0, 10)
		config.Text.Light[2] = clamp(config.Text.Light[2], 50, 100)

		// fix link colors
		config.Link.Dark[0] = clamp(config.Text.Dark[0], 180, 260)
		config.Link.Dark[2] = clamp(config.Text.Dark[2], 0, 49)

		config.Link.Light[0] = clamp(config.Text.Light[0], 180, 260)
		config.Link.Light[2] = clamp(config.Text.Light[2], 50, 100)

		// fix heading colors
		// config.Heading.Dark[0] = config.Colors.Primary[0] + clamp(config.Heading.Dark[0], -20, 20)
		// config.Heading.Light[0] = config.Colors.Primary[0] + clamp(config.Heading.Light[0], -20, 20)

		// fix strong heading colors
		if config.StrongHeading.Dark == "" || config.StrongHeading.Dark == "auto" || config.StrongHeading.Dark == "default" {
			config.StrongHeading.Dark = "linear-gradient(45deg, hsl(calc(var(--color-fg-h) - 10), var(--color-fg-s), 20), hsl(calc(var(--color-fg-h) + 10), var(--color-fg-s), 35))"
		}
		if config.StrongHeading.Light == "" || config.StrongHeading.Light == "auto" || config.StrongHeading.Light == "default" {
			config.StrongHeading.Light = "linear-gradient(45deg, hsl(calc(var(--color-fg-h) - 10), var(--color-fg-s), 75), hsl(calc(var(--color-fg-h) + 10), var(--color-fg-s), 90))"
		}
		if config.StrongHeading.Fallback == "" || config.StrongHeading.Fallback == "auto" || config.StrongHeading.Fallback == "default" {
			config.StrongHeading.Fallback = "linear-gradient(45deg, hsl(calc(var(--color-fg-h) - 10), var(--color-fg-s), 45), hsl(calc(var(--color-fg-h) + 10), var(--color-fg-s), 60))"
		}

		// fix strong heading accent colors
		if config.StrongHeadingAccent.Dark == "" || config.StrongHeadingAccent.Dark == "auto" || config.StrongHeadingAccent.Dark == "default" {
			config.StrongHeadingAccent.Dark = "linear-gradient(45deg, #4b4b4b, #312f2f)"
		}
		if config.StrongHeadingAccent.Light == "" || config.StrongHeadingAccent.Light == "auto" || config.StrongHeadingAccent.Light == "default" {
			config.StrongHeadingAccent.Light = "linear-gradient(45deg, #f3f3f3, #c4bcbc)"
		}
		if config.StrongHeadingAccent.Fallback == "" || config.StrongHeadingAccent.Fallback == "auto" || config.StrongHeadingAccent.Fallback == "default" {
			config.StrongHeadingAccent.Fallback = "linear-gradient(45deg, #bdbdbd, #9b9090)"
		}


		//todo: calculate best text contrast for background colors
		// min contrast 45
		fmt.Println(contrastRatio(config.Text.Dark, config.Colors.Primary), contrastRatio(config.Text.Light, config.Colors.Primary))

		setBestTextColor := func(setVal *string, fg cssConfigText, bg [4]int16, isLink bool) {
			if !(*setVal == "" || *setVal == "auto" || *setVal == "default") {
				return
			}
			
			textColor := [4]int16{}
			step := int16(0)
			if contrastRatio(fg.Dark, bg) > contrastRatio(fg.Light, bg) {
				textColor = fg.Dark
				step = 10
			}else{
				textColor = fg.Light
				step = 10
			}
	
			for contrastRatio(textColor, bg) < 45 && textColor[1] > 0 && textColor[1] < 100 {
				textColor[1] += step
			}
	
			tryColors := [][4]int16{
				{0, 0, 6, 100}, // black
				{0, 0, 100, 100}, // white
				{0, 0, 69, 100}, // grey
			}

			if isLink {
				tryColors = [][4]int16{
					{204, 77, 45, 100}, // dark
					{197, 81, 58, 100}, // light
					{199, 70, 42, 100}, // middle
				}
			}

			tryColorsIndex := 0
			for contrastRatio(textColor, bg) < 45 {
				textColor = tryColors[tryColorsIndex]
				tryColorsIndex++
				if tryColorsIndex >= len(tryColors) {
					break
				}
			}

			*setVal = colorful.Hsl(float64(textColor[0]), float64(textColor[1]), float64(textColor[2])).Hex()
		}

		setBestColors := func(setVal *cssConfigColorType, col [4]int16, blackAndWhiteColor bool){
			setBestTextColor(&(*setVal).Text, config.Text, col, false)
			setBestTextColor(&(*setVal).Link, config.Link, col, true)
			setBestTextColor(&(*setVal).Heading, config.Heading, col, false)
			setBestTextColor(&(*setVal).StrongHeading, cssConfigText{Dark: [4]int16{0, 0, 6, 100}, Light: [4]int16{0, 0, 100, 100}}, col, false)
			if (*setVal).StrongHeading == "#0f0f0f" {
				if blackAndWhiteColor {

				}else{
					(*setVal).StrongHeading = config.StrongHeading.Dark
				}
			}else if (*setVal).StrongHeading == "#ffffff" {
				if blackAndWhiteColor {

				}else{
					(*setVal).StrongHeading = config.StrongHeading.Light
				}
			}else{
				if blackAndWhiteColor {

				}else{
					(*setVal).StrongHeading = config.StrongHeading.Fallback
				}
			}
		}

		setBestColors(&config.Primary, config.Colors.Primary, false)
		setBestColors(&config.Accent, config.Colors.Accent, false)
		setBestColors(&config.Warn, config.Colors.Warn, false)
	}


	res := []byte(":root {\n")

	//todo: write css to result
	res = append(res, regex.JoinBytes(
		"  --primary-h: ", toString(clamp(config.Colors.Primary[0], 0, 360)), ";\n",
		"  --primary-s: ", toString(clamp(config.Colors.Primary[1], 0, 100)), "%;\n",
		"  --primary-l: ", toString(clamp(config.Colors.Primary[2], 0, 100)), "%;\n",
		"  --primary-a: ", toString(clamp(config.Colors.Primary[3], 0, 100) / 100), ";\n",
		'\n',
	)...)

	res = append(res, regex.JoinBytes(
		"  --accent-h: ", toString(clamp(config.Colors.Accent[0], 0, 360)), ";\n",
		"  --accent-s: ", toString(clamp(config.Colors.Accent[1], 0, 100)), "%;\n",
		"  --accent-l: ", toString(clamp(config.Colors.Accent[2], 0, 100)), "%;\n",
		"  --accent-a: ", toString(clamp(config.Colors.Accent[3], 0, 100) / 100), ";\n",
		'\n',
	)...)

	res = append(res, regex.JoinBytes(
		"  --warn-h: ", toString(clamp(config.Colors.Warn[0], 0, 360)), ";\n",
		"  --warn-s: ", toString(clamp(config.Colors.Warn[1], 0, 100)), "%;\n",
		"  --warn-l: ", toString(clamp(config.Colors.Warn[2], 0, 100)), "%;\n",
		"  --warn-a: ", toString(clamp(config.Colors.Warn[3], 0, 100) / 100), ";\n",
		'\n',
	)...)

	res = append(res, regex.JoinBytes(
		"  --dark-h: calc(var(--color-h) + ", toString(clamp(config.Colors.Dark[0], -20, 20)), ");\n",
		"  --dark-s: ", toString(clamp(config.Colors.Dark[1], 0, 100)), "%;\n",
		"  --dark-l: ", toString(clamp(config.Colors.Dark[2], 0, 100)), "%;\n",
		"  --dark-a: ", toString(clamp(config.Colors.Dark[3], 0, 100) / 100), ";\n",
		'\n',
	)...)

	res = append(res, regex.JoinBytes(
		"  --light-h: calc(var(--color-h) + ", toString(clamp(config.Colors.Light[0], -20, 20)), ");\n",
		"  --light-s: ", toString(clamp(config.Colors.Light[1], 0, 100)), "%;\n",
		"  --light-l: ", toString(clamp(config.Colors.Light[2], 0, 100)), "%;\n",
		"  --light-a: ", toString(clamp(config.Colors.Light[3], 0, 100) / 100), ";\n",
		'\n',
	)...)

	res = append(res, '}', '\n')

	if path, err := fs.JoinPath(dist, "config.css"); err == nil {
		os.WriteFile(path, res, 0755)
	}
}

func clamp(val int16, min int16, max int16) int16 {
	if val < min {
		val = min
	}else if val > max {
		val = max
	}
	return val
}

func toString(val int16) string {
	return strconv.Itoa(int(val))
}


func contrastRatio(fg [4]int16, bg [4]int16) int16 {
	fgColor := colorful.Hsl(float64(fg[0]), float64(fg[1]) / 100, float64(fg[2]) / 100)
	bgColor := colorful.Hsl(float64(bg[0]), float64(bg[1]) / 100, float64(bg[2]) / 100)
	return int16(math.Max(fgColor.DistanceRgb(bgColor), fgColor.DistanceLuv(bgColor)) * 100)
}
