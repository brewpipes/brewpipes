-- A Volume represents a specific quantity of liquid.
-- They can be split into multiple smaller volumes and combined with other volumes.
CREATE TABLE IF NOT EXISTS volume (
    id           serial PRIMARY KEY,
    uuid         uuid NOT NULL DEFAULT gen_random_uuid(),

    name         varchar(255) NOT NULL,
    description  text,
    amount       bigint NOT NULL,
    amount_unit  varchar(7),

    created_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at   timestamptz
);

-- A Vessel represents a container that can hold a volume of liquid.
CREATE TABLE IF NOT EXISTS vessel (
    id             serial PRIMARY KEY,
    uuid           uuid NOT NULL DEFAULT gen_random_uuid(),

    type           varchar(255) NOT NULL,
    name           varchar(255) NOT NULL,
    capacity       bigint NOT NULL,
    capacity_unit  varchar(7)  NOT NULL,
    make           varchar(255),
    model          varchar(255),

    created_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at     timestamptz
);

-- An Occupancy represents the presence of a specific volume of liquid in a vessel.
CREATE TABLE IF NOT EXISTS occupancy (
    id          serial PRIMARY KEY,
    uuid        uuid NOT NULL DEFAULT gen_random_uuid(),

    vessel_id   int NOT NULL REFERENCES vessel(id),
    volume_id   int NOT NULL REFERENCES volume(id),
    in_at       timestamptz NOT NULL DEFAULT timezone('utc', now()),

    created_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at  timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at  timestamptz
);

CREATE TABLE IF NOT EXISTS malt_lot (
    id                  serial PRIMARY KEY,
    uuid                uuid NOT NULL DEFAULT gen_random_uuid(),

    supplier_id         int NOT NULL REFERENCES supplier(id),
    maltster_id         int NOT NULL REFERENCES maltster(id),
    maltster_lot_number varchar(64),
    variety             varchar(64),
    year                int,
    acquired_at         timestamptz,
    amount              bigint NOT NULL,
    amount_unit         varchar(7),

    created_at          timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at          timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at          timestamptz
)

CREATE TABLE IF NOT EXISTS hop_lot (
    id                      serial PRIMARY KEY,
    uuid                    uuid NOT NULL DEFAULT gen_random_uuid(),

    hop_producer_id         int NOT NULL REFERENCES hop_producer(id),
    supplier_id             int NOT NULL REFERENCES supplier(id),
    variety                 varchar(64),
    year                    int,
    hop_producer_lot_number varchar(64),
    acquired_at             timestamptz,
    amount                  bigint NOT NULL,
    amount_unit             varchar(7),

    created_at              timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at              timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at              timestamptz
)

CREATE TABLE IF NOT EXISTS malt_additions (
    id             serial PRIMARY KEY,
    uuid           uuid NOT NULL DEFAULT gen_random_uuid(),

    occupancy_id   int NOT NULL REFERENCES occupancy(id),
    malt_lot_id    int NOT NULL REFERENCES malt_lot(id),
    amount         bigint NOT NULL,
    amount_unit    varchar(7),

    created_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at     timestamptz
)

CREATE TABLE IF NOT EXISTS hop_additions (
    id             serial PRIMARY KEY,
    uuid           uuid NOT NULL DEFAULT gen_random_uuid(),

    occupancy_id   int NOT NULL REFERENCES occupancy(id),
    hop_lot_id     int NOT NULL REFERENCES hop_lot(id),
    amount         bigint NOT NULL,
    amount_unit    varchar(7),

    created_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at     timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at     timestamptz
)

-- A Transfer represents the movement of a specific volume of liquid from one vessel to another.
CREATE TABLE IF NOT EXISTS transfer (
    id                   serial PRIMARY KEY,
    uuid                 uuid NOT NULL DEFAULT gen_random_uuid(),

    source_occupancy_id  int NOT NULL REFERENCES occupancy(id),
    dest_occupancy_id    int NOT NULL REFERENCES occupancy(id),
    started_at           timestamptz NOT NULL DEFAULT timezone('utc', now()),
    ended_at             timestamptz,

    created_at           timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at           timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at           timestamptz
);

-- A Batch represents a production batch of beer, tracked all the way from brewing to packaging.
-- It can be comprised of multiple worts, and after brewing, it can be transferred entirely to
-- a different vessel, split into multiple sub-batches, combined with other batches, and so forth.
CREATE TABLE IF NOT EXISTS batch (
    id           serial PRIMARY KEY,
    uuid         uuid NOT NULL DEFAULT gen_random_uuid(),

    short_name   varchar(255) NOT NULL,

    created_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at   timestamptz
);

CREATE TABLE IF NOT EXISTS wort (
    id           serial PRIMARY KEY,
    uuid         uuid NOT NULL DEFAULT gen_random_uuid(),

    batch_id     int NOT NULL REFERENCES batch(id),
    volume_id    int NOT NULL REFERENCES volume(id),

    created_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at   timestamptz
);

CREATE TABLE IF NOT EXISTS beer (
    id           serial PRIMARY KEY,
    uuid         uuid NOT NULL DEFAULT gen_random_uuid(),

    wort_id      int NOT NULL REFERENCES wort(id),
    volume_id    int NOT NULL REFERENCES volume(id),

    created_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    updated_at   timestamptz NOT NULL DEFAULT timezone('utc', now()),
    deleted_at   timestamptz
);
