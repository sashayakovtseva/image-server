package main

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

var (
	ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body>
<img src="data:image/jpg;base64,{{.Image}}">
</body>`

	CatsTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body>
{{range .}}<img src="data:image/jpg;base64,{{.}}">{{end}}
</body>`

	CatsGoodTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body>
{{range .}}<img src="{{.}}">{{end}}
</body>`
)

const (
	MAX_IMG      = math.MaxInt32
	IMG_BASE_DIR = "./images/"
)

func main() {
	http.HandleFunc("/blue/", blueHandler)
	http.HandleFunc("/red/", redHandler)
	http.HandleFunc("/cats/", catsHandler)
	http.HandleFunc("/cats/good/", catsGoodHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func blueHandler(w http.ResponseWriter, r *http.Request) {
	m := image.NewRGBA(image.Rect(0, 0, 240, 240))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	var img image.Image = m
	writeImage(w, &img)
}

func redHandler(w http.ResponseWriter, r *http.Request) {
	m := image.NewRGBA(image.Rect(0, 0, 240, 240))
	blue := color.RGBA{255, 0, 0, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	var img image.Image = m
	writeImageWithTemplate(w, &img)
}

func catsGoodHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(IMG_BASE_DIR)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	n := MAX_IMG
	if len(files) < MAX_IMG {
		n = len(files)
	}

	imgs := make([]string, n)
	for i := 0; i < n; i++ {
		imgs[i] = IMG_BASE_DIR[1:] + files[i].Name()
	}

	if tmpl, err := template.New("catsGood").Parse(CatsGoodTemplate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		if err = tmpl.Execute(w, imgs); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func catsHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(IMG_BASE_DIR)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	n := MAX_IMG
	if len(files) < MAX_IMG {
		n = len(files)
	}
	imgs := make([]image.Image, n)
	for i := 0; i < n; i++ {
		f, err := os.Open(IMG_BASE_DIR + files[i].Name())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		imgs[i], _, err = image.Decode(f)
		f.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	writeImagesWithTemplate(w, imgs)
}

// Writeimagewithtemplate encodes an image 'img' in jpeg format and writes it into ResponseWriter using a template.
func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Fatalln("unable to encode image.")
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": str}
		if err = tmpl.Execute(w, data); err != nil {
			log.Println("unable to execute template.")
		}
	}
}

func writeImagesWithTemplate(w http.ResponseWriter, imgs []image.Image) {
	b64imgs := make([]string, len(imgs))
	for i, img := range imgs {
		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, img, nil); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		b64imgs[i] = base64.StdEncoding.EncodeToString(buffer.Bytes())
	}
	if tmpl, err := template.New("cats").Parse(CatsTemplate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		if err = tmpl.Execute(w, b64imgs); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
