package contato

import (
	"database/sql"
	"errors"
)

type Contato struct {
	ID    int
	Nome  string
	Idade int
	Count int
}

func All(DB *sql.DB) ([]Contato, error) {
	stmt := `SELECT id, nome, idade FROM contato ORDER BY nome`
	rs, err := DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	contatos := []Contato{}
	for rs.Next() {
		c := Contato{}
		if err = rs.Scan(&c.ID, &c.Nome, &c.Idade); err != nil {
			return nil, err
		}
		contatos = append(contatos, c)
	}
	if err := rs.Err(); err != nil {
		return nil, err
	}
	return contatos, nil
}

func Some(DB *sql.DB, limit int, offset int) ([]Contato, error) {
	stmt := `SELECT id, nome, idade FROM contato ORDER BY nome LIMIT ? OFFSET ?`
	rs, err := DB.Query(stmt, limit, offset)
	if err != nil {
		return nil, err
	}
	contatos := []Contato{}
	for rs.Next() {
		c := Contato{}
		if err = rs.Scan(&c.ID, &c.Nome, &c.Idade); err != nil {
			return nil, err
		}
		contatos = append(contatos, c)
	}
	if err := rs.Err(); err != nil {
		return nil, err
	}
	return contatos, nil
}

func FindAll(DB *sql.DB, search string) ([]Contato, error) {
	stmt := `SELECT id, nome, idade
			FROM contato
			WHERE nome LIKE CONCAT('%', ?, '%')
			ORDER BY nome`
	rs, err := DB.Query(stmt, search)
	if err != nil {
		return nil, err
	}
	contatos := []Contato{}
	for rs.Next() {
		c := Contato{}
		if err = rs.Scan(&c.ID, &c.Nome, &c.Idade); err != nil {
			return nil, err
		}
		contatos = append(contatos, c)
	}
	if err := rs.Err(); err != nil {
		return nil, err
	}
	return contatos, nil
}

func FindSome(DB *sql.DB, search string, limit int, offset int) ([]Contato, error) {
	stmt := `SELECT id, nome, idade, (SELECT count(*) FROM telefone WHERE id_contato = contato.id) tels
			FROM contato
			WHERE nome LIKE CONCAT('%', ?, '%')
			ORDER BY nome 
			LIMIT ? OFFSET ?`
	rs, err := DB.Query(stmt, search, limit, offset)
	if err != nil {
		return nil, err
	}
	contatos := []Contato{}
	for rs.Next() {
		c := Contato{}
		if err = rs.Scan(&c.ID, &c.Nome, &c.Idade, &c.Count); err != nil {
			return nil, err
		}
		contatos = append(contatos, c)
	}
	if err := rs.Err(); err != nil {
		return nil, err
	}
	return contatos, nil
}

func FindByID(DB *sql.DB, id string) (*Contato, error) {
	r := DB.QueryRow("SELECT id, nome, idade FROM contato WHERE id = ?", id)
	c := &Contato{}
	err := r.Scan(&c.ID, &c.Nome, &c.Idade)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func Insert(DB *sql.DB, c *Contato) (*Contato, error) {
	if c == nil {
		return nil, errors.New("contact is nil")
	}
	r := DB.QueryRow(`INSERT INTO contato (
		nome,
		idade
	) VALUES (?, ?)
	RETURNING id, nome, idade`, c.Nome, c.Idade)
	err := r.Scan(&c.ID, &c.Nome, &c.Idade)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func Update(DB *sql.DB, c *Contato) (*Contato, error) {
	if c == nil {
		return nil, errors.New("contact is nil")
	}
	r := DB.QueryRow(`UPDATE contato
		SET nome = ?,
		idade = ?
		WHERE id = ?
		RETURNING id, nome, idade
		`, c.Nome, c.Idade, c.ID)
	err := r.Scan(&c.ID, &c.Nome, &c.Idade)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func Delete(DB *sql.DB, c *Contato) error {
	r := DB.QueryRow(`DELETE FROM contato WHERE id = ? RETURNING id, nome`, c.ID)
	err := r.Scan(&c.ID, &c.Nome)
	if err != nil {
		return err
	}
	return nil
}
