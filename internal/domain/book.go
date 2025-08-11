package domain

import "time"

type Book struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	Author    string    `db:"author"`
	Edition   int       `db:"edition"`
	ISBN      string    `db:"isbn"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
