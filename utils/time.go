package utils

import (
	"time"
)

// GetStartDate : Lấy start date format timestamp
func GetStartDate() int64 {
	now := GetNow()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location()).UnixNano() / int64(time.Millisecond)
}

// GetEndDate : Lấy end date format timestamp
func GetEndDate() int64 {
	now := GetNow()
	year, month, day := now.Date()
	return time.Date(year, month, day+1, 0, 0, -1, 0, now.Location()).UnixNano() / int64(time.Millisecond)
}

// GetNow : Lấy thời điểm hiện tại theo giờ Việt nam
func GetNow() time.Time {
	return time.Now()
}

// GetStartDateVN : Lấy start date format timestamp
func GetStartDateVN() int64 {
	now := GetUTCNow()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location()).UnixNano() / int64(time.Millisecond)
}

// GetEndDateVN : Lấy end date format timestamp
func GetEndDateVN() int64 {
	now := GetUTCNow()
	year, month, day := now.Date()
	return time.Date(year, month, day+1, 0, 0, -1, 0, now.Location()).UnixNano() / int64(time.Millisecond)
}

// GetStartDateVNTime : Lấy start date format time.Time
func GetStartDateVNTime() time.Time {
	now := GetUTCNow()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location())
}

// GetEndDateVNTime : Lấy end date format time.Time
func GetEndDateVNTime() time.Time {
	now := GetUTCNow()
	year, month, day := now.Date()
	return time.Date(year, month, day+1, 0, 0, -1, 0, now.Location())
}

// GetVNTime24H : Lấy cach 24h format time.Time
func GetVNTime24H() time.Time {
	now := GetUTCNow()
	year, month, day := now.Date()
	minutes := now.Minute()
	hours := now.Hour()
	return time.Date(year, month, day-1, hours, minutes, -1, 0, now.Location())
}

// GetUTCNow : Lấy thời điểm hiện tại theo giờ UTC
func GetUTCNow() time.Time {
	vietnam, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return time.Now().In(vietnam)
}

// GetDaysDuration : Lấy chênh lệch thời gian giữa beforeTime, afterTime
func GetDaysDuration(beforeTime, afterTime time.Time) (days int) {
	days = int(afterTime.Sub(beforeTime).Hours() / 24)
	return
}

func GetStartEndOfDate(date time.Time) (startOfDate *time.Time, endOfDate *time.Time) {
	currentTimestamp := time.Now()
	currentLocation := currentTimestamp.Location()

	year, month, day := date.Date()

	startTime := time.Date(year, month, day, 0, 0, 0, 0, currentLocation)
	endTime := time.Date(year, month, day+1, 0, 0, -1, 0, currentLocation)
	startOfDate = &startTime
	endOfDate = &endTime

	return
}

// Send, Use this function to send data to elastic
func ConvertDateReturn(dateString string) (date time.Time, err error) {
	//convert date string to time.Time
	if len(dateString) != 0 {
		layout := "2006-01-02T15:04:05.000Z"
		dateTime, err := time.Parse(layout, dateString)
		if err != nil {
			//dateTime, err = time.Parse("2006-01-02T17:00:00Z", dateString )
			//if err != nil {
			dateTime, err = time.Parse("02/01/2006", dateString)
			if err != nil {
				dateTime, err = time.Parse("01/02/2006", dateString)
				if err != nil {
					dateTime, err = time.Parse("02-01-2006", dateString)
					if err != nil {
						dateTime, err = time.Parse("2006-01-02", dateString)
						if err != nil {
							return date, err
						}
					}
				}
			}
			//}
		}

		date = dateTime
	}
	return
}
