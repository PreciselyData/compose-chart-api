package pic

import (
	"log"
	"strings"
)

// Config represents the configuration of the chart to be rendered.
type Config struct {
	resolver
	properties, symbols map[string]string
	fontResources       map[GUID]*FontResource
	fontStyles          map[GUID]*FontStyle
}

func newConfig(r resolver, props, syms string) *Config {
	return &Config{
		resolver:      r,
		properties:    loadSettings(props, '\n'),
		symbols:       loadSettings(syms, '\n'),
		fontResources: make(map[GUID]*FontResource),
		fontStyles:    make(map[GUID]*FontStyle),
	}
}

// NumberFormat defines how a number should be formatted for display.
func (c *Config) NumberFormat() NumberFormat {
	return c.resolver.numberFormat()
}

// Value gets the value of a property from the configuration.
func (c *Config) Value(name string) Value {
	if val, ok := c.properties[name]; ok {
		return Value(c.lookupSymbol(val))
	}
	return ""
}

// Integer gets the value of a property as an integer. Zero is returned
// if the string is empty or the conversion fails.
func (c *Config) Integer(name string) int32 {
	return c.ResolveInteger(c.Value(name))
}

// Twiplet gets the value of a property as a Twiplet. Zero is returned
// if the string is empty or the conversion fails.
func (c *Config) Twiplet(name string) Twiplet {
	return c.ResolveTwiplet(c.Value(name))
}

// Number gets the value of a property as a number. Zero is returned
// if the string is empty or the conversion fails.
func (c *Config) Number(name string) float64 {
	return c.ResolveNumber(c.Value(name))
}

// ResolveInteger converts a value to an integer. Zero is returned
// if the string is empty or the conversion fails.
func (c *Config) ResolveInteger(v Value) int32 {
	if v == "" {
		return 0
	}
	i, err := c.resolver.integer(string(v))
	if err != nil {
		log.Println(err)
		return 0
	}
	return i
}

// ResolveTwiplet converts a value to a Twiplet. Zero is returned
// if the string is empty or the conversion fails.
func (c *Config) ResolveTwiplet(v Value) Twiplet {
	return Twiplet(c.ResolveInteger(v))
}

// ResolveNumber converts a value to a number. Zero is returned
// if the string is empty or the conversion fails.
func (c *Config) ResolveNumber(v Value) float64 {
	if v == "" {
		return 0
	}
	n, err := c.resolver.number(string(v))
	if err != nil {
		log.Println(err)
		return 0
	}
	return n
}

// Color gets the value of a property as a Color.
func (c *Config) Color(name string) Color {
	if val, ok := c.properties[name]; ok && val != "" {
		return c.loadColor(val)
	}
	return DefaultColor
}

// Font gets the value of a property as a Font.
func (c *Config) Font(name string) Font {
	if val, ok := c.properties[name]; ok && val != "" {
		return c.loadFont(val)
	}
	return DefaultFont
}

// ResolveFont gets a font resource from a Font value.
func (c *Config) ResolveFont(f Font) *FontStyle {
	if f.IsStyle {
		return c.resolveFontStyle(f)
	}
	return &FontStyle{
		FontResource: c.resolveFontResource(f),
		Color:        f.Color,
		Underline:    f.Underline,
	}
}

// Dataset gets a set of data values from the configuration.
func (c *Config) Dataset(name string) Dataset {
	if val, ok := c.properties[name]; ok && val != "" {
		return c.loadDataset(val)
	}
	return Dataset{[]Value{""}}
}

// Name gets the configuration name.
func (c *Config) Name() string {
	return string(c.Value("config"))
}

// Data gets all of the data properties from the configuration.
func (c *Config) Data() *Data {
	return &Data{
		Values:  c.DataValues(),
		Titles:  c.DataTitles(),
		Colors:  c.DataColors(),
		Styles:  c.DataStyles(),
		Labels:  c.DataLabels(),
		Fonts:   c.DataFonts(),
		Formats: c.DataFormats(),
	}
}

// DataValues gets the data values dataset from the configuration.
func (c *Config) DataValues() Dataset {
	return c.Dataset("data.values")
}

// DataTitles gets the data titles dataset from the configuration.
func (c *Config) DataTitles() (titles []Value) {
	ds := c.Dataset("data.titles")
	titles = make([]Value, len(ds))
	for i, set := range ds {
		titles[i] = set[0]
	}
	return
}

// DataColors gets the data colours dataset from the configuration.
func (c *Config) DataColors() (colors []Color) {
	ds := c.Dataset("data.colors")
	if len(ds) == 1 && len(ds[0]) > 0 {
		// Get the data colours from each value in the single series.
		colors = make([]Color, len(ds[0]))
		for i, val := range ds[0] {
			colors[i] = c.loadColor(string(val))
		}
	} else {
		// Get the series colours from the first value in each set.
		colors = make([]Color, len(ds))
		for i, vals := range ds {
			colors[i] = c.loadColor(string(vals[0]))
		}
	}
	return
}

