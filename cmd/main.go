package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/tiagoinaba/davinti/internal/server"
)

func main() {
	s := server.New()
	f, err := os.OpenFile("logs/deleted_contact.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	slog.SetDefault(slog.New(slog.NewTextHandler(f, nil)))
	s.RegisterRoutes()

	log.Fatal(s.Run())
}
