package forum

import "html/template"

// LIGHT MODE
var Error404 = template.Must(template.ParseFiles("./src/templates/light/error.html"))
var ErrorUser = template.Must(template.ParseFiles("./src/templates/light/errorUser.html"))
var Login = template.Must(template.ParseFiles("./src/templates/light/login.html"))
var LoginError = template.Must(template.ParseFiles("./src/templates/light/loginerror.html"))
var Register = template.Must(template.ParseFiles("./src/templates/light/register.html"))
var RegisterError = template.Must(template.ParseFiles("./src/templates/light/registererror.html"))
var Home = template.Must(template.ParseFiles("./src/templates/light/accueil.html"))
var Profile = template.Must(template.ParseFiles("./src/templates/light/profile.html"))

// DARK MODE
var DarkError = template.Must(template.ParseFiles("./src/templates/dark/error.html"))
var DarkErrorUser = template.Must(template.ParseFiles("./src/templates/dark/errorUser.html"))
var DarkLogin = template.Must(template.ParseFiles("./src/templates/dark/login.html"))
var DarkLoginError = template.Must(template.ParseFiles("./src/templates/dark/loginerror.html"))
var DarkHome = template.Must(template.ParseFiles("./src/templates/dark/accueil.html"))
var DarkProfile = template.Must(template.ParseFiles("./src/templates/dark/profile.html"))
var DarkRegister = template.Must(template.ParseFiles("./src/templates/dark/register.html"))
var DarkRegisterError = template.Must(template.ParseFiles("./src/templates/dark/registererror.html"))
