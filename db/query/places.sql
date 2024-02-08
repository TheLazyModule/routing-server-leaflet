-- name: ListPlaces :many
SELECT name, ST_ASTEXT(location) as location
from places
order by id;

