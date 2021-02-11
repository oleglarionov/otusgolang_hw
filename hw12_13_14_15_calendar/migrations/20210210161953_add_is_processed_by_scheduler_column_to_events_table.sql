-- +goose Up
-- +goose StatementBegin
alter table events
    add is_processed_by_scheduler bool default false not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table events
    drop column is_processed_by_scheduler
-- +goose StatementEnd
