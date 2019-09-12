package pic

import "bytes"

// Builder creates a chart image.
type Builder interface {
	SetFormat(format *ImageFormat, colorSpace *ColorSpace)
	SetSize(width, height Twiplet, dpi int32)
	Render() (*bytes.Buffer, error)
}
