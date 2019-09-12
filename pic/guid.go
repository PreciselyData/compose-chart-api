package pic

import "fmt"

// GUID represents a GUID value from the chart configuration.
type GUID [16]byte

// IsZero determines whether the GUID is all zeroes.
func (g *GUID) IsZero() bool {
	for _, b := range g {
		if b != 0 {
			return false
		}
	}
	return true
}

func (g *GUID) parse(val string) error {
	if len(val) != 32 {
		return fmt.Errorf("invalid GUID format '%s'", val)
	}
	for i := 0; i < 16; i++ {
		b1, err := g.charToByte(val[2*i])
		if err != nil {
			return err
		}
		b2, err := g.charToByte(val[2*i+1])
		if err != nil {
			return err
		}
		g[i] = (b1 << 4) | b2
	}
	return nil
}

func (GUID) charToByte(c byte) (byte, error) {
	if c >= '0' && c <= '9' {
		return c - '0', nil
	}
	if c >= 'A' && c <= 'F' {
		return c - 'A' + 0x0A, nil
	}
	if c >= 'a' && c <= 'f' {
		return c - 'a' + 0x0a, nil
	}
	return 0, fmt.Errorf("invalid GUID char '%c'", c)
}
