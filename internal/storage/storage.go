package storage

import (
	"fmt"
	"net/http"
)

type gauge float64
type counter int64

type MemStorage struct {
	gaugeData   map[string]gauge
	counterData map[string]counter
}

func New() *MemStorage {
	return &MemStorage{
		gaugeData:   make(map[string]gauge),
		counterData: make(map[string]counter),
	}
}

func (m MemStorage) UpdateCounter(name string, value int64) {
	m.counterData[name] += counter(value)
}

func (m MemStorage) UpdateGauge(name string, value float64) {
	//fmt.Println(fmt.Sprintf("before update: %v", m.gaugeData))
	//fmt.Println(fmt.Sprintf("try update key: %s, value = %f", name, value))

	m.gaugeData[name] = gauge(value)

	//fmt.Println(fmt.Sprintf("after update: %v", m.gaugeData))
}

func (m MemStorage) GetValue(typeM string, name string) (string, int) {
	//fmt.Println(fmt.Sprintf("try get name: %s, type = %s", name, typeM))
	var v string
	statusCode := http.StatusOK
	switch typeM {
	case "counter":
		val, ok := m.counterData[name]
		if ok {
			v = fmt.Sprint(val)
		}
	case "gauge":
		//fmt.Println("gaugeData")
		//fmt.Println(fmt.Sprintf("before getting gauge: %v", m.gaugeData))

		val, ok := m.gaugeData[name]
		//fmt.Println(fmt.Sprintf("name=%s value=%f", name, val))
		if ok {
			v = fmt.Sprint(val)
		}
		//fmt.Println(fmt.Sprintf("after getting gauge: %v", m.gaugeData))
	default:
		statusCode = http.StatusNotFound
	}
	//fmt.Println(v, statusCode)
	return v, statusCode
}

func (m MemStorage) GetAllValues() string {
	var result string
	result += "Gauge metrics:\n"
	for n, v := range m.gaugeData {
		result += fmt.Sprintf("- %s = %f\n", n, v)
	}

	result += "Counter metrics:\n"
	for n, v := range m.counterData {
		result += fmt.Sprintf("- %s = %d\n", n, v)
	}

	return result
}
