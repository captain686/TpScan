package util

import (
	"fmt"
	"github.com/aliyun/texpr"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var randomStr = RandomString(10)

func ListDir() *[]string {
	var files []string
	basePath, err := os.Getwd()
	if IfErr(err) {
		log.Printf("can't Get PWD")
		return nil
	}
	filePath := basePath + "/exploit/"
	fileName, err := ioutil.ReadDir(filePath)
	if err != nil {
		return nil
	}
	for _, file := range fileName {
		fileNameSplit := strings.Split(file.Name(), ".")
		extName := fileNameSplit[len(fileNameSplit)-1]
		if (strings.ToUpper(extName) == "YML") || (strings.ToUpper(extName) == "YAML") {
			files = append(files, filePath+file.Name())
		}
	}
	return &files
}

func VulExist(rule Rules, host string) (*string, *string, bool) {
	method := rule.Request.Method
	path := rule.Request.Path
	headers := rule.Request.Headers
	followRedirects := rule.Request.FollowRedirects
	data := rule.Request.Body
	proxy := rule.Request.Proxy
	for key, value := range headers {
		if strings.Contains(value, "randomStr") {
			headers[key] = strings.Replace(value, "randomStr", randomStr, -1)
		}
	}
	if strings.Contains(path, "randomStr") {
		path = strings.Replace(path, "randomStr", randomStr, -1)
	}

	if strings.Contains(data, "randomStr") {
		data = strings.Replace(data, "randomStr", randomStr, -1)
	}
	var req Req
	if followRedirects == "" {
		req = Req{Method: method, Host: host, Path: path, Data: []byte(data), Header: headers, Redirect: true, ProxyUrl: proxy}
	} else {
		req = Req{Method: method, Host: host, Path: path, Data: []byte(data), Header: headers, ProxyUrl: proxy}
	}
	statusCode, response, err := req.Requests()
	if IfErr(err) {
		return nil, nil, false
	}
	//fmt.Println(string(*response))
	responseExpression := rule.Expression
	result := responseExpression["result"]
	responseStatus := responseExpression["response_status"]
	//fmt.Println(responseStatus)
	//fmt.Println(result)
	var res bool
	rulesInResponse := strings.TrimRight(responseExpression["inResponse"], "\n")
	rulesInResponse = strings.Replace(rulesInResponse, "randomStr", randomStr, -1)
	if result == "and" {
		responseStatus, err := strconv.Atoi(responseStatus)
		if err != nil {
			return nil, nil, false
		}
		res = *statusCode == responseStatus && InResponse(response, rulesInResponse)
	} else if result == "or" {
		responseStatus, err := strconv.Atoi(responseStatus)
		if err != nil {
			return nil, nil, false
		}
		res = *statusCode == responseStatus || InResponse(response, rulesInResponse)
	} else if rulesInResponse == "" {
		responseStatus, err := strconv.Atoi(responseStatus)
		if err != nil {
			return nil, nil, false
		}
		res = *statusCode == responseStatus
	} else if responseStatus == "" {
		res = InResponse(response, rulesInResponse)
	}
	vulPath := host + path
	return &vulPath, &data, res
}

func ExpExecutor(host string) {
	for _, file := range *ListDir() {
		exp, err := ReadExp(file)
		if IfErr(err) {
			return
		}
		rules := exp.Rules
		VulResults := make(map[string]bool)
		var vulResults bool
		var (
			vulPath *string
			vulData *string
		)
		for ruleName, rule := range rules {
			vulPath, vulData, vulResults = VulExist(rule, host)
			VulResults[ruleName] = vulResults
		}
		expression := exp.Expression
		var newExpression string
		for key, value := range VulResults {
			newExpression = strings.Replace(expression, key, strconv.FormatBool(value), -1)
		}
		compile, err := texpr.Compile(newExpression)
		if IfErr(err) {
			fmt.Println(err)
			return
		}
		result, err := compile.Eval(nil)
		if IfErr(err) {
			fmt.Println(err)
			return
		}
		//fmt.Println(result)
		if result == true {
			fmt.Println("[X] " + host + "存在漏洞")
			fmt.Println("[X] 漏洞路径" + *vulPath)
			if *vulData != "" {
				fmt.Println("[X] 请求数据为 " + *vulData)
			}
		} else {
			fmt.Println("[O] " + host + "不存在漏洞")
		}
	}

}
