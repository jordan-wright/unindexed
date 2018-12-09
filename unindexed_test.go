package unindexed_test

import (
	"net/http"

	"github.com/jordan-wright/unindexed"
)

// The easiest way to use unindexed is to use the Dir function, which is
// a drop-in replacement to http.Dir.
func ExampleDir() {
	http.Handle("/", http.FileServer(unindexed.Dir("./static/")))
}
