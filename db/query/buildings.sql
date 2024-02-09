-- name: ListBuildings :many
SELECT name, ST_ASTEXT(geom) as geom, ST_ASTEXT(geom_geography) as geom_geographic
from buildings
order by id;

