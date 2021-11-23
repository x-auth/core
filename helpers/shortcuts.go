/*
 * Copyright (c) 2021 X-Net Services GmbH
 * Info: https://x-net.at
 *
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"text/template"
	"x-net.at/idp/logger"
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
			logger.Warning.Println(err)
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
