package utils

import "time"

// =============================================================================================

// ternary adalah fungsi utilitas untuk melakukan operasi ternary dalam Go
func Ternary(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// =============================================================================================

const StrFormatTime string = "02 Jan 2006 15:04:05 MST"

func INDTime(inputTime time.Time) time.Time {
	// Menetapkan zona waktu ke Indonesia/Jakarta (Waktu Indonesia Barat)
	indonesiaLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Time{}
	}

	// Mengubah zona waktu ke Indonesia/Jakarta
	indonesiaTime := inputTime.In(indonesiaLocation)

	return indonesiaTime
}

func ConvertStringToTime(timeString string) (time.Time, error) {
	layout := "2006-01-02 15:04:05.999999-07" // Format sesuai dengan string yang diberikan
	parsedTime, err := time.Parse(layout, timeString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

// =============================================================================================
