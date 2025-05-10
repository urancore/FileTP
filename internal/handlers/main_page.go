package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"FileTP/internal/models"
	"FileTP/internal/utils"
)

var templateFuncMap = template.FuncMap{
	"formatSize": utils.FormatFileSize,
	"baseName":   filepath.Base,
}

type TemplateData struct {
	Files    []models.File
	CurrPath string
}

func (h *FTPHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	requestedPath := r.URL.Path

	fullPath, err := h.safeJoin(requestedPath)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		h.Log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if fileInfo.IsDir() {
		// Нормализация текущего пути
		currPath := requestedPath
		if currPath != "/" {
			currPath = strings.TrimSuffix(currPath, "/")
		}

		files, err := h.FileDB.GetAll()
		if err != nil {
			h.Log.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var filteredFiles []models.File
		for _, f := range files {
			// Проверка, находится ли файл в текущей директории
			if currPath == "/" {
				// Для корневой директории: путь должен быть вида "/file"
				if strings.Count(f.Path, "/") == 1 {
					filteredFiles = append(filteredFiles, f)
				}
			} else {
				prefix := currPath + "/"
				if strings.HasPrefix(f.Path, prefix) {
					suffix := strings.TrimPrefix(f.Path, prefix)
					if !strings.Contains(suffix, "/") {
						filteredFiles = append(filteredFiles, f)
					}
				}
			}
		}

		data := TemplateData{
			Files:    filteredFiles,
			CurrPath: currPath,
		}

		h.renderFiles(w, data)
	} else {
		// Отдача файла, если это не директория
		http.ServeFile(w, r, fullPath)
	}
}

func (h *FTPHandler) renderFiles(w http.ResponseWriter, data TemplateData) {
	tmpl := template.New("index.html").Funcs(templateFuncMap)
	tmpl, err := tmpl.ParseFiles("templates/index.html")
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		h.Log.Error(err.Error())
		http.Error(w, "Render Error", http.StatusInternalServerError)
	}
}

func (h *FTPHandler) safeJoin(subpath string) (string, error) {
	cleanedPath := filepath.Clean(subpath)
	fullPath := filepath.Join(h.Config.RootPath, cleanedPath)

	relPath, err := filepath.Rel(h.Config.RootPath, fullPath)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(relPath, "..") || strings.Contains(relPath, "..") {
		return "", http.ErrAbortHandler
	}

	return fullPath, nil
}
