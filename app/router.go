package app

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseGlob("./templates/*.tmpl"))

var views = []ViewRoute{
	{templateName: "index", pattern: regexp.MustCompile(`^/$`)},
	{templateName: "lot", pattern: regexp.MustCompile(`^/lot/(?P<lotId>[\d]+)$`)},
	{templateName: "about", pattern: regexp.MustCompile(`^/about$`)},
	{templateName: "catalog", pattern: regexp.MustCompile(`^/catalog/(?P<sectionId>[\d]+)$`)},
	{templateName: "catalog_add", pattern: regexp.MustCompile(`^/catalog/add$`)},
}

var actions = map[string]ActionRoute{
	"addLot": {
		name: "addLot",
		handlerFunc: func(r ActionParams) (ResponseData, error) {
			log.Println("addLot action processing")
			return ResponseData{}, nil
		},
	},
}

type ResponseData struct {
	data map[string]interface{}
}

type ViewRoute struct {
	pattern      *regexp.Regexp
	templateName string
	handlerFunc  func(*http.Request) (ResponseData, error)
}

type ActionRoute struct {
	name        string
	handlerFunc func(ActionParams) (ResponseData, error)
}

type ActionRequest map[string]ActionParams
type ActionParams map[string]interface{}

func AddView(view ViewRoute) {
	views = append(views, view)
	log.Printf("View handler added for [%s]", view.templateName)
}

func AddAction(action ActionRoute) {
	delete(actions, action.name)
	actions[action.name] = action
	log.Printf("Action handler added for [%s]", action.name)
}

func MainHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" { // processing template for GET request

		urlString := r.URL.String()
		log.Printf("Processing GET request url: %s", urlString)

		found := false
		for _, view := range views {
			if view.pattern.MatchString(urlString) {
				view.PageRespond(w, r)
				found = true
				break
			}
		}

		if !found {
			log.Println("Not found")
			http.Error(w, "Not found", http.StatusNotFound)
			err := templates.ExecuteTemplate(w, "404.tmpl", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	} else { // POST request

		if r.Header.Get("X-Request") == "JSON" && r.Body != nil {
			// 'action' request
			jsonResponse, err := ActionRespond(r)
			if err != nil {
				log.Println("Error processing actions")
				http.Error(w, "", http.StatusInternalServerError)
			} else {
				json.NewEncoder(w).Encode(jsonResponse)
			}

		} else {
			http.Error(w, "Wrong request format", http.StatusBadRequest)
		}

	}
}

func (view *ViewRoute) PageRespond(w http.ResponseWriter, r *http.Request) {
	pageData := ResponseData{}

	log.Printf("View [%s] selected", view.templateName)

	if view.handlerFunc != nil {
		response, err := view.handlerFunc(r)
		if err != nil {
			log.Printf("Error processing view [%s]", view.templateName)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		pageData = response
	}

	w.Header().Add("Content-Type", "text/html")
	err := templates.ExecuteTemplate(w, view.templateName+".tmpl", pageData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ActionRespond(r *http.Request) (ResponseData, error) {
	actionRequest := ActionRequest{}
	response := ResponseData{}
	err := json.NewDecoder(r.Body).Decode(&actionRequest)
	if err != nil {
		return response, err
	}

	// processing all actions
	for actionName, actionParams := range actionRequest {

		log.Println(actionName)

		if actionHanlder, ok := actions[actionName]; ok {
			actionHanlder.handlerFunc(actionParams)
		}
	}

	return response, nil
}
