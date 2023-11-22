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

func (h *handler) PostCounter(c *fiber.Ctx) error {
	value, err := strconv.ParseInt(c.Params("valueMetric"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("Невалидное значение метрики: %s", c.Params("valueMetric")))
	}
	h.store.UpdateCounter(c.Params("nameMetric"), value)
	c.Set("Content-Type", "text/plain; charset=utf-8")
	return c.Status(http.StatusOK).SendString("")
}

func (h *handler) PostGauge(c *fiber.Ctx) error {
	value, err := strconv.ParseFloat(c.Params("valueMetric"), 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("Невалидное значение метрики: %s", c.Params("valueMetric")))
	}
	h.store.UpdateGauge(c.Params("nameMetric"), value)
	c.Set("Content-Type", "text/plain; charset=utf-8")
	return c.Status(http.StatusOK).SendString("")
}

func (h *handler) GetMetric(c *fiber.Ctx) error {
	_, status := h.store.GetValue(c.Params("typeMetric"), c.Params("nameMetric"))
	//return c.Status(status).SendString(val)
	return c.Status(status).SendString("eee")
}

func (h *handler) GetAllMetrics(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	return c.Status(http.StatusOK).SendString(h.store.GetAllValues())
}
