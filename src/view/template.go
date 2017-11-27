package view

import (
	"global"
	"html/template"
	"net/http"
	"time"
)

var FuncMap = template.FuncMap{
	"noescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"formatTime": func(t time.Time, layout string) string {
		return t.Format(layout)
	},
}

func Render(w http.ResponseWriter, r *http.Request, htmlFile string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}

	tpl, err := template.New("layout.html").Funcs(FuncMap).ParseFiles(global.App.TemplateDir+"layout.html", global.App.TemplateDir+htmlFile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	err = tpl.Execute(w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
