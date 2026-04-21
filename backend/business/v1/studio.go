package businessv1

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/scheduledpublish"
	studiosvc "github.com/shekhar8352/PostEaze/services/studio"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// ErrNoTeam is returned when we can't resolve a team for the current user.
var ErrNoTeam = errors.New("user is not a member of any team")

// init wires Studio lifecycle hooks into the publishing pipeline. Using a
// callback registry avoids an import cycle (scheduledpublish -> businessv1).
func init() {
	scheduledpublish.OnPostFinalized = OnScheduledPostPublished
}

// resolveUserTeamID returns the first team UUID the user belongs to. For MVP
// each user has a single primary team; when we add explicit team switching
// this helper will accept an override.
func resolveUserTeamID(ctx context.Context, userIDStr string) (uuid.UUID, error) {
	memberships, err := repositories.GetTeamsByUserID(ctx, userIDStr)
	if err != nil {
		// The repository layer surfaces "no records found" as an error rather
		// than an empty slice. Normalize that to ErrNoTeam so callers can
		// return a friendly 400 instead of a raw 500.
		if errors.Is(err, database.ErrNoRecords) {
			return uuid.Nil, ErrNoTeam
		}
		return uuid.Nil, err
	}
	if len(memberships) == 0 {
		return uuid.Nil, ErrNoTeam
	}
	// Prefer the primary/active membership; otherwise take the first.
	pick := memberships[0]
	for _, m := range memberships {
		if m.IsPrimary && m.Status == "active" {
			pick = m
			break
		}
	}
	tid, err := uuid.Parse(pick.TeamID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid team id %q: %w", pick.TeamID, err)
	}
	return tid, nil
}

// provisionDefaultTeamForUser creates a personal team with the given user as
// both owner and primary admin member, all in a single transaction. It is
// invoked lazily by EnsureStudioForUser the first time a user opens the Studio
// without an existing team.
func provisionDefaultTeamForUser(ctx context.Context, userIDStr string) (uuid.UUID, error) {
	if _, err := uuid.Parse(userIDStr); err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id %q: %w", userIDStr, err)
	}

	user, err := repositories.GetUserByID(ctx, userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("load user: %w", err)
	}
	teamName := "My Workspace"
	if user != nil && user.Name != "" {
		teamName = user.Name + "'s Workspace"
	}

	tx, err := database.GetTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}
	team, err := repositories.CreateTeam(ctx, tx, modelsv1.Team{
		Name:       teamName,
		OwnerID:    userIDStr,
		Visibility: "private",
		Status:     "active",
	})
	if err != nil {
		database.RollbackTx(tx)
		return uuid.Nil, fmt.Errorf("create team: %w", err)
	}
	if _, err := repositories.CreateTeamMember(ctx, tx, modelsv1.TeamMember{
		TeamID:    team.ID,
		UserID:    userIDStr,
		Role:      string(modelsv1.RoleAdmin),
		Status:    "active",
		IsPrimary: true,
	}); err != nil {
		database.RollbackTx(tx)
		return uuid.Nil, fmt.Errorf("create team membership: %w", err)
	}
	if err := database.CommitTx(tx); err != nil {
		return uuid.Nil, err
	}

	tid, err := uuid.Parse(team.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid new team id %q: %w", team.ID, err)
	}
	return tid, nil
}

// assertStudioAccess verifies the user belongs to the team that owns the studio
// and returns the studio for downstream use. Returns (nil, 404) if the studio
// doesn't exist, (nil, 403) if the user isn't a member of its team.
func assertStudioAccess(ctx context.Context, userIDStr string, studioID int64) (*entities.Studio, int, error) {
	studio, err := repositories.GetStudioByID(ctx, studioID)
	if err != nil {
		return nil, 500, err
	}
	if studio == nil {
		return nil, 404, fmt.Errorf("studio not found")
	}
	teamID, err := resolveUserTeamID(ctx, userIDStr)
	if err != nil {
		if errors.Is(err, ErrNoTeam) {
			return nil, 403, err
		}
		return nil, 500, err
	}
	if studio.TeamID != teamID {
		return nil, 403, fmt.Errorf("no access to studio %d", studioID)
	}
	return studio, 200, nil
}

// assertPieceAccess loads the piece and its studio and authorizes the user.
func assertPieceAccess(ctx context.Context, userIDStr string, pieceID int64) (*entities.Piece, *entities.Studio, int, error) {
	piece, err := repositories.GetPieceByID(ctx, pieceID)
	if err != nil {
		return nil, nil, 500, err
	}
	if piece == nil {
		return nil, nil, 404, fmt.Errorf("piece not found")
	}
	studio, code, err := assertStudioAccess(ctx, userIDStr, piece.StudioID)
	if err != nil {
		return nil, nil, code, err
	}
	return piece, studio, 200, nil
}

