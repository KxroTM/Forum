package forum

import "html/template"

// LIGHT MODE
var Error = template.Must(template.ParseFiles("./src/templates/light/error.html"))
var ErrorUser = template.Must(template.ParseFiles("./src/templates/light/errorUser.html"))
var Login = template.Must(template.ParseFiles("./src/templates/light/login.html"))
var LoginError = template.Must(template.ParseFiles("./src/templates/light/loginerror.html"))
var Home = template.Must(template.ParseFiles("./src/templates/light/home.html"))
var Profile = template.Must(template.ParseFiles("./src/templates/light/profile.html"))

// DARK MODE
var DarkError = template.Must(template.ParseFiles("./src/templates/dark/error.html"))
var DarkErrorUser = template.Must(template.ParseFiles("./src/templates/dark/errorUser.html"))
var DarkLogin = template.Must(template.ParseFiles("./src/templates/dark/login.html"))
var DarkLoginError = template.Must(template.ParseFiles("./src/templates/dark/loginerror.html"))
var DarkHome = template.Must(template.ParseFiles("./src/templates/dark/home.html"))
var DarkProfile = template.Must(template.ParseFiles("./src/templates/dark/profile.html"))
