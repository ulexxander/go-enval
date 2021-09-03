package enval_test

import (
	"encoding/json"
	"testing"

	"github.com/ulexxander/enval"
)

// abc used to test custom value types
type abc struct {
	Abc int `json:"abc"`
}

// abcParseFunc parses abc type from string
func abcParseFunc(val string) (interface{}, error) {
	var out abc
	err := json.Unmarshal([]byte(val), &out)
	return out, err
}

// valType used to test multiple kinds of values
// in one table driven test
type valType int

const (
	valString valType = iota
	valInt
	valBool
	valCustom
)

func TestValuesAndErrors(t *testing.T) {
	s := enval.NewLookuper()
	vars := map[string]string{
		"STRING_PRESENT": ":80",
		// "STRING_MISSING": "actually required",

		"INT_PRESENT": "16",
		"INT_INVALID": "b4dint34",
		// "INT_MISSING":     "actually required",

		"BOOL_PRESENT": "true",
		"BOOL_INVALID": "nOTtRueOrFalsE",
		// "BOOL_MISSING": "actually required",

		"CUSTOM_PRESENT": `{"abc": 456}`,
		"CUSTOM_INVALID": `}"abc": 456{`,
		// "CUSTOM_MISSING": "actually required",
	}
	s.LookupFunc = func(key string) (string, bool) {
		val, present := vars[key]
		return val, present
	}

	tt := []struct {
		key       string
		valType   valType
		valString string
		valInt    int
		valBool   bool
		valCustom abc
		hasErr    bool
	}{
		{key: "STRING_PRESENT", valType: valString, valString: ":80"},
		{key: "STRING_MISSING", valType: valString, hasErr: true},
		{key: "INT_PRESENT", valType: valInt, valInt: 16},
		{key: "INT_INVALID", valType: valInt, hasErr: true},
		{key: "INT_MISSING", valType: valInt, hasErr: true},
		{key: "BOOL_PRESENT", valType: valBool, valBool: true},
		{key: "BOOL_INVALID", valType: valInt, hasErr: true},
		{key: "BOOL_MISSING", valType: valInt, hasErr: true},
		{key: "CUSTOM_PRESENT", valType: valCustom, valCustom: abc{Abc: 456}},
		{key: "CUSTOM_INVALID", valType: valCustom, hasErr: true},
		{key: "CUSTOM_MISSING", valType: valCustom, hasErr: true},
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

			if tc.valType == valCustom {
				v, ok := s.Custom(tc.key, abcParseFunc).(abc)
				if ok && v.Abc != tc.valCustom.Abc {
					t.Fatalf("expected custom abc value to be %d, got: %d", tc.valCustom.Abc, v.Abc)
				}
			}

			err := s.ErrByVariable[tc.key]
			if !tc.hasErr && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if tc.hasErr && err == nil {
				t.Fatalf("expected error here")
			}
		})
	}
}
