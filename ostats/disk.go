package ostats

import (
	"fmt"
	"os/exec"
	"strings"
)

type MountingPoint struct {
	FileSystem string `json:"file_system"`
	Size       string `json:"size"`
	Used       string `json:"used"`
	Avail      string `json:"avail"`
	PercentUse string `json:"percent_use"`
}

func FormatStatSlice(rawStatSlice []string) []string {
	var statSlice []string
	for _, stat := range rawStatSlice {
		if stat != "" {
			statSlice = append(statSlice, stat)
		}
	}
	return statSlice
}

func DiskUsageByPath(path string) (MountingPoint, error) {
	disk, err := diskUsage(path)
	if len(disk) > 0 {
		return disk[0], err
	}
	return MountingPoint{}, err
}

func DiskUsage() ([]MountingPoint, error) {
	return diskUsage("")
}

func diskUsage(path string) ([]MountingPoint, error) {
	var mountingPoints []MountingPoint
	cmd := "df -h"
	if path != "" {
		cmd = fmt.Sprintf("%s %s", cmd, path)
	}
	run := exec.Command("bash", "-c", cmd)
	stdout, err := run.Output()

	if err != nil {
		return nil, err
	}
	outputLines := strings.Split(string(stdout), "\n")
	for _, outputLine := range outputLines {
		mountingPointInfoSlice := FormatStatSlice(strings.Split(outputLine, " "))
		if len(mountingPointInfoSlice) > 0 {
			if strings.HasPrefix(mountingPointInfoSlice[0], "/dev") &&
				!strings.HasPrefix(mountingPointInfoSlice[0], "/dev/loop") {
				var mountingPoint MountingPoint
				mountingPoint.FileSystem = mountingPointInfoSlice[len(mountingPointInfoSlice)-1]
				mountingPoint.Size = mountingPointInfoSlice[1]
				mountingPoint.Used = mountingPointInfoSlice[2]
				mountingPoint.Avail = mountingPointInfoSlice[3]
				mountingPoint.PercentUse = mountingPointInfoSlice[4]
				mountingPoints = append(mountingPoints, mountingPoint)
			}
		}
	}
	return mountingPoints, nil
}
