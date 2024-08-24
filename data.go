package main

import (
	"sync"
	"time"
)

type DataPoint struct {
	Timestamp time.Time
	MemoryInfo
}

type DataStorage struct {
	maxPoints int
	data      []DataPoint
	mutex     sync.RWMutex
}

func NewDataStorage(maxPoints int) *DataStorage {
	return &DataStorage{
		maxPoints: maxPoints,
		data:      make([]DataPoint, 0, maxPoints),
	}
}

func (ds *DataStorage) AddDataPoint(memInfo MemoryInfo) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	if len(ds.data) >= ds.maxPoints {
		ds.data = ds.data[1:] // remove oldest dp
	}

	ds.data = append(ds.data, DataPoint{
		Timestamp:  time.Now(),
		MemoryInfo: memInfo,
	})
}

func (ds *DataStorage) GetData() []DataPoint {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	return append([]DataPoint(nil), ds.data...)
}
