package studio

import "github.com/shekhar8352/PostEaze/entities"

// DefaultPhaseTemplate describes one phase that should be created when a new
// studio is seeded. Order is significant and corresponds to OrderIndex.
type DefaultPhaseTemplate struct {
	Name      string
	Slug      string
	Kind      entities.PhaseKind
	Color     string
	IsDefault bool
}

// DefaultPhaseTemplates returns the canonical list of phases created for a
// brand-new studio. Kinds are meaningful: lifecycle automation (e.g.
// auto-transition to "Scheduled" when a scheduled_post is created) relies on
// Kind rather than Name, so users can rename phases without breaking
// automations.
func DefaultPhaseTemplates() []DefaultPhaseTemplate {
	return []DefaultPhaseTemplate{
		{Name: "Ideas", Slug: "ideas", Kind: entities.PhaseKindIdea, Color: "#6366f1", IsDefault: true},
		{Name: "Scripting", Slug: "scripting", Kind: entities.PhaseKindScript, Color: "#8b5cf6"},
		{Name: "Shooting", Slug: "shooting", Kind: entities.PhaseKindShoot, Color: "#ec4899"},
		{Name: "Editing", Slug: "editing", Kind: entities.PhaseKindEdit, Color: "#f59e0b"},
		{Name: "Review", Slug: "review", Kind: entities.PhaseKindReview, Color: "#10b981"},
		{Name: "Scheduled", Slug: "scheduled", Kind: entities.PhaseKindScheduled, Color: "#0ea5e9"},
		{Name: "Published", Slug: "published", Kind: entities.PhaseKindPublished, Color: "#22c55e"},
	}
}
