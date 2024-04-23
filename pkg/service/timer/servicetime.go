package timer

import "time"

// ServiceTimer holds the start and end timestamps.
type ServiceTimer struct {
	UEinfor   int
	StartTime time.Time
	EndTime   time.Time
}

// NORA ServiceTimer holds the start and end timestamps.
type NORAServiceTimer struct {
	UEinfor   int
	StartTime time.Time
	EndTime   time.Time
}

// NewServiceTimer creates a new ServiceTimer instance.
func NewServiceTimer(UEid int, StartTime time.Time) *ServiceTimer {
	return &ServiceTimer{
		UEinfor:   UEid,
		StartTime: StartTime,
	}
}

var ServiceTimerMap map[int]*ServiceTimer

func init() {
	ServiceTimerMap = make(map[int]*ServiceTimer)
}

func StoreTimeStamp(ue *ServiceTimer) {
	ServiceTimerMap[ue.UEinfor] = ue
}

func GetStartTime(UEinfor int) time.Time {
	ue, _ := ServiceTimerMap[UEinfor]
	return ue.StartTime
}

func SetStartTime(UEinfor int, ST time.Time) bool {
	ue, _ := ServiceTimerMap[UEinfor]
	ue.StartTime = ST
	if GetStartTime(UEinfor) != ST {
		return false
	}
	return true
}

// CalculateServiceTime calculates the service time using the start and end timestamps.
func CalculateServiceTime(ST time.Time, ET time.Time) time.Duration {
	return ET.Sub(ST)
}
