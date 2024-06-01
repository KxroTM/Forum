package forum

import "html/template"

// LIGHT MODE
var Error404 = template.Must(template.ParseFiles("./src/templates/light/error.html"))
var ErrorUser = template.Must(template.ParseFiles("./src/templates/light/errorUser.html"))
var ErrorPost = template.Must(template.ParseFiles("./src/templates/light/errorPost.html"))
var Login = template.Must(template.ParseFiles("./src/templates/light/login.html"))
var LoginError = template.Must(template.ParseFiles("./src/templates/light/loginerror.html"))
var Register = template.Must(template.ParseFiles("./src/templates/light/register.html"))
var RegisterError = template.Must(template.ParseFiles("./src/templates/light/registererror.html"))
var Home = template.Must(template.ParseFiles("./src/templates/light/accueil.html"))
var Profile = template.Must(template.ParseFiles("./src/templates/light/profile.html"))
var CreatePosts = template.Must(template.ParseFiles("./src/templates/light/createpost.html"))
var CreatePostsError = template.Must(template.ParseFiles("./src/templates/light/createposterror.html"))
var LightPost = template.Must(template.ParseFiles("./src/templates/light/post.html"))
var Populaire = template.Must(template.ParseFiles("./src/templates/light/populaire.html"))
var Posts = template.Must(template.ParseFiles("./src/templates/light/filtragePost.html"))

// DARK MODE
var DarkError = template.Must(template.ParseFiles("./src/templates/dark/error.html"))
var DarkErrorUser = template.Must(template.ParseFiles("./src/templates/dark/errorUser.html"))
var DarkErrorPost = template.Must(template.ParseFiles("./src/templates/dark/errorPost.html"))
var DarkLogin = template.Must(template.ParseFiles("./src/templates/dark/login.html"))
var DarkLoginError = template.Must(template.ParseFiles("./src/templates/dark/loginerror.html"))
var DarkHome = template.Must(template.ParseFiles("./src/templates/dark/accueil.html"))
var DarkProfile = template.Must(template.ParseFiles("./src/templates/dark/profile.html"))
var DarkRegister = template.Must(template.ParseFiles("./src/templates/dark/register.html"))
var DarkRegisterError = template.Must(template.ParseFiles("./src/templates/dark/registererror.html"))
var DarkCreatePosts = template.Must(template.ParseFiles("./src/templates/dark/createpost.html"))
var DarkCreatePostsError = template.Must(template.ParseFiles("./src/templates/dark/createposterror.html"))
var DarkPost = template.Must(template.ParseFiles("./src/templates/dark/post.html"))
var DarkPopulaire = template.Must(template.ParseFiles("./src/templates/dark/populaire.html"))
var DarkPosts = template.Must(template.ParseFiles("./src/templates/dark/filtragePost.html"))
