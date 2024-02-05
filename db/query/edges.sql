-- name: ListEdges :many
SELECT *
FROM edges;

-- name: GetEdges :one
SELECT *
FROM edges
where node_id = $1;
