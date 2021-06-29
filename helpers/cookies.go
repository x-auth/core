package helpers

import "net/http"

func GetCookie(request *http.Request,name string)string{
	for _, cookie := range  request.Cookies(){
		if cookie.Name == name{
			return cookie.Value
		}
	}
	return ""
}
