-- name: ListPlaces :many
SELECT name, ST_ASTEXT(location) as location
from places
order by id;

-- name: GetPlace :one
SELECT name, ST_ASTEXT(location) as location
from places
where name = $1;


-- name: GetPlaceGeom :one
SELECT ST_ASTEXT(location) as location
from places
where name = $1;
