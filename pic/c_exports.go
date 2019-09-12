package pic

// #include <stdlib.h>
import "C"

import (
	"log"
	"unsafe"
)

// EnchCreateImage is called by Designer/Generate to create a chart image
// from a list of configuration settings (propertiesPtr). These settings
// appear in the form name=value where value may be a constant, or refer
// to a field in the symbol table (symbolsPtr). Image information is conveyed
// via imagePtr; a C pointer represented by the Image struct. The callbackPtr
// parameter contains pointers to functions defined by Designer/Generate to
// resolve data values, fonts and locale information.
//export EnchCreateImage
func EnchCreateImage(
	callbackPtr unsafe.Pointer,
	propertiesPtr, symbolsPtr *C.char,
	imagePtr unsafe.Pointer,
) (rc ReturnCode) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Unexpected failure:", r)
			rc = Failed
		}
	}()

	if client == nil {
		log.Println("No implementation defined")
		return NotImplemented
	}

	props := C.GoString(propertiesPtr)
	syms := C.GoString(symbolsPtr)
	img := (*Image)(imagePtr)

	if options.LogInfo() {
		log.Printf(
			"INFO: Creating image: width=%v, height=%v, DPI=%v, format=%v, colorspace=%v\n"+
				"[Properties]\n%s\n"+
				"[Symbols]\n%s\n",
			img.width, img.height, img.resolution, img.format, img.colorSpace,
			props, syms,
		)
	}

	if img.width == 0 || img.height == 0 {
		log.Println("Zero dimensions supplied")
		return InvalidValue
	}

	config := newConfig(callback{callbackPtr}, props, syms)
	builder := client.NewBuilder(config)
	if builder == nil {
		log.Println("Configuration not supported")
		return NotImplemented
	}

	builder.SetFormat(&img.format, &img.colorSpace)
	builder.SetSize(img.width, img.height, img.resolution)

	buf, err := builder.Render()
	if err != nil {
		log.Println("Error rendering chart:", err)
		return Failed
	}

	img.imageDataPtr = C.CBytes(buf.Bytes())
	img.imageDataLen = uint32(buf.Len())

	if options.LogInfo() {
		log.Printf(
			"INFO: Created image: size=%d, format=%v, colorspace=%v\n",
			img.imageDataLen, img.format, img.colorSpace,
		)
	}
	return OK
}

// EnchDestroyImage is called by Designer/Generate to destroy the chart image
// data created by EnchCreateImage.
//export EnchDestroyImage
func EnchDestroyImage(imagePtr unsafe.Pointer) ReturnCode {
	img := (*Image)(imagePtr)
	if options.LogInfo() {
		log.Printf(
			"INFO: Destroying image: size=%d, format=%v, colorspace=%v\n",
			img.imageDataLen, img.format, img.colorSpace,
		)
	}
	C.free(img.imageDataPtr)
	return OK
}

// EnchTerminate is called by Generate to tidy up before the program exits.
// The Failed ReturnCode is returned here to indicate that it is unsafe for
// Generate to unload this module as the Go runtime may still need to do
// some garbage collection.
//export EnchTerminate
func EnchTerminate() ReturnCode {
	if options.LogInfo() {
		log.Println("INFO: Terminating")
	}
	return Failed
}
