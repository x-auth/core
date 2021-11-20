package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type TemplateCtx struct {
	Controller interface{}
	BasePath   string
}

func mainLang(lang string) string {
	if strings.Contains(lang, ",") {
		return strings.Split(lang, ",")[0]
	} else {
		return lang
	}
}

func Render(w http.ResponseWriter, lang string, file string, basefile string, data TemplateCtx) {
	fmt.Println(lang)
	data.BasePath = Config.BasePath
	if basefile == "" {
		tmpl := template.Must(template.ParseFiles("templates/" + mainLang(lang) + "/" + file))
		err := tmpl.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	} else {
		tmpl := template.Must(template.ParseFiles("templates/"+mainLang(lang)+"/"+file, "templates/"+basefile))
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
	Render(w, "en", "error.html", "base.html", TemplateCtx{Controller: errorInfo{code, message}})
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
