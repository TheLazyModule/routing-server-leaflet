-- name: ListBuildings :many
SELECT name,
       ST_X(st_centroid(ST_TRANSFORM(geom, 4326))) as longitude,
       ST_Y(st_centroid(ST_TRANSFORM(geom, 4326))) as latitude,
       image_urls                                  as image_urls,
       category_id
from building
order by id;

-- name: GetBuildingByID :one
SELECT name,
       ST_ASTEXT(st_centroid(geom))        as geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) as geom_geographic,
       image_urls                          as image_urls
from building
where id = $1;

-- name: GetBuildingCentroidGeom :one
SELECT ST_ASTEXT(ST_CENTROID(geom))                     as building_centroid,
       ST_ASTEXT(ST_CENTROID(ST_TRANSFORM(geom, 4326))) as building_centroid_geographic
from building
where name = $1;

