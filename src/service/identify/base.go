package identify

import (
	"fmt"
	"os"
	"syscall"
)

func BaseInfo() map[string]string {
	var finfo = make(map[string]string)
	finfo = buildWithFilename(finfo)
	finfo = buildWithStatInfo(finfo)
	finfo = buildWithAccessInfo(finfo)
	finfo = SimpleIdentify(finfo)
	return finfo
}

func buildWithStatInfo(_m map[string]string) map[string]string {
	var fstat, err = os.Stat(_m["name"])

	if err != nil {
		panic(err)
	}

	var size string
	var _size int64 = fstat.Size()
	if _size >= 0 && _size <= 1024*1024 {
		size = fmt.Sprintf("%.1f", float64(_size)/1024) + "K"
	} else if _size >= 1024*1024 || _size <= 1024*1024*1024 {
		size = fmt.Sprintf("%.1f", float64(_size)/1024/1024) + "M"
	} else {
		size = fmt.Sprintf("%.1f", float64(_size)/1024/1024/1024) + "G"
	}

	_m["size"] = size
	if fstat.IsDir() == true {
		_m["type"] = "directory"
	}
	_m["mode"] = fstat.Mode().String()[1:]
	_m["modTime"] = fstat.ModTime().Format("2006-01-02 15:04:05")
	return _m
}

func buildWithAccessInfo(_m map[string]string) map[string]string {
	err := syscall.Access(_m["name"], syscall.O_RDWR)
	if err == nil {
		_m["access"] = "writable"
	} else {
		_m["access"] = "non-writable"
	}
	return _m
}

func buildWithFilename(_m map[string]string) map[string]string {
	_m["name"] = getFilename()
	return _m
}

func getFilename() string {
	var _f string
	if len(os.Args) == 2 {
		_f = os.Args[1]
		checkExist(_f)
		return _f
	} else {
		panic("demo: ./finfo filename")
	}
}

func checkExist(filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			panic("file not exist")
		} else {
			panic(err)
		}
	}
}
