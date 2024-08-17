// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: place.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getPlace = `-- name: GetPlace :one
SELECT name,
       ST_ASTEXT(geom)                     as geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) as geom_geographic
from place
where name = $1
`

type GetPlaceRow struct {
	Name           string      `json:"name"`
	Geom           interface{} `json:"geom"`
	GeomGeographic interface{} `json:"geom_geographic"`
}

func (q *Queries) GetPlace(ctx context.Context, name string) (GetPlaceRow, error) {
	row := q.db.QueryRow(ctx, getPlace, name)
	var i GetPlaceRow
	err := row.Scan(&i.Name, &i.Geom, &i.GeomGeographic)
	return i, err
}

const getPlaceGeom = `-- name: GetPlaceGeom :one
SELECT ST_ASTEXT(geom)                     as geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) as geom_geographic
from place
where name = $1
`

type GetPlaceGeomRow struct {
	Geom           interface{} `json:"geom"`
	GeomGeographic interface{} `json:"geom_geographic"`
}

func (q *Queries) GetPlaceGeom(ctx context.Context, name string) (GetPlaceGeomRow, error) {
	row := q.db.QueryRow(ctx, getPlaceGeom, name)
	var i GetPlaceGeomRow
	err := row.Scan(&i.Geom, &i.GeomGeographic)
	return i, err
}

const listPlaces = `-- name: ListPlaces :many
SELECT name,
       ST_X(ST_TRANSFORM(geom, 4326)) as longitude,
       ST_Y(ST_TRANSFORM(geom, 4326)) as latitude,
       category_id
from place
order by id
`

type ListPlacesRow struct {
	Name       string      `json:"name"`
	Longitude  interface{} `json:"longitude"`
	Latitude   interface{} `json:"latitude"`
	CategoryID pgtype.Int4 `json:"category_id"`
}

func (q *Queries) ListPlaces(ctx context.Context) ([]ListPlacesRow, error) {
	rows, err := q.db.Query(ctx, listPlaces)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListPlacesRow{}
	for rows.Next() {
		var i ListPlacesRow
		if err := rows.Scan(
			&i.Name,
			&i.Longitude,
			&i.Latitude,
			&i.CategoryID,
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
