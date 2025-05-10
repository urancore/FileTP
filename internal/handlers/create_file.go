package handlers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"FileTP/internal/models"
	"FileTP/internal/utils"
)

func (h *FTPHandler) CreateFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	path := r.URL.Query().Get("path")
	cleanedPath := strings.ReplaceAll(path, "\\", "/")

	parsedURL, err := url.Parse(cleanedPath)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	path = parsedURL.Path

	file, handler, err := r.FormFile("fileKey")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Создаем файл в целевой директории
	dstPath := filepath.Join(h.Config.RootPath, path, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	hasher := sha256.New()
	tee := io.TeeReader(file, hasher)
	if _, err := io.Copy(dst, tee); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	hash := fmt.Sprintf("%x", hasher.Sum(nil))


	fileStat, err := os.Stat(dstPath)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}

	user_ip := utils.GetUserIP(r)

	fileObj := models.File{
		Path:        filepath.ToSlash(filepath.Join(path, handler.Filename)),
		User:        user_ip,
		Permissions: "r",
		Size:        fileStat.Size(),
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
		Type:        "file",
		LinkTarget:  filepath.ToSlash(filepath.Join(path, handler.Filename)),
		Hash:        hash,
		UploaderIP:  user_ip,
		IsDeleted:   false,
	}

	h.FileDB.Insert(fileObj)

	// Редиректим обратно в текущую директорию
	http.Redirect(w, r, path, http.StatusFound)
}
