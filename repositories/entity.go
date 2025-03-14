package repositories

import (
	"encore.dev/types/uuid"
)

type Entity struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Owner string    `db:"owner"`
}

func (e *Entity) ToModel() *Repo {
	return &Repo{
		ID:    e.ID,
		Name:  e.Name,
		Owner: e.Owner,
	}
}
