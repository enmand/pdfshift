package pdfshift

import "encoding/json"

// conversionMessage represents the message for conversion in PDFShift
type conversionMessage map[string]interface{}

func (s conversionMessage) convert() ([]byte, error) {
	return json.Marshal(s)
}

// PDFBuilder is used when building a request to PDFShift for conversion
type PDFBuilder struct {
	message conversionMessage
}

// NewPDFBuilder returns a request builder for conversion with PDFShift
func NewPDFBuilder() *PDFBuilder {
	return &PDFBuilder{
		message: make(map[string]interface{}),
	}
}

func (rb *PDFBuilder) build() conversionMessage {
	return rb.message
}

func (rb *PDFBuilder) set(k string, s interface{}) *PDFBuilder {
	rb.message[k] = s
	return rb
}

// URL sets the PDF conversion source to a given URL for conversion
func (rb *PDFBuilder) URL(url string) *PDFBuilder {
	return rb.set("source", url)
}

// Headers sets HTTP headers to use when making a URL conversion request
func (rb *PDFBuilder) Headers(h map[string]string) *PDFBuilder {
	return rb.set("headers", h)
}

// Auth is the configuration for HTTP BASIC AUTH
type Auth struct {
	Username string
	Password string
}

func (a Auth) encode() map[string]string {
	return map[string]string{
		"username": a.Username,
		"password": a.Password,
	}
}

// Auth sets BASIC Auth headers to use when making a URL conversion request
func (rb *PDFBuilder) Auth(a Auth) *PDFBuilder {
	return rb.set("auth", a.encode())
}

// Cookie represents an HTTP Cookie
type Cookie struct {
	Name     string
	Value    string
	Secure   bool
	HTTPOnly bool
}

// Cookies represents a collection of HTTP cookies
type Cookies []Cookie

func (c Cookies) encode() []map[string]interface{} {
	cookies := []map[string]interface{}{}
	for _, cookie := range c {
		cookies = append(cookies, map[string]interface{}{
			"name":      cookie.Name,
			"value":     cookie.Value,
			"secure":    cookie.Secure,
			"http_only": cookie.HTTPOnly,
		})
	}

	return cookies
}

// Cookies set HTTP cookies to use when making a URL conversion request
// Cookies should be of the form:
func (rb *PDFBuilder) Cookies(c Cookies) {
	rb.set("cookies", c.encode())
}

// CSS loads a CSS document with the HTML for styling or style overrides
// CSS input can be either a URL to a CSS file, or a string of CSS
func (rb *PDFBuilder) CSS(css string) *PDFBuilder {
	return rb.set("css", css)
}

// Javascript provides JavaScript to execute before converting the source documment
func (rb *PDFBuilder) Javascript(src string) *PDFBuilder {
	return rb.set("javascript", src)
}

// Watermark is the options for configuring a watermark
type Watermark struct {
	Image   string
	OffsetX string
	OffsetY string
	Rotate  int
}

func (w Watermark) encode() map[string]interface{} {
	return map[string]interface{}{
		"image":    w.Image,
		"offset_x": w.OffsetX,
		"offset_y": w.OffsetY,
		"rotate":   w.Rotate,
	}
}

// Watermark adds a watermark to a converted document
func (rb *PDFBuilder) Watermark(w Watermark) *PDFBuilder {
	return rb.set("watermark", w.encode())
}

// HeaderFooter represents header and footer options
type HeaderFooter struct {
	Source  string
	Spacing string
}

func (m HeaderFooter) encode() map[string]string {
	return map[string]string{
		"source":  m.Source,
		"spacing": m.Spacing,
	}
}

// Header adds a header to the PDF conversion
func (rb *PDFBuilder) Header(m HeaderFooter) *PDFBuilder {
	return rb.set("header", m.encode())
}

// Footer adds a footer to the PDF conversion
func (rb *PDFBuilder) Footer(m HeaderFooter) *PDFBuilder {
	return rb.set("footer", m.encode())
}

// Protection is the PDF production options for PDF conversion
type Protection struct {
	Author        string
	UserPassword  string
	OwnerPassword string
	NoPrint       bool
	NoCopy        bool
	NoModify      bool
}

func (p Protection) encode() map[string]interface{} {
	return map[string]interface{}{
		"author":         p.Author,
		"user_password":  p.UserPassword,
		"owner_password": p.OwnerPassword,
		"no_print":       p.NoPrint,
		"no_copy":        p.NoCopy,
		"no_modify":      p.NoModify,
	}
}

// Protection adds various protections to the PDF
func (rb *PDFBuilder) Protection(p Protection) *PDFBuilder {
	return rb.set("protection", p.encode())
}

// Sandbox sets the PDF conversion to happen in the PDFShift sandbox
func (rb *PDFBuilder) Sandbox(enabled bool) *PDFBuilder {
	return rb.set("sandbox", enabled)
}

// Encode returns base64 encoded PDF data instead of binary data for the PDF
func (rb *PDFBuilder) Encode(enabled bool) *PDFBuilder {
	return rb.set("encode", enabled)
}

// Landscape converts the PDF with a landscape orientation
func (rb *PDFBuilder) Landscape(enabled bool) *PDFBuilder {
	return rb.set("landscape", enabled)
}

// HTML sets a PDF conversion for inline HTML
func (rb *PDFBuilder) HTML(src string) *PDFBuilder {
	return rb.set("source", src)
}

// Format sets the format for the converted PDF
// Values are: Letter, Legal, Tabloid, Ledger, A0, A1, A2, A3, A4, A5 or a {width}x{height} value
func (rb *PDFBuilder) Format(f string) *PDFBuilder {
	return rb.set("format", f)
}

// Page sets the number of pages to print for the converted PDF
func (rb *PDFBuilder) Page(f string) *PDFBuilder {
	return rb.set("page", f)
}
