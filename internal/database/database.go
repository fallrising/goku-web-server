package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/fallrising/goku-api/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) SaveBookmark(info models.URLInfo) error {
	query := `
		INSERT INTO bookmarks (url, title, description, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	tags := strings.Join(info.Tags, ",")
	now := time.Now()

	_, err := d.db.Exec(query, info.URL, info.Title, info.Description, tags, now, now)
	if err != nil {
		return fmt.Errorf("error saving bookmark: %w", err)
	}

	return nil
}

func (d *Database) GetBookmarkByURL(url string) (*models.URLInfo, error) {
	query := `
		SELECT id, url, title, description, tags, created_at, updated_at
		FROM bookmarks
		WHERE url = ?
	`

	var bookmark models.URLInfo
	var tags string
	var createdAt, updatedAt time.Time

	err := d.db.QueryRow(query, url).Scan(
		&bookmark.ID,
		&bookmark.URL,
		&bookmark.Title,
		&bookmark.Description,
		&tags,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No bookmark found
		}
		return nil, fmt.Errorf("error querying bookmark: %w", err)
	}

	bookmark.Tags = strings.Split(tags, ",")

	return &bookmark, nil
}

func (d *Database) GetAllBookmarks() ([]models.URLInfo, error) {
	query := `
		SELECT id, url, title, description, tags, created_at, updated_at
		FROM bookmarks
		ORDER BY created_at DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying bookmarks: %w", err)
	}
	defer rows.Close()

	var bookmarks []models.URLInfo
	for rows.Next() {
		var bookmark models.URLInfo
		var tags string
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&bookmark.ID,
			&bookmark.URL,
			&bookmark.Title,
			&bookmark.Description,
			&tags,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning bookmark row: %w", err)
		}

		bookmark.Tags = strings.Split(tags, ",")
		bookmarks = append(bookmarks, bookmark)
	}

	return bookmarks, nil
}

func (d *Database) UpdateBookmark(info models.URLInfo) error {
	query := `
		UPDATE bookmarks
		SET url = ?, title = ?, description = ?, tags = ?, updated_at = ?
		WHERE id = ?
	`

	tags := strings.Join(info.Tags, ",")
	now := time.Now()

	_, err := d.db.Exec(query, info.URL, info.Title, info.Description, tags, now, info.ID)
	if err != nil {
		return fmt.Errorf("error updating bookmark: %w", err)
	}

	return nil
}

func (d *Database) DeleteBookmark(id int) error {
	query := `DELETE FROM bookmarks WHERE id = ?`

	_, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting bookmark: %w", err)
	}

	return nil
}
