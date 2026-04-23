package entities

import (
	"time"

	"github.com/google/uuid"
)

// Studio is a Kanban-style content pipeline board owned by a Team.
// MVP: each team has one default studio.
type Studio struct {
	ID         int64     `db:"id"`
	TeamID     uuid.UUID `db:"team_id"`
	Name       string    `db:"name"`
	PieceLabel string    `db:"piece_label"` // "Piece" | "Drop" | custom (1-32 chars)
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
