package server

import (
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/tiagoinaba/davinti/internal/models/contato"
	"github.com/tiagoinaba/davinti/internal/models/telefone"
)

func (s *Server) RegisterRoutes() {
	s.Router.Static("/public", "./assets/public")

	funcMap := template.FuncMap{
		"incr": func(a int, b int) int {
			return a + b
		},
		"decr": func(a int, b int) int {
			return a - b
		},
	}

	s.Router.SetFuncMap(funcMap)

	s.Router.LoadHTMLGlob("assets/templates/*")

	s.Router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	s.Router.GET("/contacts", func(c *gin.Context) {
		pesq := c.Query("pesquisa")
		pagina := c.Query("pagina")
		var pag int = 0
		if pagina != "" {
			var err error
			pag, err = strconv.Atoi(pagina)
			if err != nil {
				pag = 0
			}
		}
		var count float64
		r := s.DB.QueryRow(`SELECT COUNT(*) FROM contato WHERE nome LIKE CONCAT('%', ?, '%')`, pesq)
		err := r.Scan(&count)

		paginas := int(math.Ceil(count/10)) - 1

		if err != nil {
			c.Data(500, "text/html; charset=utf-8", []byte("Falha na conexão ao banco!"))
			return
		}

		cs, err := contato.FindSome(s.DB, pesq, 10, pag*10)
		if err != nil {
			c.Data(500, "text/html; charset=utf-8", []byte("Falha na conexão ao banco!"))
			return
		}

		c.HTML(http.StatusOK, "contacts.html", gin.H{
			"contacts": cs,
			"pesquisa": pesq,
			"pagina":   pag,
			"paginas":  paginas,
		})
	})

	s.Router.GET("/contacts/new", func(c *gin.Context) {
		c.HTML(200, "contact-form.html", gin.H{
			"title": "novo contato",
		})
	})

	s.Router.GET("/contacts/:id", func(c *gin.Context) {
		id := c.Param("id")
		sc := c.Query("success")
		ctt, err := contato.FindByID(s.DB, id)
		if err != nil {
			c.HTML(http.StatusNotFound, "not-found.html", gin.H{})
			return
		}

		tels, err := telefone.FindSome(s.DB, ctt.ID)
		if err != nil {
			c.HTML(http.StatusNotFound, "not-found.html", gin.H{})
			return
		}

		c.HTML(200, "contact-form.html", gin.H{
			"title":   "editar contato",
			"contato": ctt,
			"tels":    tels,
			"success": sc == "1",
		})
	})

	s.Router.POST("/contacts/:id", func(c *gin.Context) {
		ctt := &contato.Contato{}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.HTML(http.StatusBadRequest, "contact-form.html", gin.H{
				"title":   "editar contato",
				"contato": ctt,
				"error":   "contato não encontrado!",
			})
			return
		}
		ctt.ID = id
		ctt.Nome = c.PostForm("nome")
		idade, err := strconv.Atoi(c.PostForm("idade"))
		if err != nil {
			c.HTML(http.StatusBadRequest, "contact-form.html", gin.H{
				"title":   "editar contato",
				"contato": ctt,
				"error":   "a idade deve ser um número inteiro!",
			})
			return
		}
		ctt.Idade = idade
		if len(ctt.Nome) == 0 {
			c.HTML(http.StatusBadRequest, "contact-form.html", gin.H{
				"title":   "editar contato",
				"contato": ctt,
				"error":   "a idade deve ser um número inteiro!",
			})
			return
		}
		if len(ctt.Nome) > 100 {
			c.HTML(http.StatusBadRequest, "contact-form.html", gin.H{
				"title":   "editar contato",
				"contato": ctt,
				"error":   "o nome deve ter menos de 100 caracteres!",
			})
			return
		}
		ctt, err = contato.Update(s.DB, ctt)
		if err != nil {
			c.HTML(http.StatusBadRequest, "contact-form.html", gin.H{
				"title":   "editar contato",
				"contato": ctt,
				"error":   "ops, algo deu errado!",
			})
			return
		}

		tels, err := telefone.FindSome(s.DB, ctt.ID)
		if err != nil {
			c.HTML(http.StatusNotFound, "not-found.html", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "contact-form.html", gin.H{
			"title":   "editar contato",
			"contato": ctt,
			"tels":    tels,
			"success": true,
		})
	})

	s.Router.POST("/contacts/new", func(c *gin.Context) {
		ctt := &contato.Contato{}
		ctt.Nome = c.PostForm("nome")
		idade, err := strconv.Atoi(c.PostForm("idade"))
		if err != nil {
			c.HTML(http.StatusBadRequest, "contact-form.html", gin.H{
				"title":   "novo contato",
				"contato": ctt,
				"error":   "a idade deve ser um número inteiro!",
			})
			return
		}
		ctt.Idade = idade
		if len(ctt.Nome) == 0 {
			c.HTML(http.StatusBadRequest, "contact-form.html", gin.H{
				"title":   "novo contato",
				"contato": ctt,
				"error":   "a idade deve ser um número inteiro!",
			})
			return
		}
		if len(ctt.Nome) > 100 {
			c.HTML(http.StatusBadRequest, "contact-form.html", gin.H{
				"title":   "novo contato",
				"contato": ctt,
				"error":   "o nome deve ter menos de 100 caracteres!",
			})
			return
		}
		ctt, err = contato.Insert(s.DB, ctt)
		if err != nil {
			c.HTML(http.StatusBadRequest, "contact-form.html", gin.H{
				"title":   "editar contato",
				"contato": ctt,
				"error":   "ops, algo deu errado!",
			})
			return
		}
		c.Header("HX-Redirect", fmt.Sprintf("/contacts/%d?success=1", ctt.ID))
	})

	s.Router.DELETE("/contacts/:id", func(c *gin.Context) {
		ctt := &contato.Contato{}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.HTML(http.StatusBadRequest, "contact-form.html", gin.H{
				"title":   "editar contato",
				"contato": ctt,
				"error":   "contato não encontrado!",
			})
			return
		}
		ctt.ID = id
		err = contato.Delete(s.DB, ctt)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		slog.Info("contato excluído", "id", ctt.ID, "nome", ctt.Nome)

		c.Header("HX-Redirect", "/")
	})

	s.Router.DELETE("/delete", func(c *gin.Context) {
		c.Status(200)
	})

	s.Router.GET("/contact/:id/phone-number/new", func(c *gin.Context) {
		contactID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Data(400, "text/html", []byte("contact not found!"))
			return
		}
		c.HTML(http.StatusOK, "tel-form.html", gin.H{
			"telefone": &telefone.Telefone{
				ContatoID: contactID,
			},
		})
	})

	s.Router.POST("/contact/:id/phone-number", func(c *gin.Context) {
		contactID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Data(400, "text/html", []byte("contact not found!"))
			return
		}
		t := &telefone.Telefone{}
		t.ContatoID = contactID
		t.Descricao = c.PostForm("descricao")
		t.Numero = c.PostForm("numero")
		telefone.Insert(s.DB, t)

		c.HTML(http.StatusOK, "tel-component.html", gin.H{
			"tels": []telefone.Telefone{*t},
		})
	})

	s.Router.GET("/phone-number/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Data(400, "text/html", []byte("contact not found!"))
			return
		}
		t, err := telefone.FindByID(s.DB, id)
		if err != nil {
			c.Data(404, "text/html", []byte("não encontrado!"))
			return
		}

		c.HTML(http.StatusOK, "tel-form.html", gin.H{
			"telefone": t,
		})
	})

	s.Router.POST("/phone-number/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Data(400, "text/html", []byte("contact not found!"))
			return
		}
		c.PostForm("descricao")
		t := &telefone.Telefone{}
		t.ID = id
		t.Descricao = c.PostForm("descricao")
		t.Numero = c.PostForm("numero")

		t, err = telefone.Update(s.DB, t)
		if err != nil {
			c.Data(500, "text/html", []byte("não foi possível atualizar o telefone!"))
			return
		}

		c.HTML(http.StatusOK, "tel-component.html", gin.H{
			"tels": []telefone.Telefone{*t},
		})
	})

	s.Router.DELETE("/phone-number/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Data(400, "text/html", []byte("contact not found!"))
			return
		}
		err = telefone.Delete(s.DB, &telefone.Telefone{ID: id})
		if err != nil {
			c.Status(500)
			return
		}
		c.Status(http.StatusOK)
	})
}
