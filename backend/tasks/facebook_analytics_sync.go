package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/provider/meta"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

const (
	facebookPageInsightsLookbackDays = 90
	facebookFeedMaxPosts             = 80
	facebookFeedPageSize             = 25
)

// SyncFacebookChannelAnalytics pulls Page and post insights from the Graph API into facebook_* tables.
func SyncFacebookChannelAnalytics(ctx context.Context, channel entities.Channel) error {
	if channel.Provider != "facebook" {
		return fmt.Errorf("channel %d is not a facebook channel", channel.ID)
	}
	token, err := repositories.GetLatestTokenByChannelID(ctx, channel.ID)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	pageToken, err := encryption.Decrypt(token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to decrypt access token: %w", err)
	}

	pageID := channel.ProviderChannelID
	if pageID == "" {
		return fmt.Errorf("empty provider_channel_id for facebook channel %d", channel.ID)
	}

	prov := meta.NewMetaProvider()

	now := time.Now().UTC()
	since := now.AddDate(0, 0, -facebookPageInsightsLookbackDays).Unix()
	until := now.Unix()

	pageResp, err := prov.GetPageDailyInsights(pageToken, pageID, since, until, meta.PageDailyInsightMetrics)
	if err != nil {
		return fmt.Errorf("page insights: %w", err)
	}

	byDate := meta.InsightValueByMetricEndDate(pageResp)
	for dateKey, m := range byDate {
		d, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}
		d = d.UTC().Truncate(24 * time.Hour)

		row := &entities.FacebookPageAnalytics{
			ChannelID: channel.ID,
			Date:      d,
		}
		if v, ok := m["page_fans"]; ok {
			i := int(v)
			row.Followers = &i
		}
		if v, ok := m["page_fan_adds"]; ok {
			i := int(v)
			row.Likes = &i
		}
		if v, ok := m["page_impressions_unique"]; ok {
			i := int(v)
			row.Reach = &i
		}
		if v, ok := m["page_impressions"]; ok {
			i := int(v)
			row.Impressions = &i
		}
		if v, ok := m["page_engaged_users"]; ok {
			i := int(v)
			row.PageEngagedUsers = &i
		}
		raw, _ := json.Marshal(m)
		row.Raw = raw

		if err := repositories.UpsertFacebookPageAnalytics(ctx, row); err != nil {
			utils.Logger.Error(ctx, fmt.Sprintf("facebook page analytics upsert channel %d date %s: %v", channel.ID, dateKey, err))
		}
	}

	if err := syncFacebookPostInsights(ctx, prov, pageToken, pageID, channel); err != nil {
		return fmt.Errorf("post insights: %w", err)
	}

	utils.Logger.Info(ctx, fmt.Sprintf("Successfully synced Facebook analytics for channel %d", channel.ID))
	return nil
}

func syncFacebookPostInsights(ctx context.Context, prov *meta.MetaProviderImpl, pageToken, pageID string, channel entities.Channel) error {
	var after string
	synced := 0
	for synced < facebookFeedMaxPosts {
		feed, err := prov.ListPageFeedPosts(pageToken, pageID, facebookFeedPageSize, after)
		if err != nil {
			return err
		}
		if len(feed.Data) == 0 {
			break
		}
		for _, fp := range feed.Data {
			if synced >= facebookFeedMaxPosts {
				break
			}
			if err := syncSingleFacebookPost(ctx, prov, pageToken, channel, fp); err != nil {
				utils.Logger.Warn(ctx, fmt.Sprintf("facebook post %s channel %d: %v", fp.ID, channel.ID, err))
			}
			synced++
		}
		if feed.Paging == nil || feed.Paging.Cursors == nil || feed.Paging.Cursors.After == "" {
			break
		}
		after = feed.Paging.Cursors.After
	}
	return nil
}

func syncSingleFacebookPost(ctx context.Context, prov *meta.MetaProviderImpl, pageToken string, channel entities.Channel, fp meta.FeedPost) error {
	postRow, err := repositories.GetPostByProviderKey(ctx, channel.ID, "facebook", fp.ID)
	if err != nil {
		return err
	}
	if postRow == nil {
		publishedAt, perr := parseFacebookTime(fp.CreatedTime)
		if perr != nil {
			publishedAt = time.Now().UTC()
		}
		ids, _ := json.Marshal(map[string]string{"facebook": fp.ID})
		ownerStr := channel.OwnerUserID.String()
		pt := "post"
		caption := fp.Message
		postRow = &entities.Post{
			ChannelIDs:      pq.Int64Array{channel.ID},
			OwnerID:         &ownerStr,
			Providers:       pq.StringArray{"facebook"},
			ProviderPostIDs: ids,
			Source:          "native",
			PostType:        &pt,
			Caption:         &caption,
			PublishedAt:     &publishedAt,
		}
		if err := repositories.CreatePost(ctx, postRow); err != nil {
			return fmt.Errorf("create post: %w", err)
		}
	}

	ins, err := prov.GetPostInsights(pageToken, fp.ID, meta.PostInsightMetrics)
	if err != nil {
		return fmt.Errorf("post insights: %w", err)
	}

	m := postInsightsFirstValue(ins)
	raw, _ := json.Marshal(m)

	a := &entities.FacebookPostAnalytics{
		ChannelID: channel.ID,
		PostID:    postRow.ID,
		Date:      utils.ToUTCDate(time.Now().UTC()),
		Raw:       raw,
	}
	if v, ok := m["post_impressions"]; ok {
		i := int(v)
		a.Impressions = &i
	}
	if v, ok := m["post_impressions_unique"]; ok {
		i := int(v)
		a.Reach = &i
	}
	if v, ok := m["post_reactions_like_total"]; ok {
		i := int(v)
		a.Likes = &i
	}
	if v, ok := m["post_video_views"]; ok {
		i := int(v)
		a.VideoViews = &i
	}

	return repositories.UpsertFacebookPostAnalytics(ctx, a)
}

func postInsightsFirstValue(resp *meta.FBInsightsResponse) map[string]int64 {
	out := make(map[string]int64)
	if resp == nil {
		return out
	}
	for _, series := range resp.Data {
		if len(series.Values) == 0 {
			continue
		}
		if v, ok := meta.ParseScalarInsightValue(series.Values[0].Value); ok {
			out[series.Name] = v
		}
	}
	return out
}

func parseFacebookTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("empty time")
	}
	t, err := time.Parse("2006-01-02T15:04:05-0700", s)
	if err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, s)
}
