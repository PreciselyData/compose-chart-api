package pic

import (
	"fmt"
	"strconv"
)

// Color represents a colour value from the chart configuration.
type Color struct {
	R, G, B    uint8
	C, M, Y, K uint8
}

// DefaultColor represents the default colour value (black).
var DefaultColor = Color{K: 100}

func (c *Color) parse(v []Value) error {
	if len(v) == 1 {
		// The value just contains the named colour index.
		c.parseNamedColor(v[0])
		return nil
	}
	if len(v) == 4 {
		// Not interested in the first two attributes:
		// - Master colour type (Named, RGB or CMYK).
		// - Named colour index.
		rgb := c.parseAttribute(v[2])
		cmyk := c.parseAttribute(v[3])
		c.R = uint8((rgb >> 16) & 0xff)
		c.G = uint8((rgb >> 8) & 0xff)
		c.B = uint8(rgb & 0xff)
		c.C = uint8((cmyk >> 24) & 0xff)
		c.M = uint8((cmyk >> 16) & 0xff)
		c.Y = uint8((cmyk >> 8) & 0xff)
		c.K = uint8(cmyk & 0xff)
		return nil
	}
	return fmt.Errorf("invalid color set %v", v)
}

func (c *Color) parseNamedColor(v Value) {
	index := c.parseAttribute(v)
	if index > 15 {
		index = 0
	}
	c.R, c.G, c.B = c.namedRgb(index)
	c.C, c.M, c.Y, c.K = c.namedCmyk(index)
}

func (Color) parseAttribute(v Value) int32 {
	i, err := strconv.ParseInt(string(v), 10, 32)
	if err != nil {
		i = 0
	}
	return int32(i)
}

var namedRgbColors = [][]uint8{
	{0, 0, 0},       // Black
	{0, 0, 255},     // Blue
	{144, 48, 0},    // Brown
	{0, 255, 0},     // Green
	{255, 0, 255},   // Magenta
	{255, 0, 0},     // Red
	{0, 255, 255},   // Cyan
	{255, 255, 0},   // Yellow
	{0, 0, 170},     // Dark Blue
	{0, 146, 0},     // Dark Green
	{0, 146, 170},   // Teal
	{131, 131, 131}, // Gray
	{196, 160, 32},  // Mustard
	{255, 128, 0},   // Orange
	{170, 0, 170},   // Purple
	{255, 255, 255}, // White
}

func (Color) namedRgb(index int32) (r, g, b uint8) {
	rgb := namedRgbColors[index]
	r, b, g = rgb[0], rgb[1], rgb[2]
	return
}

var namedCmykColors = [][]uint8{
	{0, 0, 0, 100},   // Black
	{100, 100, 0, 0}, // Blue
	{0, 38, 57, 43},  // Brown
	{100, 0, 100, 0}, // Green
	{0, 100, 0, 0},   // Magenta
	{0, 100, 100, 0}, // Red
	{100, 0, 0, 0},   // Cyan
	{0, 0, 100, 0},   // Yellow
	{67, 67, 0, 33},  // Dark Blue
	{58, 0, 58, 42},  // Dark Green
	{67, 9, 0, 33},   // Teal
	{0, 0, 0, 48},    // Gray
	{0, 14, 64, 23},  // Mustard
	{0, 49, 100, 0},  // Orange
	{0, 67, 0, 33},   // Purple
	{0, 0, 0, 0},     // White
}

func (Color) namedCmyk(index int32) (c, m, y, k uint8) {
	cmyk := namedCmykColors[index]
	c, m, y, k = cmyk[0], cmyk[1], cmyk[2], cmyk[3]
	return
}
