package pic

import (
	"io/ioutil"

	"github.com/golang/freetype/truetype"
)

// FontResource represents a font resource from Designer/Generate.
type FontResource struct {
	Typeface     string
	PointSize    float64
	Attributes   FontAttribute
	Filename     string
	TruetypeFont *truetype.Font
}

// Bold determines whether the font has the Bold attribute.
func (fr *FontResource) Bold() bool {
	return (fr.Attributes & Bold) != 0
}

// Italic determines whether the font has the Italic attribute.
func (fr *FontResource) Italic() bool {
	return (fr.Attributes & Italic) != 0
}

// EmulateBold determines whether the font has the EmulateBold attribute.
func (fr *FontResource) EmulateBold() bool {
	return (fr.Attributes & EmulateBold) != 0
}

// EmulateItalic determines whether the font has the EmulateItalic attribute.
func (fr *FontResource) EmulateItalic() bool {
	return (fr.Attributes & EmulateItalic) != 0
}

// EmulateTypeface determines whether the font has the EmulateTypeface attribute.
func (fr *FontResource) EmulateTypeface() bool {
	return (fr.Attributes & EmulateTypeface) != 0
}

func (fr *FontResource) loadTruetype() error {
	if fr.TruetypeFont == nil {
		if fr.Filename == "" {
			return nil
		}
		data, err := ioutil.ReadFile(fr.Filename)
		if err != nil {
			return err
		}
		fr.TruetypeFont, err = truetype.Parse(data)
		if err != nil {
			return err
		}
	}
	return nil
}
