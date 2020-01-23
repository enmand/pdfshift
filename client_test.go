package pdfshift

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_New(t *testing.T) {
	p := New("test key")

	if p == nil {
		t.Error("invalid PDFShift client")
	}
}

func Test_Convert(t *testing.T) {
	badBuilder := NewPDFBuilder()
	badBuilder.message["s"] = func() {}

	type Error struct {
		Error string `json:"error"`
	}
	badResp, _ := json.Marshal(struct {
		Error Error `json:"error"`
	}{
		Error: Error{
			Error: "bad request",
		},
	})

	type server struct {
		statusCode int
		resp       []byte
	}
	type expectation struct {
		want []byte
		err  error
	}
	tcs := map[string]struct {
		builder  *PDFBuilder
		server   server
		expected expectation
	}{
		"ok": {},
		"bad message": {
			builder: badBuilder,
			expected: expectation{
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
				err: fmt.Errorf("unable to decode error response: EOF"),
			},
		},
		"server error": {
			server: server{
				statusCode: 500,
				resp:       badResp,
			},
			expected: expectation{
				err: fmt.Errorf("bad request"),
			},
		},
	}
	for _, tc := range tcs {
		if tc.builder == nil {
			tc.builder = NewPDFBuilder()
		}
		c := New("test key")
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
