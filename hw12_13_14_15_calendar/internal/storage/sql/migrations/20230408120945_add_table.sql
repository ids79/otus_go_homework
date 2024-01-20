-- +goose Up
-- +goose StatementBegin
CREATE table events (
    id              text,
    title           text,
    description     text,
	date_time       timestamp,
	year            numeric,
	day             numeric,
	week            numeric,
	month           numeric,
	duration        numeric,
	user_id         numeric,
	time_before     numeric);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events;
-- +goose StatementEnd

