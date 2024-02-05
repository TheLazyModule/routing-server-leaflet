-- name: ListNodes :many
SELECT name, ST_ASTEXT(point_geom) as point_geom
FROM nodes
ORDER BY id;