package api

import (
	"fmt"
	"github.com/lionslon/metrics-learning/internal/storage"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	store storageUpdater
}

type storageUpdater interface {
	UpdateCounter(string, int64)
	UpdateGauge(string, float64)
}

func New(s *storage.MemStorage) *Server {
	return &Server{
		store: s,
	}
}

func (s *Server) Start() error {
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	return http.ListenAndServe(`:8080`, s.mainPage())
}

func (s *Server) mainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			// разрешаем только POST-запросы
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		parsedPath := strings.Split(r.URL.Path, "/")
		if len(parsedPath) != 5 || parsedPath[1] != "update" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Println(r.URL.Path)
		metricType := parsedPath[2]
		metricName := parsedPath[3]
		metricValue := parsedPath[4]

		if metricType == "counter" {
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err == nil {
				s.store.UpdateCounter(metricName, value)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("К метрике %s добавлено значение %s", metricName, metricValue)))
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Ошибка преобразования к int64 значения: %s", metricValue)))
				return
			}
		} else if metricType == "gauge" {
			value, err := strconv.ParseFloat(metricValue, 64)
			if err == nil {
				s.store.UpdateGauge(metricName, value)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("Метрике %s установлено значение %s", metricName, metricValue)))
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Ошибка преобразования к float64 значения: %s", metricValue)))
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Нет такого типа метрики: %s", metricType)))
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		//fmt.Println(fmt.Printf("%v", s.store))
	}
}
