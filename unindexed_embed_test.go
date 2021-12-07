package unindexed

import (
	"bytes"
	"embed"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// The easiest way to use unindexed is to use the Dir function, which is
// a drop-in replacement to http.Dir.
func ExampleEmbedDir() {
	http.Handle("/", http.FileServer(http.FileSystem(EmbedFS{FS: embed.FS{}})))
}

//go:embed test_embed
var embedTest embed.FS

func TestEmbedNotFound(t *testing.T) {
	r := httptest.NewRequest("GET", "/doesntexist", nil)
	w := httptest.NewRecorder()
	handler := http.FileServer(http.FileSystem(EmbedFS{FS: embedTest}))
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("unexpected status received. expected %d got %d", http.StatusNotFound, w.Code)
	}
}

func TestEmbedNoIndex(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler := http.FileServer(EmbedFS{FS: embedTest})
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("unexpected status received. expected %d got %d", http.StatusForbidden, w.Code)
	}
}

func TestEmbedWithIndex(t *testing.T) {
	r := httptest.NewRequest("GET", "/test_embed/unindexed-static/", nil)
	w := httptest.NewRecorder()
	handler := http.FileServer(EmbedFS{FS: embedTest})
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status received. expected %d got %d", http.StatusOK, w.Code)
	}

	resp := w.Result()
	defer resp.Body.Close()
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading from response body: %s", err)
	}

	expectedBody := []byte("testing")
	if !bytes.Equal(expectedBody, got) {
		t.Fatalf("unexpected response received. expected %s got %s", expectedBody, got)
	}
}

func TestEmbedAccessSecret(t *testing.T) {
	r := httptest.NewRequest("GET", "/test_embed/unindexed-static/secret.txt", nil)
	w := httptest.NewRecorder()
	handler := http.FileServer(EmbedFS{FS: embedTest})
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status received. expected %d got %d", http.StatusOK, w.Code)
	}

	resp := w.Result()
	defer resp.Body.Close()
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading from response body: %s", err)
	}

	expectedBody := []byte("This is a secret!")
	if !bytes.Equal(expectedBody, got) {
		t.Fatalf("unexpected response received. expected %s got %s", expectedBody, got)
	}
}
