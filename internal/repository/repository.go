package repository

import (
	"database/sql"
	"discordbot/internal/models"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(userID, userName, reason string, createdAt, endAt time.Time) error {
	_, err := r.db.Exec("INSERT INTO inactive (user_id, user_name, reason, created_at, end_at) VALUES (?, ?, ?, ?, ?)",
		userID, userName, reason, createdAt, endAt)
	return err
}

func (r *Repository) List() ([]models.Inactive, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, user_name, reason, created_at, end_at
		FROM inactive`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.Inactive
	for rows.Next() {
		var user models.Inactive
		if err := rows.Scan(
			&user.ID,
			&user.UserID,
			&user.UserName,
			&user.Reason,
			&user.CreatedAt,
			&user.EndAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM inactive WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ListWithEndAt(now time.Time) ([]models.Inactive, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, user_name, reason, created_at, end_at
		FROM inactive
		WHERE end_at < ?`, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.Inactive
	for rows.Next() {
		var user models.Inactive
		if err := rows.Scan(
			&user.ID,
			&user.UserID,
			&user.UserName,
			&user.Reason,
			&user.CreatedAt,
			&user.EndAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
