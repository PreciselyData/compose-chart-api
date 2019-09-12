package pic

// DataStyle represents a style to apply to a data value.
type DataStyle struct {
	Type     string
	Settings map[string]Value
}

// DataStyles represents a set of styles to apply to the data values.
type DataStyles [][]DataStyle

// At returns the data style of a particular data value.
func (ds DataStyles) At(row, col int) *DataStyle {
	if row < 0 || row >= len(ds) || col < 0 || col >= len(ds[row]) {
		return nil
	}
	return &ds[row][col]
}

// Setting determines a data style setting for a particular data value.
func (ds DataStyles) Setting(row, col int, name string) Value {
	style := ds.At(row, col)
	if style == nil {
		return ""
	}
	val, ok := style.Settings[name]
	if !ok {
		return ""
	}
	return val
}

// CustomFormat determines the special custom format setting of a data value.
func (ds DataStyles) CustomFormat(row, col int) Value {
	style := ds.At(row, col)
	if style == nil || style.Type != "custom" {
		return ""
	}
	val, ok := style.Settings["customFmt"]
	if !ok {
		return ""
	}
	return val
}
