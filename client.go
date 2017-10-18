package gosseract

// #cgo LDFLAGS: -llept -ltesseract
// #include <stdlib.h>
// #include "tessbridge.h"
import "C"
import (
	"fmt"
	"strings"
)

// Version returns the version of Tesseract-OCR
func Version() string {
	api := C.Create()
	defer C.Free(api)
	version := C.Version(api)
	return C.GoString(version)
}

// Client is argument builder for tesseract::TessBaseAPI.
type Client struct {
	api            C.TessBaseAPI
	Trim           bool
	TessdataPrefix *string
	Languages      []string
	ImagePath      string
}

// NewClient construct new Client. It's due to caller to Close this client.
func NewClient() *Client {
	client := &Client{
		api: C.Create(),
	}
	return client
}

// Close frees allocated API.
func (c *Client) Close() (err error) {
	// defer func() {
	// 	if e := recover(); e != nil {
	// 		err = fmt.Errorf("%v", e)
	// 	}
	// }()
	C.Free(c.api)
	return err
}

// SetImage sets image to execute OCR.
func (c *Client) SetImage(imagepath string) *Client {
	c.ImagePath = imagepath
	return c
}

// Text finally initalize tesseract::TessBaseAPI, execute OCR and extract text detected as string.
func (c *Client) Text() (string, error) {

	// Defer recover and make error
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	// Initialize tesseract::TessBaseAPI
	if len(c.Languages) == 0 {
		C.Init(c.api, nil, nil)
	} else {
		langs := strings.Join(c.Languages, "+")
		C.Init(c.api, nil, C.CString(langs))
	}

	// Set Image by giving path
	C.SetImage(c.api, C.CString(c.ImagePath))

	// Get text by execuitng
	out := C.GoString(C.UTF8Text(c.api))

	// Trim result if needed
	if c.Trim {
		out = strings.Trim(out, "\n")
	}

	return out, err
}
