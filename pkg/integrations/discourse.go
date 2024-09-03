package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/harish-dalal/feedback-ingestion-system/pkg/models"
)

type DiscourseIntegration struct{}

func NewDiscourseStrategy() *DiscourseIntegration {
	return &DiscourseIntegration{}
}

func (s *DiscourseIntegration) Pull(ctx context.Context, sub *models.Subscription) ([]*models.Feedback, error) {
	lastPulled := sub.LastPulled.Format("2006-01-02")
	now := time.Now().Format("2006-01-02")

	url := fmt.Sprintf("https://meta.discourse.org/search.json?page=1&q=after%%3A%s+before%%3A%s", lastPulled, now)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var searchResults struct {
		Posts []struct {
			ID         int    `json:"id"`
			TopicID    int    `json:"topic_id"`
			CreatedAt  string `json:"created_at"`
			Blurb      string `json:"blurb"`
			Username   string `json:"username"`
			TopicTitle string `json:"topic_title_headline"`
		} `json:"posts"`
	}

	if err := json.Unmarshal(body, &searchResults); err != nil {
		return nil, fmt.Errorf("failed to unmarshal discourse search results: %v", err)
	}

	searchResults.Posts = searchResults.Posts[:3]

	var (
		feedbacks []*models.Feedback
		mutex     sync.Mutex
		wg        sync.WaitGroup
	)

	for _, post := range searchResults.Posts {
		wg.Add(1)
		go func(post struct {
			ID         int    `json:"id"`
			TopicID    int    `json:"topic_id"`
			CreatedAt  string `json:"created_at"`
			Blurb      string `json:"blurb"`
			Username   string `json:"username"`
			TopicTitle string `json:"topic_title_headline"`
		}) {
			defer wg.Done()

			// Create a new context with a specific timeout for each post
			postCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			feedback, err := s.processPullPost(postCtx, post.ID, post.TopicID, sub)
			if err != nil {
				fmt.Printf("Failed to process post ID %d: %v\n", post.ID, err)
				return
			}

			// Safely append to feedbacks slice
			mutex.Lock()
			feedbacks = append(feedbacks, feedback...)
			mutex.Unlock()
		}(post)
	}

	wg.Wait()
	fmt.Println("Successfully pulled and processed data from Discourse")
	return feedbacks, nil
}

func (s *DiscourseIntegration) processPullPost(ctx context.Context, postID, topicID int, sub *models.Subscription) ([]*models.Feedback, error) {
	url := fmt.Sprintf("https://meta.discourse.org/t/%d/posts.json?post_ids[]=%d", topicID, postID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch post, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var postResponse struct {
		PostStream struct {
			Posts []struct {
				ID        int    `json:"id"`
				CreatedAt string `json:"created_at"`
				Cooked    string `json:"cooked"`
				Username  string `json:"username"`
				TopicID   int    `json:"topic_id"`
				TopicSlug string `json:"topic_slug"`
			} `json:"posts"`
		} `json:"post_stream"`
	}

	if err := json.Unmarshal(body, &postResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal post data: %v", err)
	}

	if len(postResponse.PostStream.Posts) == 0 {
		return nil, fmt.Errorf("no posts found in the response")
	}

	post := postResponse.PostStream.Posts[0]

	feedback := &models.Feedback{
		ID:          fmt.Sprintf("%d", post.ID),
		TenantID:    sub.TenantID,
		SubSourceID: sub.SubSourceId,
		Source:      s.GetSourceName(),
		SourceType:  s.GetSourceType(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content: map[string]interface{}{
			"body": post.Cooked,
		},
		Metadata: map[string]interface{}{
			"username":   post.Username,
			"topic_id":   post.TopicID,
			"topic_slug": post.TopicSlug,
			"created_at": post.CreatedAt,
		},
	}

	return []*models.Feedback{feedback}, nil
}

func (s *DiscourseIntegration) Push(ctx context.Context, r *http.Request, body []byte) ([]*models.Feedback, error) {
	// implement webhook logic here
	return nil, fmt.Errorf("discourse push method not implemented yet")
}

func (a *DiscourseIntegration) GetSourceName() models.Source {
	return models.SourceDiscourse
}

func (a *DiscourseIntegration) GetSourceType() models.SourceType {
	return models.STFeedback
}
