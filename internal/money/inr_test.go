package money

import "testing"

func TestFormatLakhs(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{120, "₹1.2 Cr"},
		{100, "₹1 Cr"},
		{250, "₹2.5 Cr"},
		{36.5, "₹36.5 L"},
		{1, "₹1 L"},
		{99.9, "₹99.9 L"},
		{0.125, "₹12,500"},
		{0.5, "₹50,000"},
	}
	for _, c := range cases {
		if got := FormatLakhs(c.in); got != c.want {
			t.Errorf("FormatLakhs(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatRupeesIndianGrouping(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{500, "₹500"},
		{12500, "₹12,500"},
		{100000, "₹1,00,000"},
		{1234567, "₹12,34,567"},
		{123456789, "₹12,34,56,789"},
		{-12500, "-₹12,500"},
	}
	for _, c := range cases {
		if got := FormatRupees(c.in); got != c.want {
			t.Errorf("FormatRupees(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}
