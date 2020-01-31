package pdfshift

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_New(t *testing.T) {
	p := New("test key", true)

	if p == nil {
		t.Error("invalid PDFShift client")
	}
}

func Test_Convert(t *testing.T) {
	badBuilder := NewPDFBuilder()
	badBuilder.message["s"] = func() {}

	badResp, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: "bad request",
	})

	type fields struct {
		sandbox bool
	}
	type server struct {
		statusCode int
		resp       []byte
	}
	type expectation struct {
		want    []byte
		reqBody string
		err     error
	}
	tcs := map[string]struct {
		fields   fields
		builder  *PDFBuilder
		server   server
		expected expectation
	}{
		"ok": {
			expected: expectation{
				reqBody: `{"sandbox":false}`,
			},
		},
		"ok - sandbox": {
			fields: fields{
				sandbox: true,
			},
			expected: expectation{
				reqBody: `{"sandbox":true}`,
			},
		},
		"bad message": {
			builder: badBuilder,
			expected: expectation{
				reqBody: `{"sandbox":false}`,
				err: fmt.Errorf(
					"unable to marshal conversion message: json: unsupported type: func()",
				),
			},
		},
		"bad response": {
			server: server{
				statusCode: 401,
			},
			expected: expectation{
				reqBody: `{"sandbox":false}`,
				err:     fmt.Errorf("unable to decode error response: EOF"),
			},
		},
		"server error": {
			server: server{
				statusCode: 500,
				resp:       badResp,
			},
			expected: expectation{
				reqBody: `{"sandbox":false}`,
				err:     fmt.Errorf("bad request"),
			},
		},
	}
	for _, tc := range tcs {
		if tc.builder == nil {
			tc.builder = NewPDFBuilder()
		}
		c := New("test key", tc.fields.sandbox)
		ctx := context.Background()

		s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			status := 200
			if tc.server.statusCode != 0 {
				status = tc.server.statusCode
			}
			rw.WriteHeader(status)
			if tc.server.resp != nil {
				rw.Write(tc.server.resp)
			}

			bs, _ := ioutil.ReadAll(req.Body)
			if string(bs) != tc.expected.reqBody {
				t.Errorf("request bodies do not match: want %s, got %s", tc.expected.reqBody, bs)
			}
		}))
		defer s.Close()
		c.url = s.URL
		c.client = s.Client()

		out, err := c.Convert(ctx, tc.builder)
		if err != nil {
			if tc.expected.err == nil {
				t.Errorf("unexpected error: %s", err)
				return
			}

			if err.Error() != tc.expected.err.Error() {
				t.Errorf("unmatched error: %s", err)
			}
		}

		if !bytes.Equal(out, tc.expected.want) {
			t.Error("output does not match expected value")
		}
	}
}
