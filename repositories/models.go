package repositories

import (
	"encore.dev/types/uuid"
)

type AddParams struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type DeleteParams struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type Repo struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Owner string    `json:"owner"`
}

type Response struct {
	Repos []Repo `json:"repos"`
}
