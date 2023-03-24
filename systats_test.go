package systats

import (
	"fmt"
	pprint "github.com/NubeIO/lib-system/print"
	"testing"
)

func TestSyStats_GetSystem(t *testing.T) {
	s := New()
	disks, err := s.GetSystem()
	if err != nil {
		return
	}
	pprint.PrintJOSN(disks)
}

func TestSyStats_GetMemoryUsage(t *testing.T) {
	s := New()
	disks, err := s.GetMemoryUsage()
	if err != nil {
		return
	}
	pprint.PrintJOSN(disks)
}

func TestSyStats_GetDisksPretty(t *testing.T) {
	s := New()
	disks, err := s.GetDisksPretty()
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.PrintJOSN(disks)
}

func TestSyStats_GetTopProcesses(t *testing.T) {
	s := New()
	disks, err := s.GetTopProcesses(2, "mem")
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.PrintJOSN(disks)
}
