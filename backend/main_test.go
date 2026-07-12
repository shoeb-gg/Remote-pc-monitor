package main

import (
	"strings"
	"testing"
)

func TestParseFloat(t *testing.T) {
	cases := []struct {
		name   string
		value  string
		unit   string
		want   float64
		wantOK bool
	}{
		{"valid celsius", "72.3 °C", "°C", 72.3, true},
		{"valid watts", "65.7 W", "W", 65.7, true},
		{"na string", "N/A", "°C", 0, false},
		{"dash", "-", "°C", 0, false},
		{"empty", "", "°C", 0, false},
	}
	for _, c := range cases {
		got, ok := parseFloat(c.value, c.unit)
		if ok != c.wantOK || (ok && got != c.want) {
			t.Errorf("%s: parseFloat(%q, %q) = (%v, %v), want (%v, %v)",
				c.name, c.value, c.unit, got, ok, c.want, c.wantOK)
		}
	}
}

// Reproduces the "silent 0" bug: a missing sensor must report not-found,
// never a 0 that is indistinguishable from a real reading.
func TestExtractMetricByPath_MissingSensorIsNotFound(t *testing.T) {
	computer := HardwareNode{
		Text: "PC",
		Children: []HardwareNode{
			{Text: "AMD Ryzen", Children: []HardwareNode{
				{Text: "Temperatures", Children: []HardwareNode{
					{Text: "Core (Tctl)", Value: "72.3 °C"},
				}},
			}},
		},
	}

	if got, ok := extractMetricByPath(computer, []string{"AMD Ryzen", "Temperatures", "Core|Tctl"}, "°C"); !ok || got != 72.3 {
		t.Errorf("present metric: got (%v, %v), want (72.3, true)", got, ok)
	}

	if _, ok := extractMetricByPath(computer, []string{"AMD Ryzen", "Temperatures", "CCD1 (Tdie)"}, "°C"); ok {
		t.Error("missing sensor should return found=false, not a silent 0")
	}
}

func TestFindNodeByPattern(t *testing.T) {
	nodes := []HardwareNode{
		{Text: "AMD Ryzen 5 7600X"},
		{Text: "NVIDIA GeForce RTX"},
	}
	if n := findNodeByPattern(nodes, "Ryzen|AMD"); n == nil || !strings.Contains(n.Text, "Ryzen") {
		t.Error("should match Ryzen via alternative pattern")
	}
	if n := findNodeByPattern(nodes, "Intel"); n != nil {
		t.Error("no match should return nil")
	}
}
