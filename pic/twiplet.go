package pic

// Twiplet represents a unit of measure in 1/100ths of a twip.
type Twiplet int32

// Inch defines the number of units in one inch.
const Inch = 144000.0

// Inches converts the measurement to inches.
func (t Twiplet) Inches() float64 {
	return float64(t) / Inch
}

// Pixels converts the measurement to pixels.
func (t Twiplet) Pixels(dpi int32) int {
	f := t.Inches() * float64(dpi)
	if f < 0 {
		return int(f - 0.5)
	}
	return int(f + 0.5)
}
