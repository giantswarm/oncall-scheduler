package opsgenieclient

import (
	"strconv"
	"testing"
	"time"

	"github.com/giantswarm/micrologger/microloggertest"
)

func Test_GetUnixTime(t *testing.T) {
	testCases := []struct {
		name           string
		input_time     time.Time
		input_shift    int
		expectedOutput int64
	}{
		{
			name:           "case 0: normal date, no shift",
			input_time:     time.Date(2019, 2, 10, 9, 0, 0, 0, time.UTC),
			input_shift:    0,
			expectedOutput: 1549789200000,
		},
		{
			name:           "case 0: normal date, 1 day shift",
			input_time:     time.Date(2019, 2, 10, 9, 0, 0, 0, time.UTC),
			input_shift:    1,
			expectedOutput: 1549702800000,
		},
		{
			name:           "case 0: normal date, large shift",
			input_time:     time.Date(2019, 2, 10, 9, 0, 0, 0, time.UTC),
			input_shift:    30,
			expectedOutput: 1547197200000,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client, err := New(Config{
				Logger: microloggertest.New(),
				APIKey: "test",
			})
			if err != nil {
				t.Fatalf("could not create client: %#v", err)
			}

			output := client.getUnixTime(tc.input_time, tc.input_shift)

			if output != tc.expectedOutput {
				t.Fatalf("wanted: %v, got: %v", tc.expectedOutput, output)
			}
		})
	}
}

func Test_CalculatePercentageChange(t *testing.T) {
	testCases := []struct {
		name           string
		input_a        int
		input_b        int
		expectedOutput int
	}{
		{
			name:           "case 0: 0 change to 0",
			input_a:        0,
			input_b:        0,
			expectedOutput: 0,
		},
		// Note: There is no 'right' output here,
		// as any percentage increase from 0 is infinite.
		// Choosing 0 as the least 'wrong' option.
		{
			name:           "case 1: 0 change to 1",
			input_a:        0,
			input_b:        1,
			expectedOutput: 100,
		},
		{
			name:           "case 2: 1 change to 0",
			input_a:        1,
			input_b:        0,
			expectedOutput: -100,
		},
		{
			name:           "case 3: 1 change to 2",
			input_a:        1,
			input_b:        2,
			expectedOutput: 100,
		},
		{
			name:           "case 4: 5 change to 20",
			input_a:        5,
			input_b:        20,
			expectedOutput: 300,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client, err := New(Config{
				Logger: microloggertest.New(),
				APIKey: "test",
			})
			if err != nil {
				t.Fatalf("could not create client: %#v", err)
			}

			output := client.calculatePercentageChange(tc.input_a, tc.input_b)

			if output != tc.expectedOutput {
				t.Fatalf("wanted: %v, got: %v", tc.expectedOutput, output)
			}
		})
	}
}
