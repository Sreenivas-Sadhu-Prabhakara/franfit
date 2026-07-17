// Package money formats Indian Rupee amounts the way Indians read them:
// crores, lakhs and Indian digit grouping (₹12,34,567).
package money

import (
	"math"
	"strconv"
	"strings"
)

// FormatLakhs formats an amount expressed in ₹ lakhs.
// 120 -> "₹1.2 Cr", 36.5 -> "₹36.5 L", 0.125 -> "₹12,500".
func FormatLakhs(l float64) string {
	if l >= 100 {
		return "₹" + trimFloat(l/100) + " Cr"
	}
	if l >= 1 {
		return "₹" + trimFloat(l) + " L"
	}
	return FormatRupees(int64(math.Round(l * 100000)))
}

// FormatRupees formats a plain rupee amount with Indian grouping.
// 12500 -> "₹12,500", 1234567 -> "₹12,34,567".
func FormatRupees(n int64) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return sign + "₹" + s
	}
	head := s[:len(s)-3]
	tail := s[len(s)-3:]
	var groups []string
	for len(head) > 2 {
		groups = append([]string{head[len(head)-2:]}, groups...)
		head = head[:len(head)-2]
	}
	if head != "" {
		groups = append([]string{head}, groups...)
	}
	return sign + "₹" + strings.Join(groups, ",") + "," + tail
}

func trimFloat(f float64) string {
	s := strconv.FormatFloat(math.Round(f*10)/10, 'f', 1, 64)
	return strings.TrimSuffix(s, ".0")
}
