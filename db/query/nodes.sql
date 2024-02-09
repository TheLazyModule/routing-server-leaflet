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


-- name: GetNodePointGeom :one
SELECT ST_ASTEXT(point_geom) as point_geom
FROM nodes
WHERE id = $1;


-- name: ListNodePointGeoms :many
SELECT ST_ASTEXT(point_geom) as point_geom
FROM nodes;


-- name: GetClosestPointToQueryLocation :one
SELECT id,name,
       ST_AsText(point_geom) AS closest_geom
FROM nodes
ORDER BY point_geom <-> ST_GeomFromText($1, 3857)
LIMIT 1;

