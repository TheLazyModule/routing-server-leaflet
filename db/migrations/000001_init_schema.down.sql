-- Drop indexes for the 'classroom.sql' table
DROP INDEX IF EXISTS "idx_classroom_room_code";
DROP INDEX IF EXISTS "idx_classroom_building_id";

-- Drop indexes for the 'building' table
DROP INDEX IF EXISTS "idx_building_geom";
DROP INDEX IF EXISTS "idx_building_name";

-- Drop indexes for the 'place' table
DROP INDEX IF EXISTS "idx_place_location";
DROP INDEX IF EXISTS "idx_place_name";

-- Drop indexes for the 'edge' table
DROP INDEX IF EXISTS "idx_edge_weight";
DROP INDEX IF EXISTS "idx_edge_from_to";

-- Drop indexes for the 'node' table
DROP INDEX IF EXISTS "idx_node_geom";

-- Drop the 'classroom.sql' table
DROP TABLE IF EXISTS "classroom.sql";

-- Drop the 'building' table
DROP TABLE IF EXISTS "building";

-- Drop the 'place' table
DROP TABLE IF EXISTS "place";

-- Drop the 'edge' table
DROP TABLE IF EXISTS "edge";

-- Drop the 'node' table
DROP TABLE IF EXISTS "node";