// ----------------------------------------------------------------------------
// Studio
// ----------------------------------------------------------------------------

// EnsureStudioForUser returns the user's default studio, creating one (plus
// default phases) if none exists. This is the entry point the frontend calls
// when opening the Studio tab.
//
// If the user has no team yet (new Firebase sign-up that never went through the
// explicit team-creation flow), we lazily provision a personal team and admin
// membership so the Studio is usable out of the box.
func EnsureStudioForUser(ctx context.Context, userIDStr string) (*modelsv1.StudioResponse, int, error) {
	teamID, err := resolveUserTeamID(ctx, userIDStr)
	if err != nil {
		if errors.Is(err, ErrNoTeam) {
			teamID, err = provisionDefaultTeamForUser(ctx, userIDStr)
			if err != nil {
				return nil, 500, err
			}
		} else {
			return nil, 500, err
		}
	}
	existing, err := repositories.GetStudioByTeamID(ctx, teamID)
	if err != nil {
		return nil, 500, err
	}
	if existing != nil {
		return mapStudioToResponse(existing), 200, nil
	}

	studio := &entities.Studio{
		TeamID:     teamID,
		Name:       "Studio",
		PieceLabel: "Piece",
	}
	if err := repositories.CreateStudio(ctx, studio); err != nil {
		return nil, 500, err
	}
	if err := seedDefaultPhases(ctx, studio.ID); err != nil {
		return nil, 500, err
	}
	return mapStudioToResponse(studio), 200, nil
}

// UpdateStudio edits the studio's name or piece label.
func UpdateStudio(ctx context.Context, userIDStr string, studioID int64, req *modelsv1.UpdateStudioRequest) (*modelsv1.StudioResponse, int, error) {
	studio, code, err := assertStudioAccess(ctx, userIDStr, studioID)
	if err != nil {
		return nil, code, err
	}
	name := req.Name
	if name == "" {
		name = studio.Name
	}
	label := req.PieceLabel
	if label == "" {
		label = studio.PieceLabel
	}
	if err := repositories.UpdateStudio(ctx, studioID, name, label); err != nil {
		return nil, 500, err
	}
	studio.Name, studio.PieceLabel = name, label
	return mapStudioToResponse(studio), 200, nil
}

// GetStudioBoard returns the full board payload: studio metadata, ordered
// phases, and active pieces.
func GetStudioBoard(ctx context.Context, userIDStr string, studioID int64) (*modelsv1.StudioBoardResponse, int, error) {
	studio, code, err := assertStudioAccess(ctx, userIDStr, studioID)
	if err != nil {
		return nil, code, err
	}
	phases, err := repositories.ListPhasesByStudio(ctx, studioID)
	if err != nil {
		return nil, 500, err
	}
	pieces, err := repositories.ListPiecesByStudio(ctx, repositories.ListPiecesFilters{
		StudioID: studioID,
		Status:   string(entities.PieceStatusActive),
	})
	if err != nil {
		return nil, 500, err
	}

	phaseResps := make([]modelsv1.PhaseResponse, 0, len(phases))
	for i := range phases {
		phaseResps = append(phaseResps, *mapPhaseToResponse(&phases[i]))
	}
	pieceResps := make([]modelsv1.PieceResponse, 0, len(pieces))
	for i := range pieces {
		pieceResps = append(pieceResps, *mapPieceToResponse(&pieces[i]))
	}
	return &modelsv1.StudioBoardResponse{
		Studio: *mapStudioToResponse(studio),
		Phases: phaseResps,
		Pieces: pieceResps,
	}, 200, nil
}

// ----------------------------------------------------------------------------
// Phase
// ----------------------------------------------------------------------------

// seedDefaultPhases creates the canonical idea→published pipeline for a new studio.
func seedDefaultPhases(ctx context.Context, studioID int64) error {
	tpls := studiosvc.DefaultPhaseTemplates()
	for i, t := range tpls {
		p := &entities.Phase{
			StudioID:   studioID,
			Name:       t.Name,
			Slug:       t.Slug,
			OrderIndex: i,
			Kind:       string(t.Kind),
			IsDefault:  t.IsDefault,
			Color:      t.Color,
		}
		if err := repositories.CreatePhase(ctx, p); err != nil {
			return fmt.Errorf("seed phase %q: %w", t.Name, err)
		}
	}
	return nil
}

func ListPhases(ctx context.Context, userIDStr string, studioID int64) ([]modelsv1.PhaseResponse, int, error) {
	if _, code, err := assertStudioAccess(ctx, userIDStr, studioID); err != nil {
		return nil, code, err
	}
	phases, err := repositories.ListPhasesByStudio(ctx, studioID)
	if err != nil {
		return nil, 500, err
	}
	out := make([]modelsv1.PhaseResponse, 0, len(phases))
	for i := range phases {
		out = append(out, *mapPhaseToResponse(&phases[i]))
	}
	return out, 200, nil
}

