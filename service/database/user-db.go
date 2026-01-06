package database

import (
    "database/sql"
    "time"
)
// Database function that gets the stream of a user (photos of people that are followed by the latter)
func (db *appdbimpl) GetStream(user User) ([]Photo, error) {

	rows, err := db.c.Query(`SELECT * FROM photos WHERE id_user IN (SELECT followed FROM followers WHERE follower = ?) ORDER BY date DESC`,
		user.IdUser)
	if err != nil {
		return nil, err
	}
	// Wait for the function to finish before closing rows
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset
	var res []Photo
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.PhotoId, &photo.Owner, &photo.Date) //  &photo.Comments, &photo.Likes,
		if err != nil {
			return nil, err
		}
		res = append(res, photo)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return res, nil
}

// Database function that retrieves all users identifiers
func (db *appdbimpl) GetAllUsers() ([]User, error) {
	rows, err := db.c.Query("SELECT id_user FROM users")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.IdUser); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return users, nil
}

// CreateMessage inserts a direct message between two users
func (db *appdbimpl) CreateMessage(from User, to User, body string) (int64, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	now := time.Now().UTC()
	res, err := tx.Exec("INSERT INTO messages (sender, receiver, body, date) VALUES (?,?,?,?)", from.IdUser, to.IdUser, body, now)
	if err != nil {
		return 0, err
	}
	messageID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Mark as received for receiver (read_at will be set when receiver reads the conversation).
	_, err = tx.Exec(
		"INSERT OR REPLACE INTO direct_message_receipts (message_id, receiver_id, received_at, read_at) VALUES (?,?,?,NULL)",
		messageID, to.IdUser, now,
	)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return messageID, nil
}

// ListMessages returns messages between a and b ordered by date descending (reverse chronological)
func (db *appdbimpl) ListMessages(a User, b User, limit int, offset int) ([]Message, error) {
    rows, err := db.c.Query(
		"SELECT id, sender, receiver, body, date FROM messages "+
			"WHERE ((sender=? AND receiver=?) OR (sender=? AND receiver=?)) "+
			directMessageNotDeletedClause +
			"ORDER BY date DESC LIMIT ? OFFSET ?",
		a.IdUser, b.IdUser, b.IdUser, a.IdUser, limit, offset)
    if err != nil {
        return nil, err
    }
    defer func() { _ = rows.Close() }()
    var msgs []Message
    for rows.Next() {
        var m Message
        var dt time.Time
        if err := rows.Scan(&m.Id, &m.Sender, &m.Receiver, &m.Body, &dt); err != nil {
            return nil, err
        }
        m.Date = dt
        msgs = append(msgs, m)
    }
    if rows.Err() != nil {
        return nil, rows.Err()
    }
    return msgs, nil
}

// Database function that adds a new user in the database upon registration
func (db *appdbimpl) CreateUser(u User) error {

	_, err := db.c.Exec("INSERT INTO users (id_user,nickname) VALUES (?, ?)",
		u.IdUser, u.IdUser)

	if err != nil {
		return err
	}

	return nil
}

// FindUserByNickname returns the user identifier of the user that owns the given nickname.
// It returns found=false if the nickname does not exist.
func (db *appdbimpl) FindUserByNickname(nickname string) (User, bool, error) {
	var id string
	err := db.c.QueryRow("SELECT id_user FROM users WHERE nickname = ?", nickname).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, false, nil
		}
		return User{}, false, err
	}
	return User{IdUser: id}, true, nil
}

// [EXTRA] Database function that checks if a user exists
func (db *appdbimpl) CheckUser(targetUser User) (bool, error) {

	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE id_user = ?",
		targetUser.IdUser).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return true, err
	}

	// If the counter is 1 then the user exists
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}
