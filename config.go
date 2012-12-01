package main

import "plong"

// For future operating restrictions.
//
// r = read
// w = write
// x = exchange
//
// E.g. rw = no exchange, r = read-only, etc
var Mode string = "rwx"

// Plong configuration.
var PlongConfig plong.Config = plong.Config{1800}

// Connection buffer size
var BufferSize int = 256
