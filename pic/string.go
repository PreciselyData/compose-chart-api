package pic

import "fmt"

var returnCodes = map[ReturnCode]string{
	OK:                   "OK",
	Failed:               "Failed",
	NotImplemented:       "NotImplemented",
	InvalidFilePath:      "InvalidFilePath",
	MissingProperty:      "MissingProperty",
	InvalidValue:         "InvalidValue",
	UnresolvedFont:       "UnresolvedFont",
	InvalidDataString:    "InvalidDataString",
	EmptyDataString:      "EmptyDataString",
	JavaException:        "JavaException",
	FailedToCreateJavaVM: "FailedToCreateJavaVM",
}

func (rc ReturnCode) String() string {
	if s, ok := returnCodes[rc]; ok {
		return s
	}
	return fmt.Sprintf("Unknown (%d)", rc)
}

var imageFormats = map[ImageFormat]string{
	BMP: "BMP",
	PNG: "PNG",
	JPG: "JPG",
	SVG: "SVG",
}

func (f ImageFormat) String() string {
	if s, ok := imageFormats[f]; ok {
		return s
	}
	return fmt.Sprintf("Unknown (%d)", f)
}

var colorSpaces = map[ColorSpace]string{
	Named: "Named",
	RGB:   "RGB",
	CMYK:  "CMYK",
}

func (cs ColorSpace) String() string {
	if s, ok := colorSpaces[cs]; ok {
		return s
	}
	return fmt.Sprintf("Unknown (%d)", cs)
}
