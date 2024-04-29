package forum

import "html/template"

var Error = template.Must(template.ParseFiles("./src/templates/error.html"))
