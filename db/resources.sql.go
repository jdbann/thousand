// Code generated by sqlc. DO NOT EDIT.
// source: resources.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createResource = `-- name: CreateResource :one
INSERT INTO resources (vampire_id, description, stationary)
    VALUES ($1, $2, $3)
RETURNING
    id, vampire_id, description, stationary, created_at, updated_at
`

type CreateResourceParams struct {
	VampireID   uuid.UUID
	Description string
	Stationary  bool
}

func (q *Queries) CreateResource(ctx context.Context, arg CreateResourceParams) (Resource, error) {
	row := q.db.QueryRow(ctx, createResource, arg.VampireID, arg.Description, arg.Stationary)
	var i Resource
	err := row.Scan(
		&i.ID,
		&i.VampireID,
		&i.Description,
		&i.Stationary,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getResourcesForVampire = `-- name: GetResourcesForVampire :many
SELECT
    resources.id, resources.vampire_id, resources.description, resources.stationary, resources.created_at, resources.updated_at
FROM
    resources
WHERE
    resources.vampire_id = $1
`

func (q *Queries) GetResourcesForVampire(ctx context.Context, vampireID uuid.UUID) ([]Resource, error) {
	rows, err := q.db.Query(ctx, getResourcesForVampire, vampireID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Resource
	for rows.Next() {
		var i Resource
		if err := rows.Scan(
			&i.ID,
			&i.VampireID,
			&i.Description,
			&i.Stationary,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
