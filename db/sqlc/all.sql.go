// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: all.sql

package db

import (
	"context"
)

const fuzzyFindPlaceOrBuilding = `-- name: FuzzyFindPlaceOrBuilding :many
WITH Combined AS (SELECT name, ST_AsText(geom) AS geom
                  FROM place
                  WHERE name ILIKE '%' || $1::text || '%'
UNION
SELECT name, ST_AsText(ST_Centroid(geom)) AS geom
FROM building
WHERE name ILIKE '%' || $1::text || '%'
)
SELECT name, geom
FROM Combined
`

type FuzzyFindPlaceOrBuildingRow struct {
	Name string      `json:"name"`
	Geom interface{} `json:"geom"`
}

func (q *Queries) FuzzyFindPlaceOrBuilding(ctx context.Context, text string) ([]FuzzyFindPlaceOrBuildingRow, error) {
	rows, err := q.db.Query(ctx, fuzzyFindPlaceOrBuilding, text)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FuzzyFindPlaceOrBuildingRow{}
	for rows.Next() {
		var i FuzzyFindPlaceOrBuildingRow
		if err := rows.Scan(&i.Name, &i.Geom); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBuildingOrPlace = `-- name: GetBuildingOrPlace :one
WITH Combined AS (SELECT name, st_astext(geom) as geom
                  FROM place
                  WHERE place.name = $1::text

                  UNION

                  SELECT name, st_astext(st_centroid(geom)) as geom
                  FROM building
                  WHERE building.name = $1::text)
SELECT name, geom
FROM Combined limit 1
`

type GetBuildingOrPlaceRow struct {
	Name string      `json:"name"`
	Geom interface{} `json:"geom"`
}

func (q *Queries) GetBuildingOrPlace(ctx context.Context, name string) (GetBuildingOrPlaceRow, error) {
	row := q.db.QueryRow(ctx, getBuildingOrPlace, name)
	var i GetBuildingOrPlaceRow
	err := row.Scan(&i.Name, &i.Geom)
	return i, err
}

const getClosestPointToQueryLocationByLatLngGeom = `-- name: GetClosestPointToQueryLocationByLatLngGeom :one
SELECT id,
       name,
       ST_ASTEXT(geom) AS closest_geom
FROM node
ORDER BY geom <-> st_transform(ST_GEOMFROMTEXT($1, 4326), 3857) LIMIT 1
`

type GetClosestPointToQueryLocationByLatLngGeomRow struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	ClosestGeom interface{} `json:"closest_geom"`
}

func (q *Queries) GetClosestPointToQueryLocationByLatLngGeom(ctx context.Context, stGeomfromtext interface{}) (GetClosestPointToQueryLocationByLatLngGeomRow, error) {
	row := q.db.QueryRow(ctx, getClosestPointToQueryLocationByLatLngGeom, stGeomfromtext)
	var i GetClosestPointToQueryLocationByLatLngGeomRow
	err := row.Scan(&i.ID, &i.Name, &i.ClosestGeom)
	return i, err
}
