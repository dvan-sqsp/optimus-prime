package repositories

import (
	"context"
	"errors"

	"encore.app/client"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"encore.dev/types/uuid"
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
	ghClient client.Client
	dao      *Dao
}

func initRepoService() (*RepoService, error) {
	ghClient := client.NewGithubClient()
	dao := NewDao(db.Stdlib())
	return &RepoService{
		ghClient: ghClient,
		dao:      dao,
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

	_, err = s.ghClient.GetRepository(ctx, p.Owner, p.Name)
	if err != nil {
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