func CreatePhase(ctx context.Context, userIDStr string, studioID int64, req *modelsv1.CreatePhaseRequest) (*modelsv1.PhaseResponse, int, error) {
	if _, code, err := assertStudioAccess(ctx, userIDStr, studioID); err != nil {
		return nil, code, err
	}
	existing, err := repositories.ListPhasesByStudio(ctx, studioID)
	if err != nil {
		return nil, 500, err
	}
	kind := req.Kind
	if kind == "" {
		kind = string(entities.PhaseKindCustom)
	}
	slug := req.Slug
	if slug == "" {
		slug = slugify(req.Name)
	}
	color := req.Color
	if color == "" {
		color = "#64748b"
	}
	p := &entities.Phase{
		StudioID:   studioID,
		Name:       req.Name,
		Slug:       slug,
		OrderIndex: len(existing),
		Kind:       kind,
		WIPLimit:   req.WIPLimit,
		Color:      color,
	}
	if err := repositories.CreatePhase(ctx, p); err != nil {
		return nil, 500, err
	}
	return mapPhaseToResponse(p), 200, nil
}

func UpdatePhase(ctx context.Context, userIDStr string, phaseID int64, req *modelsv1.UpdatePhaseRequest) (*modelsv1.PhaseResponse, int, error) {
	phase, err := repositories.GetPhaseByID(ctx, phaseID)
	if err != nil {
		return nil, 500, err
	}
	if phase == nil {
		return nil, 404, fmt.Errorf("phase not found")
	}
	if _, code, err := assertStudioAccess(ctx, userIDStr, phase.StudioID); err != nil {
		return nil, code, err
	}
	name := req.Name
	if name == "" {
		name = phase.Name
	}
	slug := req.Slug
	if slug == "" {
		slug = phase.Slug
	}
	kind := req.Kind
	if kind == "" {
		kind = phase.Kind
	}
	color := req.Color
	if color == "" {
		color = phase.Color
	}
	wip := phase.WIPLimit
	if req.WIPLimit != nil {
		wip = req.WIPLimit
	}
	if err := repositories.UpdatePhase(ctx, phaseID, name, slug, kind, color, wip); err != nil {
		return nil, 500, err
	}
	phase.Name, phase.Slug, phase.Kind, phase.Color, phase.WIPLimit = name, slug, kind, color, wip
	return mapPhaseToResponse(phase), 200, nil
}

func ReorderPhases(ctx context.Context, userIDStr string, studioID int64, req *modelsv1.ReorderPhasesRequest) (int, error) {
	if _, code, err := assertStudioAccess(ctx, userIDStr, studioID); err != nil {
		return code, err
	}
	if err := repositories.ReorderPhases(ctx, studioID, req.PhaseIDs); err != nil {
		return 500, err
	}
	return 200, nil
}

func DeletePhase(ctx context.Context, userIDStr string, phaseID int64) (int, error) {
	phase, err := repositories.GetPhaseByID(ctx, phaseID)
	if err != nil {
		return 500, err
	}
	if phase == nil {
		return 404, fmt.Errorf("phase not found")
	}
	if _, code, err := assertStudioAccess(ctx, userIDStr, phase.StudioID); err != nil {
		return code, err
	}
	count, err := repositories.CountPiecesInPhase(ctx, phaseID)
	if err != nil {
		return 500, err
	}
	if count > 0 {
		return 409, fmt.Errorf("cannot delete phase with %d active pieces; move them first", count)
	}
	if err := repositories.DeletePhase(ctx, phaseID); err != nil {
		return 500, err
	}
	return 200, nil
}

// ----------------------------------------------------------------------------
// Piece
// ----------------------------------------------------------------------------

// CreatePiece adds a new piece to a phase, appending it to the end of the column.
func CreatePiece(ctx context.Context, userIDStr string, studioID int64, req *modelsv1.CreatePieceRequest) (*modelsv1.PieceResponse, int, error) {
	if _, code, err := assertStudioAccess(ctx, userIDStr, studioID); err != nil {
		return nil, code, err
	}
	phase, err := repositories.GetPhaseByID(ctx, req.PhaseID)
	if err != nil {
		return nil, 500, err
	}
	if phase == nil || phase.StudioID != studioID {
		return nil, 400, fmt.Errorf("invalid phase_id for this studio")
	}

	creator, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	var assignee *uuid.UUID
	if req.AssigneeUserID != nil && *req.AssigneeUserID != "" {
		u, err := uuid.Parse(*req.AssigneeUserID)
		if err != nil {
			return nil, 400, fmt.Errorf("invalid assignee_user_id")
		}
		assignee = &u
	}
	var dueAt *time.Time
	if req.DueAt != nil && *req.DueAt != "" {
		t, err := time.Parse(time.RFC3339, *req.DueAt)
		if err != nil {
			return nil, 400, fmt.Errorf("invalid due_at; expected RFC3339")
		}
		dueAt = &t
	}

	contentType := req.ContentType
	if contentType == "" {
		contentType = string(entities.PieceContentTypePost)
	}

	position, err := nextPositionInPhase(ctx, req.PhaseID)
	if err != nil {
		return nil, 500, err
	}

	p := &entities.Piece{
		StudioID:       studioID,
		PhaseID:        req.PhaseID,
		Title:          req.Title,
		Description:    req.Description,
		ContentType:    contentType,
		Status:         string(entities.PieceStatusActive),
		AssigneeUserID: assignee,
		DueAt:          dueAt,
		Position:       position,
		CreatedBy:      creator,
	}
	if err := repositories.CreatePiece(ctx, p); err != nil {
		return nil, 500, err
	}
	logPieceActivity(ctx, p.ID, &creator, entities.PieceActivityCreated, map[string]any{
		"phase_id": p.PhaseID,
		"title":    p.Title,
	})
	return mapPieceToResponse(p), 200, nil
}

