package sql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"FileTP/internal/models"
)

type FileDB struct {
	db *sql.DB
}

func NewFileDB(dataSourceName string) (*FileDB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := createTable(db); err != nil {
		return nil, err
	}

	return &FileDB{db: db}, nil
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS files (
		path TEXT PRIMARY KEY,
		user_name TEXT,
		permissions TEXT,
		size INTEGER,
		created_at DATETIME,
		modified_at DATETIME,
		type TEXT CHECK(type IN ('file', 'directory', 'symlink')),
		link_target TEXT,
		hash TEXT,
		uploader_ip TEXT,
		is_deleted INTEGER
	)`

	_, err := db.Exec(query)
	return err
}

func (fdb *FileDB) Close() error {
	return fdb.db.Close()
}

func (fdb *FileDB) Insert(file models.File) error {
	query := `
	INSERT INTO files (
		path, user_name, permissions, size, created_at,
		modified_at, type, link_target, hash, uploader_ip, is_deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	isDeleted := 0
	if file.IsDeleted {
		isDeleted = 1
	}

	_, err := fdb.db.Exec(query,
		file.Path,
		file.User,
		file.Permissions,
		file.Size,
		file.CreatedAt,
		file.ModifiedAt,
		file.Type,
		file.LinkTarget,
		file.Hash,
		file.UploaderIP,
		isDeleted,
	)

	return err
}

func (fdb *FileDB) Get(path string) (*models.File, error) {
	query := `SELECT * FROM files WHERE path = ?`
	row := fdb.db.QueryRow(query, path)

	var file models.File
	var isDeleted int
	var createdAt, modifiedAt string

	err := row.Scan(
		&file.Path,
		&file.User,
		&file.Permissions,
		&file.Size,
		&createdAt,
		&modifiedAt,
		&file.Type,
		&file.LinkTarget,
		&file.Hash,
		&file.UploaderIP,
		&isDeleted,
	)

	if err != nil {
		return nil, err
	}

	// Конвертация времени из строки
	file.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %v", err)
	}

	file.ModifiedAt, err = time.Parse(time.RFC3339, modifiedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse modified_at: %v", err)
	}

	file.IsDeleted = isDeleted == 1

	return &file, nil
}

func (fdb *FileDB) GetByPath(dirPath string) ([]models.File, error) {
	query := `SELECT * FROM files WHERE (path = ? OR path LIKE ? || '/%') AND is_deleted = 0`
	rows, err := fdb.db.Query(query, dirPath, dirPath)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var files []models.File

	for rows.Next() {
		var file models.File
		var isDeleted int
		var createdAt, modifiedAt string

		err := rows.Scan(
			&file.Path,
			&file.User,
			&file.Permissions,
			&file.Size,
			&createdAt,
			&modifiedAt,
			&file.Type,
			&file.LinkTarget,
			&file.Hash,
			&file.UploaderIP,
			&isDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		file.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		file.ModifiedAt, _ = time.Parse(time.RFC3339, modifiedAt)
		file.IsDeleted = isDeleted == 1

		files = append(files, file)
	}
	return files, nil
}

func (fdb *FileDB) GetAll() ([]models.File, error) {
	query := `SELECT
        path,
        user_name,
        permissions,
        size,
        created_at,
        modified_at,
        type,
        link_target,
        hash,
        uploader_ip,
        is_deleted
    FROM files`

	rows, err := fdb.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		var isDeleted int
		var createdAt, modifiedAt string

		err := rows.Scan(
			&file.Path,
			&file.User,
			&file.Permissions,
			&file.Size,
			&createdAt,
			&modifiedAt,
			&file.Type,
			&file.LinkTarget,
			&file.Hash,
			&file.UploaderIP,
			&isDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		// Парсинг времени создания
		file.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("created_at parsing error: %v", err)
		}

		// Парсинг времени изменения
		file.ModifiedAt, err = time.Parse(time.RFC3339, modifiedAt)
		if err != nil {
			return nil, fmt.Errorf("modified_at parsing error: %v", err)
		}

		file.IsDeleted = isDeleted == 1
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return files, nil
}

func (fdb *FileDB) Update(file models.File) error {
	query := `
	UPDATE files SET
		user_name = ?,
		permissions = ?,
		size = ?,
		created_at = ?,
		modified_at = ?,
		type = ?,
		link_target = ?,
		hash = ?,
		uploader_ip = ?,
		is_deleted = ?
	WHERE path = ?`

	isDeleted := 0
	if file.IsDeleted {
		isDeleted = 1
	}

	_, err := fdb.db.Exec(query,
		file.User,
		file.Permissions,
		file.Size,
		file.CreatedAt,
		file.ModifiedAt,
		file.Type,
		file.LinkTarget,
		file.Hash,
		file.UploaderIP,
		isDeleted,
		file.Path,
	)

	return err
}

func (fdb *FileDB) SoftDelete(path string) error {
	query := `UPDATE files SET is_deleted = 1 WHERE path = ?`
	_, err := fdb.db.Exec(query, path)
	return err
}
