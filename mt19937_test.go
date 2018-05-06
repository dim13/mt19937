package mt19937

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestSource(t *testing.T) {
	s := new(Source)
	s.SeedByArray([]uint64{0x12345, 0x23456, 0x34567, 0x45678})
	testCases := []struct {
		title  string
		golden string
		format func(s *Source) string
	}{
		{
			title:  "uint64",
			golden: "testdata/uint64.golden",
			format: func(s *Source) string {
				return fmt.Sprintf("%20d", s.Uint64())
			},
		},
		{
			title:  "float64",
			golden: "testdata/float64.golden",
			format: func(s *Source) string {
				return fmt.Sprintf("%10.8f", s.Float64B())
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			fd, err := os.Open(tc.golden)
			if err != nil {
				t.Fatal(err)
			}
			defer fd.Close()
			scanner := bufio.NewScanner(fd)
			for scanner.Scan() {
				if got, want := tc.format(s), scanner.Text(); got != want {
					t.Errorf("got %v; want %v", got, want)
				}
			}
			if err := scanner.Err(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func BenchmarkSource(b *testing.B) {
	s := new(Source)
	for i := 0; i < b.N; i++ {
		s.Uint64()
	}
}
