package app

import (
	"log"
	"net/http"
)

type App struct {
	//render  Renderer
	//storage Strorage
	config Config
}

type Config struct {
	addr string
}

//type Page struct {
//	Title string
//	Body  []string
//}

func Init() (*App, error) {

	app := &App{}

	// check configuration
	app.config = Config{addr: ":8080"}

	// load handlers
	//InitHandlers()
	//http.HandleFunc("/", MainHandler)

	// init storage
	return app, nil
}

func (app *App) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", MainHandler)

	jsFiles := http.FileServer(http.Dir("./static/js"))
	jsStyles := http.FileServer(http.Dir("./static/css"))
	imgFiles := http.FileServer(http.Dir("./static/img"))

	mux.Handle("/js/", http.StripPrefix("/js", jsFiles))
	mux.Handle("/css/", http.StripPrefix("/css", jsStyles))
	mux.Handle("/img/", http.StripPrefix("/img", imgFiles))

	log.Fatal(http.ListenAndServe(app.config.addr, mux))
}

//func renderTemplate(w http.ResponseWriter, tmpl string, p Page) {
//	err := templates.ExecuteTemplate(w, tmpl+".html", p)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//	}
//}
