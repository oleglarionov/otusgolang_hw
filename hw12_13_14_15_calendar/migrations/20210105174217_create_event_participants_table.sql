-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table event_participants
(
    event_id uuid    not null
        constraint event_participants_events_id_fk
            references events
            on delete cascade,
    uid      varchar not null,
    constraint event_participants_pk
        primary key (event_id, uid)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table event_participants