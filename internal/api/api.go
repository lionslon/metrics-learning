package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lionslon/metrics-learning/internal/handlers"
	"github.com/lionslon/metrics-learning/internal/storage"
	"log"
)

type Server struct {
	storage *storage.MemStorage
}

//func New() *Server {
//	return &Server{
//		store: s,
//	}
//}

func Start() {
	app := fiber.New()
	st := storage.New()
	handl := handlers.New(st)
	app.Post("/update/:typeMetric/:nameMetric/:valueMetric", handl.PostMetric)
	app.Get("/value/:typeMetric/:nameMetric", handl.GetMetric)
	app.Get("/", handl.GetAllMetrics)
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	log.Fatal(app.Listen(":8080"))
}
