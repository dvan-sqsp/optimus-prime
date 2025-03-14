package pull_requests

import (
	"context"
	"database/sql"

	ghclient "encore.app/gh_client"
	"encore.app/repositories"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

//encore:service
type PRService struct {
	client ghclient.Client
	db     *sql.DB
}

func initPRService() (*PRService, error) {
	ghClient := ghclient.NewClient()
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

	prs, err := s.client.GetPullRequests(ctx, owner, name)
	if err != nil {
		rlog.Error("error fetching PRs", "err", err)
		return nil, errs.B().Code(errs.Internal).Cause(err).Err()
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

	return &Response{PRs: prsToSave}, nil
}

// this has to be called package level
var db = sqldb.NewDatabase("pull_requests", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
