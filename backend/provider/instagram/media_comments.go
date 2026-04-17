package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// MediaCommentsResponse is returned by GET /{ig-media-id}/comments.
type MediaCommentsResponse struct {
	Data   []MediaCommentItem `json:"data"`
	Paging *Paging            `json:"paging,omitempty"`
}

// MediaCommentItem is one comment on a media object.
type MediaCommentItem struct {
	ID        string         `json:"id"`
	Text      string         `json:"text"`
	Timestamp string         `json:"timestamp"`
	Username  string         `json:"username"`
	From      *MediaFromUser `json:"from"`
	ParentID  *string        `json:"parent_id"`
}

// MediaFromUser is the comment author.
type MediaFromUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// GetMediaComments lists top-level comments on an IG media object (paginated with `after`).
func (p *InstagramProviderImpl) GetMediaComments(accessToken, mediaID, after string) (*MediaCommentsResponse, error) {
	q := url.Values{}
	q.Set("fields", "id,text,timestamp,username,from{id,username},parent_id")
	q.Set("access_token", accessToken)
	if after != "" {
		q.Set("after", after)
	}
	reqURL := fmt.Sprintf("https://graph.instagram.com/v24.0/%s/comments?%s", mediaID, q.Encode())

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if ge, ok := GraphAPIErrorFromBody(body); ok {
			return nil, ge
		}
		return nil, fmt.Errorf("get media comments: %s: %s", resp.Status, string(body))
	}

	var out MediaCommentsResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("parse media comments: %w", err)
	}
	return &out, nil
}
