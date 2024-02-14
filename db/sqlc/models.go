// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twpayne/go-geom"
	dto "routing/db/dto"
)

type Building struct {
	ID                int64       `json:"id"`
	Name              pgtype.Text `json:"name"`
	Geom              geom.Point  `json:"geom"`
	GeomGeography     interface{} `json:"geom_geography"`
	Centroid          geom.Point  `json:"centroid"`
	CentroidGeography interface{} `json:"centroid_geography"`
}

type Classroom struct {
	ID         int64  `json:"id"`
	BuildingID int64  `json:"building_id"`
	Name       string `json:"name"`
}

type Edge struct {
	ID        int64         `json:"id"`
	NodeID    int64         `json:"node_id"`
	Neighbors dto.EdgesData `json:"neighbors"`
}

type Node struct {
	ID                 int64       `json:"id"`
	Name               string      `json:"name"`
	PointGeom          geom.Point  `json:"point_geom"`
	PointGeomGeography interface{} `json:"point_geom_geography"`
}

type Place struct {
	ID                int64       `json:"id"`
	Name              pgtype.Text `json:"name"`
	Location          geom.Point  `json:"location"`
	LocationGeography interface{} `json:"location_geography"`
}

type Weight struct {
	FromNodeID int64   `json:"from_node_id"`
	ToNodeID   int64   `json:"to_node_id"`
	Distance   float64 `json:"distance"`
}
