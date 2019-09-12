package pic

import "log"

// Value represents a value of a property from the configuration.
type Value string

// Type returns the value type.
func (v Value) Type() DataType {
	if len(v) > 1 && v[0] == ascESC {
		switch v[1] {
		case 'i':
			return Integer
		case 'n':
			return Number
		case 'd':
			return Date
		case 't':
			return Time
		case '$':
			return Currency
		default:
			log.Printf("Unrecognised value type '%c'\n", v[1])
			return NotSet
		}
	}
	return Neutral
}

// Text converts the value to displayable text by removing the escape
// sequence identifying the value type.
func (v Value) Text() string {
	if len(v) > 1 && v[0] == ascESC {
		return string(v)[2:]
	}
	return string(v)
}

// True determines whether a bool value is true.
func (v Value) True() bool {
	return v == "true"
}
