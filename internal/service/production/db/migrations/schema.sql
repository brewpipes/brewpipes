CREATE TYPE volume_unit AS ENUM ('ml', 'oz');

-- A Volume represents a specific quantity of liquid.
CREATE TABLE IF NOT EXISTS volume (
    id           serial PRIMARY KEY,
    name         varchar(255) NOT NULL,
    description  text,
    amount       bigint NOT NULL,
    amount_unit  volume_unit  NOT NULL,
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
    target_vessel_id     int NOT NULL REFERENCES occupancy(id),
    amount               bigint NOT NULL,
    amount_unit          volume_unit  NOT NULL,
    started_at           timestamptz NOT NULL DEFAULT timezone('utc', now())
    ended_at             timestamptz
)
