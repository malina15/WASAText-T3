package database

import (
	"database/sql"
	"errors"
	"time"
)

func (db *appdbimpl) MarkDirectConversationRead(reader User, peer User) error {
	// Mark as read all messages sent by peer to reader.
	_, err := db.c.Exec(
		"UPDATE direct_message_receipts SET read_at = ? "+
			"WHERE receiver_id = ? AND read_at IS NULL "+
			"AND message_id IN (SELECT id FROM messages WHERE sender = ? AND receiver = ?)",
		time.Now().UTC(), reader.IdUser, peer.IdUser, reader.IdUser,
	)
	return err
}

func (db *appdbimpl) MarkGroupConversationRead(groupId int64, reader User) error {
	_, err := db.c.Exec(
		"UPDATE group_message_receipts SET read_at = ? "+
			"WHERE id_group = ? AND id_user = ? AND read_at IS NULL",
		time.Now().UTC(), groupId, reader.IdUser,
	)
	return err
}

func (db *appdbimpl) GetDirectMessageCheckmarks(messageId int64) (int, error) {
	var readAt sql.NullTime
	err := db.c.QueryRow("SELECT read_at FROM direct_message_receipts WHERE message_id = ?", messageId).Scan(&readAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Backward compatibility: if receipt row is missing, treat as received (1 check).
			return 1, nil
		}
		return 0, err
	}
	if readAt.Valid {
		return 2, nil
	}
	return 1, nil
}

func (db *appdbimpl) GetGroupMessageCheckmarks(groupId int64, messageId int64) (int, error) {
	var total int
	var readCnt int
	err := db.c.QueryRow(
		"SELECT COUNT(*), SUM(CASE WHEN read_at IS NOT NULL THEN 1 ELSE 0 END) "+
			"FROM group_message_receipts WHERE id_group = ? AND message_id = ?",
		groupId, messageId,
	).Scan(&total, &readCnt)
	if err != nil {
		return 0, err
	}
	if total == 0 {
		// No recipients (group only has sender) -> consider it read.
		return 2, nil
	}
	if readCnt == total {
		return 2, nil
	}
	return 1, nil
}
