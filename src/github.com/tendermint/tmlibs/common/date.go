package common

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

var TimeLayout = "2006-01-02"

func ParseDateRange(dateRange string) (startDate, endDate time.Time, err error) {
	dates := strings.Split(dateRange, ":")
	if len(dates) != 2 {
		err = errors.New("bad date range, must be in format date:date")
		return
	}
	parseDate := func(date string) (out time.Time, err error) {
		if len(date) == 0 {
			return
		}
		out, err = time.Parse(TimeLayout, date)
		return
	}
	startDate, err = parseDate(dates[0])
	if err != nil {
		return
	}
	endDate, err = parseDate(dates[1])
	if err != nil {
		return
	}
	return
}
