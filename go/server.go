package forum

import "net/http"

var ResetPasswordMap = make(map[string]string)
var URL string

func CreateRoute(w http.ResponseWriter, r *http.Request, url string) {
	URL = url + "/"
	ResetPasswordMap[URL] = "valid"
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	IPsLog(clientIP + "  ==>  " + r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	p := "Page not found"
	err := Error.ExecuteTemplate(w, "error.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