func GetPieceDetail(ctx context.Context, userIDStr string, pieceID int64) (*modelsv1.PieceDetailResponse, int, error) {
	piece, _, code, err := assertPieceAccess(ctx, userIDStr, pieceID)
	if err != nil {
		return nil, code, err
	}
	assets, err := repositories.ListPieceAssets(ctx, pieceID)
	if err != nil {
		return nil, 500, err
	}
	sched, err := repositories.ListPieceScheduledPosts(ctx, pieceID)
	if err != nil {
		return nil, 500, err
	}
	acts, err := repositories.ListPieceActivities(ctx, pieceID, 50)
	if err != nil {
		return nil, 500, err
	}
	comments, err := repositories.ListPieceComments(ctx, pieceID)
	if err != nil {
		return nil, 500, err
	}

	out := &modelsv1.PieceDetailResponse{
		Piece:          *mapPieceToResponse(piece),
		Assets:         make([]modelsv1.PieceAssetResponse, 0, len(assets)),
		ScheduledPosts: make([]modelsv1.PieceScheduledPostResponse, 0, len(sched)),
		Activities:     make([]modelsv1.PieceActivityResponse, 0, len(acts)),
		Comments:       make([]modelsv1.PieceCommentResponse, 0, len(comments)),
	}
	for _, a := range assets {
		out.Assets = append(out.Assets, modelsv1.PieceAssetResponse{
			MediaAssetID: a.MediaAssetID, Role: a.Role, AddedAt: modelsv1.FormatTime(a.AddedAt),
		})
	}
	for _, s := range sched {
		out.ScheduledPosts = append(out.ScheduledPosts, modelsv1.PieceScheduledPostResponse{
			ScheduledPostID: s.ScheduledPostID, Role: s.Role, CreatedAt: modelsv1.FormatTime(s.CreatedAt),
		})
	}
	for _, a := range acts {
		out.Activities = append(out.Activities, mapActivityToResponse(&a))
	}
	for _, c := range comments {
		out.Comments = append(out.Comments, mapCommentToResponse(&c))
	}
	return out, 200, nil
}

func UpdatePiece(ctx context.Context, userIDStr string, pieceID int64, req *modelsv1.UpdatePieceRequest) (*modelsv1.PieceResponse, int, error) {
	piece, _, code, err := assertPieceAccess(ctx, userIDStr, pieceID)
	if err != nil {
		return nil, code, err
	}
	title := req.Title
	if title == "" {
		title = piece.Title
	}
	description := piece.Description
	if req.Description != nil {
		description = *req.Description
	}
	contentType := req.ContentType
	if contentType == "" {
		contentType = piece.ContentType
	}
	assignee := piece.AssigneeUserID
	if req.AssigneeUserID != nil {
		if *req.AssigneeUserID == "" {
			assignee = nil
		} else {
			u, err := uuid.Parse(*req.AssigneeUserID)
			if err != nil {
				return nil, 400, fmt.Errorf("invalid assignee_user_id")
			}
			assignee = &u
		}
	}
	var due sql.NullTime
	if req.DueAt != nil {
		if *req.DueAt == "" {
			due = sql.NullTime{Valid: false}
		} else {
			t, err := time.Parse(time.RFC3339, *req.DueAt)
			if err != nil {
				return nil, 400, fmt.Errorf("invalid due_at; expected RFC3339")
			}
			due = sql.NullTime{Time: t, Valid: true}
		}
	} else if piece.DueAt != nil {
		due = sql.NullTime{Time: *piece.DueAt, Valid: true}
	}

	if err := repositories.UpdatePieceFields(ctx, pieceID, title, description, contentType, assignee, &due); err != nil {
		return nil, 500, err
	}
	actor := parseUserID(userIDStr)
	logPieceActivity(ctx, pieceID, actor, entities.PieceActivityUpdated, map[string]any{
		"title": title, "content_type": contentType,
	})
	piece.Title, piece.Description, piece.ContentType, piece.AssigneeUserID = title, description, contentType, assignee
	if due.Valid {
		t := due.Time
		piece.DueAt = &t
	} else {
		piece.DueAt = nil
	}
	return mapPieceToResponse(piece), 200, nil
}

