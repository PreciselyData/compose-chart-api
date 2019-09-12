package pic

import (
	"fmt"
	"strconv"
)

// Font represents a font value from the chart configuration.
type Font struct {
	IsStyle bool
	GUID
	Color
	Underline bool
}

// DefaultFont can be used to query the default font from Designer/Generate.
var DefaultFont = Font{}

func (f *Font) parse(ds Dataset) error {
	v := ds[0][0]
	if len(v) < 2 || v[0] != ascESC || v[1] != 'f' {
		return fmt.Errorf("invalid font value '%s'", v)
	}
	s := string(v[2:])

	if len(s) > 1 && s[0] == '$' {
		f.IsStyle = true
		s = s[1:]
	} else {
		f.IsStyle = false
	}

	if (f.IsStyle && len(ds) != 1) || (!f.IsStyle && len(ds) != 3) {
		return fmt.Errorf("invalid font set %v", ds)
	}

	err := f.GUID.parse(s)
	if err != nil {
		return err
	}

	if !f.IsStyle {
		if err = f.Color.parse(ds[1]); err != nil {
			return err
		}
		if f.Underline, err = strconv.ParseBool(string(ds[2][0])); err != nil {
			return err
		}
	}

	return nil
}
