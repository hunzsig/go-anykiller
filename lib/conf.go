package lib

import (
	"os"
	"regexp"
	"strings"
)

type Conf struct {
	init  bool
	Port  string
	Proxy string
}

var (
	data Conf
)

func GetConf() Conf {
	if data.init == true {
		return data
	}
	if !IsFile("./conf.ini") {
		return data
	}
	c, err := os.ReadFile("./conf.ini")
	if err != nil {
		Panic(err)
	}
	content := string(c)
	reg, _ := regexp.Compile("#(.*)")
	content = reg.ReplaceAllString(content, "")
	content = strings.Replace(content, "\r\n", "\n", -1)
	content = strings.Replace(content, "\r", "\n", -1)
	split := strings.Split(content, "\n")
	conf := make(map[string]string)
	for _, iniItem := range split {
		if len(iniItem) > 0 {
			itemSplit := strings.Split(iniItem, "=")
			itemKey := strings.Trim(itemSplit[0], " ")
			itemKey = strings.ToLower(strings.Trim(itemSplit[0], " "))
			conf[itemKey] = strings.Trim(itemSplit[1], " ")
		}
	}
	if conf["port"] != "" {
		data.Port = conf["port"]
	} else {
		data.Port = "9999"
	}
	if conf["proxy"] != "" {
		data.Proxy = conf["proxy"]
	}
	return data
}
