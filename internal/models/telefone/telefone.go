package telefone

import (
	"database/sql"
	"errors"
)

type Telefone struct {
	ID        int
	ContatoID int
	Descricao string
	Numero    string
}

func FindSome(DB *sql.DB, contatoID int) ([]Telefone, error) {
	rs, err := DB.Query(`SELECT id, id_contato, descricao, numero FROM telefone WHERE id_contato = ?`, contatoID)
	if err != nil {
		return nil, err
	}
	tels := []Telefone{}
	for rs.Next() {
		tel := Telefone{}
		err = rs.Scan(&tel.ID, &tel.ContatoID, &tel.Descricao, &tel.Numero)
		if err != nil {
			return nil, err
		}
		tels = append(tels, tel)
	}
	return tels, nil
}

func FindByID(DB *sql.DB, id int) (*Telefone, error) {
	rs, err := DB.Query(`SELECT id, id_contato, descricao, numero FROM telefone WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	tel := Telefone{}
	for rs.Next() {
		err = rs.Scan(&tel.ID, &tel.ContatoID, &tel.Descricao, &tel.Numero)
		if err != nil {
			return nil, err
		}
	}
	return &tel, nil
}

func Insert(DB *sql.DB, t *Telefone) (*Telefone, error) {
	if t == nil {
		return nil, errors.New("contact is nil")
	}
	r := DB.QueryRow(`INSERT INTO telefone (
			id_contato, descricao, numero
		) VALUES (
			?, ?, ?
		) RETURNING id`, t.ContatoID, t.Descricao, t.Numero)
	err := r.Scan(&t.ID)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func Update(DB *sql.DB, t *Telefone) (*Telefone, error) {
	if t == nil {
		return nil, errors.New("contact is nil")
	}
	r := DB.QueryRow(`UPDATE telefone
		SET descricao = ?,
			numero = ?
		WHERE id = ?
		RETURNING id, id_contato, descricao, numero`, t.Descricao, t.Numero, t.ID)
	err := r.Scan(&t.ID, &t.ContatoID, &t.Descricao, &t.Numero)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func Delete(DB *sql.DB, t *Telefone) error {
	_, err := DB.Exec(`DELETE FROM telefone WHERE id = ?`, t.ID)
	return err
}
