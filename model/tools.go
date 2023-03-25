package model

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/oswaldoooo/octools/toolsbox"
)

func generateid(langth int) (res string) {
	res = ""
	for i := 0; i < langth; i++ {
		rand.Seed(time.Now().UnixNano())
		res += string(byte(rand.Intn(26) + 97))
	}
	return
}
func ParseDbContent(content []byte) (restr []string, err error) {
	strcont := string(content)
	firarr := strings.Split(strcont, "\n")
	restr = BinarySearch(strcont, firarr)
	if len(restr) == 0 {
		err = errors.New(strcont + " dont existed")
	}
	return
}

func BinarySearch(content string, origin_arr []string) (resarr []string) {
	resarr = []string{}
	if id, ok := compareid(content, origin_arr[0]); ok {
		resarr = append(resarr, id)
	}
	resarr = append(resarr, binarysearch(content, origin_arr, 0)...)
	return
}
func binarysearch(content string, origin_arr []string, pos int) (resarr []string) {
	resarr = []string{}
	//search left child if exist
	if 2*pos+1 < len(origin_arr) {
		if id, ok := compareid(origin_arr[2*pos+1], content); ok {
			resarr = append(resarr, id)
		}
	} else {
		return
	}
	//search right child if exist
	if 2*pos+2 < len(origin_arr) {
		if id, ok := compareid(origin_arr[2*pos+2], content); ok {
			resarr = append(resarr, id)
		}
	} else {
		return
	}
	leftstrarr := binarysearch(content, origin_arr, pos*2+1)
	rightstrarr := binarysearch(content, origin_arr, 2*pos+2)
	if len(leftstrarr) > 0 {
		resarr = append(resarr, leftstrarr...)
	}
	if len(rightstrarr) > 0 {
		resarr = append(resarr, rightstrarr...)
	}
	return
}

// target equal to dataid****value
func compareid(content string, target string) (id string, res bool) {
	targetarr := strings.Split(target, "****")
	if len(targetarr) != 2 {
		res = false
	} else {
		if content == targetarr[1] {
			id = targetarr[0]
			res = true
		} else {
			res = false
		}
	}
	return
}
func ParseDb(actor *os.File) (res map[string]string, err error) {
	buffer := make([]byte, wrbuffersize)
	read := bufio.NewReader(actor)
	lang, err := read.Read(buffer)
	if err == nil {
		res, err = toolsbox.ParseListUltra(buffer[:lang], "****")
	}
	return
}
func FormatDb(originmap map[string]string, actor *os.File) (err error) {
	res := toolsbox.FormatListUltra(originmap, "****")
	_, err = actor.Write([]byte(res))
	return
}

// delete the target vale
func deletetargetvale(idarray []string, actor *os.File) (err error) {
	resmap, err := ParseDb(actor)
	if err == nil {
		//delete the columns that include the target id
		for _, id := range idarray {
			if _, ok := resmap[id]; ok {
				delete(resmap, id)
			}
		}
		err = FormatDb(resmap, actor)
	}
	return
}
