package repositories

import (
	"context"
	"database/sql"
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
	db     *sql.DB
}

func initRepoService() (*RepoService, error) {
	token := os.Getenv("GITHUB_TOKEN_ENCORE")
	if token == "" {
		rlog.Error("GITHUB_TOKEN_ENCORE not set")
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	ghClient := github.NewClient(oauth2.NewClient(context.Background(), tokenSource))
	return &RepoService{
		client: ghClient,
		db:     db.Stdlib(),
	}, nil
}

//encore:api public method=GET path=/repos
func (s *RepoService) List(ctx context.Context) (*Response, error) {

	query := `
		SELECT id, name, owner FROM repositories;
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errs.B().Code(errs.Internal).Msg("error retrieving data").Err()
	}

	defer rows.Close()

	var repos []Repo
	for rows.Next() {
		var entity Entity
		if err = rows.Scan(&entity.ID, &entity.Name, &entity.Owner); err != nil {
			return nil, errs.B().Code(errs.Internal).Msg("error scanning row").Cause(err).Err()
		}
		repos = append(repos, *entity.ToModel())
	}

	return &Response{Repos: repos}, nil
}

//encore:api public method=GET path=/repos/:owner/:name
func (s *RepoService) Get(ctx context.Context, owner string, name string) (*Repo, error) {
	repo, err := s.getRepoByNameAndOwner(ctx, name, owner)
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
	repo, err := s.getRepoByNameAndOwner(ctx, p.Name, p.Owner)
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

	newRepo := &Repo{Name: p.Name, Owner: p.Owner}

	query := `
		INSERT INTO repositories (name, owner)
		VALUES ($1, $2) RETURNING id
	`

	args := []interface{}{p.Name, p.Owner}
	err = s.db.QueryRowContext(ctx, query, args...).Scan(&newRepo.ID)
	if err != nil {
		rlog.Error("error inserting repo", "err", err)
		return nil, err
	}
	return newRepo, nil
}

//encore:api public method=DELETE path=/repos/:id
func (s *RepoService) Delete(ctx context.Context, id uuid.UUID) (*Repo, error) {
	_, err := s.getRepoByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, errs.B().Code(errs.NotFound).Cause(err).Err()
		}
		return nil, errs.B().Code(errs.Internal).Cause(err).Err()
	}

	query := `
		DELETE FROM repositories WHERE id = $1;
	`

	_, err = s.db.ExecContext(ctx, query, id)
	if err != nil {
		rlog.Error("error deleting repo", "err", err)
		return nil, err
	}
	return nil, nil
}

func (s *RepoService) getRepoByNameAndOwner(ctx context.Context, name string, owner string) (*Repo, error) {
	query := `SELECT * FROM repositories WHERE name = $1 AND owner = $2;`

	args := []interface{}{name, owner}

	var entity Entity
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&entity.ID, &entity.Name, &entity.Owner)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return entity.ToModel(), nil
}

func (s *RepoService) getRepoByID(ctx context.Context, id uuid.UUID) (*Repo, error) {
	query := `SELECT * FROM repositories WHERE id = $1;`

	var entity Entity
	err := s.db.QueryRowContext(ctx, query, id).Scan(&entity.ID, &entity.Name, &entity.Owner)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return entity.ToModel(), nil
}
