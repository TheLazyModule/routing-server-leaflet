-- name: GetBuildingOrPlace :one
WITH Combined AS
         (
             SELECT name, ST_AsText(geom) as geom
             FROM place
             WHERE place.name = @name::text

UNION

SELECT name, ST_AsText(ST_Centroid(geom)) as geom
FROM building
WHERE building.name = @name::text

UNION

SELECT CASE
           WHEN classroom.room_code IS NOT NULL THEN classroom.room_code
           ELSE classroom.name
           END || ' (' || building.name || ')' as name,
       ST_AsText(ST_Centroid(building.geom)) as geom
FROM classroom
         JOIN building ON classroom.building_id = building.id
WHERE classroom.room_code = @name::text OR classroom.name = @name::text
)
SELECT name, geom
FROM Combined
         LIMIT 1;


-- name: GetClosestPointToQueryLocationByLatLngGeom :one
SELECT id,
       name,
       ST_ASTEXT(geom) AS closest_geom
FROM node
ORDER BY geom <-> st_transform(ST_GEOMFROMTEXT($1, 4326), 3857) LIMIT 1;

-- name: FuzzyFindPlaceBuildingClassroom :many
WITH Combined AS
         (
             SELECT id, name, geom, category_id, NULL::TEXT[] AS image_urls
             FROM place
             WHERE name ILIKE '%' || @text::text || '%'

UNION

SELECT id, name, ST_Centroid(geom) as geom, category_id, image_urls
FROM building
WHERE name ILIKE '%' || @text::text || '%'

UNION

SELECT classroom.id,
       COALESCE(classroom.room_code, classroom.name) as name,
       ST_Centroid(building.geom) as geom,
       classroom.category_id,
       NULL::TEXT[] AS image_urls
FROM classroom
         JOIN building ON classroom.building_id = building.id
WHERE classroom.room_code ILIKE '%' || @text::text || '%'
       OR classroom.name ILIKE '%' || @text::text || '%'
)
SELECT id,
       name,
       category_id,
       ST_AsText(ST_Transform(geom, 4326)) as geom,
       image_urls
FROM Combined;
