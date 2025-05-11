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
		currentPath := strings.Trim(currPath, "/")
		currentComponents := strings.Split(currentPath, "/")

		for _, f := range files {
			filePath := strings.Trim(f.Path, "/")
			fileComponents := strings.Split(filePath, "/")

			// Проверяем, что текущий путь является префиксом пути файла
			if len(fileComponents) < len(currentComponents) {
				continue
			}

			match := true
			for i := 0; i < len(currentComponents); i++ {
				if currentComponents[i] != fileComponents[i] {
					match = false
					break
				}
			}

			// Проверяем уровень вложенности
			if match && len(fileComponents) == len(currentComponents)+1 {
				filteredFiles = append(filteredFiles, f)
			} else if currentPath == "" && len(fileComponents) == 1 {
				// Корневая директория
				filteredFiles = append(filteredFiles, f)
			}
		}
		data := TemplateData{
			Files:    filteredFiles,
			CurrPath: currPath,
		}

		h.renderFiles(w, data)
	} else {
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
