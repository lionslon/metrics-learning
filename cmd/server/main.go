package main

import "github.com/lionslon/metrics-learning/internal/api"

func main() {
	//s := storage.New()
	//a := api.New(s)
	//if err := a.Start(); err != nil {
	//	panic(err)
	//}
	api.Start()
}
