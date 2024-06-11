package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestGzip(t *testing.T) {

	bb := []byte("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTYyNTk0NzQsImZpbGVJZCI6ImFkZTEwM2Q4NzVlZjQ5YWNhMjc4MDUwNjE3MWEyM2RiODE3OTQ0OTU2IiwiY291bnQiOjMsImV4dHJhIjoie1wiT3JpZ2luXCI6W1wiaHR0cHM6Ly8xOTIuMTY4LjgwLjExNToxMDAwMlwiXSxcIlNlYy1DaC1VYS1QbGF0Zm9ybVwiOltcIlxcXCJXaW5kb3dzXFxcIlwiXSxcIlVzZXItQWdlbnRcIjpbXCJNb3ppbGxhLzUuMCAoV2luZG93cyBOVCAxMC4wOyBXaW42NDsgeDY0KSBBcHBsZVdlYktpdC81MzcuMzYgKEtIVE1MLCBsaWtlIEdlY2tvKSBDaHJvbWUvMTI0LjAuMC4wIFNhZmFyaS81MzcuMzZcIl0sXCJTZWMtRmV0Y2gtRGVzdFwiOltcImVtcHR5XCJdLFwiQWNjZXB0LUVuY29kaW5nXCI6W1wiZ3ppcCwgZGVmbGF0ZSwgYnIsIHpzdGRcIl0sXCJTZWMtRmV0Y2gtTW9kZVwiOltcImNvcnNcIl0sXCJBdXRob3JpemF0aW9uXCI6W1wiQmVhcmVyIGV5SmhiR2NpT2lKSVV6STFOaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpwY0NJNklqRTVNaTR4TmpndU1TNHlNek1pTENKbWFYaGxaQ0k2SWlJc0ltVjRjQ0k2TVRjeE5qSTJNamc0T0gwLjZ3NU1xWml2bWM0bW91V3BsWDFqV25zcnN4aEpud3RnYWJrM0lyQ0FfcmNcIl0sXCJ1c2VyUmlnaHRcIjowLFwiZmxvb3JFZGl0VXJsUHJlZml4XCI6XCJodHRwczovLzE5Mi4xNjguODAuMTE1OjEwMDAyL3MvXCIsXCJDb250ZW50LUxlbmd0aFwiOltcIjBcIl0sXCJtb2JpbGVGbGFnXCI6ZmFsc2UsXCJTZWMtQ2gtVWFcIjpbXCJcXFwiQ2hyb21pdW1cXFwiO3Y9XFxcIjEyNFxcXCIsIFxcXCJHb29nbGUgQ2hyb21lXFxcIjt2PVxcXCIxMjRcXFwiLCBcXFwiTm90LUEuQnJhbmRcXFwiO3Y9XFxcIjk5XFxcIlwiXSxcIlgtUmVxdWVzdGVkLVdpdGhcIjpbXCJYTUxIdHRwUmVxdWVzdFwiXSxcIkFjY2VwdFwiOltcImFwcGxpY2F0aW9uL2pzb24sIHRleHQvcGxhaW4sICovKlwiXSxcIlByaW9yaXR5XCI6W1widT0xLCBpXCJdLFwiWW96by1BZ2VudFwiOltcInlvem9zb2Z0LmNvbVwiXSxcImZpbGVQYXRoXCI6XCIv5Yqg5a-GLTEyMy54bHN4XCIsXCJSZWZlcmVyXCI6W1wiaHR0cHM6Ly8xOTIuMTY4LjgwLjExNToxMDAwMi9kZW1vL2QvXCJdLFwiWC1Gb3J3YXJkZWQtUHJvdG9cIjpbXCJodHRwc1wiXSxcIlNlYy1GZXRjaC1TaXRlXCI6W1wic2FtZS1vcmlnaW5cIl0sXCJ1c2VySWRcIjpcIjE5Mi4xNjguMS4yMzNcIixcImZhbGxiYWNrVXJsXCI6XCJodHRwczovLzE5Mi4xNjguODAuMTE1OjEwMDAyL2RlbW8vZC8jL2luZGV4XCIsXCJzYXZlRmxhZ1wiOnRydWUsXCJYLUZvcndhcmRlZC1Gb3JcIjpbXCIxOTIuMTY4LjEuMjMzXCJdLFwiQWNjZXB0LUxhbmd1YWdlXCI6W1wiemgtQ04semg7cT0wLjlcIl0sXCJTZWMtQ2gtVWEtTW9iaWxlXCI6W1wiPzBcIl0sXCJYLVJlYWwtSXBcIjpbXCIxOTIuMTY4LjEuMjMzXCJdLFwiZmlsZUlkXCI6XCJhZGUxMDNkODc1ZWY0OWFjYTI3ODA1MDYxNzFhMjNkYjgxNzk0NDk1NlwifSJ9.OapYK2x8sQrLXmyIEyzKiVr6X1jJTATEVqu3avd0FA6-BAhFwRavn17Skp32bIZQpYdujyi8IdYs7m-F0JVMSw")
	str := writer(bb)
	fmt.Println("压缩后的字符串：", str)
	str2 := writer([]byte(str))
	fmt.Println("压缩后的字符串：", str2)
	str3 := writer([]byte(str2))
	fmt.Println("压缩后的字符串：", str3)
	str4 := writer([]byte(str3))
	fmt.Println("压缩后的字符串：", str4)
	/*data, _ := base64.StdEncoding.DecodeString(str)
	rdata := bytes.NewReader(data)
	r, _ := gzip.NewReader(rdata)
	s, _ := ioutil.ReadAll(r)
	fmt.Println("解压后的字符串：", string(s))*/
	decodeString, err := base64.URLEncoding.DecodeString(string("eyJtZXRob2QiOjMsImRhdGEiOnsidGV4dCI6IkVkaXRvcklkKDE5Mi4xNjguMS4yMzMkMCkgaXMgY29uc3VtZWQgYnkgb2ZmaWNlKDE5Mi4xNjguMS4yMzM6MTkwMDApLiJ9fQ=="))
	if err != nil {
		return
	}
	println(string(decodeString))

}

func writer(bb []byte) string {
	var b bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&b, 9)
	if _, err := gz.Write(bb); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}

	str := base64.StdEncoding.EncodeToString(b.Bytes())
	return str
}
