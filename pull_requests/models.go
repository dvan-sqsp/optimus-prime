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

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type PR struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	AvatarURL string    `json:"avatar_url"`
	Status    string    `json:"status"`
	HTMLURL   string    `json:"html_url"`
	Labels    []Label   `json:"labels"`
	Draft     bool      `json:"draft"`
	CreatedAt time.Time `json:"created_at"`
}

func entityToModel(pr *github.PullRequest) (*PR, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	m := &PR{
		ID:        uuid,
		Title:     *pr.Title,
		Author:    *pr.User.Login,
		AvatarURL: *pr.User.AvatarURL,
		HTMLURL:   *pr.HTMLURL,
		Draft:     *pr.Draft,
		Status:    *pr.State,
		CreatedAt: pr.CreatedAt.Time,
	}

	for _, l := range pr.Labels {
		label := Label{
			Name:  *l.Name,
			Color: *l.Color,
		}
		m.Labels = append(m.Labels, label)
	}
	return m, nil
}

type Response struct {
	PRs []PR `json:"pull_requests"`
}
