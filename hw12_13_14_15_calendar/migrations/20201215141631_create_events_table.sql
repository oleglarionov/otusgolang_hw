-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table events
(
    id          uuid not null
        constraint events_pk
            primary key,
    title       text not null,
    description text,
    begin_date  timestamp with time zone,
    end_date    timestamp with time zone
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table events;