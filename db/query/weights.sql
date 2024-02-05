-- name: ListWeights :many
SELECT *
FROM weights;


-- name: GetWeight :one
SELECT *
FROM weights
where from_node_id = $1
  and to_node_id = $2;

