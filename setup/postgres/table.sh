#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$GOIDDD_USERNAME" "$GOIDDD_DATABASE" <<-EOSQL
create table if not exists eventstore
(
	id serial not null,
	stream_id varchar(255) not null,
	stream_version integer default 0 not null,
	event_name varchar(255) not null,
	payload jsonb default '{}'::jsonb not null,
	occurred_at timestamp with time zone not null
);

create unique index if not exists id_unique
	on eventstore (id);

create unique index if not exists stream_unique
	on eventstore (stream_id, stream_version);

create index if not exists event_name_idx
	on eventstore (event_name);

create index if not exists occurred_at_idx
	on eventstore (occurred_at);
EOSQL