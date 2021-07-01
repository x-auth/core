package helpers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func Render(w http.ResponseWriter, file string, basefile string, data interface{}) {
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

func Error(w http.ResponseWriter, code int, message string) {
	log.Println(message)
	type errorInfo struct {
		Code    int
		Message string
	}
	Render(w, "error.html", "base.html", errorInfo{code, message})
}
