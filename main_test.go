package main

import "testing"

func TestRun(t *testing.T) {
	testCases := []struct {
		src    string
		strict bool
		ok     bool
	}{
		{"testdata/t1.json", false, true},
		{"testdata/t2.json", false, true},
		{"testdata/t3.json", false, false},
		{"testdata/t4.json", false, true},
		// Strict mode
		{"testdata/t1.json", true, true},
		{"testdata/t2.json", true, false}, // array of integers
		{"testdata/t3.json", true, false},
		{"testdata/t4.json", true, false}, // camelcase
	}

	for _, tc := range testCases {
		if err := run("testdata/example.proto", "foo.Bar", nil, tc.src, tc.strict); (err == nil) != tc.ok {
			t.Errorf("%q expected ok? %v, error: %+v", tc.src, tc.ok, err)
		}
	}
}
