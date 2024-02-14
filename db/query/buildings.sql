-- name: ListBuildings :many
SELECT name,
       ST_ASTEXT(geom)           as geom,
       ST_ASTEXT(geom_geography) as geom_geographic
from buildings
order by id;

-- name: GetBuildingCentroidGeom :one
SELECT ST_ASTEXT(centroid)           as building_centroid,
       ST_ASTEXT(centroid_geography) as building_centroid_geographic
from buildings
where name = $1;

