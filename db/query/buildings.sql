-- name: ListBuildings :many
SELECT name, ST_ASTEXT(geom) as geom
from buildings
order by id;

