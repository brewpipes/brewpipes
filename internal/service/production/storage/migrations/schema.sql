CREATE TYPE volume_unit AS ENUM ('ml', 'oz');

-- A Volume represents a specific quantity of liquid.
-- They can be split into multiple smaller volumes and combined with other volumes.
CREATE TABLE IF NOT EXISTS volume (
    id           serial PRIMARY KEY,
    name         varchar(255) NOT NULL,
    description  text,
    
)

-- A Vessel represents a container that can hold a volume of liquid.
CREATE TABLE IF NOT EXISTS vessel (
    id             serial PRIMARY KEY,
    type           varchar(255) NOT NULL,
    name           varchar(255) NOT NULL,
    capacity       bigint NOT NULL,
    capacity_unit  volume_unit  NOT NULL
)

-- An Occupancy represents the presence of a specific volume of liquid in a vessel.
CREATE TABLE IF NOT EXISTS occupancy (
    id           serial PRIMARY KEY,
    vessel_id    int NOT NULL REFERENCES vessel(id),
    volume_id    int NOT NULL REFERENCES volume(id),
    in_at        timestamptz NOT NULL DEFAULT timezone('utc', now())
)

-- A Transfer represents the movement of a specific volume of liquid from one vessel to another.
CREATE TABLE IF NOT EXISTS transfer (
    id                   serial PRIMARY KEY,
    source_occupancy_id  int NOT NULL REFERENCES occupancy(id),
    dest_occupancy_id    int NOT NULL REFERENCES occupancy(id),
    amount               bigint NOT NULL,
    amount_unit          volume_unit  NOT NULL,
    started_at           timestamptz NOT NULL DEFAULT timezone('utc', now()),
    ended_at             timestamptz
)

-- A Batch represents a production batch of beer, tracked all the way from brewing to packaging.
-- It can be comprised of multiple worts, and after brewing, it can be transferred entirely to
-- a different vessel, split into multiple sub-batches, combined with other batches, and so forth.
CREATE TABLE IF NOT EXISTS batch (
    id           serial PRIMARY KEY,
    short_name   varchar(255) NOT NULL
)

CREATE TABLE IF NOT EXISTS wort (
    id           serial PRIMARY KEY,
    batch_id     int NOT NULL REFERENCES batch(id),
    volume_id    int NOT NULL REFERENCES volume(id)
)

CREATE TABLE IF NOT EXISTS beer (
    id           serial PRIMARY KEY,
    wort_id      int NOT NULL REFERENCES wort(id),
    volume_id    int NOT NULL REFERENCES volume(id)
)
