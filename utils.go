package tests

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
)

func makeAPIRequest(baseURL string, uri string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("%s/%s", baseURL, uri))

	resp := fasthttp.AcquireResponse()
	fasthttp.Do(req, resp)

	return resp
}

func readDataFile(filename string) []string{
	data, err := os.ReadFile("testdata/" + filename + ".txt")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("File does not exist")  
		} else if errors.Is(err, os.ErrPermission) {
			fmt.Println("Permission denied")
		} else {
			fmt.Printf("Unhandled error %v occurred\n", err)
			panic(err)
		}
	}
	return strings.Split(string(data), "\n")
}