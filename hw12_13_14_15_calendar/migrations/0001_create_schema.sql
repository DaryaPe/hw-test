-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE SCHEMA calendar;
CREATE TABLE calendar.event
(
id               varchar      not null,
user_id          bigint       not null,
title            varchar(100) not null,
start_date       timestamp    not null,
end_date         timestamp    not null,
notification     integer,
description      bytea,
PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE calendar.event;
DROP SCHEMA calendar CASCADE ;
-- +goose StatementEnd