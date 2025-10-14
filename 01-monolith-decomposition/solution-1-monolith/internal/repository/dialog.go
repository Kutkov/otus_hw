package repository

import (
	"database/sql"
)

type DialogRepository struct {
	db *sql.DB
}

func NewDialogRepository(database *sql.DB) *DialogRepository {
	return &DialogRepository{db: database}
}

type DialogMessage struct {
	ID         string
	FromUserID string
	ToUserID   string
	Text       string
	CreatedAt  string
}

func (r *DialogRepository) CreateMessage(message *DialogMessage) error {
	_, err := r.db.Exec(
		`INSERT INTO dialog_messages(id, from_user_id, to_user_id, text, created_at) VALUES(?,?,?,?,?)`,
		message.ID, message.FromUserID, message.ToUserID, message.Text, message.CreatedAt,
	)
	return err
}

func (r *DialogRepository) GetMessagesBetweenUsers(userID1, userID2 string) ([]DialogMessage, error) {
	rows, err := r.db.Query(`
		SELECT id, from_user_id, to_user_id, text, created_at 
		FROM dialog_messages 
		WHERE (from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)
		ORDER BY created_at ASC
	`, userID1, userID2, userID2, userID1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []DialogMessage
	for rows.Next() {
		var msg DialogMessage
		if err := rows.Scan(&msg.ID, &msg.FromUserID, &msg.ToUserID, &msg.Text, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
