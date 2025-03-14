package pull_requests

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"

	"encore.app/repositories"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

//encore:service
type PRService struct {
	client *github.Client
	db     *sql.DB
}

func initPRService() (*PRService, error) {
	token := os.Getenv("GITHUB_TOKEN_ENCORE")
	if token == "" {
		rlog.Error("GITHUB_TOKEN_ENCORE not set")
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	ghClient := github.NewClient(oauth2.NewClient(context.Background(), tokenSource))
	return &PRService{
		client: ghClient,
		db:     db.Stdlib(),
	}, nil
}

//encore:api public method=GET path=/pull_requests/:owner/:name
func (s *PRService) List(ctx context.Context, owner string, name string) (*Response, error) {
	_, err := repositories.Get(ctx, owner, name)
	if err != nil {
		if errs.Code(err) == errs.NotFound {
			return nil, errs.B().Code(errs.NotFound).Cause(err).Err()
		}
		return nil, errs.B().Code(errs.Internal).Cause(err).Err()
	}

	prs, resp, err := s.client.PullRequests.List(ctx, owner, name, nil)
	if err != nil {
		rlog.Error("error fetching PRs", "err", err)
		return nil, errs.B().Code(errs.Internal).Cause(err).Err()
	}

	if resp.StatusCode != http.StatusOK {
		rlog.Error("status code not ok", "status", resp.Status)
		return nil, errs.B().Code(errs.Internal).Cause(errors.New("status code not ok from github")).Meta(resp.StatusCode).Err()
	}

	prsToSave := make([]PR, 0, len(prs))
	for _, pr := range prs {
		model, mapErr := entityToModel(pr)
		if mapErr != nil {
			rlog.Error("error mapping entity to model", "err", mapErr)
			continue
		}
		prsToSave = append(prsToSave, *model)
	}

	rlog.Info("github model", "model", prsToSave)

	return &Response{PRs: prsToSave}, nil
}

// this has to be called package level
var db = sqldb.NewDatabase("pull_requests", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
