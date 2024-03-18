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

-- name: GetPlaceByNameOrGeom :one
select name            as placeName,
       st_astext(geom) as PlaceGeom
from place
where place.name = $1
   or place.geom = (select geom as Geom
                    from place
                    order by geom <-> st_geomfromtext($2, 3857)
                    limit 1);
