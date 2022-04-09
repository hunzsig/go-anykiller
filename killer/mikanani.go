package killer

import (
	"fmt"
	"go-anykiller/lib"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func Mikanani(search string) {
	net := "http://mikanani.me"
	root := "download/"

	if len(search) < 4 {
		lib.Panic("Too short!")
	}

	content, er := Maho("http://mikanani.me/Home/Search?searchstr=" + url.QueryEscape(search))
	if er != nil {
		lib.Panic(er)
	}
	content = strings.Replace(content, "\r", "", -1)
	content = strings.Replace(content, "\n", "", -1)

	dir := root + search + "/"
	lib.DirCheck(dir)
	dir, _ = filepath.Abs(dir)

	var swg sync.WaitGroup
	srcs := make(map[string]string)
	reg, er := regexp.Compile(`<tr class="js-search-results-row".+>(.*?)</tr>`)
	if er == nil {
		tds := reg.FindAllString(content, -1)
		for _, v := range tds {
			reg, er = regexp.Compile(`<a href=".*?" target="_blank" class="magnet-link-wrap">(.*?)</a>.*?<a href="(.*?)\.torrent"><img`)
			if er == nil {
				torrent := reg.FindAllStringSubmatch(v, -1)
				if len(torrent) > 0 {
					for i := 0; i < len(torrent); i++ {
						t := torrent[i]
						name := lib.ClearName(t[1]) + ".torrent"
						src := net + t[2] + ".torrent"
						srcs[name] = src
					}
				}
			}
		}
	}
	i := 0
	total := len(srcs)
	for n, s := range srcs {
		no := func(name string) string {
			i = i + 1
			return "[" + strconv.Itoa(i) + "/" + strconv.Itoa(total) + "]" + name + " "
		}
		file := dir + "/" + n
		src := s
		if lib.IsFile(file) {
			fmt.Println("[skip]" + no(n))
			continue
		}
		swg.Add(1)
		go func() {
			data, errd := Maho(src)
			if errd != nil {
				fmt.Println(errd.Error())
				fmt.Println("[fail]" + no(src))
				swg.Done()
				return
			}
			errp := lib.FilePutContents(file, data, os.ModePerm)
			if errp != nil {
				fmt.Println(errp.Error())
				fmt.Println("[fail]" + no(src))
				swg.Done()
				return
			}
			fmt.Println("[OK]" + no(src))
			swg.Done()
		}()
	}
	swg.Wait()
}
