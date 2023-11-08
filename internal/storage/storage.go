package storage

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
	m.gaugeData[name] = gauge(value)
}
