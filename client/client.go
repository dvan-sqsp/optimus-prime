package client

import (
	"context"
	"errors"
	"net/http"
	"os"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

var ErrGithubFailure = errors.New("request from github was a non 200")

type Client interface {
	GetPullRequests(ctx context.Context, owner string, repo string) ([]*github.PullRequest, error)
	GetRepository(ctx context.Context, owner string, repo string) (*github.Repository, error)
}

type GithubClient struct {
	ghClient *github.Client
}

func NewGithubClient() *GithubClient {
	token := os.Getenv("GITHUB_TOKEN_ENCORE")
	if token == "" {
		rlog.Error("GITHUB_TOKEN_ENCORE not set")
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	ghClient := github.NewClient(oauth2.NewClient(context.Background(), tokenSource))

	return &GithubClient{
		ghClient: ghClient,
	}
}

func (c *GithubClient) GetPullRequests(ctx context.Context, owner string, repoName string) ([]*github.PullRequest, error) {
	prs, resp, err := c.ghClient.PullRequests.List(ctx, owner, repoName, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrGithubFailure
	}

	return prs, nil
}

func (c *GithubClient) GetRepository(ctx context.Context, owner string, repoName string) (*github.Repository, error) {
	repo, resp, err := c.ghClient.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errs.B().Code(errs.NotFound).Msg("repo not found in github").Err()
		}
		rlog.Error("error fetching repo", "err", err)
		return nil, err
	}

	return repo, nil
}
