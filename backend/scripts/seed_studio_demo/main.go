// Command seed_studio_demo populates a user's Studio with a realistic set of
// demo Pieces distributed across the default phases. It is intended for local
// development and demos only.
//
// Prerequisites:
//   - The target user must already exist and belong to a team.
//   - Database migrations must be applied (including the Studio pipeline
//     migration 009).
//
// Usage (run from the backend directory):
//
//	go run scripts/seed_studio_demo/main.go --user-id <uuid>
//	go run scripts/seed_studio_demo/main.go --user-id <uuid> --reset
//
// Flags:
//
//	--user-id  UUID of the user whose primary team's Studio will be seeded.
//	--reset    Archive all existing active pieces in the Studio before seeding.
//	--count    Number of demo pieces to create (default: 12, min: 1).
//
// The script is idempotent in the sense that it will happily add more pieces
// on repeated runs. Pass --reset to get a clean board each time.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils/configs"
	"github.com/shekhar8352/PostEaze/utils/database"
	"github.com/shekhar8352/PostEaze/utils/encryption"
	"github.com/shekhar8352/PostEaze/utils/env"
	"github.com/shekhar8352/PostEaze/utils/flags"
)

// demoPiece describes one seeded piece. PhaseSlug is matched against the
// phases returned from the database so the script keeps working even if phase
// IDs drift between environments.
type demoPiece struct {
	PhaseSlug   string
	Title       string
	Description string
	ContentType entities.PieceContentType
	DueInDays   *int // nil = no due date; negative = overdue
}

func ptrInt(i int) *int { return &i }

// demoPieces is intentionally hand-crafted so the board tells a coherent story
// rather than looking like lorem ipsum.
var demoPieces = []demoPiece{
	// Ideas
	{PhaseSlug: "ideas", Title: "Behind-the-scenes: studio tour", Description: "Quick 30s walkthrough of the recording setup.", ContentType: entities.PieceContentTypeReel},
	{PhaseSlug: "ideas", Title: "Top 5 productivity apps 2025", Description: "Carousel with our favourite picks + one surprise.", ContentType: entities.PieceContentTypeCarousel},
	{PhaseSlug: "ideas", Title: "Community Q&A thread", Description: "Collect follower questions for Friday's live.", ContentType: entities.PieceContentTypePost},

	// Scripting
	{PhaseSlug: "scripting", Title: "Launch week recap video", Description: "Write voiceover script covering day 1–5 highlights.", ContentType: entities.PieceContentTypeVideo, DueInDays: ptrInt(2)},
	{PhaseSlug: "scripting", Title: "'How we ship' explainer", Description: "Long-form explainer on our release process.", ContentType: entities.PieceContentTypeVideo, DueInDays: ptrInt(5)},

	// Shooting
	{PhaseSlug: "shooting", Title: "Product demo reel — v2", Description: "Reshoot hero shots with new lighting rig.", ContentType: entities.PieceContentTypeReel, DueInDays: ptrInt(1)},

	// Editing
	{PhaseSlug: "editing", Title: "Customer testimonial montage", Description: "Edit 6 clips into a 60s highlight reel.", ContentType: entities.PieceContentTypeVideo, DueInDays: ptrInt(3)},
	{PhaseSlug: "editing", Title: "Monthly digest carousel", Description: "Design 10-slide carousel summarising the month.", ContentType: entities.PieceContentTypeCarousel},

	// Review
	{PhaseSlug: "review", Title: "Feature teaser post", Description: "Pending legal + brand review.", ContentType: entities.PieceContentTypePost, DueInDays: ptrInt(-1)},

	// Scheduled
	{PhaseSlug: "scheduled", Title: "Weekly tip: keyboard shortcuts", Description: "Queued for Tuesday 09:00.", ContentType: entities.PieceContentTypePost, DueInDays: ptrInt(4)},

	// Published
	{PhaseSlug: "published", Title: "PostEaze launch announcement", Description: "Shipped last Monday — crossposted everywhere.", ContentType: entities.PieceContentTypePost},
	{PhaseSlug: "published", Title: "Founder intro reel", Description: "Published two weeks ago, great engagement.", ContentType: entities.PieceContentTypeReel},
}

