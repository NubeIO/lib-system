package systats

import (
	"github.com/NubeIO/lib-system/exec"
	"github.com/NubeIO/lib-system/internal/fileops"
	"github.com/rvflash/elapsed"
	"path"
	"regexp"
	"strings"
	"time"
)

// System holds operating system information
type System struct {
	HostName      string    `json:"host_name"`
	OS            string    `json:"os"`
	Kernel        string    `json:"kernel"`
	UpTime        string    `json:"up_time"`
	LastBootDate  time.Time `json:"last_boot_date"`
	LoggedInUsers []User    `json:"logged_in_users"`
	Time          time.Time `json:"time"`
	Timezone      string    `json:"timezone"`
}

// User holds logged in user information
type User struct {
	Username     string    `json:"username"`
	RemoteHost   string    `json:"remote_host"`
	LoggedInTime time.Time `json:"logged_in_time"`
}

func getSystem(systats *SyStats) (System, error) {
	system := System{}

	systemOS, err := getOperatingSystem(systats)
	if err != nil {
		return system, err
	}
	system.OS = systemOS
	system.Kernel = getKernel(systats.VersionPath)
	system.HostName = strings.TrimSpace(fileops.ReadFile(systats.EtcPath + "hostname"))

	err = processSystemBootTimes(&system)
	if err != nil {
		return system, err
	}
	processLoggedInUsers(&system)
	system.Time = time.Now()
	return system, nil
}

func getOperatingSystem(systats *SyStats) (string, error) {
	p := path.Join(systats.EtcPath, "os-release")
	content, err := fileops.ReadFileWithError(p)
	if err != nil {
		p, err = fileops.FindFileWithNameLike(systats.EtcPath, "-release")
		if err != nil {
			return "", err
		}
		content = fileops.ReadFile(p)
	}

	split := strings.Split(content, "\n")
	var systemOS string
	for _, line := range split {
		r, _ := regexp.Compile("^(PRETTY_NAME=\")(.+)(\")")
		matches := r.FindAllStringSubmatch(line, -1)
		if len(matches) > 0 && len(matches[0]) >= 3 {
			systemOS = matches[0][2]
		}
	}

	return systemOS, nil
}

func processSystemBootTimes(system *System) error {
	tt, err := getBootTime()
	if err != nil {
		return err
	}
	system.UpTime = timeSince(tt)
	system.LastBootDate = tt
	tz, err := getTZ()
	if err != nil {
		return err
	}
	system.Timezone = tz
	return nil
}

func processLoggedInUsers(system *System) {
	split := strings.Split(exec.Execute("who"), "\n")
	system.LoggedInUsers = []User{}
	for _, line := range split {
		loggedInInfo := strings.Fields(line)
		if len(loggedInInfo) >= 5 {
			timeAdd := loggedInInfo[2] + " " + loggedInInfo[3] + ":00"
			loggedInTime, _ := parseTimeWithTimezone(timeLayout, timeAdd, system.Timezone)
			system.LoggedInUsers = append(system.LoggedInUsers, User{
				Username:     loggedInInfo[0],
				LoggedInTime: loggedInTime,
				RemoteHost:   loggedInInfo[4],
			})
		}
	}
}

func getKernel(versionPath string) string {
	var out string
	split := strings.Fields(fileops.ReadFile(versionPath))
	if len(split) >= 3 {
		out = strings.TrimSpace(split[2])
	}
	return out
}

const timeLayout = "2006-01-02 15:04:05"

func getBootTime() (time.Time, error) {
	up, err := exec.ExecuteWithError("uptime", "-s")
	if err != nil {
		return time.Time{}, err
	}
	tz, err := getTZ()
	if err != nil {
		return time.Time{}, err
	}
	timeWithZone, err := parseTimeWithTimezone(timeLayout, strings.TrimSpace(up), strings.TrimSpace(tz))
	if err != nil {
		return time.Time{}, err
	}
	return timeWithZone, err
}

func getTZ() (string, error) {
	tz, err := exec.ExecuteWithError("cat", "/etc/timezone")
	return strings.TrimSpace(tz), err
}

func parseTimeWithTimezone(layout, value, zone string) (time.Time, error) {
	tt, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	loc := tt.Location()
	loc, err = time.LoadLocation(zone)
	if err != nil {
		return time.Time{}, err // or `return tt, nil` if you more prefer
	}
	return time.ParseInLocation(layout, value, loc)
}

// TimeSince returns in a human-readable format the elapsed time
// e.g. 12 hours, 12 days
func timeSince(t time.Time) string {
	return elapsed.Time(t)
}
