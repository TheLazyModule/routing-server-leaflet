-- name: ListBuildings :many
SELECT name,
       ST_ASTEXT(geom)                     as geom,
       st_astext(st_transform(geom, 4326)) as geom_geographic
from buildings
order by id;

-- name: GetBuildingNames: many
SELECT name
from buildings;

-- name: GetBuildingCentroidGeom :one
SELECT st_centroid(geom)                   as building_centroid,
       st_astext(st_transform(geom, 4326)) as building_centroid_geographic
from buildings
where name = $1;

