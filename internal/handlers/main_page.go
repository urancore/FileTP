package handlers

import (
	"html/template"
	"net/http"

	"FileTP/internal/models"
	"FileTP/internal/utils"
)

var templateFuncMap = template.FuncMap{
	"formatSize": utils.FormatFileSize,
	"formatPath": utils.FormatFilePath,
}

func (h *FTPHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("index.html").Funcs(templateFuncMap)
	tmpl, err := tmpl.ParseFiles("templates/index.html")
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	files, err := h.FileDB.GetAll()
	if err != nil {
		h.Log.Error(err.Error())
	}

	data := struct {
		Files []models.File
	}{
		Files: files,
	}

	if err := tmpl.Execute(w, data); err != nil {
		h.Log.Error(err.Error())
		http.Error(w, "Render Error", http.StatusInternalServerError)
	}
}
