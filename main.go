package main

import (
	"encoding/hex"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

const (
	defaultWidth = 300
	defaultFg    = "ffffff"
	defaultBg    = "969696"
	appName      = "go-imgplaceholder"
	version      = "0.1"
)

var (
	templates = template.Must(template.ParseFiles("home.html"))
)

type canvas struct {
	img *image.RGBA
	w   int
	h   int
	fg  color.RGBA
	bg  color.RGBA
}

type data struct {
	Title string
	Req   *http.Request
}

func (c *canvas) fillBackground() {
	draw.Draw(c.img, c.img.Bounds(), &image.Uniform{c.bg}, image.ZP, draw.Src)
}

func (c *canvas) fillString(s string) {
	drawString(c, s)
}

func main() {
	r := mux.NewRouter()

	// Fixes for no trailing slash
	r.HandleFunc("/{width:[0-9]+}", imgHandler)
	r.HandleFunc("/{width:[0-9]+}x{height:[0-9]+}", imgHandler)

	// Subrouter with image's width specified
	s1 := r.PathPrefix("/{width:[0-9]+}").Subrouter()
	subroute(s1)

	// Subrouter image's width x height specified
	s2 := r.PathPrefix("/{width:[0-9]+}x{height:[0-9]+}").Subrouter()
	subroute(s2)

	// Home
	r.HandleFunc("/", home)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func subroute(s *mux.Router) {
	s.HandleFunc("/", imgHandler)
	s.HandleFunc("/{bg:[a-fA-F0-9]+}", imgHandler)
	s.HandleFunc("/{bg:[a-fA-F0-9]+}/{fg:[a-fA-F0-9]+}", imgHandler)
}

func home(w http.ResponseWriter, r *http.Request) {
	d := &data{
		Title: fmt.Sprintf("%s/%s", appName, version),
		Req:   r,
	}

	if err := templates.ExecuteTemplate(w, "home.html", d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func imgHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var val, fgs, bgs string
	var ok bool

	width := defaultWidth
	if val, ok = vars["width"]; ok {
		width, _ = strconv.Atoi(val)
	}
	height := width
	if val, ok = vars["height"]; ok {
		height, _ = strconv.Atoi(val)
	}

	if fgs, ok = vars["fg"]; !ok {
		fgs = defaultFg
	}
	if bgs, ok = vars["bg"]; !ok {
		bgs = defaultBg
	}

	c := &canvas{
		img: image.NewRGBA(image.Rect(0, 0, width, height)),
		w:   width,
		h:   height,
		fg:  stringToColor(fgs),
		bg:  stringToColor(bgs),
	}
	c.fillBackground()
	c.fillString(fmt.Sprintf("%dx%d", width, height))

	err := png.Encode(w, c.img)
	if err != nil {
		errStr := fmt.Sprintf("Error encoding: %v", err)
		http.Error(w, errStr, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "image/png")
}

func stringToColor(s string) color.RGBA {
	switch len(s) {
	case 1:
		s = strings.Repeat(s, 6)
	case 2:
		s = strings.Repeat(s, 3)
	case 3:
		s = strings.Repeat(s, 2)
	case 4:
		fallthrough
	case 5:
		s = s[:4] + "00"
	}

	h, _ := hex.DecodeString(s)

	return color.RGBA{h[0], h[1], h[2], 255}
}
