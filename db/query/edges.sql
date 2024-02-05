-- name: ListEdges :many
SELECT node_id, neighbors
FROM edges;

-- name: GetEdges :one
SELECT *
FROM edges
where node_id = $1;
