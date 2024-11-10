-- +goose Up
-- +goose StatementBegin
CREATE TABLE contato (
	id INTEGER PRIMARY KEY,
	nome VARCHAR(100),
	idade NUMBER(3)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE contato;
-- +goose StatementEnd
