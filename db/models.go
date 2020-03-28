// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"time"

	"github.com/google/uuid"
)

type AnidbIndexFile struct {
	ID        uuid.UUID `json:"id"`
	Hash      string    `json:"hash"`
	Source    int32     `json:"source"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
