-- Remove Foreign Keys
ALTER TABLE "classrooms" DROP CONSTRAINT IF EXISTS classrooms_building_id_fkey;
ALTER TABLE "weights" DROP CONSTRAINT IF EXISTS weights_to_node_id_fkey;
ALTER TABLE "weights" DROP CONSTRAINT IF EXISTS weights_from_node_id_fkey;
ALTER TABLE "edges" DROP CONSTRAINT IF EXISTS edges_node_id_fkey;

-- Drop indexes
DROP INDEX IF EXISTS ix_classrooms_building_id;
DROP INDEX IF EXISTS ix_buildings_centroid_geography;
DROP INDEX IF EXISTS ix_buildings_centroid;
DROP INDEX IF EXISTS ix_buildings_geom_geography;
DROP INDEX IF EXISTS ix_buildings_geom;
DROP INDEX IF EXISTS ix_places_location_geography;
DROP INDEX IF EXISTS ix_places_location;
DROP INDEX IF EXISTS ix_weights_to_node_id;
DROP INDEX IF EXISTS ix_weights_from_node_id;
DROP INDEX IF EXISTS ix_edges_node_id;
DROP INDEX IF EXISTS ix_nodes_point_geom_geography;
DROP INDEX IF EXISTS ix_nodes_point_geom;
DROP INDEX IF EXISTS ix_nodes_id;

-- Drop tables
DROP TABLE IF EXISTS "classrooms";
DROP TABLE IF EXISTS "buildings";
DROP TABLE IF EXISTS "places";
DROP TABLE IF EXISTS "weights";
DROP TABLE IF EXISTS "edges";
DROP TABLE IF EXISTS "nodes";

-- -- Drop extension
-- DROP EXTENSION IF EXISTS postgis;
