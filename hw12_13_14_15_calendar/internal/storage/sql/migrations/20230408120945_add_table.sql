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

--SELECT TOP(10) cust.id, SUM(car amount) as ammount 
--FROM public.customeer as cust
--   LEFT OUTER JOIN public.carts as car 
--   ON cust.id = car.customer_id
--GROUP BY cust.id
--ORDER BY SUM(car.amount) desc   

