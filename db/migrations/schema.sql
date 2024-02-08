CREATE
EXTENSION IF NOT EXISTS postgis;

CREATE TABLE "nodes"
(
    "id"         bigserial PRIMARY KEY,
    "name"       varchar UNIQUE NOT NULL,
    "point_geom" geometry(Point, 3857) NOT NULL
);

CREATE TABLE "edges"
(
    "id"        bigserial PRIMARY KEY,
    "node_id"   bigint UNIQUE NOT NULL,
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
    "id"       bigserial PRIMARY KEY,
    "name"     varchar,
    "location" geometry(point, 3857)
);

CREATE TABLE "buildings"
(
    "id"       bigserial PRIMARY KEY,
    "name"     varchar,
    "geom"     geometry(polygon, 3857),
    "centroid" geometry(point, 3857)
);

CREATE TABLE "classrooms"
(
    "id"          bigserial PRIMARY KEY,
    "building_id" bigint  NOT NULL,
    "name"        varchar NOT NULL
);


-- indexes
CREATE INDEX ix_nodes_id
    ON "nodes" ("id");

CREATE INDEX ix_nodes_point_geom
    ON "nodes" using gist ("point_geom");

CREATE INDEX ix_edges_node_id
    ON "edges" ("node_id");

CREATE INDEX ix_weights_from_node_id
    ON "weights" ("from_node_id");

CREATE INDEX ix_weights_to_node_id
    ON "weights" ("to_node_id");


-- Foreign Keys
ALTER TABLE "edges"
    ADD FOREIGN KEY ("node_id")
        REFERENCES "nodes" ("id")
        ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE "weights"
    ADD FOREIGN KEY ("from_node_id") REFERENCES "nodes" ("id")
        ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE "weights"
    ADD FOREIGN KEY ("to_node_id") REFERENCES "nodes" ("id")
        ON UPDATE CASCADE ON DELETE CASCADE;



CREATE INDEX ON "places" ("id");

CREATE INDEX ON "places" ("name");

CREATE INDEX ON "places" ("location");

CREATE INDEX ON "buildings" ("id");

CREATE INDEX ON "buildings" ("name");

CREATE INDEX ON "buildings" ("geom");

CREATE INDEX ON "buildings" ("centroid");

CREATE INDEX ON "classrooms" ("id");

CREATE INDEX ON "classrooms" ("building_id");

CREATE INDEX ON "classrooms" ("name");

ALTER TABLE "classrooms"
    ADD FOREIGN KEY ("id") REFERENCES "buildings" ("id");
