package goplex

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func NewReadCloser(body string) *ReadCloser {
	return &ReadCloser{
		Reader: strings.NewReader(body),
	}
}

type ReadCloser struct {
	io.Reader
}

func (f ReadCloser) Close() error {
	return nil
}

type fakeRoundTripper struct {
	t        *testing.T
	resp     *http.Response
	err      error
	expected *http.Request
}

func (f fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if !reflect.DeepEqual(f.expected, req) {
		f.t.Errorf("Unexpected Request:\nExpected: %+v\nGot: %+v\n", f.expected, req)
	}
	return f.resp, f.err
}

func makeFakeClient(t *testing.T, statusCode int, resp string, expected *http.Request) *http.Client {
	r := &http.Response{
		StatusCode: statusCode,
		Body:       NewReadCloser(resp),
	}

	client := &http.Client{
		Transport: fakeRoundTripper{
			t:        t,
			resp:     r,
			expected: expected,
		},
	}

	return client
}
