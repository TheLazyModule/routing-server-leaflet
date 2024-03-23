-- name: ListPlaces :many
SELECT name
     , ST_ASTEXT(geom)                     as geom
     , ST_ASTEXT(ST_TRANSFORM(geom, 4326)) as geom_geographic
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

-- name: GetBuildingOrPlace :one
WITH Combined AS (
    SELECT name, st_astext(geom) as geom
    FROM place
    WHERE place.name = $1

    UNION

    SELECT name, st_astext(geom) as geom
    FROM building
    WHERE building.name = $1
)
SELECT name, geom
FROM Combined;



