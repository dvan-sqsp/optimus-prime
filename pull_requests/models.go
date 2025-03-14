package pull_requests

import (
	"time"

	"encore.dev/types/uuid"
	"github.com/google/go-github/v69/github"
)

type AddParams struct {
	Owner    string `json:"owner"`
	RepoName string `json:"repo_name"`
}

type PR struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	AvatarURL string    `json:"avatar_url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func entityToModel(pr *github.PullRequest) (*PR, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &PR{
		ID:        uuid,
		Title:     *pr.Title,
		Author:    *pr.User.Login,
		AvatarURL: *pr.User.AvatarURL,
		Status:    *pr.State,
		CreatedAt: pr.CreatedAt.Time,
	}, nil
}

type Response struct {
	PRs []PR `json:"pull_requests"`
}
