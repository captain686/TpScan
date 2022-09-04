package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Req struct {
	Method   string
	Host     string
	Path     string
	Data     []byte
	ProxyUrl string
	Header   map[string]string
	Redirect bool
}

func randomUa() string {
	ua := browser.Random()
	return ua
}

func HeaderMap(jsonData ...string) map[string]string {
	var headerMap = make(map[string]string)
	headerMap["User-Agent"] = randomUa()
	headerMap["Connection"] = "close"
	headerMap["Content-Type"] = "application/x-www-form-urlencoded"
	if len(jsonData) > 0 {
		for _, arg := range jsonData {
			err := json.Unmarshal([]byte(arg), &headerMap)
			if IfErr(err) {
				fmt.Println(err)
				return nil
			}
		}
	}
	return headerMap
}

func (requests Req) Requests() (*int, *[]byte, error) {

	//代理
	//proxyUrl := "http://127.0.0.1:8080"
	var client = &http.Client{Timeout: time.Second * 5}

	if requests.ProxyUrl != "" {
		proxy, _ := url.Parse(requests.ProxyUrl)
		tr := &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client = &http.Client{
			Transport: tr,
			Timeout:   time.Second * 5, //超时时间
		}
	}
	url := HostFormat(requests.Host) + PathFormat(requests.Path)
	//fmt.Println(host)
	req, err := http.NewRequest(requests.Method, url, bytes.NewReader(requests.Data))
	if IfErr(err) {
		fmt.Println("[X] Error:  " + HostFormat(requests.Host) + " 无法访问目标站点")
		return nil, nil, err
	}

	req.Header.Set("User-Agent", randomUa())
	req.Header.Set("Connection", "close")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if !requests.Redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return fmt.Errorf("\n[O] disable Redirect")
		}
	}
	resp, err := client.Do(req)
	if IfErr(err) {
		if !strings.Contains(fmt.Sprint(err), "disable Redirect") {
			fmt.Println("[X] Error:  " + HostFormat(requests.Host) + " 无法访问目标站点")
		}
		return nil, nil, err
	}
	if resp.StatusCode != 200 {
		fmt.Println("[x] Target: "+HostFormat(requests.Host)+" StatusCode: ", resp.StatusCode)
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if IfErr(err) {
		return nil, nil, err
	}
	return &resp.StatusCode, &response, nil
}
