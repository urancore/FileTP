package handlers

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"FileTP/internal/models"
	"FileTP/internal/utils"
)

func (h *FTPHandler) CreateDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим форму
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Получаем путь из query параметра
	path := r.URL.Query().Get("path")
	cleanedPath := strings.ReplaceAll(path, "\\", "/")

	parsedURL, err := url.Parse(cleanedPath)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	path = parsedURL.Path

	// Получаем имя директории из формы
	dirName := strings.TrimSpace(r.FormValue("dirname"))
	if dirName == "" {
		http.Error(w, "Directory name is required", http.StatusBadRequest)
		return
	}

	// Создаем полный путь к новой директории
	fullPath := filepath.Join(h.Config.RootPath, path, dirName)
	dirPath := filepath.Join(path, dirName)
	if !strings.HasPrefix(dirPath, "/") {
		dirPath = "/" + dirPath
	}
	dirPath = filepath.Clean(dirPath)
	dirPath = filepath.ToSlash(dirPath)

	// Проверяем безопасность пути
	if !strings.HasPrefix(fullPath, h.Config.RootPath) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Создаем директорию
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		h.Log.Error("Failed to create directory: " + err.Error())
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	userIP := utils.GetUserIP(r)

	// Создаем запись о директории
	dirObj := models.File{
		Path:        filepath.ToSlash(filepath.Join(path, dirName)),
		User:        userIP,
		Permissions: "rw",
		Size:        0,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
		Type:        "directory",
		LinkTarget:  filepath.ToSlash(filepath.Join("/", path, dirName)) + "/",
		Hash:        "",
		UploaderIP:  userIP,
		IsDeleted:   false,
	}

	// Сохраняем в базе данных
	if err := h.FileDB.Insert(dirObj); err != nil {
		h.Log.Error("Database error: " + err.Error())
		http.Error(w, "Error saving directory info", http.StatusInternalServerError)
		return
	}

	// Редирект обратно в текущую директорию
	http.Redirect(w, r, path, http.StatusFound)
}
