package scheduledpublish

import "context"

// OnPostFinalized is an optional hook invoked after a scheduled post reaches a
// terminal state (published or failed). Higher layers (e.g. business/v1) can
// register a callback here to propagate lifecycle events to linked entities
// such as Studio Pieces without creating an import cycle back into this
// package.
//
// The hook signature intentionally mirrors businessv1.OnScheduledPostPublished
// so callers can assign it directly from an init() function.
var OnPostFinalized func(ctx context.Context, scheduledPostID int64, success bool) error
