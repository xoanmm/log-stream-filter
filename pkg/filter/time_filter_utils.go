package filter

import (
	"fmt"
	"log"
	"time"
)

const shortDateFormat = "01/02/2006 15:04:05"

//ShortDateFromString parse shot date from string
func ShortDateFromString(ds string) (time.Time, error) {
	t, err := time.Parse(shortDateFormat, ds)
	if err != nil {
		return t, err
	}
	return t, nil
}

//CheckDataBoundariesStr checks is startdate <= enddate
func CheckDataBoundariesStr(startdate, enddate string) (bool, error) {

	tstart, err := ShortDateFromString(startdate)
	if err != nil {
		return false, fmt.Errorf("cannot parse startdate: %v", err)
	}
	tend, err := ShortDateFromString(enddate)
	if err != nil {
		return false, fmt.Errorf("cannot parse enddate: %v", err)
	}

	if tstart.After(tend) {
		return false, fmt.Errorf("startdate > enddate - please set proper data boundaries")
	}
	return true, err
}

func getTimeStampUnixFromDate(startDate string) int64 {
	layout := "01/02/2006 15:04:05"
	t, err := time.Parse(layout, startDate)
	if err != nil {
		log.Fatalf("Error converting date to timestamp %s", err)
	}
	return t.UnixNano() / int64(time.Millisecond)
}

func getTimeInUTCFromMilliseconds(timestamp int64) string {
	return time.Unix(0, timestamp*int64(1000000)).UTC().Format(shortDateFormat)
}
