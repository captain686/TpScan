package util

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
)

// RandomString 与Docker ContainerID 生成方法相同
func RandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func IfErr(err error) bool {
	if err != nil {
		return true
	}
	return false
}

func HostFormat(host string) string {
	subTime := strings.Count(host, "/")
	if subTime == 2 {
		host += "/"
	} else {
		re := regexp.MustCompile("http*.://.*?/")
		host = re.FindString(host)
	}
	return host
}

func PathFormat(path string) string {
	if strings.HasPrefix(path, "/") {
		return strings.TrimPrefix(path, "/")
	}
	return path
}

func FindinResponse(html *[]byte, pattern string) string {
	strHtml := string(*html)
	re := regexp.MustCompile(pattern)
	result := re.FindString(strHtml)
	if result != "" {
		return result
	}
	return ""
}
func InResponse(html *[]byte, pattern string) bool {
	if FindinResponse(html, pattern) != "" {
		return true
	}
	return false
}

//func VulExist(result string, statusCode string, inResponse string) bool {
//	if result == "and" {
//		return
//	}
//}
