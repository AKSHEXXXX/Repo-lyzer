package github

import (
	"fmt"
	"time"
)

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Number            int        `json:"number"`
	Title             string     `json:"title"`
	State             string     `json:"state"`
	CreatedAt         time.Time  `json:"created_at"`
	MergedAt          *time.Time `json:"merged_at"`
	ClosedAt          *time.Time `json:"closed_at"`
	User              User       `json:"user"`
	Draft             bool       `json:"draft"`
	Additions         int        `json:"additions"`
	Deletions         int        `json:"deletions"`
	ChangedFiles      int        `json:"changed_files"`
	AuthorAssociation string     `json:"author_association"`
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
// Handles pagination to retrieve all reviews (GitHub default is 30 per page)
func (c *Client) GetPullRequestReviews(owner, repo string, prNumber int) ([]Review, error) {
	var allReviews []Review

	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/pulls/%d/reviews?per_page=%d&page=%d",
			owner, repo, prNumber, perPage, page,
		)

		var reviews []Review
		err := c.get(url, &reviews)
		if err != nil {
			return nil, err
		}

		// Stop when no more reviews
		if len(reviews) == 0 {
			break
		}

		allReviews = append(allReviews, reviews...)

		// Stop when fewer than per_page (last page)
		if len(reviews) < perPage {
			break
		}

		page++
	}

	return allReviews, nil
}

// GetPullRequestDetails fetches detailed information for a specific PR
// This endpoint includes additions, deletions, and changed_files which are
// not available in the list endpoint
func (c *Client) GetPullRequestDetails(owner, repo string, prNumber int) (*PullRequest, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/pulls/%d",
		owner, repo, prNumber,
	)

	var pr PullRequest
	err := c.get(url, &pr)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}
