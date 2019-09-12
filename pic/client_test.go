package pic

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

type mockClient struct{}
type mockBuilder struct{}

func (mockClient) NewBuilder(c *Config) Builder {
	return mockBuilder{}
}

func (mockBuilder) SetFormat(format *ImageFormat, colorSpace *ColorSpace) {
}

func (mockBuilder) SetSize(width, height Twiplet, dpi int32) {
}

func (mockBuilder) Render() (*bytes.Buffer, error) {
	return &bytes.Buffer{}, nil
}

func TestSetClient(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file")
	}
	tmpfile.Close()

	defer func() {
		log.SetOutput(ioutil.Discard)
		os.Remove(tmpfile.Name())
	}()

	mc := mockClient{}
	SetClient(
		mc,
		Options{
			LogLevel:    LogInfo,
			LogFileName: tmpfile.Name(),
		},
	)

	assertEqual(t, mc, client)
	assertEqual(t, options.LogInfo(), true)
}
