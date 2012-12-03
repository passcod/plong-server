package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"plong"
	"strings"
)

// For future operating restrictions.
//
// h = http
// w = websocket
// x = passthru
var Mode string = "hwx"

// Plong configuration.
var PlongConfig plong.Config = plong.Config{1800}

// Connection buffer size
var BufferSize int = 256

func readConfig() {
	var i interface{}
	cfg, err := ioutil.ReadFile("config.json")
	if err != nil {
		e := fmt.Sprintf("Fatal: can’t read config.json (%s).", err)
		println(e)
		panic(e)
	}

	err = json.Unmarshal(cfg, &i)
	if err != nil {
		e := fmt.Sprintf("Fatal: JSON encoding error (%s).", err)
		println(e)
		panic(e)
	}

	m := i.(map[string]interface{})
	Mode = m["mode"].(string)
	PlongConfig = plong.Config{int64(m["identity_timeout"].(float64))}
	BufferSize = int(m["buffer_size"].(float64))
}

func whatMode(s string) bool {
	return strings.Contains(Mode, s)
}
