-- name: GetBuildingOrPlace :one
WITH Combined AS (SELECT name, st_astext(geom) as geom
                  FROM place
                  WHERE place.name = @name::text

UNION

SELECT name, st_astext(st_centroid(geom)) as geom
FROM building
WHERE building.name = @name::text)
SELECT name, geom
FROM Combined limit 1;


-- name: GetClosestPointToQueryLocationByLatLngGeom :one
SELECT id,
       name,
       ST_ASTEXT(geom) AS closest_geom
FROM node
ORDER BY geom <-> st_transform(ST_GEOMFROMTEXT($1, 4326), 3857) LIMIT 1;

-- name: FuzzyFindPlaceOrBuilding :many
WITH Combined AS
         (SELECT id, name, geom as geom
          FROM place
          WHERE name ILIKE '%' || @text::text || '%'

UNION

SELECT id, name, st_centroid(geom) as geom
FROM building
WHERE name ILIKE '%' || @text::text || '%'
)
SELECT id,
       name,
       ST_ASTEXT(ST_TRANSFORM(geom, 4326)) as geom

FROM Combined;
