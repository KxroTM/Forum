package forum

import "html/template"

var Error = template.Must(template.ParseFiles("./src/templates/error.html"))
var Login = template.Must(template.ParseFiles("./src/templates/login.html"))
var LoginError = template.Must(template.ParseFiles("./src/templates/loginerror.html"))
