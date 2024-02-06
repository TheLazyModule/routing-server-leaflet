package api

type getWeightRequest struct {
	FromNodeID int64 `json:"from_node_id" binding:"required"`
	ToNodeID   int64 `json:"to_node_id"`
}

type routeRequest struct {
	FromNode string `json:"from_node" binding:"required"`
	ToNode   string `json:"to_node" binding:"required"`
}

type routeRequestByID struct {
	FromNodeID int64 `json:"from_node_id" binding:"required,min=1"`
	ToNodeID   int64 `json:"to_node_id" binding:"required,min=1"`
}

type EdgesData []string
