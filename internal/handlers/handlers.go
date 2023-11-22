package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lionslon/metrics-learning/internal/storage"
	"net/http"
	"strconv"
)

type handler struct {
	store storageUpdater
}

type storageUpdater interface {
	UpdateCounter(string, int64)
	UpdateGauge(string, float64)
	GetValue(string, string) (string, int)
	GetAllValues() string
}

func New(stor *storage.MemStorage) *handler {
	return &handler{
		store: stor,
	}
}

func (h *handler) PostMetric(c *fiber.Ctx) error {
	switch c.Params("typeMetric") {
	case "counter":
		value, err := strconv.ParseInt(c.Params("valueMetric"), 10, 64)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("Невалидное значение метрики: %s", c.Params("valueMetric")))
		}
		h.store.UpdateCounter(c.Params("nameMetric"), value)
		c.Set("Content-Type", "text/plain; charset=utf-8")
	case "gauge":
		value, err := strconv.ParseFloat(c.Params("valueMetric"), 64)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("Невалидное значение метрики: %s", c.Params("valueMetric")))
		}
		h.store.UpdateGauge(c.Params("nameMetric"), value)
		c.Set("Content-Type", "text/plain; charset=utf-8")
	default:
		return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("Невалидное тип метрики: %s", c.Params("typeMetric")))
	}
	c.Set("Content-Type", "text/plain; charset=utf-8")
	return c.Status(http.StatusOK).SendString("")
}

func (h *handler) GetMetric(c *fiber.Ctx) error {
	val, status := h.store.GetValue(c.Params("typeMetric"), c.Params("nameMetric"))
	return c.Status(status).SendString(val)
	//return c.Status(status).SendString("eee")
}

func (h *handler) GetAllMetrics(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	return c.Status(http.StatusOK).SendString(h.store.GetAllValues())
}
