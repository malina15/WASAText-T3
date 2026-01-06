package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

const directMessageNotDeletedClause = "AND id NOT IN (SELECT message_id FROM direct_message_deletions) "

func (db *appdbimpl) ListConversations(user User) ([]Conversation, error) {
	groups, err := db.ListGroupsForUser(user)
	if err != nil {
		return nil, err
	}

	conversations := make([]Conversation, 0, len(groups)+8)

	// Direct conversations: derive peers from messages
	rows, err := db.c.Query(
		"SELECT CASE WHEN sender = ? THEN receiver ELSE sender END AS peer_id, "+
			"MAX(CAST(strftime('%s', date) AS INTEGER)) AS last_ts "+
			"FROM messages "+
			"WHERE (sender = ? OR receiver = ?) " +
			directMessageNotDeletedClause +
			"GROUP BY peer_id",
		user.IdUser, user.IdUser, user.IdUser,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var peerID string
		var lastTs int64
		if err := rows.Scan(&peerID, &lastTs); err != nil {
			return nil, err
		}
		lastDate := time.Unix(lastTs, 0).UTC()

		// last message preview
		var lastBody string
		_ = db.c.QueryRow(
			"SELECT body FROM messages "+
				"WHERE ((sender=? AND receiver=?) OR (sender=? AND receiver=?)) "+
				directMessageNotDeletedClause +
				"ORDER BY date DESC LIMIT 1",
			user.IdUser, peerID, peerID, user.IdUser,
		).Scan(&lastBody)

		nickname, err := db.GetNickname(User{IdUser: peerID})
		if err != nil {
			// If nickname can't be retrieved, fallback to identifier.
			nickname = peerID
		}

		conversations = append(conversations, Conversation{
			Peer:               peerID,
			IsGroup:            false,
			Name:               nickname,
			PhotoURL:           fmt.Sprintf("/users/%s/photo", peerID),
			LastMessageAt:      lastDate,
			LastMessagePreview: snippet(lastBody, 40),
		})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	// Group conversations: list groups for user and attach last activity from group_messages
	for _, g := range groups {
		var lastBody string
		var lastDate time.Time
		err := db.c.QueryRow(
			"SELECT body, date FROM group_messages "+
				"WHERE id_group = ? "+
				"AND id NOT IN (SELECT message_id FROM group_message_deletions WHERE id_group = ?) "+
				"ORDER BY date DESC LIMIT 1",
			g.Id, g.Id,
		).Scan(&lastBody, &lastDate)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				lastDate = g.CreatedAt
				lastBody = ""
			} else {
				return nil, err
			}
		}

		conversations = append(conversations, Conversation{
			Peer:               fmt.Sprintf("g-%d", g.Id),
			IsGroup:            true,
			Name:               g.Name,
			PhotoURL:           fmt.Sprintf("/groups/%d/photo", g.Id),
			LastMessageAt:      lastDate,
			LastMessagePreview: snippet(lastBody, 40),
		})
	}

	sort.Slice(conversations, func(i, j int) bool {
		return conversations[i].LastMessageAt.After(conversations[j].LastMessageAt)
	})
	return conversations, nil
}

func snippet(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max]
}


