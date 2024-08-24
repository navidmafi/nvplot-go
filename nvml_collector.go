package main

import (
	"errors"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

type NVMLCollector struct {
	device nvml.Device
}

type MemoryInfo struct {
	UsedMB  int
	TotalMB int
	FreeMB  int
}

func NewNVMLCollector() (*NVMLCollector, error) {
	if err := nvml.Init(); err != nvml.SUCCESS {
		return nil, err
	}

	count, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New("no NVIDIA GPU found")
	}

	device, err := nvml.DeviceGetHandleByIndex(0)
	if err != nvml.SUCCESS {
		return nil, err
	}

	return &NVMLCollector{device: device}, nil
}

func (nc *NVMLCollector) GetVRAMUsage() (MemoryInfo, error) {
	memory, err := nc.device.GetMemoryInfo()
	if err != nvml.SUCCESS {
		return MemoryInfo{}, err
	}

	byteToMB := func(n uint64) int {
		return int(n / (1 << 20))
	}

	return MemoryInfo{
		UsedMB:  byteToMB(memory.Used),
		TotalMB: byteToMB(memory.Total),
		FreeMB:  byteToMB(memory.Free),
	}, nil
}

func (nc *NVMLCollector) Close() {
	nvml.Shutdown()
}