// DataStyles gets the data styles dataset from the configuration.
func (c *Config) DataStyles() DataStyles {
	return c.loadDataStyles(c.Dataset("data.styles"))
}

// DataLabels gets the data labels dataset from the configuration.
func (c *Config) DataLabels() []Value {
	ds := c.Dataset("data.labels")
	return ds[0]
}

// DataFonts gets the data fonts dataset from the configuration.
func (c *Config) DataFonts() (fonts []Font) {
	ds := c.Dataset("data.fonts")
	fonts = make([]Font, len(ds[0]))
	for i, val := range ds[0] {
		fonts[i] = c.loadFont(string(val))
	}
	return
}

// DataFormats gets the data formats dataset from the configuration.
func (c *Config) DataFormats() DataStyles {
	return c.loadDataStyles(c.Dataset("data.formats"))
}

func loadSettings(input string, sep byte) map[string]string {
	out := make(map[string]string)
	for _, line := range strings.Split(input, string(sep)) {
		line = strings.Trim(line, "\r")
		setting := strings.SplitN(line, "=", 2)
		if len(setting) == 2 {
			out[setting[0]] = setting[1]
		}
	}
	return out
}

func (c *Config) lookupSymbol(val string) string {
	if val == "" || val[0] != ascDLE {
		return val
	}
	if res, ok := c.symbols[val[1:]]; ok {
		return res
	}
	return ""
}

func (c *Config) loadDataset(input string) (ds Dataset) {
	// If the first character is SOH then RS and US are being used as the
	// set and value separators, otherwise '|' and ',' are being used.
	var setSep, valSep string
	if input != "" && input[0] == ascSOH {
		setSep = string(ascRS)
		valSep = string(ascUS)
		input = input[1:]
	} else {
		setSep = "|"
		valSep = ","
	}

	sets := strings.Split(input, setSep)
	ds = make(Dataset, len(sets))
	for i, set := range sets {
		vals := strings.Split(set, valSep)
		ds[i] = make([]Value, 0, len(vals))
		for _, val := range vals {
			// Each value could be a symbol reference to a list of values.
			for _, s := range strings.Split(c.lookupSymbol(val), valSep) {
				ds[i] = append(ds[i], Value(s))
			}
		}
	}
	return
}

func (c *Config) loadDataStyles(ds Dataset) (styles DataStyles) {
	styles = make(DataStyles, len(ds))
	for i, set := range ds {
		styles[i] = make([]DataStyle, len(set))
		for j, val := range set {
			styles[i][j] = c.loadDataStyle(val)
		}
	}
	return
}

func (c *Config) loadDataStyle(v Value) (style DataStyle) {
	options := strings.SplitN(string(v), ":", 2)
	if len(options) != 2 {
		return
	}
	style.Type = options[0]
	style.Settings = make(map[string]Value)
	if len(options[1]) > 0 {
		settings := loadSettings(options[1][1:], options[1][0])
		for k, s := range settings {
			style.Settings[k] = Value(c.lookupSymbol(s))
		}
	}
	return
}

func (c *Config) loadColor(val string) (color Color) {
	// A colour value is represented by a single-row dataset.
	ds := c.loadDataset(val)
	if err := color.parse(ds[0]); err != nil {
		log.Println(err)
		return DefaultColor
	}
	return
}

func (c *Config) loadFont(val string) (f Font) {
	// A font value is represented by a multi-row dataset.
	ds := c.loadDataset(val)
	if err := f.parse(ds); err != nil {
		log.Println(err)
		return DefaultFont
	}
	return
}

func (c *Config) resolveFontResource(f Font) *FontResource {
	fr, ok := c.fontResources[f.GUID]
	if ok {
		return fr
	}
	var err error
	if fr, err = c.resolver.fontResource(f.GUID); err != nil {
		log.Println(err)
		fr = &FontResource{}
	} else if err = fr.loadTruetype(); err != nil {
		log.Println(err)
	}
	c.fontResources[f.GUID] = fr
	return fr
}

func (c *Config) resolveFontStyle(f Font) *FontStyle {
	fs, ok := c.fontStyles[f.GUID]
	if ok {
		return fs
	}
	var err error
	if fs, err = c.resolver.fontStyle(f.GUID); err != nil {
		log.Println(err)
		fs = &FontStyle{
			FontResource: &FontResource{},
			Color:        DefaultColor,
		}
	} else if err = fs.FontResource.loadTruetype(); err != nil {
		log.Println(err)
	}
	c.fontStyles[f.GUID] = fs
	return fs
}