// MovePiece relocates a piece to a target phase/position. The caller provides
// either BeforePieceID (drop above that piece) or AfterPieceID (drop below);
// if neither is given the piece is appended to the end of the target phase.
func MovePiece(ctx context.Context, userIDStr string, pieceID int64, req *modelsv1.MovePieceRequest) (*modelsv1.PieceResponse, int, error) {
	piece, studio, code, err := assertPieceAccess(ctx, userIDStr, pieceID)
	if err != nil {
		return nil, code, err
	}
	targetPhase, err := repositories.GetPhaseByID(ctx, req.PhaseID)
	if err != nil {
		return nil, 500, err
	}
	if targetPhase == nil || targetPhase.StudioID != studio.ID {
		return nil, 400, fmt.Errorf("invalid phase_id for this studio")
	}

	// WIP-limit enforcement (only when moving into a new phase).
	if piece.PhaseID != req.PhaseID && targetPhase.WIPLimit != nil && *targetPhase.WIPLimit > 0 {
		count, err := repositories.CountPiecesInPhase(ctx, req.PhaseID)
		if err != nil {
			return nil, 500, err
		}
		if count >= *targetPhase.WIPLimit {
			return nil, 409, fmt.Errorf("phase %q is at WIP limit (%d)", targetPhase.Name, *targetPhase.WIPLimit)
		}
	}

	prev, next, code, err := resolveMoveBounds(ctx, pieceID, req)
	if err != nil {
		return nil, code, err
	}
	position, err := studiosvc.GeneratePositionBetween(prev, next)
	if err != nil {
		return nil, 500, fmt.Errorf("compute position: %w", err)
	}

	if err := repositories.MovePiece(ctx, pieceID, req.PhaseID, position); err != nil {
		return nil, 500, err
	}
	actor := parseUserID(userIDStr)
	logPieceActivity(ctx, pieceID, actor, entities.PieceActivityMoved, map[string]any{
		"from_phase_id": piece.PhaseID,
		"to_phase_id":   req.PhaseID,
		"position":      position,
	})
	piece.PhaseID, piece.Position = req.PhaseID, position
	return mapPieceToResponse(piece), 200, nil
}

// resolveMoveBounds determines the (prev, next) position strings flanking the
// intended insertion slot. The moving piece is excluded from the neighbors so
// it can be repositioned within its current phase without conflicting with
// itself.
func resolveMoveBounds(ctx context.Context, movingPieceID int64, req *modelsv1.MovePieceRequest) (prev, next string, code int, err error) {
	pieces, err := repositories.ListPiecesByStudio(ctx, repositories.ListPiecesFilters{
		StudioID: 0, // filled below
		PhaseID:  &req.PhaseID,
		Status:   string(entities.PieceStatusActive),
	})
	// The repo requires studio_id; fall back to fetching by phase id via a
	// quick lookup if the above query returned nothing or errored with studio=0.
	// To keep the API simple we use a direct listing helper: find the phase's
	// studio first.
	if err != nil || len(pieces) == 0 {
		phase, ferr := repositories.GetPhaseByID(ctx, req.PhaseID)
		if ferr != nil {
			return "", "", 500, ferr
		}
		if phase == nil {
			return "", "", 400, fmt.Errorf("invalid phase_id")
		}
		pieces, err = repositories.ListPiecesByStudio(ctx, repositories.ListPiecesFilters{
			StudioID: phase.StudioID,
			PhaseID:  &req.PhaseID,
			Status:   string(entities.PieceStatusActive),
		})
		if err != nil {
			return "", "", 500, err
		}
	}
	// Drop the moving piece from the neighbor list.
	filtered := pieces[:0]
	for _, p := range pieces {
		if p.ID != movingPieceID {
			filtered = append(filtered, p)
		}
	}

	switch {
	case req.BeforePieceID != nil:
		for i, p := range filtered {
			if p.ID == *req.BeforePieceID {
				next = p.Position
				if i > 0 {
					prev = filtered[i-1].Position
				}
				return prev, next, 200, nil
			}
		}
		return "", "", 400, fmt.Errorf("before_piece_id not found in target phase")
	case req.AfterPieceID != nil:
		for i, p := range filtered {
			if p.ID == *req.AfterPieceID {
				prev = p.Position
				if i+1 < len(filtered) {
					next = filtered[i+1].Position
				}
				return prev, next, 200, nil
			}
		}
		return "", "", 400, fmt.Errorf("after_piece_id not found in target phase")
	default:
		if len(filtered) > 0 {
			prev = filtered[len(filtered)-1].Position
		}
		return prev, "", 200, nil
	}
}

