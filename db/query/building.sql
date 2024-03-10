-- name: ListBuildings :many
SELECT name, ST_ASTEXT(geom) as geom,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) as geom_geographic
from building
order by id;

-- name: GetBuildingNames :many
SELECT name
from building;

-- name: GetBuildingCentroidGeom :one
SELECT ST_CENTROID(geom)                   as building_centroid,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) as building_centroid_geographic
from building
where name = $1;

