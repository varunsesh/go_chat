package views

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {

	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, create new empty Page and save it
			p := &Page{Title: title, Body: []byte{}}
			if err := p.save(); err != nil {
				return nil, err
			}
			return p, nil

		}
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	log.Println("Viewing page:", title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Chat Application Server!")
}
