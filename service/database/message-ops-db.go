package database

import (
	"database/sql"
	"errors"
	"time"
)

// CreateGroupMessage inserts a message into a group conversation.
func (db *appdbimpl) CreateGroupMessage(groupId int64, from User, body string) (int64, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	now := time.Now().UTC()
	res, err := tx.Exec(
		"INSERT INTO group_messages (id_group, sender, body, date) VALUES (?,?,?,?)",
		groupId, from.IdUser, body, now,
	)
	if err != nil {
		return 0, err
	}
	messageID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Create receipts for all members except the sender.
	rows, err := tx.Query("SELECT id_user FROM group_members WHERE id_group = ? AND id_user <> ?", groupId, from.IdUser)
	if err != nil {
		return 0, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return 0, err
		}
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO group_message_receipts (message_id, id_group, id_user, received_at, read_at) VALUES (?,?,?,?,NULL)",
			messageID, groupId, uid, now,
		)
		if err != nil {
			return 0, err
		}
	}
	if rows.Err() != nil {
		return 0, rows.Err()
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return messageID, nil
}

// ListGroupMessages returns group messages ordered by date descending (reverse chronological), excluding deleted messages.
func (db *appdbimpl) ListGroupMessages(groupId int64, limit int, offset int) ([]GroupMessage, error) {
	rows, err := db.c.Query(
		"SELECT id, id_group, sender, body, date FROM group_messages "+
			"WHERE id_group = ? "+
			"AND id NOT IN (SELECT message_id FROM group_message_deletions WHERE id_group = ?) "+
			"ORDER BY date DESC LIMIT ? OFFSET ?",
		groupId, groupId, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var msgs []GroupMessage
	for rows.Next() {
		var m GroupMessage
		var dt time.Time
		if err := rows.Scan(&m.Id, &m.GroupID, &m.Sender, &m.Body, &dt); err != nil {
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

func (db *appdbimpl) GetDirectMessageInConversation(a User, b User, messageId int64) (Message, error) {
	var m Message
	var dt time.Time
	err := db.c.QueryRow(
		"SELECT id, sender, receiver, body, date FROM messages "+
			"WHERE id = ? AND ((sender=? AND receiver=?) OR (sender=? AND receiver=?))",
		messageId, a.IdUser, b.IdUser, b.IdUser, a.IdUser,
	).Scan(&m.Id, &m.Sender, &m.Receiver, &m.Body, &dt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Message{}, ErrMessageNotFound
		}
		return Message{}, err
	}
	m.Date = dt
	return m, nil
}

func (db *appdbimpl) GetGroupMessageInGroup(groupId int64, messageId int64) (GroupMessage, error) {
	var m GroupMessage
	var dt time.Time
	err := db.c.QueryRow(
		"SELECT id, id_group, sender, body, date FROM group_messages WHERE id = ? AND id_group = ?",
		messageId, groupId,
	).Scan(&m.Id, &m.GroupID, &m.Sender, &m.Body, &dt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GroupMessage{}, ErrMessageNotFound
		}
		return GroupMessage{}, err
	}
	m.Date = dt
	return m, nil
}

func (db *appdbimpl) DeleteDirectMessage(messageId int64, deletedBy User) error {
	var sender string
	err := db.c.QueryRow("SELECT sender FROM messages WHERE id = ?", messageId).Scan(&sender)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrMessageNotFound
		}
		return err
	}
	if sender != deletedBy.IdUser {
		return ErrForbiddenMessageAction
	}

	_, err = db.c.Exec(
		"INSERT OR IGNORE INTO direct_message_deletions (message_id, deleted_at, deleted_by) VALUES (?,?,?)",
		messageId, time.Now().UTC(), deletedBy.IdUser,
	)
	return err
}

func (db *appdbimpl) DeleteGroupMessage(groupId int64, messageId int64, deletedBy User) error {
	var sender string
	err := db.c.QueryRow("SELECT sender FROM group_messages WHERE id = ? AND id_group = ?", messageId, groupId).Scan(&sender)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrMessageNotFound
		}
		return err
	}
	if sender != deletedBy.IdUser {
		return ErrForbiddenMessageAction
	}

	_, err = db.c.Exec(
		"INSERT OR IGNORE INTO group_message_deletions (message_id, id_group, deleted_at, deleted_by) VALUES (?,?,?,?)",
		messageId, groupId, time.Now().UTC(), deletedBy.IdUser,
	)
	return err
}

func (db *appdbimpl) SetDirectMessageReaction(messageId int64, user User, reaction string) error {
	_, err := db.c.Exec(
		"INSERT OR REPLACE INTO direct_message_reactions (message_id, id_user, reaction, created_at) VALUES (?,?,?,?)",
		messageId, user.IdUser, reaction, time.Now().UTC(),
	)
	return err
}

func (db *appdbimpl) RemoveDirectMessageReaction(messageId int64, user User) error {
	_, err := db.c.Exec("DELETE FROM direct_message_reactions WHERE message_id = ? AND id_user = ?", messageId, user.IdUser)
	return err
}

func (db *appdbimpl) ListDirectMessageReactions(messageId int64) ([]MessageReaction, error) {
	rows, err := db.c.Query("SELECT id_user, reaction FROM direct_message_reactions WHERE message_id = ?", messageId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []MessageReaction
	for rows.Next() {
		var mr MessageReaction
		if err := rows.Scan(&mr.UserID, &mr.Reaction); err != nil {
			return nil, err
		}
		out = append(out, mr)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return out, nil
}

func (db *appdbimpl) SetGroupMessageReaction(groupId int64, messageId int64, user User, reaction string) error {
	// Validate membership of the message to this group
	if _, err := db.GetGroupMessageInGroup(groupId, messageId); err != nil {
		return err
	}
	_, err := db.c.Exec(
		"INSERT OR REPLACE INTO group_message_reactions (message_id, id_user, reaction, created_at) VALUES (?,?,?,?)",
		messageId, user.IdUser, reaction, time.Now().UTC(),
	)
	return err
}

func (db *appdbimpl) RemoveGroupMessageReaction(groupId int64, messageId int64, user User) error {
	// Validate membership of the message to this group
	if _, err := db.GetGroupMessageInGroup(groupId, messageId); err != nil {
		return err
	}
	_, err := db.c.Exec("DELETE FROM group_message_reactions WHERE message_id = ? AND id_user = ?", messageId, user.IdUser)
	return err
}

func (db *appdbimpl) ListGroupMessageReactions(groupId int64, messageId int64) ([]MessageReaction, error) {
	// Validate membership of the message to this group
	if _, err := db.GetGroupMessageInGroup(groupId, messageId); err != nil {
		return nil, err
	}

	rows, err := db.c.Query("SELECT id_user, reaction FROM group_message_reactions WHERE message_id = ?", messageId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []MessageReaction
	for rows.Next() {
		var mr MessageReaction
		if err := rows.Scan(&mr.UserID, &mr.Reaction); err != nil {
			return nil, err
		}
		out = append(out, mr)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return out, nil
}
