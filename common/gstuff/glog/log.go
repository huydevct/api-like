package glog

import (
	"encoding/json"

	"github.com/kr/pretty"
)

// Pretty func
func Pretty(data ...interface{}) {
	pretty.Println(data)
}

// Send ..
func Send(args map[string]interface{}) {

	data := map[string]interface{}{
		"log":  args,
		"keep": "true",
	}

	jsonData, _ := json.Marshal(data)
	logFile.Write(string(jsonData))
}
