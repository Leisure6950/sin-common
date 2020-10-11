package main

import "github.com/sin-z/sin-common/sinhttp"

func main() {
	sinhttp.Client(nil, "").Post("", nil).Do()

	sinhttp.Client(nil, "").Get("").Do()

}