// nextPositionInPhase returns a position that appends to the end of a phase.
func nextPositionInPhase(ctx context.Context, phaseID int64) (string, error) {
	positions, err := repositories.GetPhasePositions(ctx, phaseID)
	if err != nil {
		return "", err
	}
	var last string
	if len(positions) > 0 {
		last = positions[len(positions)-1]
	}
	return studiosvc.GeneratePositionBetween(last, "")
}

func SetPieceStatus(ctx context.Context, userIDStr string, pieceID int64, status string) (int, error) {
	piece, _, code, err := assertPieceAccess(ctx, userIDStr, pieceID)
	if err != nil {
		return code, err
	}
	if piece.Status == status {
		return 200, nil
	}
	if err := repositories.SetPieceStatus(ctx, pieceID, status); err != nil {
		return 500, err
	}
	actor := parseUserID(userIDStr)
	kind := entities.PieceActivityArchived
	if status == string(entities.PieceStatusActive) {
		kind = entities.PieceActivityRestored
	}
	logPieceActivity(ctx, pieceID, actor, kind, map[string]any{"status": status})
	return 200, nil
}

func DeletePiece(ctx context.Context, userIDStr string, pieceID int64) (int, error) {
	if _, _, code, err := assertPieceAccess(ctx, userIDStr, pieceID); err != nil {
		return code, err
	}
	if err := repositories.DeletePiece(ctx, pieceID); err != nil {
		return 500, err
	}
	return 200, nil
}

// ----------------------------------------------------------------------------
// Piece <-> MediaAsset linking
// ----------------------------------------------------------------------------

func LinkPieceAsset(ctx context.Context, userIDStr string, pieceID int64, req *modelsv1.LinkAssetRequest) (int, error) {
	if _, _, code, err := assertPieceAccess(ctx, userIDStr, pieceID); err != nil {
		return code, err
	}
	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return 400, fmt.Errorf("invalid user")
	}
	asset, err := repositories.GetMediaAssetByID(ctx, req.MediaAssetID, ownerID)
	if err != nil {
		return 500, err
	}
	if asset == nil {
		return 404, fmt.Errorf("media asset not found")
	}
	role := req.Role
	if role == "" {
		role = string(entities.PieceAssetRoleAttachment)
	}
	if err := repositories.LinkPieceAsset(ctx, pieceID, req.MediaAssetID, role); err != nil {
		return 500, err
	}
	actor := parseUserID(userIDStr)
	logPieceActivity(ctx, pieceID, actor, entities.PieceActivityAssetLinked, map[string]any{
		"media_asset_id": req.MediaAssetID, "role": role,
	})
	return 200, nil
}

func UnlinkPieceAsset(ctx context.Context, userIDStr string, pieceID, mediaAssetID int64, role string) (int, error) {
	if _, _, code, err := assertPieceAccess(ctx, userIDStr, pieceID); err != nil {
		return code, err
	}
	if err := repositories.UnlinkPieceAsset(ctx, pieceID, mediaAssetID, role); err != nil {
		return 500, err
	}
	actor := parseUserID(userIDStr)
	logPieceActivity(ctx, pieceID, actor, entities.PieceActivityAssetUnlinked, map[string]any{
		"media_asset_id": mediaAssetID, "role": role,
	})
	return 200, nil
}

// ----------------------------------------------------------------------------
// Piece <-> ScheduledPost linking
// ----------------------------------------------------------------------------

func LinkPieceScheduledPost(ctx context.Context, userIDStr string, pieceID int64, req *modelsv1.LinkScheduledPostRequest) (int, error) {
	piece, studio, code, err := assertPieceAccess(ctx, userIDStr, pieceID)
	if err != nil {
		return code, err
	}
	role := req.Role
	if role == "" {
		role = string(entities.PieceScheduledPostRolePrimary)
	}
	if err := repositories.LinkPieceScheduledPost(ctx, pieceID, req.ScheduledPostID, role); err != nil {
		return 500, err
	}
	actor := parseUserID(userIDStr)
	logPieceActivity(ctx, pieceID, actor, entities.PieceActivityScheduled, map[string]any{
		"scheduled_post_id": req.ScheduledPostID, "role": role,
	})

	// Auto-transition the piece to the studio's "scheduled" phase if one exists
	// and the piece isn't already there (or further down the pipeline).
	scheduledPhase, err := repositories.GetPhaseByStudioKind(ctx, studio.ID, string(entities.PhaseKindScheduled))
	if err == nil && scheduledPhase != nil && piece.PhaseID != scheduledPhase.ID {
		pos, perr := nextPositionInPhase(ctx, scheduledPhase.ID)
		if perr == nil {
			_ = repositories.MovePiece(ctx, pieceID, scheduledPhase.ID, pos)
			logPieceActivity(ctx, pieceID, actor, entities.PieceActivityMoved, map[string]any{
				"from_phase_id": piece.PhaseID,
				"to_phase_id":   scheduledPhase.ID,
				"reason":        "auto_scheduled",
			})
		}
	}
	return 200, nil
}

