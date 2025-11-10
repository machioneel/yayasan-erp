package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateCode generates a unique code with prefix
func GenerateCode(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	timestamp := time.Now().Format("20060102")
	random := rand.Intn(9999)
	return fmt.Sprintf("%s-%s-%04d", prefix, timestamp, random)
}

// FormatCurrency formats number as Indonesian Rupiah
func FormatCurrency(amount float64) string {
	return fmt.Sprintf("Rp %.2f", amount)
}

// CalculateAge calculates age from birth date
func CalculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

// Pagination helper
type PaginationParams struct {
	Page     int
	PageSize int
	Total    int64
}

func (p *PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationParams) TotalPages() int {
	if p.PageSize == 0 {
		return 0
	}
	return int((p.Total + int64(p.PageSize) - 1) / int64(p.PageSize))
}

// DIPERBAIKI: Fungsi ParseDate ditambahkan
// ParseDate mengurai string tanggal (YYYY-MM-DD atau RFC3339)
func ParseDate(dateString string) (time.Time, error) {
	// Coba format YYYY-MM-DD
	t, err := time.Parse("2006-01-02", dateString)
	if err == nil {
		return t, nil
	}

	// Coba format timestamp penuh RFC3339
	t, err = time.Parse(time.RFC3339, dateString)
	if err == nil {
		return t, nil
	}

	// Jika keduanya gagal, kembalikan eror
	return time.Time{}, fmt.Errorf("format tanggal tidak valid untuk '%s', gunakan YYYY-MM-DD atau RFC3339", dateString)
}