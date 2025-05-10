package models

import "time"

type File struct {
	Path        string
	User        string // Владелец
	Permissions string // Пример: "r" "w"
	Size        int64  // Размер в байтах
	CreatedAt   time.Time
	ModifiedAt  time.Time
	Type        string // "file", "directory", "symlink"
	LinkTarget  string
	Hash        string // Хэш содержимого
	UploaderIP  string // IP-адрес загрузившего
	IsDeleted   bool
}