func UnlinkPieceScheduledPost(ctx context.Context, userIDStr string, pieceID, scheduledPostID int64) (int, error) {
	if _, _, code, err := assertPieceAccess(ctx, userIDStr, pieceID); err != nil {
		return code, err
	}
	if err := repositories.UnlinkPieceScheduledPost(ctx, pieceID, scheduledPostID); err != nil {
		return 500, err
	}
	return 200, nil
}

// OnScheduledPostPublished is called by the publishing pipeline when a
// scheduled_post finishes publishing. It finds the linked piece (if any) and
// transitions it into the studio's "published" phase, recording an activity.
func OnScheduledPostPublished(ctx context.Context, scheduledPostID int64, success bool) error {
	pieceID, err := repositories.GetPieceIDByScheduledPost(ctx, scheduledPostID)
	if err != nil {
		return err
	}
	if pieceID == 0 {
		return nil // not linked to any piece; nothing to do
	}
	piece, err := repositories.GetPieceByID(ctx, pieceID)
	if err != nil || piece == nil {
		return err
	}

	if !success {
		logPieceActivity(ctx, pieceID, nil, entities.PieceActivityPublishFailed, map[string]any{
			"scheduled_post_id": scheduledPostID,
		})
		return nil
	}

	publishedPhase, err := repositories.GetPhaseByStudioKind(ctx, piece.StudioID, string(entities.PhaseKindPublished))
	if err != nil {
		return err
	}
	if publishedPhase != nil && piece.PhaseID != publishedPhase.ID {
		pos, perr := nextPositionInPhase(ctx, publishedPhase.ID)
		if perr == nil {
			_ = repositories.MovePiece(ctx, pieceID, publishedPhase.ID, pos)
		}
	}
	logPieceActivity(ctx, pieceID, nil, entities.PieceActivityPublished, map[string]any{
		"scheduled_post_id": scheduledPostID,
	})
	return nil
}

// ----------------------------------------------------------------------------
// Comments
// ----------------------------------------------------------------------------

func CreatePieceComment(ctx context.Context, userIDStr string, pieceID int64, req *modelsv1.CreatePieceCommentRequest) (*modelsv1.PieceCommentResponse, int, error) {
	if _, _, code, err := assertPieceAccess(ctx, userIDStr, pieceID); err != nil {
		return nil, code, err
	}
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, 400, fmt.Errorf("invalid user")
	}
	c := &entities.PieceComment{
		PieceID:  pieceID,
		UserID:   uid,
		Body:     req.Body,
		ParentID: req.ParentID,
	}
	if err := repositories.CreatePieceComment(ctx, c); err != nil {
		return nil, 500, err
	}
	logPieceActivity(ctx, pieceID, &uid, entities.PieceActivityCommented, map[string]any{
		"comment_id": c.ID,
	})
	resp := mapCommentToResponse(c)
	return &resp, 200, nil
}

func ListPieceComments(ctx context.Context, userIDStr string, pieceID int64) ([]modelsv1.PieceCommentResponse, int, error) {
	if _, _, code, err := assertPieceAccess(ctx, userIDStr, pieceID); err != nil {
		return nil, code, err
	}
	comments, err := repositories.ListPieceComments(ctx, pieceID)
	if err != nil {
		return nil, 500, err
	}
	out := make([]modelsv1.PieceCommentResponse, 0, len(comments))
	for i := range comments {
		out = append(out, mapCommentToResponse(&comments[i]))
	}
	return out, 200, nil
}

func DeletePieceComment(ctx context.Context, userIDStr string, pieceID, commentID int64) (int, error) {
	if _, _, code, err := assertPieceAccess(ctx, userIDStr, pieceID); err != nil {
		return code, err
	}
	if err := repositories.DeletePieceComment(ctx, commentID); err != nil {
		return 500, err
	}
	return 200, nil
}

// ----------------------------------------------------------------------------
// Activity log
// ----------------------------------------------------------------------------

func ListPieceActivities(ctx context.Context, userIDStr string, pieceID int64, limit int) ([]modelsv1.PieceActivityResponse, int, error) {
	if _, _, code, err := assertPieceAccess(ctx, userIDStr, pieceID); err != nil {
		return nil, code, err
	}
	acts, err := repositories.ListPieceActivities(ctx, pieceID, limit)
	if err != nil {
		return nil, 500, err
	}
	out := make([]modelsv1.PieceActivityResponse, 0, len(acts))
	for i := range acts {
		out = append(out, mapActivityToResponse(&acts[i]))
	}
	return out, 200, nil
}

