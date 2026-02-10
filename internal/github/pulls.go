package github

import (
	"fmt"
	"time"
)

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Number       int        `json:"number"`
	Title        string     `json:"title"`
	State        string     `json:"state"`
	CreatedAt    time.Time  `json:"created_at"`
	MergedAt     *time.Time `json:"merged_at"`
	ClosedAt     *time.Time `json:"closed_at"`
	User         User       `json:"user"`
	Draft        bool       `json:"draft"`
	Additions    int        `json:"additions"`
	Deletions    int        `json:"deletions"`
	ChangedFiles int        `json:"changed_files"`
}

// Review represents a pull request review
type Review struct {
	User        User      `json:"user"`
	State       string    `json:"state"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// GetPullRequests fetches all pull requests for a repository with pagination
// state can be "open", "closed", or "all"
func (c *Client) GetPullRequests(owner, repo, state string) ([]PullRequest, error) {
	var allPRs []PullRequest

	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/pulls?state=%s&per_page=%d&page=%d&sort=created&direction=desc",
			owner, repo, state, perPage, page,
		)

		var prs []PullRequest
		err := c.get(url, &prs)
		if err != nil {
			return nil, err
		}

		// Stop when no more pull requests
		if len(prs) == 0 {
			break
		}

		allPRs = append(allPRs, prs...)

		// Stop when fewer than per_page (last page)
		if len(prs) < perPage {
			break
		}

		page++
	}

	return allPRs, nil
}

// GetPullRequestsWithLimit fetches pull requests with a maximum limit
func (c *Client) GetPullRequestsWithLimit(owner, repo, state string, limit int) ([]PullRequest, error) {
	var allPRs []PullRequest

	page := 1
	perPage := 100
	if limit < perPage {
		perPage = limit
	}

	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/pulls?state=%s&per_page=%d&page=%d&sort=created&direction=desc",
			owner, repo, state, perPage, page,
		)

		var prs []PullRequest
		err := c.get(url, &prs)
		if err != nil {
			return nil, err
		}

		// Stop when no more pull requests
		if len(prs) == 0 {
			break
		}

		allPRs = append(allPRs, prs...)

		// Stop when we've reached the limit
		if len(allPRs) >= limit {
			if len(allPRs) > limit {
				allPRs = allPRs[:limit]
			}
			break
		}

		// Stop when fewer than per_page (last page)
		if len(prs) < perPage {
			break
		}

		page++
	}

	return allPRs, nil
}

// GetPullRequestReviews fetches all reviews for a specific pull request
func (c *Client) GetPullRequestReviews(owner, repo string, prNumber int) ([]Review, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/pulls/%d/reviews",
		owner, repo, prNumber,
	)

	var reviews []Review
	err := c.get(url, &reviews)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
