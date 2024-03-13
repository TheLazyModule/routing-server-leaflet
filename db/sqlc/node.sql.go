// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: node.sql

package db

import (
	"context"
)

const getClosestPointToQueryLocation = `-- name: GetClosestPointToQueryLocation :one
SELECT id,
       name,
       ST_ASTEXT(geom)                     AS closest_geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS closest_geom_geographic
FROM node
ORDER BY geom <-> ST_GEOMFROMTEXT($1, 3857)
LIMIT 1
`

type GetClosestPointToQueryLocationRow struct {
	ID                    int64       `json:"id"`
	Name                  string      `json:"name"`
	ClosestGeom           interface{} `json:"closest_geom"`
	ClosestGeomGeographic interface{} `json:"closest_geom_geographic"`
}

func (q *Queries) GetClosestPointToQueryLocation(ctx context.Context, stGeomfromtext interface{}) (GetClosestPointToQueryLocationRow, error) {
	row := q.db.QueryRow(ctx, getClosestPointToQueryLocation, stGeomfromtext)
	var i GetClosestPointToQueryLocationRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ClosestGeom,
		&i.ClosestGeomGeographic,
	)
	return i, err
}

const getNodeByID = `-- name: GetNodeByID :one
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS geom_geographic
FROM node
WHERE id = $1
`

type GetNodeByIDRow struct {
	Name           string      `json:"name"`
	Geom           interface{} `json:"geom"`
	GeomGeographic interface{} `json:"geom_geographic"`
}

func (q *Queries) GetNodeByID(ctx context.Context, id int64) (GetNodeByIDRow, error) {
	row := q.db.QueryRow(ctx, getNodeByID, id)
	var i GetNodeByIDRow
	err := row.Scan(&i.Name, &i.Geom, &i.GeomGeographic)
	return i, err
}

const getNodePointGeom = `-- name: GetNodePointGeom :one
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS geom_geographic
FROM node
WHERE id = $1
`

type GetNodePointGeomRow struct {
	Name           string      `json:"name"`
	Geom           interface{} `json:"geom"`
	GeomGeographic interface{} `json:"geom_geographic"`
}

func (q *Queries) GetNodePointGeom(ctx context.Context, id int64) (GetNodePointGeomRow, error) {
	row := q.db.QueryRow(ctx, getNodePointGeom, id)
	var i GetNodePointGeomRow
	err := row.Scan(&i.Name, &i.Geom, &i.GeomGeographic)
	return i, err
}

const getNodesByIds = `-- name: GetNodesByIds :many
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS geom_geographic
FROM node
WHERE id = ANY ($1)
ORDER BY ARRAY_POSITION($1, id)
`

type GetNodesByIdsRow struct {
	Name           string      `json:"name"`
	Geom           interface{} `json:"geom"`
	GeomGeographic interface{} `json:"geom_geographic"`
}

func (q *Queries) GetNodesByIds(ctx context.Context, id []int64) ([]GetNodesByIdsRow, error) {
	rows, err := q.db.Query(ctx, getNodesByIds, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetNodesByIdsRow{}
	for rows.Next() {
		var i GetNodesByIdsRow
		if err := rows.Scan(&i.Name, &i.Geom, &i.GeomGeographic); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listNodePointGeoms = `-- name: ListNodePointGeoms :many
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS geom_geographic
FROM node
`

type ListNodePointGeomsRow struct {
	Name           string      `json:"name"`
	Geom           interface{} `json:"geom"`
	GeomGeographic interface{} `json:"geom_geographic"`
}

func (q *Queries) ListNodePointGeoms(ctx context.Context) ([]ListNodePointGeomsRow, error) {
	rows, err := q.db.Query(ctx, listNodePointGeoms)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListNodePointGeomsRow{}
	for rows.Next() {
		var i ListNodePointGeomsRow
		if err := rows.Scan(&i.Name, &i.Geom, &i.GeomGeographic); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listNodes = `-- name: ListNodes :many
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS geom_geographic
FROM node
ORDER BY id
`

type ListNodesRow struct {
	Name           string      `json:"name"`
	Geom           interface{} `json:"geom"`
	GeomGeographic interface{} `json:"geom_geographic"`
}

func (q *Queries) ListNodes(ctx context.Context) ([]ListNodesRow, error) {
	rows, err := q.db.Query(ctx, listNodes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListNodesRow{}
	for rows.Next() {
		var i ListNodesRow
		if err := rows.Scan(&i.Name, &i.Geom, &i.GeomGeographic); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
