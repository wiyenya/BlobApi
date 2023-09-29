-- +migrate Up

CREATE TABLE my_table
(
    index BIGSERIAL PRIMARY KEY,
    data TEXT
);

-- +migrate Down
DROP TABLE my_table;
