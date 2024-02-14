-- name: ListEdges :many
SELECT node_id, neighbors
FROM edges;

-- name: GetEdgeByID :one
SELECT *
FROM edges
where node_id = $1;
