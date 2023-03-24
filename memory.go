package systats

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/NubeIO/lib-system/internal/fileops"
	"github.com/NubeIO/lib-system/internal/strops"
	"github.com/NubeIO/lib-system/internal/unitconv"
)

// Memory holds information on system memory usage
type Memory struct {
	PercentageUsed float64 `json:"percentage_used"`
	Available      uint64  `json:"available"`
	Free           uint64  `json:"free"`
	Used           uint64  `json:"used"`
	Time           int64   `json:"time"`
	Total          uint64  `json:"total"`
	Unit           string  `json:"unit"`
}

func getMemory(systats *SyStats, unit string) (Memory, error) {
	output := Memory{}
	output.Unit = unit

	meminfoStr, err := fileops.ReadFileWithError(systats.MemInfoPath)
	if err != nil {
		return output, err
	}

	meminfoSplit := strings.Split(meminfoStr, "\n")
	var buffers, cached uint64

	for _, line := range meminfoSplit {
		lineArr := strings.Fields(line)
		if len(lineArr) == 0 {
			continue
		}
		if lineArr[0] == "MemTotal:" {
			output.Total = strops.ToUint64(lineArr[1])
		}
		if lineArr[0] == "MemFree:" {
			output.Free = strops.ToUint64(lineArr[1])
		}
		if lineArr[0] == "MemAvailable:" {
			output.Available = strops.ToUint64(lineArr[1])
		}
		if lineArr[0] == "Buffers:" {
			buffers = strops.ToUint64(lineArr[1])
		}
		if lineArr[0] == "Cached:" {
			cached = strops.ToUint64(lineArr[1])
		}
	}

	if output.Total > 0 {
		output.Used = output.Total - (output.Free + buffers + cached)
		percentage := float64(output.Used) / float64(output.Total) * 100
		output.PercentageUsed = percentage
	}

	output.Time = time.Now().Unix()
	if unit == Byte {
	} else if unit == Kilobyte {
		output.Available = unitconv.KibToKB(output.Available)
		output.Total = unitconv.KibToKB(output.Total)
		output.Used = unitconv.KibToKB(output.Used)
		output.Free = unitconv.KibToKB(output.Free)
	} else if unit == Megabyte {
		output.Available = unitconv.KibToMB(output.Available)
		output.Total = unitconv.KibToMB(output.Total)
		output.Used = unitconv.KibToMB(output.Used)
		output.Free = unitconv.KibToMB(output.Free)
	} else {
		return output, errors.New(unit + " is not supported")
	}

	return output, nil
}

type MemoryUsage struct {
	MemoryPercentageUsed float64
	MemoryPercentage     string
	MemoryAvailable      string
	MemoryFree           string
	MemoryUsed           string
	MemoryTotal          string
	SwapPercentageUsed   float64
	SwapPercentage       string
	SwapFree             string
	SwapUsed             string
	SwapTotal            string
}

func getMemoryUsage(systats *SyStats) (*MemoryUsage, error) {
	m, err := systats.GetMemory(Kilobyte)
	if err != nil {
		return nil, err
	}
	s, err := systats.GetSwap(Kilobyte)
	if err != nil {
		return nil, err
	}

	return &MemoryUsage{
		MemoryPercentageUsed: m.PercentageUsed,
		MemoryPercentage:     fmt.Sprintf("%s", format(float32(m.PercentageUsed))) + "%",
		MemoryAvailable:      bytePretty(kbToByte(m.Available)),
		MemoryFree:           bytePretty(kbToByte(m.Free)),
		MemoryUsed:           bytePretty(kbToByte(m.Used)),
		MemoryTotal:          bytePretty(kbToByte(m.Total)),
		SwapPercentageUsed:   s.PercentageUsed,
		SwapPercentage:       fmt.Sprintf("%s", format(float32(s.PercentageUsed))) + "%",
		SwapFree:             bytePretty(kbToByte(s.Free)),
		SwapUsed:             bytePretty(kbToByte(s.Used)),
		SwapTotal:            bytePretty(kbToByte(s.Total)),
	}, nil
}

func kbToByte(input uint64) uint64 {
	return uint64(float64(input) * 1024)
}

func format(num float32) string {
	s := fmt.Sprintf("%.2f", num)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)

func bytePretty(size uint64) string {
	unit := ""
	value := float32(size)

	switch {
	case size >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case size >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case size >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case size >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case size >= BYTE:
		unit = "B"
	case size == 0:
		return "0"
	}
	stringValue := fmt.Sprintf("%.2f", value)
	stringValue = strings.TrimSuffix(stringValue, ".00")
	return fmt.Sprintf("%s%s", stringValue, unit)
}
