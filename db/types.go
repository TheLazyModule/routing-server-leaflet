package db

import (
	"github.com/jackc/pgx/v5/pgtype"
	db "routing/db/sqlc"
)

type routeRequest struct {
	FromNode string `json:"from_node" binding:"required"`
	ToNode   string `json:"to_node" binding:"required"`
}

type ReqID struct {
	ID int64 `uri:"id" binding:"required"`
}

type RouteRequestByID struct {
	FromNodeID int64 `json:"from_node_id" binding:"required,min=1"`
	ToNodeID   int64 `json:"to_node_id" binding:"required,min=1"`
}

type RouteRequestByPlaceJSON struct {
	From pgtype.Text `json:"from" binding:"required"`
	To   pgtype.Text `json:"to" binding:"required"`
}

type PlaceOrGeomRequest struct {
	Name string      `json:"name"`
	Geom interface{} `json:"geom"`
}

type RouteRequestByBuildingJSON struct {
	From pgtype.Text `json:"from" binding:"required"`
	To   pgtype.Text `json:"to" binding:"required"`
}

type RouteRequestByBuildingOrPlace struct {
	From pgtype.Text `json:"from" binding:"required"`
	To   pgtype.Text `json:"to" binding:"required"`
}

type ClosestNodeResult struct {
	Node db.GetClosestPointToQueryLocationRow
	Err  error
}

type DijkstraResult struct {
	Paths    []int64
	Distance float64
	Err      error
}

type PlacesResult struct {
	Places []db.ListPlacesRow
	Err    error
}

type BuildingResult struct {
	Building []db.ListBuildingsRow
	Err      error
}

type Nodes struct {
	Nodes []db.GetNodesByIdsRow
	Err   error
}

type Neighbours []int64

type Edge struct {
	FromNodeID int64
	ToNodeID   int64
	Weight     float64
}
