// Package unindexed provides an HTTP filesystem that disables directory
// indexing
//
// Motivation
//
// By default, the "http.Dir" filesystem has directory indexing enabled, which
// means that if a directory is requested that doesn't include an index.html
// file, a list of files in the directory is returned. This could leak
// sensitive information and should be avoided unless needed.
//
// Usage
//
// The easiest way to use this package is through unindexed.Dir, which is a
// drop-in replacement for http.Dir. If a directory is requested that doesn't
// have an index.html file, this package returns a http.StatusNotFound response.
package unindexed
