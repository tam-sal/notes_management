package date

import "time"

func ArgentinaTimeNow() *time.Time {
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		return nil
	}
	currentTime := time.Now().In(loc)
	return &currentTime
}
