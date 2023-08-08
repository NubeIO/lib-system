package ostats

import (
	"fmt"
	pprint "github.com/NubeIO/lib-system/print"
	"testing"
)

func TestDiskUsageRootInGB(t *testing.T) {
	usage, err := DiskUsageByPath("/")
	if err != nil {
		return
	}
	pprint.PrintJOSN(usage)
}

func TestDiskGetMemory(t *testing.T) {
	data, err := getMem()
	if err != nil {
		fmt.Println(err)
		return
	}
	pprint.PrintJOSN(data)
}
