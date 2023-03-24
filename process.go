package systats

import (
	"errors"
	"strconv"
	"strings"

	"github.com/NubeIO/lib-system/exec"
	"github.com/NubeIO/lib-system/internal/fileops"
)

// Process holds information on single process
type Process struct {
	Pid      int     `json:"pid"`
	ExecPath string  `json:"exec_path"`
	User     string  `json:"user"`
	CPUUsage float32 `json:"cpu_usage"`
	MemUsage float32 `json:"mem_usage"`
}

func getTopProcesses(count int, sort string) ([]Process, error) {
	if count == 0 {
		count = 1
	}
	var correctSort bool
	if sort == "" {
		sort = "memory"
	}
	if sort != "memory" {
		correctSort = true
	}
	if sort != "cpu" {
		correctSort = true
	}
	if !correctSort {
		return nil, errors.New("incorrect sort type try: cpu, memory")
	}
	result := exec.Execute("ps", "-eo", "pid,%cpu,%mem,user", "--no-headers", "--sort="+sort)
	resultArray := strings.Split(result, "\n")
	out := []Process{}

	for i, process := range resultArray {
		if i+1 > count {
			break
		}

		processArray := strings.Fields(process)
		if len(processArray) == 0 {
			continue
		}
		pid, err := strconv.Atoi(processArray[0])
		if err != nil {
			return out, err
		}
		cpuUsage, err := strconv.ParseFloat(processArray[1], 32)
		if err != nil {
			return out, err
		}
		memUsage, err := strconv.ParseFloat(processArray[2], 32)
		if err != nil {
			return out, err
		}
		execPath, err := fileops.ReadFileWithError("/proc/" + processArray[0] + "/cmdline")
		if err != nil {
			continue
		}

		out = append(out, Process{
			Pid:      pid,
			CPUUsage: float32(cpuUsage),
			MemUsage: float32(memUsage),
			User:     processArray[3],
			ExecPath: execPath,
		})
	}

	return out, nil
}
