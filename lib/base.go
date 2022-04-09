package lib

import (
	"embed"
	"fmt"
	"html"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func stack() string {
	var buf [2 << 10]byte
	res := string(buf[:runtime.Stack(buf[:], true)])
	res = strings.Replace(res, "Z:/Workspace/web/hunzsig/go-anykiller", "", -1)
	return res
}

func Panic(what interface{}) {
	t := reflect.TypeOf(what)
	switch t.Kind() {
	case reflect.String:
		fmt.Println("<EStr>", what)
	case reflect.Ptr:
		fmt.Println("<EPtr>", what)
		fmt.Println("<STACK>", stack())
	default:
		fmt.Println("<KIND>", t.Kind())
	}
	os.Exit(555)
}

// IsFile is_file()
func IsFile(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// IsDir is_dir()
func IsDir(filename string) bool {
	fd, err := os.Stat(filename)
	if err != nil {
		return false
	}
	fm := fd.Mode()
	return fm.IsDir()
}

// FileSize filesize()
func FileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return 0, err
	}
	return info.Size(), nil
}

// FilePutContents file_put_contents()
func FilePutContents(filename string, data string, mode os.FileMode) error {
	return os.WriteFile(filename, []byte(data), mode)
}

// FileGetContents file_get_contents()
func FileGetContents(filename string) (string, error) {
	d, err := os.ReadFile(filename)
	return string(d), err
}

// FileGetContentsEmbeds file_get_contents()
func FileGetContentsEmbeds(f embed.FS, filename string) (string, error) {
	d, err := f.ReadFile(filename)
	return string(d), err
}

// GetModTime 获取文件(架)修改时间 返回unix时间戳
func GetModTime(path string) int64 {
	modTime := int64(0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		u := info.ModTime().Unix()
		if u > modTime {
			modTime = u
		}
		return nil
	})
	if err != nil {
		return 0
	}
	return modTime
}

// Rand rand()
// Range: [0, 2147483647]
func Rand(min, max int) int {
	if min > max {
		Panic("min: min cannot be greater than max")
	}
	// PHP: getrandmax()
	if int31 := 1<<31 - 1; max > int31 {
		Panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

// InArray in_array()
// haystack supported types: slice, array or map
func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		Panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}

func ExeRunningQty(names []string) int {
	cmd := exec.Command("cmd", "/C", "tasklist")
	output, _ := cmd.Output()
	n := strings.Index(string(output), "System")
	if n == -1 {
		return 0
	}
	qty := 0
	sd := string(output)[n:]
	fields := strings.Fields(sd)
	for _, v := range fields {
		for _, n := range names {
			if v == n {
				qty += 1
			}
		}
	}
	return qty
}

func StringsUnique(arr []string) []string {
	var res []string
	tmp := make(map[string]bool)
	for _, v := range arr {
		if v == "" || tmp[v] == true {
			continue
		}
		res = append(res, v)
		tmp[v] = true
	}
	return res
}

func Shuffle(s []string) []string {
	ns := []string{}
	for _, v := range s {
		ns = append(ns, v)
	}
	length := len(ns) - 1
	times := length
	for {
		if times < 0 {
			break
		}
		random := Rand(0, length)
		temp := ns[times]
		ns[times] = ns[random]
		ns[random] = temp
		times -= 1
	}
	return ns
}

func RandIP() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

func ClearName(name string) string {
	name = html.UnescapeString(name)
	name = strings.Trim(name, " ")
	name = strings.Replace(name, "\\", "_", -1)
	name = strings.Replace(name, "/", "_", -1)
	name = strings.Replace(name, ":", "_", -1)
	name = strings.Replace(name, "*", "_", -1)
	name = strings.Replace(name, "?", "_", -1)
	name = strings.Replace(name, "\"", "_", -1)
	name = strings.Replace(name, "'", "_", -1)
	name = strings.Replace(name, "<", "_", -1)
	name = strings.Replace(name, ">", "_", -1)
	name = strings.Replace(name, "|", "_", -1)
	return name
}
