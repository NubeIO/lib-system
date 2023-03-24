package systats

import (
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
