-- name: ListPlaces :many
SELECT name, ST_ASTEXT(location) as location, ST_ASTEXT(location_geography) as location_geographic
from places
order by id;

-- name: GetPlace :one
SELECT name, ST_ASTEXT(location) as location, ST_ASTEXT(location_geography) as location_geographic
from places
where name = $1;


-- name: GetPlaceGeom :one
SELECT ST_ASTEXT(location) as location, ST_ASTEXT(location_geography) as location_geographic
from places
where name = $1;
