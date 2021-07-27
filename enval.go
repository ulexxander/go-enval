package enval

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Lookuper struct {
	ErrByVariable map[string]error
	LookupFunc    func(key string) (string, bool)
}

func NewLookuper() *Lookuper {
	return &Lookuper{
		ErrByVariable: map[string]error{},
		LookupFunc:    os.LookupEnv,
	}
}

func (l *Lookuper) String(key string) string {
	val, present := l.LookupFunc(key)
	if !present {
		l.addError(key, fmt.Errorf("missing"))
		return ""
	}
	return val
}

func (l *Lookuper) StringWithDefault(key string, def string) string {
	val, present := l.LookupFunc(key)
	if !present {
		return def
	}
	return val
}

func (l *Lookuper) Int(key string) int {
	val, present := l.LookupFunc(key)
	if !present {
		l.addError(key, fmt.Errorf("missing"))
		return 0
	}
	return l.parseInt(key, val)
}

func (l *Lookuper) IntWithDefault(key string, def int) int {
	val, present := l.LookupFunc(key)
	if !present {
		return def
	}
	return l.parseInt(key, val)
}

func (l *Lookuper) parseInt(key, val string) int {
	valInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		l.addError(key, fmt.Errorf("unparsable: %s", err))
	}
	return int(valInt)
}

func (l *Lookuper) Bool(key string) bool {
	val, present := l.LookupFunc(key)
	if !present {
		l.addError(key, fmt.Errorf("missing"))
		return false
	}
	return l.parseBool(key, val)
}

func (l *Lookuper) BoolWithDefault(key string, def bool) bool {
	val, present := l.LookupFunc(key)
	if !present {
		return def
	}
	return l.parseBool(key, val)
}

func (l *Lookuper) parseBool(key, val string) bool {
	valBool, err := strconv.ParseBool(val)
	if err != nil {
		l.addError(key, fmt.Errorf("unparsable: %s", err))
	}
	return valBool
}

type ParseFunc func(val string, present bool) error

func (l *Lookuper) Custom(key string, pf ParseFunc) {
	val, present := l.LookupFunc(key)
	if !present {
		l.addError(key, fmt.Errorf("missing"))
		return
	}
	if err := pf(val, present); err != nil {
		l.addError(key, err)
	}
}

func (l *Lookuper) CustomWithDefault(key string, pf ParseFunc) {
	val, present := l.LookupFunc(key)
	if err := pf(val, present); err != nil {
		l.addError(key, err)
	}
}

func (l *Lookuper) addError(key string, err error) {
	l.ErrByVariable[key] = err
}

func (l *Lookuper) Err() error {
	if l.ErrByVariable == nil {
		return nil
	}

	var errStr string
	for key, varErr := range l.ErrByVariable {
		errStr += key + ":" + varErr.Error() + "\n"
	}

	return errors.New(errStr)
}