func main() {
	userIDFlag := flag.String("user-id", "", "UUID of the user whose Studio should be seeded (required)")
	resetFlag := flag.Bool("reset", false, "Archive all existing active pieces before seeding")
	countFlag := flag.Int("count", 0, "Max number of demo pieces to create (0 = all)")
	flag.Parse()

	if *userIDFlag == "" {
		flag.Usage()
		log.Fatal("--user-id is required")
	}
	if _, err := uuid.Parse(*userIDFlag); err != nil {
		log.Fatalf("--user-id must be a valid UUID: %v", err)
	}

	ctx := context.Background()
	initBackend(ctx)

	// 1. Ensure the user's studio exists (auto-seeds default phases on creation).
	studio, status, err := businessv1.EnsureStudioForUser(ctx, *userIDFlag)
	if err != nil {
		log.Fatalf("ensure studio: status=%d err=%v", status, err)
	}
	log.Printf("Studio ready: id=%d team=%s piece_label=%q", studio.ID, studio.TeamID, studio.PieceLabel)

	// 2. Load phases so we can map slugs -> ids.
	phases, err := repositories.ListPhasesByStudio(ctx, studio.ID)
	if err != nil {
		log.Fatalf("list phases: %v", err)
	}
	if len(phases) == 0 {
		log.Fatalf("studio %d has no phases; default phase seeding appears broken", studio.ID)
	}
	slugToPhaseID := make(map[string]int64, len(phases))
	for _, p := range phases {
		slugToPhaseID[p.Slug] = p.ID
	}
	log.Printf("Found %d phases (%s)", len(phases), describePhases(phases))

	// 3. Optionally archive existing active pieces for a clean slate.
	if *resetFlag {
		existing, err := repositories.ListPiecesByStudio(ctx, repositories.ListPiecesFilters{
			StudioID: studio.ID,
			Status:   string(entities.PieceStatusActive),
		})
		if err != nil {
			log.Fatalf("list existing pieces: %v", err)
		}
		for _, p := range existing {
			if err := repositories.SetPieceStatus(ctx, p.ID, string(entities.PieceStatusArchived)); err != nil {
				log.Fatalf("archive piece %d: %v", p.ID, err)
			}
		}
		log.Printf("Archived %d existing pieces (--reset)", len(existing))
	}

	// 4. Create demo pieces.
	created := 0
	skipped := 0
	for _, d := range demoPieces {
		if *countFlag > 0 && created >= *countFlag {
			break
		}
		phaseID, ok := slugToPhaseID[d.PhaseSlug]
		if !ok {
			log.Printf("  skip %q: phase slug %q not found in studio", d.Title, d.PhaseSlug)
			skipped++
			continue
		}
		req := &modelsv1.CreatePieceRequest{
			PhaseID:     phaseID,
			Title:       d.Title,
			Description: d.Description,
			ContentType: string(d.ContentType),
		}
		if d.DueInDays != nil {
			due := time.Now().Add(time.Duration(*d.DueInDays) * 24 * time.Hour).UTC().Format(time.RFC3339)
			req.DueAt = &due
		}
		piece, status, err := businessv1.CreatePiece(ctx, *userIDFlag, studio.ID, req)
		if err != nil {
			log.Fatalf("create piece %q (phase=%s): status=%d err=%v", d.Title, d.PhaseSlug, status, err)
		}
		log.Printf("  ✓ [%s] %s  (id=%d)", d.PhaseSlug, piece.Title, piece.ID)
		created++
	}

	log.Printf("Done. created=%d skipped=%d (studio_id=%d)", created, skipped, studio.ID)
}

// initBackend brings up the minimal subset of the backend needed to call the
// business/repository layers: configs, database, and encryption (some flows
// pull encrypted settings even during seed).
func initBackend(ctx context.Context) {
	env.InitEnv()

	configNames := []string{constants.DatabaseConfig}
	if err := configs.InitDev(flags.BaseConfigPath(), configNames...); err != nil {
		log.Fatalf("init configs: %v", err)
	}

	driverName, err := configs.Get().GetString(constants.DatabaseConfig, constants.DatabaseDriverNameConfigKey)
	if err != nil {
		log.Fatalf("get db driver: %v", err)
	}
	urlString, err := configs.Get().GetString(constants.DatabaseConfig, constants.DatabaseURLConfigKey)
	if err != nil {
		log.Fatalf("get db url: %v", err)
	}
	dbURL := env.ApplyEnvironmentToString(urlString)

	if err := database.Init(ctx, database.Config{
		DriverName: driverName,
		URL:        dbURL,
	}); err != nil {
		log.Fatalf("init database: %v", err)
	}

	// Encryption init is cheap and keeps us compatible with any repository
	// that lazily touches encrypted columns.
	if err := encryption.Init(); err != nil {
		log.Printf("warning: encryption init failed (continuing): %v", err)
	}
}

func describePhases(phases []entities.Phase) string {
	out := ""
	for i, p := range phases {
		if i > 0 {
			out += " → "
		}
		out += p.Slug
	}
	return out
}
