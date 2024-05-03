package forum

import "html/template"

var Error = template.Must(template.ParseFiles("./src/templates/error.html"))
var ErrorUser = template.Must(template.ParseFiles("./src/templates/errorUser.html"))
var Login = template.Must(template.ParseFiles("./src/templates/login.html"))
var LoginError = template.Must(template.ParseFiles("./src/templates/loginerror.html"))
var Home = template.Must(template.ParseFiles("./src/templates/home.html"))
var Profile = template.Must(template.ParseFiles("./src/templates/profile.html"))
