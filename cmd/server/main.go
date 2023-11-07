package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type gauge float64
type counter int64

type MemStorage struct {
	gaugeData   map[string]gauge
	counterData map[string]counter
}

type storageUpdater interface {
	updateCounter(string, int64)
	updateGauge(string, float64)
}

var storage = MemStorage{
	gaugeData:   make(map[string]gauge),
	counterData: make(map[string]counter),
}

func (m MemStorage) updateCounter(name string, value int64) {
	m.counterData[name] += counter(value)
}

func (m MemStorage) updateGauge(name string, value float64) {
	m.gaugeData[name] = gauge(value)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	parsedPath := strings.Split(r.URL.Path, "/")
	if len(parsedPath) != 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println(r.URL.Path)
	metricType := parsedPath[2]
	metricName := parsedPath[3]
	metricValue := parsedPath[4]

	var stor storageUpdater = storage

	if metricType == "counter" {
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err == nil {
			stor.updateCounter(metricName, value)
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
			stor.updateGauge(metricName, value)
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

	fmt.Println(storage)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, mainPage)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
