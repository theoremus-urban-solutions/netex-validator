package utils

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOptimizedHTTPClientGetGzipSchema(t *testing.T) {
	const schema = `<xsd:schema xmlns:xsd="http://www.w3.org/2001/XMLSchema"/>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			t.Error("request did not advertise gzip support")
			http.Error(w, "gzip expected", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Content-Encoding", "gzip")
		z := gzip.NewWriter(w)
		if _, err := io.WriteString(z, schema); err != nil {
			t.Error(err)
		}
		if err := z.Close(); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(server.Close)

	client := NewOptimizedHTTPClient(nil)
	t.Cleanup(client.Close)
	resp, err := client.Get(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Error(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != schema {
		t.Errorf("body = %q, want %q", body, schema)
	}
}
