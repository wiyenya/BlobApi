-- +migrate Up

CREATE TABLE my_table
(
    index BIGSERIAL PRIMARY KEY,
    user_id INT, 
    data jsonb
);

-- +migrate Down
DROP TABLE my_table;
