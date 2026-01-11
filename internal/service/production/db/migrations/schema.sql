CREATE TYPE volume_unit AS ENUM ('ml', 'oz');

CREATE TABLE IF NOT EXISTS vessel (
    id             serial PRIMARY KEY,
    type           varchar(255) NOT NULL,
    name           varchar(255) NOT NULL,
    capacity       bigint NOT NULL,
    capacity_unit  volume_unit  NOT NULL
)

CREATE TABLE IF NOT EXIST volume (
    id           serial PRIMARY KEY,
    name         varchar(255) NOT NULL,
    amount       bigint NOT NULL,
    amount_unit  volume_unit  NOT NULL,
    vessel_id    int NOT NULL REFERENCES vessel(id)
)

CREATE TABLE IF NOT EXISTS transfer (
    id                serial PRIMARY KEY,
    source_volume_id  int NOT NULL REFERENCES volume(id),
    target_volume_id  int NOT NULL REFERENCES volume(id),
    amount            bigint NOT NULL,
    amount_unit       volume_unit  NOT NULL,
    transferred_at    timestamp NOT NULL DEFAULT timezone('utc', now())
)
