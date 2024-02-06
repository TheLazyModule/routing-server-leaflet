-- name: ListNodes :many
SELECT name, ST_ASTEXT(point_geom) as point_geom
FROM nodes
ORDER BY id;

-- name: GetNodeAndEdges :one
select *
from nodes
         join edges
              on nodes.id = edges.node_id
where nodes.id = $1;


-- name: GetNodesByIds :many
SELECT name, ST_ASTEXT(point_geom) as point_geom
FROM nodes
WHERE id = ANY ($1);


-- name: GetNodeByID :one
SELECT name, ST_ASTEXT(point_geom) as point_geom
FROM nodes
WHERE id = $1;



