package time

import (
	gotime "time"
)

//GetUnix fetches the current unixtime
func GetUnix() int64 {
	return gotime.Now().Unix()
}
