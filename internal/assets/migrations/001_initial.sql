-- +migrate Up

CREATE TABLE my_table
(
    index BIGSERIAL PRIMARY KEY,
    user_id INTEGER,
    data TEXT
);

-- +migrate Down
DROP TABLE my_table;
