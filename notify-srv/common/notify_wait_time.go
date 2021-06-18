package common

const (
	NotifyWaitTimesOne      = 1
	NotifyWaitTimesTwo      = 2
	NotifyWaitTimesThree    = 3
	NotifyWaitTimesFour     = 4
	NotifyWaitTimesFive     = 5
	NotifyWaitTimesSex      = 6
	NotifyWaitTimesSeven    = 7
	NotifyWaitTimesEight    = 8
	NotifyWaitTimesNine     = 9
	NotifyWaitTimesTen      = 10
	NotifyWaitTimesEleven   = 11
	NotifyWaitTimesTwelve   = 12
	NotifyWaitTimesThirteen = 13
	NotifyWaitTimesFourteen = 14
	NotifyWaitTimesFifteen  = 15
)

var NotifyWaitTimes = make(map[int]int64)

func init() {
	NotifyWaitTimes[NotifyWaitTimesOne] = 15         // 15s
	NotifyWaitTimes[NotifyWaitTimesTwo] = 15         // 15s
	NotifyWaitTimes[NotifyWaitTimesThree] = 30       // 30s
	NotifyWaitTimes[NotifyWaitTimesFour] = 60        // 60s
	NotifyWaitTimes[NotifyWaitTimesFive] = 180       // 3m
	NotifyWaitTimes[NotifyWaitTimesSex] = 600        // 10m
	NotifyWaitTimes[NotifyWaitTimesSeven] = 600      // 10m
	NotifyWaitTimes[NotifyWaitTimesEight] = 1800     // 30m
	NotifyWaitTimes[NotifyWaitTimesNine] = 3600      // 1h
	NotifyWaitTimes[NotifyWaitTimesTen] = 7200       // 2h
	NotifyWaitTimes[NotifyWaitTimesEleven] = 10800   // 3h
	NotifyWaitTimes[NotifyWaitTimesTwelve] = 10800   // 3h
	NotifyWaitTimes[NotifyWaitTimesThirteen] = 10800 // 3h
	NotifyWaitTimes[NotifyWaitTimesFourteen] = 21600 // 6h
	NotifyWaitTimes[NotifyWaitTimesFifteen] = 21600  // 6h
}

func GetNotifyWaitTimeById(id int) int64 {
	if v, ok := NotifyWaitTimes[id]; ok {
		return v
	}

	return -1
}
