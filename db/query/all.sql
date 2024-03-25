-- name: GetBuildingOrPlace :one
WITH Combined AS (SELECT name, st_astext(geom) as geom
                  FROM place
                  WHERE place.name = $1

                  UNION

                  SELECT name, st_astext(st_centroid(geom)) as geom
                  FROM building
                  WHERE building.name = $1)
SELECT name, geom
FROM Combined
limit 1;


-- name: GetClosestPointToQueryLocationByLatLngGeom :one
SELECT id,
       name,
       ST_ASTEXT(geom) AS closest_geom
FROM node
ORDER BY geom <-> st_transform(ST_GEOMFROMTEXT($1, 4326), 3857)
LIMIT 1;

