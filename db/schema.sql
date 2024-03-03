CREATE
EXTENSION IF NOT EXISTS postgis;

CREATE TABLE "nodes"
(
    "id"                   bigserial PRIMARY KEY,
    "name"                 varchar UNIQUE NOT NULL,
    "point_geom"           geometry(Point, 3857) NOT NULL,
    "point_geom_geography" geography(Point, 4326) NOT NULL
);

CREATE TABLE "edges"
(
    "id"        bigserial PRIMARY KEY,
    "node_id"   bigint NOT NULL,
    "neighbors" jsonb -- '["A", "B", "C"]'
);

CREATE TABLE "weights"
(
    "from_node_id" bigint           NOT NULL,
    "to_node_id"   bigint           NOT NULL,
    "distance"     double precision NOT NULL,
    PRIMARY KEY ("from_node_id", "to_node_id")
);

CREATE TABLE "places"
(
    "id"                 bigserial PRIMARY KEY,
    "name"               varchar,
    "location"           geometry(Point, 3857),
    "location_geography" geography(Point, 4326)
);

CREATE TABLE "buildings"
(
    "id"                 bigserial PRIMARY KEY,
    "name"               varchar,
    "geom"               geometry(Polygon, 3857),
    "geom_geography"     geography(Polygon, 4326),
    "centroid"           geometry(Point, 3857),
    "centroid_geography" geography(Point, 4326) -- Corrected SRID
);

CREATE TABLE "classrooms"
(
    "id"          bigserial PRIMARY KEY,
    "building_id" bigint  NOT NULL,
    "name"        varchar NOT NULL
);

-- indexes
CREATE INDEX ix_nodes_id ON "nodes" USING btree ("id");
CREATE INDEX ix_nodes_point_geom ON "nodes" USING gist ("point_geom");
CREATE INDEX ix_nodes_point_geom_geography ON "nodes" USING gist ("point_geom_geography");
CREATE INDEX ix_edges_node_id ON "edges" USING btree ("node_id");
CREATE INDEX ix_weights_from_node_id ON "weights" USING btree ("from_node_id");
CREATE INDEX ix_weights_to_node_id ON "weights" USING btree ("to_node_id");

CREATE INDEX ix_places_location ON "places" USING gist ("location");
CREATE INDEX ix_places_location_geography ON "places" USING gist ("location_geography");

CREATE INDEX ix_buildings_geom ON "buildings" USING gist ("geom");
CREATE INDEX ix_buildings_geom_geography ON "buildings" USING gist ("geom_geography");
CREATE INDEX ix_buildings_centroid ON "buildings" USING gist ("centroid");
CREATE INDEX ix_buildings_centroid_geography ON "buildings" USING gist ("centroid_geography");

CREATE INDEX ix_classrooms_building_id ON "classrooms" USING btree ("building_id");

-- Foreign Keys
ALTER TABLE "edges"
    ADD FOREIGN KEY ("node_id") REFERENCES "nodes" ("id") ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE "weights"
    ADD FOREIGN KEY ("from_node_id") REFERENCES "nodes" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
    ADD FOREIGN KEY ("to_node_id") REFERENCES "nodes" ("id") ON
UPDATE CASCADE
ON
DELETE
CASCADE;

ALTER TABLE "classrooms"
    ADD FOREIGN KEY ("building_id") REFERENCES "buildings" ("id") ON UPDATE CASCADE ON DELETE RESTRICT; -- Corrected foreign key
