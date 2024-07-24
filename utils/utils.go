package utils

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

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

func GetDataFingerprint(c *gin.Context) map[string]string {
	var ip string

	// Mendapatkan IP dari header "X-Forwarded-For"
	XForwardedFor := c.GetHeader("X-Forwarded-For")
	if XForwardedFor != "" {
		ip = XForwardedFor
	}

	// Jika "X-Forwarded-For" tidak tersedia, gunakan "X-Real-IP"
	XRealIP := c.GetHeader("X-Real-IP")
	if ip == "" && XRealIP != "" {
		ip = XRealIP
	}

	// Jika header tidak tersedia, gunakan metode bawaan ClientIP [Carrier-Grade NAT (CGNAT)]
	ClientIP := c.ClientIP()
	if ip == "" && ClientIP != "" {
		ip = ClientIP
	}

	userAgent := c.GetHeader("User-Agent")
	acceptLanguage := c.GetHeader("Accept-Language")
	referer := c.GetHeader("Referer")

	fingerprintData := map[string]string{
		"UserAgent":      userAgent,
		"AcceptLanguage": acceptLanguage,
		"Referer":        referer,
		"IP":             ip,
	}

	return fingerprintData
}

func GenerateFingerprint(data map[string]string) string {
	fingerprintSource := ""
	for key, value := range data {
		fingerprintSource += key + ":" + value + ";"
	}
	hash := md5.Sum([]byte(fingerprintSource))
	return hex.EncodeToString(hash[:])
}

// =============================================================================================
