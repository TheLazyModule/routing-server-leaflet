package db

import (
	"github.com/jackc/pgx/v5/pgtype"
	db "routing/db/sqlc"
)

type ReqID struct {
	ID int64 `uri:"id" binding:"required"`
}

type RouteRequestByID struct {
	FromNodeID int64 `json:"from_node_id" binding:"required,min=1"`
	ToNodeID   int64 `json:"to_node_id" binding:"required,min=1"`
}

type RouteRequest struct {
	From         pgtype.Text `json:"from" binding:"required"`
	FromLocation pgtype.Text `json:"from_location"`
	To           pgtype.Text `json:"to" binding:"required"`
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

type PlacesResult struct {
	Places []db.ListPlacesRow
	Err    error
}

type BuildingsResult struct {
	Buildings []db.ListBuildingsRow
	Err       error
}

type ClosestNodeResult struct {
	Node db.GetClosestPointToQueryLocationRow
	Err  error
}

type ClosestNodeResultToUserLocation struct {
	Node db.GetClosestPointToQueryLocationByLatLngGeomRow
	Err  error
}

type DijkstraResult struct {
	Paths    []int64
	Distance float64
	Err      error
}
