package enval_test

import (
	"fmt"
	"testing"

	"github.com/ulexxander/enval"
)

func TestValuesAndErrors(t *testing.T) {
	s := enval.NewLookuper()
	vars := map[string]string{
		// strings
		"STRING_PRESENT": ":80",
		// "STRING_MISSING": "actually required",

		// ints
		"INT_PRESENT": "16",
		// "INT_MISSING":     "actually required",
		"INT_INVALID": "475bads",

		// bools
		"BOOL_PRESENT": "true",
		// "BOOL_MISSING": "missing",
		"BOOL_INVALID": "nOTtRueOrFalsE",

		// custom
		"CUSTOM_PRESENT": "TRACE",
		// "CUSTOM_MISSING": "TRACE",
	}
	s.LookupFunc = func(key string) (string, bool) {
		val, present := vars[key]
		return val, present
	}

	const (
		valString int = iota
		valInt
		valBool
		valCustom
	)

	tt := []struct {
		key       string
		valType   int
		valString string
		valInt    int
		valBool   bool
		hasErr    bool
	}{
		{key: "STRING_PRESENT", valType: valString, valString: ":80"},
		{key: "STRING_MISSING", valType: valString, hasErr: true},
		{key: "INT_PRESENT", valType: valInt, valInt: 16},
		{key: "INT_MISSING", valType: valInt, hasErr: true},
		{key: "INT_INVALID", valType: valInt, hasErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.key, func(t *testing.T) {
			if tc.valType == valString {
				s := s.String(tc.key)
				if s != tc.valString {
					t.Fatalf("expected string value to be %s, got: %s", tc.valString, s)
				}
			}

			if tc.valType == valInt {
				i := s.Int(tc.key)
				if i != tc.valInt {
					t.Fatalf("expected int value to be %d, got: %d", tc.valInt, i)
				}
			}

			if tc.valType == valBool {
				b := s.Bool(tc.key)
				if b != tc.valBool {
					t.Fatalf("expected bool value to be %t, got: %t", tc.valBool, b)
				}
			}

			err := s.ErrByVariable[tc.key]
			if !tc.hasErr && err != nil {
				fmt.Printf("unexpected error: %s", err)
			}

			if tc.hasErr && err == nil {
				fmt.Printf("expected error here")
			}
		})
	}
}
