package database

import (
	"database/sql"
	"errors"
	"time"
)

// CreateGroup creates a new group and adds creator + members to it.
func (db *appdbimpl) CreateGroup(creator User, name string, members []User) (int64, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.Exec("INSERT INTO groups (name, photo_path, created_at) VALUES (?,?,?)", name, "", time.Now().UTC())
	if err != nil {
		return 0, err
	}
	groupID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Always add creator
	_, err = tx.Exec("INSERT INTO group_members (id_group, id_user) VALUES (?,?)", groupID, creator.IdUser)
	if err != nil {
		return 0, err
	}

	// Add other members (ignore duplicates)
	for _, m := range members {
		if m.IdUser == "" || m.IdUser == creator.IdUser {
			continue
		}
		_, err = tx.Exec("INSERT OR IGNORE INTO group_members (id_group, id_user) VALUES (?,?)", groupID, m.IdUser)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return groupID, nil
}

func (db *appdbimpl) AddUserToGroup(groupId int64, user User) error {
	// Ensure group exists
	if _, err := db.GetGroup(groupId); err != nil {
		return err
	}

	inGroup, err := db.IsUserInGroup(groupId, user)
	if err != nil {
		return err
	}
	if inGroup {
		return ErrUserAlreadyInGroup
	}

	_, err = db.c.Exec("INSERT INTO group_members (id_group, id_user) VALUES (?,?)", groupId, user.IdUser)
	return err
}

func (db *appdbimpl) RemoveUserFromGroup(groupId int64, user User) error {
	// Ensure group exists
	if _, err := db.GetGroup(groupId); err != nil {
		return err
	}

	res, err := db.c.Exec("DELETE FROM group_members WHERE id_group = ? AND id_user = ?", groupId, user.IdUser)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrUserNotInGroup
	}
	return nil
}

func (db *appdbimpl) IsUserInGroup(groupId int64, user User) (bool, error) {
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM group_members WHERE id_group = ? AND id_user = ?", groupId, user.IdUser).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (db *appdbimpl) SetGroupName(groupId int64, name string) error {
	res, err := db.c.Exec("UPDATE groups SET name = ? WHERE id_group = ?", name, groupId)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrGroupNotFound
	}
	return nil
}

func (db *appdbimpl) SetGroupPhotoPath(groupId int64, photoPath string) error {
	res, err := db.c.Exec("UPDATE groups SET photo_path = ? WHERE id_group = ?", photoPath, groupId)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrGroupNotFound
	}
	return nil
}

func (db *appdbimpl) GetGroup(groupId int64) (Group, error) {
	var g Group
	var created time.Time
	err := db.c.QueryRow("SELECT id_group, name, photo_path, created_at FROM groups WHERE id_group = ?", groupId).
		Scan(&g.Id, &g.Name, &g.PhotoPath, &created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Group{}, ErrGroupNotFound
		}
		return Group{}, err
	}
	g.CreatedAt = created
	return g, nil
}

func (db *appdbimpl) ListGroupsForUser(user User) ([]Group, error) {
	rows, err := db.c.Query(`
		SELECT g.id_group, g.name, g.photo_path, g.created_at
		FROM groups g
		INNER JOIN group_members gm ON gm.id_group = g.id_group
		WHERE gm.id_user = ?
		ORDER BY g.created_at DESC
	`, user.IdUser)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var groups []Group
	for rows.Next() {
		var g Group
		var created time.Time
		if err := rows.Scan(&g.Id, &g.Name, &g.PhotoPath, &created); err != nil {
			return nil, err
		}
		g.CreatedAt = created
		groups = append(groups, g)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return groups, nil
}


