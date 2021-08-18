package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type TemplateCtx struct {
	Controller interface{}
	BasePath   string
}

func Render(w http.ResponseWriter, file string, basefile string, data TemplateCtx) {
	data.BasePath = Config.BasePath
	if basefile == "" {
		tmpl := template.Must(template.ParseFiles("templates/" + file))
		err := tmpl.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	} else {
		tmpl := template.Must(template.ParseFiles("templates/"+file, "templates/"+basefile))
		err := tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func JsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	type errorInfo struct {
		Code    int
		Message string
	}
	if !Config.Debug {
		if code == 400 {
			message = "Bad Request"
		} else {
			message = "Internal Server Error\n Please try again later"
		}
	}

	w.WriteHeader(code)
	Render(w, "error.html", "base.html", TemplateCtx{Controller: errorInfo{code, message}})
}

func JsonError(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	type errorInfo struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	w.WriteHeader(code)
	JsonResponse(w, errorInfo{code, message})
}
