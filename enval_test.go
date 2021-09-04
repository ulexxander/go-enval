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

// testVariables contains mock env variables to test
// it should cover all possible cases
var testVariables = map[string]string{
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

func testVariablesLookupFunc(key string) (string, bool) {
	val, present := testVariables[key]
	return val, present
}

func TestValuesAndErrors(t *testing.T) {
	l := enval.NewLookuper()
	l.LookupFunc = testVariablesLookupFunc

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
				s := l.String(tc.key)
				if s != tc.valString {
					t.Fatalf("expected string value to be %s, got: %s", tc.valString, s)
				}
			}

			if tc.valType == valInt {
				i := l.Int(tc.key)
				if i != tc.valInt {
					t.Fatalf("expected int value to be %d, got: %d", tc.valInt, i)
				}
			}

			if tc.valType == valBool {
				b := l.Bool(tc.key)
				if b != tc.valBool {
					t.Fatalf("expected bool value to be %t, got: %t", tc.valBool, b)
				}
			}

			if tc.valType == valCustom {
				v, ok := l.Custom(tc.key, abcParseFunc).(abc)
				if ok && v.Abc != tc.valCustom.Abc {
					t.Fatalf("expected custom abc value to be %d, got: %d", tc.valCustom.Abc, v.Abc)
				}
			}

			err := l.ErrByVariable[tc.key]
			if !tc.hasErr && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if tc.hasErr && err == nil {
				t.Fatalf("expected error here")
			}
		})
	}
}

func TestDefaults(t *testing.T) {
	l := enval.NewLookuper()
	l.LookupFunc = testVariablesLookupFunc

	type valStringDef struct{ expected, def string }
	type valIntDef struct{ expected, def int }
	type valBoolDef struct{ expected, def bool }
	type valCustomDef struct{ expected, def abc }

	tt := []struct {
		key       string
		valType   valType
		valString valStringDef
		valInt    valIntDef
		valBool   valBoolDef
		valCustom valCustomDef
		hasErr    bool
	}{
		{key: "STRING_PRESENT", valType: valString, valString: valStringDef{":80", ":10030"}},
		{key: "STRING_MISSING", valType: valString, valString: valStringDef{":10030", ":10030"}},
		{key: "INT_PRESENT", valType: valInt, valInt: valIntDef{16, 550}},
		{key: "INT_INVALID", valType: valInt, hasErr: true},
		{key: "INT_MISSING", valType: valInt, valInt: valIntDef{550, 550}},
		{key: "BOOL_PRESENT", valType: valBool, valBool: valBoolDef{true, false}},
		{key: "BOOL_INVALID", valType: valInt, hasErr: true},
		{key: "BOOL_MISSING", valType: valInt, valBool: valBoolDef{true, true}},
		// {key: "CUSTOM_PRESENT", valType: valCustom, valCustom: valCustomDef{abc{Abc: 456}, abc{Abc: 999}},
		{key: "CUSTOM_INVALID", valType: valCustom, hasErr: true},
		{key: "CUSTOM_MISSING", valType: valCustom, valCustom: valCustomDef{abc{Abc: 999}, abc{Abc: 999}}},
	}

	for _, tc := range tt {
		t.Run(tc.key, func(t *testing.T) {
			if tc.valType == valString {
				s := l.StringWithDefault(tc.key, tc.valString.def)
				if s != tc.valString.expected {
					t.Fatalf("expected string value to be %s, got: %s", tc.valString.expected, s)
				}
			}

			if tc.valType == valInt {
				i := l.IntWithDefault(tc.key, tc.valInt.def)
				if i != tc.valInt.expected {
					t.Fatalf("expected int value to be %d, got: %d", tc.valInt.expected, i)
				}
			}

			if tc.valType == valBool {
				b := l.BoolWithDefault(tc.key, tc.valBool.def)
				if b != tc.valBool.expected {
					t.Fatalf("expected bool value to be %t, got: %t", tc.valBool.expected, b)
				}
			}

			if tc.valType == valCustom {
				v, ok := l.CustomWithDefault(tc.key, tc.valCustom.def, abcParseFunc).(abc)
				if ok && v.Abc != tc.valCustom.expected.Abc {
					t.Fatalf("expected custom abc value to be %d, got: %d", tc.valCustom.expected.Abc, v.Abc)
				}
			}

			err := l.ErrByVariable[tc.key]
			if !tc.hasErr && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if tc.hasErr && err == nil {
				t.Fatalf("expected error here")
			}
		})
	}
}

func TestErr(t *testing.T) {
	l := enval.NewLookuper()
	l.LookupFunc = testVariablesLookupFunc

	tt := []struct {
		key          string
		valType      valType
		errTextDelta string
	}{
		{key: "STRING_PRESENT", valType: valString},
		{key: "STRING_MISSING", valType: valString, errTextDelta: "STRING_MISSING: key missing"},
		{key: "INT_PRESENT", valType: valInt},
		{key: "INT_INVALID", valType: valInt, errTextDelta: `, INT_INVALID: unparsable int: strconv.ParseInt: parsing "b4dint34": invalid syntax`},
		{key: "INT_MISSING", valType: valInt, errTextDelta: `, INT_MISSING: key missing`},
		{key: "BOOL_PRESENT", valType: valBool},
		{key: "BOOL_INVALID", valType: valBool, errTextDelta: `, BOOL_INVALID: unparsable bool: strconv.ParseBool: parsing "nOTtRueOrFalsE": invalid syntax`},
		{key: "BOOL_MISSING", valType: valBool, errTextDelta: `, BOOL_MISSING: key missing`},
		{key: "CUSTOM_PRESENT", valType: valCustom},
		{key: "CUSTOM_INVALID", valType: valCustom, errTextDelta: `, CUSTOM_INVALID: invalid character '}' looking for beginning of value`},
		{key: "CUSTOM_MISSING", valType: valCustom, errTextDelta: `, CUSTOM_MISSING: key missing`},
	}

	var errText string
	for _, tc := range tt {
		t.Run(tc.key, func(t *testing.T) {
			switch tc.valType {
			case valString:
				l.String(tc.key)
			case valInt:
				l.Int(tc.key)
			case valBool:
				l.Bool(tc.key)
			case valCustom:
				l.Custom(tc.key, abcParseFunc)
			}

			errText += tc.errTextDelta

			err := l.Err()
			if err == nil && errText != "" {
				t.Log("expected:\n", errText)
				t.Fatalf("expected to have error, got nil")
			}

			if err != nil && errText != err.Error() {
				t.Log("expected:\n", errText)
				t.Log("actual:\n", err)
				t.Fatalf("error texts are not equal")
			}
		})
	}
}
