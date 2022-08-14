package reward

import (
	"testing"

	"github.com/PositionExchange/posichain/numeric"
)

func TestPercentageForTimeStamp(t *testing.T) {
	testCases := []struct {
		time     string
		expected string
	}{
		{"2022-Jul-01", "1.000000000000000"},
		{"2022-Jul-31", "1.000000000000000"},
		{"2023-Jul-30", "1.000000000000000"},
		{"2026-Apr-29", "1.000000000000000"},
		{"2023-Apr-30", "1.000000000000000"},
		{"2028-May-31", "1.000000000000000"},
	}

	for _, tc := range testCases {
		result := PercentageForTimeStamp(mustParse(tc.time))
		expect := numeric.MustNewDecFromStr(tc.expected)
		if !result.Equal(expect) {
			t.Errorf("Time: %s, Chosen bucket percent: %s, Expected: %s",
				tc.time, result, expect)
		}
	}
}
