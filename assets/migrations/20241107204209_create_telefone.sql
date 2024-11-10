-- +goose Up
-- +goose StatementBegin
CREATE TABLE telefone (
	id INTEGER PRIMARY KEY,
	id_contato NUMBER(14),
	descricao varchar(50) NOT NULL,
	numero varchar(16) NOT NULL,
	FOREIGN KEY (id_contato) REFERENCES contato(id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE telefone;
-- +goose StatementEnd
