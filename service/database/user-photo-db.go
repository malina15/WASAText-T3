package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) SetUserPhotoPath(user User, photoPath string) error {
	// SQLite: replace to upsert by primary key.
	_, err := db.c.Exec("INSERT OR REPLACE INTO user_photos (id_user, photo_path) VALUES (?,?)", user.IdUser, photoPath)
	return err
}

func (db *appdbimpl) GetUserPhotoPath(user User) (string, error) {
	var path string
	err := db.c.QueryRow("SELECT photo_path FROM user_photos WHERE id_user = ?", user.IdUser).Scan(&path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserPhotoNotFound
		}
		return "", err
	}
	return path, nil
}