// logPieceActivity is a best-effort audit-log writer. Failures are intentionally
// swallowed so they never break the caller's primary operation; real failures
// still surface via logs at the infrastructure layer.
func logPieceActivity(ctx context.Context, pieceID int64, actor *uuid.UUID, kind entities.PieceActivityKind, payload map[string]any) {
	data, _ := json.Marshal(payload)
	if len(data) == 0 {
		data = []byte("{}")
	}
	_ = repositories.CreatePieceActivity(ctx, &entities.PieceActivity{
		PieceID:     pieceID,
		ActorUserID: actor,
		Kind:        string(kind),
		Payload:     data,
	})
}

// ----------------------------------------------------------------------------
// helpers
// ----------------------------------------------------------------------------

func parseUserID(s string) *uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &u
}

// slugify produces a lowercase, hyphenated slug from a display name. The DB has
// a uniqueness constraint on (studio_id, slug); collision handling is deferred
// to callers (currently we rely on validation errors to bubble up).
func slugify(s string) string {
	out := make([]byte, 0, len(s))
	prevDash := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'A' && c <= 'Z':
			out = append(out, c+32)
			prevDash = false
		case (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'):
			out = append(out, c)
			prevDash = false
		default:
			if !prevDash && len(out) > 0 {
				out = append(out, '-')
				prevDash = true
			}
		}
	}
	// Trim trailing dash
	for len(out) > 0 && out[len(out)-1] == '-' {
		out = out[:len(out)-1]
	}
	if len(out) == 0 {
		return "phase"
	}
	return string(out)
}

func mapStudioToResponse(s *entities.Studio) *modelsv1.StudioResponse {
	return &modelsv1.StudioResponse{
		ID:         s.ID,
		TeamID:     s.TeamID.String(),
		Name:       s.Name,
		PieceLabel: s.PieceLabel,
		CreatedAt:  modelsv1.FormatTime(s.CreatedAt),
		UpdatedAt:  modelsv1.FormatTime(s.UpdatedAt),
	}
}

func mapPhaseToResponse(p *entities.Phase) *modelsv1.PhaseResponse {
	return &modelsv1.PhaseResponse{
		ID:         p.ID,
		StudioID:   p.StudioID,
		Name:       p.Name,
		Slug:       p.Slug,
		OrderIndex: p.OrderIndex,
		Kind:       p.Kind,
		WIPLimit:   p.WIPLimit,
		IsDefault:  p.IsDefault,
		Color:      p.Color,
		CreatedAt:  modelsv1.FormatTime(p.CreatedAt),
		UpdatedAt:  modelsv1.FormatTime(p.UpdatedAt),
	}
}

func mapPieceToResponse(p *entities.Piece) *modelsv1.PieceResponse {
	resp := &modelsv1.PieceResponse{
		ID:          p.ID,
		StudioID:    p.StudioID,
		PhaseID:     p.PhaseID,
		Title:       p.Title,
		Description: p.Description,
		ContentType: p.ContentType,
		Status:      p.Status,
		Position:    p.Position,
		CreatedBy:   p.CreatedBy.String(),
		CreatedAt:   modelsv1.FormatTime(p.CreatedAt),
		UpdatedAt:   modelsv1.FormatTime(p.UpdatedAt),
	}
	if p.AssigneeUserID != nil {
		s := p.AssigneeUserID.String()
		resp.AssigneeUserID = &s
	}
	if p.DueAt != nil {
		s := modelsv1.FormatTime(*p.DueAt)
		resp.DueAt = &s
	}
	if p.ArchivedAt != nil {
		s := modelsv1.FormatTime(*p.ArchivedAt)
		resp.ArchivedAt = &s
	}
	return resp
}

func mapActivityToResponse(a *entities.PieceActivity) modelsv1.PieceActivityResponse {
	var payload any
	_ = json.Unmarshal(a.Payload, &payload)
	r := modelsv1.PieceActivityResponse{
		ID:        a.ID,
		PieceID:   a.PieceID,
		Kind:      a.Kind,
		Payload:   payload,
		CreatedAt: modelsv1.FormatTime(a.CreatedAt),
	}
	if a.ActorUserID != nil {
		s := a.ActorUserID.String()
		r.ActorUserID = &s
	}
	return r
}

func mapCommentToResponse(c *entities.PieceComment) modelsv1.PieceCommentResponse {
	return modelsv1.PieceCommentResponse{
		ID:        c.ID,
		PieceID:   c.PieceID,
		UserID:    c.UserID.String(),
		Body:      c.Body,
		ParentID:  c.ParentID,
		CreatedAt: modelsv1.FormatTime(c.CreatedAt),
		UpdatedAt: modelsv1.FormatTime(c.UpdatedAt),
	}
}
