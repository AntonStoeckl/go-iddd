BEGIN;

CREATE TABLE IF NOT EXISTS eventstore
(
    id serial not null,
    stream_id varchar(255) not null,
    stream_version integer default 0 not null,
    event_name varchar(255) not null,
    payload jsonb default '{}'::jsonb not null,
    occurred_at timestamp with time zone not null
);

CREATE UNIQUE INDEX IF NOT EXISTS id_unique
    on eventstore (id);

CREATE UNIQUE INDEX IF NOT EXISTS stream_unique
    on eventstore (stream_id, stream_version);

CREATE INDEX IF NOT EXISTS event_name_idx
    on eventstore (event_name);

CREATE INDEX IF NOT EXISTS occurred_at_idx
    on eventstore (occurred_at);

COMMIT;