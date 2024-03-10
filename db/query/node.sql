-- name: ListNodes :many
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_AsText(ST_Transform(geom, 4326)) AS geom_geographic
FROM node
ORDER BY id;

-- name: GetNodesByIds :many
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_AsText(ST_Transform(geom, 4326)) AS geom_geographic
FROM node
WHERE id = ANY ($1);

-- name: GetNodeByID :one
SELECT name,
       ST_ASTEXT(geom)                     AS geom,
       ST_AsText(ST_Transform(geom, 4326)) AS geom_geographic
FROM node
WHERE id = $1;

-- name: GetNodePointGeom :one
SELECT ST_ASTEXT(geom)                     AS geom,
       ST_AsText(ST_Transform(geom, 4326)) AS geom_geographic
FROM node
WHERE id = $1;

-- name: ListNodePointGeoms :many
SELECT ST_ASTEXT(geom)                     AS geom,
       ST_AsText(ST_Transform(geom, 4326)) AS geom_geographic
FROM node;

-- name: GetClosestPointToQueryLocation :one
SELECT id,
       name,
       ST_AsText(geom)                     AS closest_geom,
       ST_AsText(ST_Transform(geom, 4326)) AS closest_geom_geographic
FROM node
ORDER BY geom <-> ST_GeomFromText($1, 3857)
LIMIT 1;
