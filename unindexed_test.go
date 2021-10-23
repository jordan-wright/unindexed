package unindexed_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/jordan-wright/unindexed"
)

// The easiest way to use unindexed is to use the Dir function, which is
// a drop-in replacement to http.Dir.
func ExampleDir() {
	http.Handle("/", http.FileServer(unindexed.Dir("./static/")))
}

func createTemporaryDirectory(t *testing.T) string {
	// Create the static directory
	dir, err := ioutil.TempDir(os.TempDir(), "unindexed-static")
	if err != nil {
		t.Fatalf("unable to create temp directory: %s", err)
	}
	// Add a secret file
	secret, err := ioutil.TempFile(dir, "secret.txt")
	if err != nil {
		t.Fatalf("unable to create secret file: %s", err)
	}
	fmt.Fprintf(secret, "This is a secret!")
	return dir
}

func tearDown(dir string, t *testing.T) {
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("unable to remove temp directory: %s", err)
	}
}

func TestNotFound(t *testing.T) {
	dir := createTemporaryDirectory(t)
	defer tearDown(dir, t)
	r := httptest.NewRequest("GET", "/doesntexist", nil)
	w := httptest.NewRecorder()
	handler := http.FileServer(unindexed.Dir(dir))
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("unexpected status received. expected %d got %d", http.StatusNotFound, w.Code)
	}
}

func TestNoIndex(t *testing.T) {
	dir := createTemporaryDirectory(t)
	defer tearDown(dir, t)
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler := http.FileServer(unindexed.Dir(dir))
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("unexpected status received. expected %d got %d", http.StatusForbidden, w.Code)
	}
}

func TestWithIndex(t *testing.T) {
	dir := createTemporaryDirectory(t)
	defer tearDown(dir, t)

	expectedBody := []byte("testing")

	err := ioutil.WriteFile(filepath.Join(dir, "index.html"), expectedBody, 0644)
	if err != nil {
		t.Fatalf("error creating index.html in temp directory: %s", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler := http.FileServer(unindexed.Dir(dir))
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
	if !bytes.Equal(expectedBody, got) {
		t.Fatalf("unexpected response received. expected %s got %s", expectedBody, got)
	}
}
