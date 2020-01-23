package pdfshift

import (
	"testing"
)

func TestPDFBuilder_conversionMessage(t *testing.T) {
	m := conversionMessage{}
	m["key"] = "value"
	expect := "{\"key\":\"value\"}"

	out, err := m.convert()
	if err != nil {
		t.Errorf("conversion failed, err: %s", err)
	}

	if string(out) != expect {
		t.Error("conversion value unexpected")
	}
}

func Test_NewPDFBuilder(t *testing.T) {
	p := NewPDFBuilder()

	if p == nil {
		t.Error("no PDF builder given")
	}
}

func TestPDFBuilder_build(t *testing.T) {
	p := NewPDFBuilder()
	p.set("key", "value")

	b := p.build()

	if b["key"] != "value" {
		t.Error("invalid PDF message")
	}
}

func TestPDFBuilder_URL(t *testing.T) {
	p := NewPDFBuilder()
	p.URL("google.ca")
	b := p.build()

	if b["source"] != "google.ca" {
		t.Error("invalid PDF message")
	}
}

func TestPDFBuilder_Header(t *testing.T) {
	p := NewPDFBuilder()
	headers := map[string]string{
		"Content-Type": "application/test",
	}
	p.Headers(headers)
	b := p.build()

	if b["headers"].(map[string]string)["Content-Type"] != "application/test" {
		t.Error("invalid headers")
	}
}

func TestPDFBuilder_Auth(t *testing.T) {
	p := NewPDFBuilder()
	a := Auth{
		Username: "username",
		Password: "password",
	}
	p.Auth(a)
	b := p.build()

	if b["auth"].(map[string]string)["username"] != a.Username {
		t.Error("invalid auth username")
	}
	if b["auth"].(map[string]string)["password"] != a.Password {
		t.Error("invalid auth password")
	}
}

func TestPDFBuilder_Cookies(t *testing.T) {
	p := NewPDFBuilder()
	c := Cookies{
		{
			Name:  "cookie",
			Value: "secret",
		},
	}
	p.Cookies(c)
	b := p.build()

	if b["cookies"].([]map[string]interface{})[0]["name"] != "cookie" {
		t.Error("invalid auth username")
	}
	if b["cookies"].([]map[string]interface{})[0]["value"] != "secret" {
		t.Error("invalid auth password")
	}
}

func TestPDFBuilder_CSS(t *testing.T) {
	p := NewPDFBuilder()
	css := "a { font-color: red; }"
	p.CSS(css)
	b := p.build()

	if b["css"] != css {
		t.Error("invalid PDF message")
	}
}

func TestPDFBuilder_Javascript(t *testing.T) {
	p := NewPDFBuilder()
	js := "alert(\"hello\")"
	p.Javascript(js)
	b := p.build()

	if b["javascript"] != js {
		t.Error("invalid PDF message")
	}
}

func TestPDFBuilder_Watermark(t *testing.T) {
	p := NewPDFBuilder()
	w := Watermark{
		Image:   "manifold.co/logo.png",
		OffsetX: "50px",
		OffsetY: "45px",
		Rotate:  45,
	}
	p.Watermark(w)
	b := p.build()
	bw := b["watermark"].(map[string]interface{})
	if bw["image"] != w.Image {
		t.Error("invalid watermark image")
	}

	if bw["offset_x"] != w.OffsetX {
		t.Error("invalid offset_x")
	}

	if bw["offset_y"] != w.OffsetY {
		t.Error("invalid offset_y")
	}

	if bw["rotate"] != w.Rotate {
		t.Error("invalid rotate")
	}
}

func TestPDFBuilder_HeaderFooter(t *testing.T) {
	p := NewPDFBuilder()
	hf := HeaderFooter{
		Source:  "document!",
		Spacing: "100px",
	}
	p.Header(hf)
	b := p.build()
	bh := b["header"].(map[string]string)

	if bh["source"] != hf.Source {
		t.Error("bad header source")
	}
	if bh["spacing"] != hf.Spacing {
		t.Error("bad header spacing")
	}

	p.Footer(hf)
	b = p.build()
	bf := b["header"].(map[string]string)

	if bf["source"] != hf.Source {
		t.Error("bad header source")
	}
	if bf["spacing"] != hf.Spacing {
		t.Error("bad header spacing")
	}
}

func TestPDFBuilder_Protection(t *testing.T) {
	p := NewPDFBuilder()
	hf := Protection{
		Author:        "author",
		UserPassword:  "user_password",
		OwnerPassword: "owner_password",
		NoPrint:       true,
		NoCopy:        true,
		NoModify:      true,
	}
	p.Protection(hf)
	b := p.build()
	ptc := b["protection"].(map[string]interface{})

	if ptc["author"] != hf.Author {
		t.Error("bad author")
	}
	if ptc["user_password"] != hf.UserPassword {
		t.Error("bad user_password")
	}
	if ptc["owner_password"] != hf.OwnerPassword {
		t.Error("bad owner_password")
	}
	if ptc["no_print"].(bool) != hf.NoPrint {
		t.Error("bad no_print")
	}
	if ptc["no_copy"].(bool) != hf.NoCopy {
		t.Error("bad no_copy")
	}
	if ptc["no_modify"].(bool) != hf.NoModify {
		t.Error("bad no_modify")
	}
}

func TestPDFBuilder_Sandbox(t *testing.T) {
	p := NewPDFBuilder()
	p.Sandbox(true)
	b := p.build()

	if b["sandbox"].(bool) != true {
		t.Error("invalid sandbox request")
	}
}

func TestPDFBuilder_Encode(t *testing.T) {
	p := NewPDFBuilder()
	p.Encode(true)
	b := p.build()

	if b["encode"].(bool) != true {
		t.Error("invalid encode request")
	}
}

func TestPDFBuilder_Landscape(t *testing.T) {
	p := NewPDFBuilder()
	p.Landscape(true)
	b := p.build()

	if b["landscape"].(bool) != true {
		t.Error("invalid encode request")
	}
}

func TestPDFBuilder_HTML(t *testing.T) {
	p := NewPDFBuilder()
	html := "<blink>important!</blink>"
	p.HTML(html)
	b := p.build()

	if b["source"] != html {
		t.Error("invalid html source request")
	}
}

func TestPDFBuilder_Format(t *testing.T) {
	p := NewPDFBuilder()
	p.Format("A4")
	b := p.build()

	if b["format"] != "A4" {
		t.Error("invalid format request")
	}
}

func TestPDFBuilder_Page(t *testing.T) {
	p := NewPDFBuilder()
	p.Page("1,2,3")
	b := p.build()

	if b["page"] != "1,2,3" {
		t.Error("invalid format request")
	}
}
