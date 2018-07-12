package time

import (
	gotime "time"
)

//GetUnix fetches the current unixtime
func GetUnix() int {
	return int(gotime.Now().Unix())
}
