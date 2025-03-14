package repositories

import (
	"context"
	"errors"
	"net/http"
	"os"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
	"encore.dev/types/uuid"
	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

var (
	ErrNotFound      = errors.New("repo not found")
	ErrAlreadyExists = errors.New("repo already exists")
)

// this has to be called package level
var db = sqldb.NewDatabase("repositories", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

//encore:service
type RepoService struct {
	client *github.Client
	dao    *Dao
}

func initRepoService() (*RepoService, error) {
	token := os.Getenv("GITHUB_TOKEN_ENCORE")
	if token == "" {
		rlog.Error("GITHUB_TOKEN_ENCORE not set")
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	ghClient := github.NewClient(oauth2.NewClient(context.Background(), tokenSource))

	dao := NewDao(db.Stdlib())
	return &RepoService{
		client: ghClient,
		dao:    dao,
	}, nil
}

//encore:api public method=GET path=/repos
func (s *RepoService) List(ctx context.Context) (*Response, error) {
	entityRepos, err := s.dao.list(ctx)
	if err != nil {
		return nil, errs.B().Code(errs.Internal).Cause(err).Err()
	}

	var repos []Repo
	for _, e := range entityRepos {
		repos = append(repos, *e.ToModel())
	}

	return &Response{Repos: repos}, nil
}

//encore:api public method=GET path=/repos/:owner/:name
func (s *RepoService) Get(ctx context.Context, owner string, name string) (*Repo, error) {
	repo, err := s.dao.getRepoByNameAndOwner(ctx, name, owner)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, errs.B().Code(errs.NotFound).Cause(err).Err()
		}
		return nil, err
	}

	return repo, nil
}

//encore:api public method=POST path=/repos
func (s *RepoService) Add(ctx context.Context, p *AddParams) (*Repo, error) {
	repo, err := s.dao.getRepoByNameAndOwner(ctx, p.Name, p.Owner)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, errs.B().Code(errs.Internal).Cause(err).Err()
	}

	if repo != nil {
		return nil, errs.B().Code(errs.AlreadyExists).Cause(ErrAlreadyExists).Err()
	}

	_, resp, err := s.client.Repositories.Get(ctx, p.Owner, p.Name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errs.B().Code(errs.NotFound).Msg("repo not found in github").Err()
		}
		rlog.Error("error fetching repo", "err", err)
		return nil, errs.B().Code(errs.Internal).Msg("error requesting github").Cause(err).Err()
	}

	e, err := s.dao.add(ctx, p.Owner, p.Name)
	if err != nil {
		return nil, errs.B().Code(errs.Internal).Cause(err).Err()
	}

	return e.ToModel(), nil
}

//encore:api public method=DELETE path=/repos/:id
func (s *RepoService) Delete(ctx context.Context, id uuid.UUID) (*Repo, error) {
	_, err := s.dao.getRepoByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, errs.B().Code(errs.NotFound).Cause(err).Err()
		}
		return nil, errs.B().Code(errs.Internal).Cause(err).Err()
	}

	err = s.dao.delete(ctx, id)
	if err != nil {
		return nil, errs.B().Code(errs.Internal).Cause(err).Err()
	}

	return nil, nil
}
