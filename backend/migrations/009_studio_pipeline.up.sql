-- Studio Content Pipeline: Kanban-style board for content creation workflow.
-- A Studio belongs to a Team and contains ordered Phases (columns). Pieces
-- (content items, optionally labeled "Drop") flow through Phases and can be
-- linked to media_assets and scheduled_posts.

-- ============================================================
-- studios: one board per team (MVP: single default studio per team)
-- ============================================================
CREATE TABLE IF NOT EXISTS studios (
    id          BIGSERIAL    PRIMARY KEY,
    team_id     UUID         NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    name        TEXT         NOT NULL DEFAULT 'Studio',
    piece_label TEXT         NOT NULL DEFAULT 'Piece'
                             CHECK (piece_label IN ('Piece', 'Drop') OR char_length(piece_label) BETWEEN 1 AND 32),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_studios_team ON studios (team_id);

-- ============================================================
-- phases: ordered columns within a studio
-- ============================================================
CREATE TABLE IF NOT EXISTS phases (
    id           BIGSERIAL    PRIMARY KEY,
    studio_id    BIGINT       NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
    name         TEXT         NOT NULL,
    slug         TEXT         NOT NULL,
    order_index  INT          NOT NULL DEFAULT 0,
    kind         TEXT         NOT NULL DEFAULT 'custom'
                              CHECK (kind IN ('idea', 'script', 'shoot', 'edit', 'review', 'scheduled', 'published', 'custom')),
    wip_limit    INT,
    is_default   BOOLEAN      NOT NULL DEFAULT FALSE,
    color        TEXT         NOT NULL DEFAULT '#64748b',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_phases_studio_slug UNIQUE (studio_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_phases_studio_order ON phases (studio_id, order_index);
CREATE INDEX IF NOT EXISTS idx_phases_studio_kind  ON phases (studio_id, kind);

-- ============================================================
-- pieces: individual content items on the board
-- position is a fractional-index string for cheap reorders
-- ============================================================
CREATE TABLE IF NOT EXISTS pieces (
    id               BIGSERIAL    PRIMARY KEY,
    studio_id        BIGINT       NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
    phase_id         BIGINT       NOT NULL REFERENCES phases(id) ON DELETE RESTRICT,
    title            TEXT         NOT NULL,
    description      TEXT         NOT NULL DEFAULT '',
    content_type     TEXT         NOT NULL DEFAULT 'post'
                                  CHECK (content_type IN ('post', 'reel', 'story', 'video', 'carousel', 'other')),
    status           TEXT         NOT NULL DEFAULT 'active'
                                  CHECK (status IN ('active', 'archived')),
    assignee_user_id UUID         REFERENCES users(id) ON DELETE SET NULL,
    due_at           TIMESTAMPTZ,
    position         TEXT         NOT NULL DEFAULT 'a0',
    created_by       UUID         NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    archived_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_pieces_studio_phase_pos ON pieces (studio_id, phase_id, position);
CREATE INDEX IF NOT EXISTS idx_pieces_phase            ON pieces (phase_id);
CREATE INDEX IF NOT EXISTS idx_pieces_assignee         ON pieces (assignee_user_id);
CREATE INDEX IF NOT EXISTS idx_pieces_due_at           ON pieces (due_at);
CREATE INDEX IF NOT EXISTS idx_pieces_status           ON pieces (status);

-- ============================================================
-- piece_assets: join to media_assets with a role label
-- ============================================================
CREATE TABLE IF NOT EXISTS piece_assets (
    piece_id       BIGINT       NOT NULL REFERENCES pieces(id) ON DELETE CASCADE,
    media_asset_id BIGINT       NOT NULL REFERENCES media_assets(id) ON DELETE CASCADE,
    role           TEXT         NOT NULL DEFAULT 'attachment'
                                CHECK (role IN ('script', 'raw', 'edit', 'thumbnail', 'final', 'attachment')),
    added_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (piece_id, media_asset_id, role)
);

CREATE INDEX IF NOT EXISTS idx_piece_assets_piece ON piece_assets (piece_id);
CREATE INDEX IF NOT EXISTS idx_piece_assets_media ON piece_assets (media_asset_id);

-- ============================================================
-- piece_scheduled_posts: join to scheduled_posts
-- ============================================================
CREATE TABLE IF NOT EXISTS piece_scheduled_posts (
    piece_id           BIGINT  NOT NULL REFERENCES pieces(id) ON DELETE CASCADE,
    scheduled_post_id  BIGINT  NOT NULL REFERENCES scheduled_posts(id) ON DELETE CASCADE,
    role               TEXT    NOT NULL DEFAULT 'primary'
                               CHECK (role IN ('primary', 'cross_post', 'repost')),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (piece_id, scheduled_post_id)
);

CREATE INDEX IF NOT EXISTS idx_piece_scheduled_posts_piece ON piece_scheduled_posts (piece_id);
CREATE INDEX IF NOT EXISTS idx_piece_scheduled_posts_post  ON piece_scheduled_posts (scheduled_post_id);

-- ============================================================
-- piece_activities: audit log (phase moves, asset links, publishes, etc.)
-- ============================================================
CREATE TABLE IF NOT EXISTS piece_activities (
    id              BIGSERIAL    PRIMARY KEY,
    piece_id        BIGINT       NOT NULL REFERENCES pieces(id) ON DELETE CASCADE,
    actor_user_id   UUID         REFERENCES users(id) ON DELETE SET NULL,
    kind            TEXT         NOT NULL
                                 CHECK (kind IN (
                                     'created', 'updated', 'moved', 'assigned',
                                     'asset_linked', 'asset_unlinked',
                                     'scheduled', 'published', 'publish_failed',
                                     'commented', 'archived', 'restored'
                                 )),
    payload         JSONB        NOT NULL DEFAULT '{}'::jsonb,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_piece_activities_piece ON piece_activities (piece_id, created_at DESC);

-- ============================================================
-- piece_comments: threaded discussion on a piece
-- ============================================================
CREATE TABLE IF NOT EXISTS piece_comments (
    id          BIGSERIAL    PRIMARY KEY,
    piece_id    BIGINT       NOT NULL REFERENCES pieces(id) ON DELETE CASCADE,
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body        TEXT         NOT NULL,
    parent_id   BIGINT       REFERENCES piece_comments(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_piece_comments_piece  ON piece_comments (piece_id, created_at);
CREATE INDEX IF NOT EXISTS idx_piece_comments_parent ON piece_comments (parent_id);
