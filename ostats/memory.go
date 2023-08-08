package ostats

import (
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func GetMemory() (*Memory, error) {
	return getMem()
}

const (
	filePath = "/proc/meminfo"
)

type Memory struct {
	Buffers    uint64 `json:"buffers"`
	Cached     uint64 `json:"cached"`
	PageTables uint64 `json:"page_tables"`

	Total      uint64  `json:"total"`
	Free       uint64  `json:"free"`
	Used       uint64  `json:"used"`
	TotalInMB  float64 `json:"total_in_mb"`
	FreeInMB   float64 `json:"free_in_mb"`
	UsedInMB   float64 `json:"used_in_mb"`
	Percentage float64 `json:"percentage"`

	SwapTotal      uint64  `json:"swap_total"`
	SwapFree       uint64  `json:"swap_free"`
	SwapUsed       uint64  `json:"swap_used"`
	SwapTotalInMB  float64 `json:"swap_total_in_mb"`
	SwapFreeInMB   float64 `json:"swap_free_in_mb"`
	SwapUsedInMB   float64 `json:"swap_used_in_mb"`
	SwapPercentage float64 `json:"swap_percentage"`
}

func getMem() (*Memory, error) {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	out := &Memory{}

	s := strings.Replace(string(f), ":", "", -1)

	for _, v := range strings.Split(s, "\n") {
		x := strings.Fields(v)
		if len(x) < 2 {
			continue
		}
		switch x[0] {
		case "MemTotal":
			out.Total, err = strconv.ParseUint(x[1], 10, 64)
			if err != nil {
				return nil, err
			}
		case "MemFree":
			out.Free, err = strconv.ParseUint(x[1], 10, 64)
			if err != nil {
				return nil, err
			}
		case "Buffers":
			out.Buffers, err = strconv.ParseUint(x[1], 10, 64)
			if err != nil {
				return nil, err
			}
		case "Cached":
			out.Cached, err = strconv.ParseUint(x[1], 10, 64)
			if err != nil {
				return nil, err
			}
		case "PageTables":
			out.PageTables, err = strconv.ParseUint(x[1], 10, 64)
			if err != nil {
				return nil, err
			}
		case "SwapTotal":
			out.SwapTotal, err = strconv.ParseUint(x[1], 10, 64)
			if err != nil {
				return nil, err
			}
		case "SwapFree":
			out.SwapFree, err = strconv.ParseUint(x[1], 10, 64)
			if err != nil {
				return nil, err
			}

		}
	}

	if out.SwapTotal > 0 {
		out.SwapUsed = out.SwapTotal - out.SwapFree
		out.SwapPercentage = toFixed(100*float64(out.SwapUsed)/float64(out.SwapTotal), 2)

	}

	if out != nil {
		out.TotalInMB = bToMb(out.Total)
		out.FreeInMB = bToMb(out.Free)
		out.Used = calcUsed(out)
		out.UsedInMB = calcUsedMb(out)
		out.Percentage = calcPercentage(out)

		out.SwapTotalInMB = bToMb(out.SwapTotal)
		out.SwapFreeInMB = bToMb(out.SwapFree)
		out.SwapUsedInMB = bToMb(out.SwapUsed)
	}
	return out, err
}

func calcUsed(d *Memory) uint64 {
	used := d.Total - d.Free - d.Buffers - d.Cached + d.PageTables
	return used
}

func calcUsedMb(d *Memory) float64 {
	used := calcUsed(d)
	return float64(used) / 1024
}

func calcPercentage(d *Memory) float64 {
	used := calcUsed(d)
	return toFixed(100*float64(used)/float64(d.Total), 2)
}

func bToMb(b uint64) float64 {
	return float64(b) / (1 << 10)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
