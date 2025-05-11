package handlers

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"FileTP/internal/utils"
)

func (h *FTPHandler) ChangePermissionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data, checking for errors
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Retrieve form values from the POST data
	path := r.PostForm.Get("path")
	perm := r.PostForm.Get("perm")
	redirectUrl := strings.ReplaceAll(filepath.Dir(path), "\\", "/")

	// Декодируем и очищаем путь
	decodedPath, err := url.PathUnescape(path)

	if err != nil {
		h.Log.Error("Path unescape error: " + err.Error())
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Проверяем валидность прав
	if !(perm == "r" || perm == "w" || perm == "rw") {
		http.Error(w, "Invalid permissions", http.StatusBadRequest)
		return
	}

	file, err := h.FileDB.Get(decodedPath)
	if err != nil {
		h.Log.Error("File not found: " + decodedPath)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	currentUser := utils.GetUserIP(r)
	if file.User != currentUser && currentUser != "localhost" {
		h.Log.Error("Permission denied for user: " + currentUser)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.FileDB.ChangePermissions(perm, decodedPath); err != nil {
		h.Log.Error("Change permissions error: " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}
