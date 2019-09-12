package pic

import (
	"fmt"
	"testing"
)

func TestReturnCodeString(t *testing.T) {
	assertEqual(t, fmt.Sprintf("%v", OK), "OK")
	assertEqual(t, fmt.Sprintf("%v", FailedToCreateJavaVM), "FailedToCreateJavaVM")
	assertEqual(t, fmt.Sprintf("%v", ReturnCode(999)), "Unknown (999)")
}

func TestImageFormatString(t *testing.T) {
	assertEqual(t, fmt.Sprintf("%v", BMP), "BMP")
	assertEqual(t, fmt.Sprintf("%v", SVG), "SVG")
	assertEqual(t, fmt.Sprintf("%v", ImageFormat(999)), "Unknown (999)")
}

func TestColorSpaceString(t *testing.T) {
	assertEqual(t, fmt.Sprintf("%v", Named), "Named")
	assertEqual(t, fmt.Sprintf("%v", CMYK), "CMYK")
	assertEqual(t, fmt.Sprintf("%v", ColorSpace(999)), "Unknown (999)")
}
