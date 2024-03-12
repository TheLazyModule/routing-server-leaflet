-- name: ListNodes :many
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS geom_geographic
FROM node
ORDER BY id;

-- name: GetNodesByIds :many
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS geom_geographic
FROM node
WHERE id = ANY ($1);

-- name: GetNodeByID :one
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS geom_geographic
FROM node
WHERE id = $1;

-- name: GetNodePointGeom :one
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS geom_geographic
FROM node
WHERE id = $1;

-- name: ListNodePointGeoms :many
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS geom_geographic
FROM node;

-- name: GetClosestPointToQueryLocation :one
SELECT id,
       name,
       ST_ASTEXT(geom)                     AS closest_geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) AS closest_geom_geographic
FROM node
ORDER BY geom <-> ST_GEOMFROMTEXT($1, 3857)
LIMIT 1;

