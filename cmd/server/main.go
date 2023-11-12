package main

import (
	"github.com/lionslon/metrics-learning/internal/api"
	"github.com/lionslon/metrics-learning/internal/storage"
)

func main() {
	s := storage.New()
	a := api.New(s)
	if err := a.Start(); err != nil {
		panic(err)
	}
}
