package killer

import (
	"fmt"
	"go-anykiller/lib"
	"golang.org/x/net/proxy"
	"html"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func Maho(target string) (string, error) {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: nil,
	})
	trans := &http.Transport{}
	//proxy
	proxyAddr := lib.GetConf().Proxy
	if proxyAddr != "" {
		if proxyAddr[0:4] == "http" {
			_, errPX := http.Get(proxyAddr)
			if errPX != nil {
			} else {
				proxyURL, _ := url.Parse(proxyAddr)
				trans.Proxy = http.ProxyURL(proxyURL)
			}
		} else if proxyAddr[0:8] == "sock5://" {
			proxyAddr = proxyAddr[8:]
			fmt.Println("proxyAddr", proxyAddr)
			dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
			if err != nil {
			} else {
				trans.Dial = dialer.Dial
			}
		}
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Jar:       jar,
		Transport: trans,
	}
	request, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return "", err
	}
	if strings.Index(target, "e-hentai.org") != -1 {
		cookie := &http.Cookie{Name: "nw", Value: "1"}
		request.AddCookie(cookie)
	}
	ip := lib.RandIP()
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	request.Header.Add("X-FORWARDED-FOR", ip)
	request.Header.Add("CLIENT-IP", ip)

	resp, err2 := client.Do(request)
	if err2 != nil {
		return "", err2
	}
	defer resp.Body.Close()
	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return "", err3
	}
	content := string(body)
	content = html.UnescapeString(content)
	return content, nil
}
