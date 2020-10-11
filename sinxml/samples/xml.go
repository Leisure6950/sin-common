package main

import (
	"encoding/json"
	"fmt"
	"github.com/sin-z/sin-common/sinlog"
	"github.com/sin-z/sin-common/sinxml"
	"strings"
)

func main() {
	sinlog.Init(sinlog.Config{
		Path:   ".",
		Prefix: "",
		Level:  "debug",
		Rotate: "day",
	})
	xmlStr := `<?xml version="1.0" encoding="utf-8"?>
<message>
    <monitor>1</monitor>
    <subwin>2</subwin>
    <camera>3</camera>
    <command>4</command>
    <param>asdf</param>
    <direction>6</direction>
    <speedx>63</speedx>
    <speedy>63</speedy>
</message>`
	jsonStr := `{
	    "message":{
	    "monitor":0,
	    "subwin":0,
	    "camera":0,
	    "command":0,
	    "param":"",
	    "direction":0,
	    "speedx":0,
	    "speedy":0
	}
}`
	v := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &v)
	if err != nil {
		fmt.Println("unmarshal json error:", err)
		return
	}
	err = sinxml.Unmarshal(strings.NewReader(xmlStr), &v)
	if err != nil {
		fmt.Println("unmarshal xml error:", err)
		return
	}

	vDst, _ := json.Marshal(v)
	fmt.Println(string(vDst))
}
