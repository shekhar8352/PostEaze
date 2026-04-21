package entities

import "time"

// PhaseKind enumerates the canonical phase categories used for lifecycle automation
// (e.g. auto-transition to "scheduled" when a Piece's scheduled_post is queued).
type PhaseKind string

const (
	PhaseKindIdea      PhaseKind = "idea"
	PhaseKindScript    PhaseKind = "script"
	PhaseKindShoot     PhaseKind = "shoot"
	PhaseKindEdit      PhaseKind = "edit"
	PhaseKindReview    PhaseKind = "review"
	PhaseKindScheduled PhaseKind = "scheduled"
	PhaseKindPublished PhaseKind = "published"
	PhaseKindCustom    PhaseKind = "custom"
)

// Phase is an ordered column within a Studio. Users can CRUD and reorder phases.
type Phase struct {
	ID         int64     `db:"id"`
	StudioID   int64     `db:"studio_id"`
	Name       string    `db:"name"`
	Slug       string    `db:"slug"`
	OrderIndex int       `db:"order_index"`
	Kind       string    `db:"kind"`
	WIPLimit   *int      `db:"wip_limit"`
	IsDefault  bool      `db:"is_default"`
	Color      string    `db:"color"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
