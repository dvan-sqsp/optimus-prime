package repositories

import (
	"context"
	"database/sql"
	"errors"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"encore.dev/types/uuid"
)

type Dao struct {
	db *sql.DB
}

func NewDao(db *sql.DB) *Dao {
	return &Dao{db: db}
}

func (s *Dao) list(ctx context.Context) ([]Entity, error) {
	query := `
		SELECT id, name, owner FROM repositories;
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errs.B().Code(errs.Internal).Msg("error retrieving data").Err()
	}

	defer rows.Close()

	var repos []Entity
	for rows.Next() {
		var entity Entity
		if err = rows.Scan(&entity.ID, &entity.Name, &entity.Owner); err != nil {
			return nil, errs.B().Code(errs.Internal).Msg("error scanning row").Cause(err).Err()
		}
		repos = append(repos, entity)
	}

	return repos, nil
}

func (s *Dao) add(ctx context.Context, owner string, repoName string) (*Entity, error) {
	query := `
		INSERT INTO repositories (name, owner)
		VALUES ($1, $2) RETURNING id
	`

	entity := &Entity{
		Name:  repoName,
		Owner: owner,
	}

	args := []interface{}{repoName, owner}
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&entity.ID)
	if err != nil {
		rlog.Error("error inserting repo", "err", err)
		return nil, err
	}

	return entity, err
}

func (s *Dao) getRepoByNameAndOwner(ctx context.Context, name string, owner string) (*Repo, error) {
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

func (s *Dao) getRepoByID(ctx context.Context, id uuid.UUID) (*Repo, error) {
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

func (s *Dao) delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM repositories WHERE id = $1;
	`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		rlog.Error("error deleting repo", "err", err)
		return err
	}

	return nil
}
