package api

import "github.com/jackc/pgx/v5/pgtype"

type getWeightRequest struct {
	FromNodeID int64 `json:"from_node_id" binding:"required"`
	ToNodeID   int64 `json:"to_node_id"`
}

type routeRequest struct {
	FromNode string `json:"from_node" binding:"required"`
	ToNode   string `json:"to_node" binding:"required"`
}

type ReqID struct {
	ID int64 `uri:"id" binding:"required"`
}

type routeRequestByID struct {
	FromNodeID int64 `json:"from_node_id" binding:"required,min=1"`
	ToNodeID   int64 `json:"to_node_id" binding:"required,min=1"`
}

type routeRequestByPlaceForm struct {
	From pgtype.Text `form:"from" binding:"required"`
	To   pgtype.Text `form:"to" binding:"required"`
}

type routeRequestByPlaceOrBuildingJSON struct {
	From pgtype.Text `json:"from" binding:"required"`
	To   pgtype.Text `json:"to" binding:"required"`
}

type EdgesData []string
