package pic

// #include "c_interface.h"
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// ReturnCode represents the API return codes.
type ReturnCode int32

// API return codes understood by Designer/Generate.
const (
	OK                   ReturnCode = C.ENCHRC_OK
	Failed               ReturnCode = C.ENCHRC_Failed
	NotImplemented       ReturnCode = C.ENCHRC_NotImplemented
	InvalidFilePath      ReturnCode = C.ENCHRC_InvalidFilePath
	MissingProperty      ReturnCode = C.ENCHRC_MissingProperty
	InvalidValue         ReturnCode = C.ENCHRC_InvalidValue
	UnresolvedFont       ReturnCode = C.ENCHRC_UnresolvedFont
	InvalidDataString    ReturnCode = C.ENCHRC_InvalidDataString
	EmptyDataString      ReturnCode = C.ENCHRC_EmptyDataString
	JavaException        ReturnCode = C.ENCHRC_JavaException
	FailedToCreateJavaVM ReturnCode = C.ENCHRC_FailedToCreateJavaVM
)

// FontAttribute represents the attributes of a font.
type FontAttribute uint16

// Font attributes understood by Designer/Generate.
const (
	Bold            FontAttribute = C.ENCH_FontBold
	Italic          FontAttribute = C.ENCH_FontItalic
	EmulateBold     FontAttribute = C.ENCH_FontEmulateBold
	EmulateItalic   FontAttribute = C.ENCH_FontEmulateItalic
	EmulateTypeface FontAttribute = C.ENCH_FontEmulateTypeface
)

// DataType represents a data field type.
type DataType int32

// Data types understood by Designer/Generate.
const (
	NotSet   DataType = C.ENCH_DataNotSet
	Neutral  DataType = C.ENCH_DataNeutral
	Integer  DataType = C.ENCH_DataInteger
	Number   DataType = C.ENCH_DataNumber
	Date     DataType = C.ENCH_DataDate
	Time     DataType = C.ENCH_DataTime
	Currency DataType = C.ENCH_DataCurrency
)

// ImageFormat represents the preferred image format.
type ImageFormat int32

// Image formats understood by Designer/Generate.
const (
	BMP ImageFormat = C.ENCH_ImageBmp
	PNG ImageFormat = C.ENCH_ImagePng
	JPG ImageFormat = C.ENCH_ImageJpg
	SVG ImageFormat = C.ENCH_ImageSvg
)

// ColorSpace represents the preferred image colour space.
type ColorSpace int32

// Image colour spaces understood by Designer/Generate.
const (
	Named ColorSpace = C.ENCH_ColorNamed
	RGB   ColorSpace = C.ENCH_ColorRgb
	CMYK  ColorSpace = C.ENCH_ColorCmyk
)

// NumberFormat defines how a number should be formatted for display.
type NumberFormat struct {
	ThousandsSeparator, DecimalPoint rune
}

// Image represents the chart image requirements.
type Image struct {
	format       ImageFormat    // [In/Out] Preferred image format.
	colorSpace   ColorSpace     // [In/Out] Preferred image colour space.
	width        Twiplet        // [In] Image width in 1/100ths of a twip.
	height       Twiplet        // [In] Image height in 1/100ths of a twip.
	resolution   int32          // [In] Image resolution (DPI).
	handle       unsafe.Pointer // [Out] Reference handle (not used).
	imageDataPtr unsafe.Pointer // [Out] Physical image data.
	imageDataLen uint32         // [Out] Number of bytes in the image data.
}

type callback struct {
	p unsafe.Pointer
}

func (c callback) integer(s string) (int32, error) {
	var i C.int
	cs := C.CString(s)
	rc := C.EnchGetInteger(c.p, cs, &i)
	C.free(unsafe.Pointer(cs))
	if rc != C.ENCHRC_OK {
		err := fmt.Errorf(
			"invalid integer format '%s', error %v",
			s, ReturnCode(rc),
		)
		return 0, err
	}
	return int32(i), nil
}

func (c callback) number(s string) (float64, error) {
	var d C.double
	cs := C.CString(s)
	rc := C.EnchGetNumber(c.p, cs, &d)
	C.free(unsafe.Pointer(cs))
	if rc != C.ENCHRC_OK {
		err := fmt.Errorf(
			"invalid number format '%s', error %v",
			s, ReturnCode(rc),
		)
		return 0, err
	}
	return float64(d), nil
}

func (c callback) numberFormat() NumberFormat {
	var f C.EnchNumberFormat
	if C.EnchGetNumberFormat(c.p, &f) == 0 {
		return NumberFormat{
			ThousandsSeparator: ',',
			DecimalPoint:       '.',
		}
	}
	return NumberFormat{
		ThousandsSeparator: rune(f.chThousandsSeparator),
		DecimalPoint:       rune(f.chDecimalPoint),
	}
}

func (c callback) fontResource(guid GUID) (*FontResource, error) {
	var cfr C.EnchFontResourceUtf8
	if guid.IsZero() {
		if C.EnchGetFont(c.p, nil, &cfr) == 0 {
			return nil, errors.New("default font not found")
		}
	} else {
		if C.EnchGetFont(c.p, (*C.uchar)(&guid[0]), &cfr) == 0 {
			return nil, fmt.Errorf("font not found for guid %v", guid)
		}
	}
	fr := &FontResource{
		Typeface:   C.GoString(cfr.pszTypeface),
		PointSize:  float64(cfr.nDeciPointSize) / 10.0,
		Attributes: FontAttribute(cfr.fsFlags),
		Filename:   C.GoString(cfr.pszFileName),
	}
	C.EnchFreeFont(&cfr)
	return fr, nil
}

func (c callback) fontStyle(guid GUID) (*FontStyle, error) {
	var csr C.EnchStyleResourceUtf8
	if C.EnchGetStyle(c.p, (*C.uchar)(&guid[0]), &csr) == 0 {
		return nil, fmt.Errorf("style not found for guid %v", guid)
	}
	fs := &FontStyle{
		FontResource: &FontResource{
			Typeface:   C.GoString(csr.fontResource.pszTypeface),
			PointSize:  float64(csr.fontResource.nDeciPointSize) / 10.0,
			Attributes: FontAttribute(csr.fontResource.fsFlags),
			Filename:   C.GoString(csr.fontResource.pszFileName),
		},
		Color: Color{
			R: uint8(csr.color.red),
			G: uint8(csr.color.green),
			B: uint8(csr.color.blue),
			C: uint8(csr.color.cyan),
			M: uint8(csr.color.magenta),
			Y: uint8(csr.color.yellow),
			K: uint8(csr.color.keyBlack),
		},
		Underline: (csr.fsFlags & C.ENCH_StyleUnderline) != 0,
	}
	C.EnchFreeFont(&csr.fontResource)
	return fs, nil
}
