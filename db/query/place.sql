-- name: ListPlaces :many
SELECT name
     , ST_ASTEXT(location)                     as location
     , ST_ASTEXT(ST_TRANSFORM(location, 4326)) as location_geographic
from place
order by id;

-- name: GetPlace :one
SELECT name,
       ST_ASTEXT(location)                     as location,
       ST_ASTEXT(ST_TRANSFORM(location, 4326)) as location_geographic
from place
where name = $1;


-- name: GetPlaceGeom :one
SELECT ST_ASTEXT(location)                     as location,
       ST_ASTEXT(ST_TRANSFORM(location, 4326)) as location_geographic
from place
where name = $1;
