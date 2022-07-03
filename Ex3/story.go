package storybook

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template
var renderStoryTemplate = `
<!doctype html>
<html>
    <head>
        <title>Our Funky HTML Page</title>
        <meta name="description" content="Go lang excerise">
        <meta name="keywords" content="html tutorial template">
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Stories}}
        <p>{{.}}</p>
        {{end}}
        <ol>
        {{range .Options}}
        <li><a href="/{{.Arc}}">{{.Text}}</a></li>
        {{end}}
        </ol>
    </body>
</html>
`

func init() {
	tpl = template.Must(template.New("").Parse(renderStoryTemplate))
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	if arc, ok := h.s[path]; ok {
		err := tpl.Execute(w, arc)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}
	log.Println("Chapter is not found")
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

func StoryHandler(s Story) http.Handler {
	return handler{s}
}

func JSONStory(r io.Reader) (Story, error) {
	var story Story
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type handler struct {
	s Story
}

type Story map[string]Arc

type Arc struct {
	Title   string   `json:"title"`
	Stories []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
