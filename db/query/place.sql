-- name: ListPlaces :many
SELECT name,
       ST_X(ST_TRANSFORM(geom, 4326)) as longitude,
       ST_Y(ST_TRANSFORM(geom, 4326)) as latitude,
       category_id
from place
order by id;

-- name: GetPlace :one
SELECT name,
       ST_ASTEXT(geom)                     as geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) as geom_geographic
from place
where name = $1;

-- name: GetPlaceGeom :one
SELECT ST_ASTEXT(geom)                     as geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) as geom_geographic
from place
where name = $1;