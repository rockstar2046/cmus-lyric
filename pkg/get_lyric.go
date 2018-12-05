// Copyright 2018 ra <rockagen@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	LyricApi   = "http://music.163.com/weapi/song/lyric?csrf_token="
	SearchApi  = "http://music.163.com/weapi/cloudsearch/get/web?csrf_token="
	CommentApi = "http://music.163.com/weapi/v1/resource/comments/R_SO_4_%v/?csrf_token="
	Cookie     = "os=pc; osver=Microsoft-Windows-10-Professional-build-10586-64bit; appver=2.0.3.131777; channel=netease; __remember_me=true"
	UserAgent  = "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"
)

type ReInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type SongInfo struct {
	Id    int        `json:"id"`
	Name  string     `json:"name"`
	Dt    int        `json:"dt"`
	Songs []SongInfo `json:"songs"`
}

type SearchInfo struct {
	Result DesInfo `json:"result"`
	ReInfo
}

type DesInfo struct {
	Songs []SongInfo `json:"songs"`
	ReInfo
}

type LrcInfo struct {
	Lyric string `json:"lyric"`
}

type LyricInfo struct {
	Lrc    LrcInfo `json:"lrc"`
	Tlyric LrcInfo `json:"tlyric"`
	ReInfo
}

type HotCommentInfo struct {
	Content    string `json:"content"`
	LikedCount int    `json:"likedCount"`
}

type CommentInfo struct {
	HotCommentInfo []HotCommentInfo `json:"hotComments"`
	CommentInfo    []HotCommentInfo `json:"comments"`
	ReInfo
}

func FetchLyric(dir string, name string, duration int) {
	sid := FindId(name, duration)
	if len(sid) > 0 {
		lyrc, tlyrc := GetLyric(sid)
		if len(lyrc) > 0 {
			path := dir + "/" + name + ".lyric"
			save(path, strings.NewReader(lyrc))
		}
		if len(tlyrc) > 0 {
			path := dir + "/" + name + ".t.lyric"
			save(path, strings.NewReader(tlyrc))
		}
	}
}

func FetchLyricCmus(file string, dt int) {
	pathIdx := strings.LastIndexAny(file, ".")
	titleIdx := strings.LastIndexAny(file, "/")
	dir := file[:titleIdx]
	title := file[titleIdx+1 : pathIdx]
	FetchLyric(dir, title, dt)
}

func FindId(name string, duration int) string {

	m := make(map[string]interface{})

	m["s"] = name
	m["type"] = 1
	m["limit"] = 10
	m["offset"] = 0
	m["total"] = true
	m["csrf_token"] = ""

	req, _ := json.Marshal(m)
	params, encSecKey, _ := EncParams(string(req))

	resp, e := post(SearchApi, params, encSecKey)
	if e != nil {
		log.Println(e)
		return ""
	}

	ret := &SearchInfo{}

	e2 := json.Unmarshal(resp, ret)

	if e2 != nil {
		log.Println(e2)
		return ""
	}
	code := ret.Code

	if 200 != code {
		log.Printf("code: %v, msg: %v \n", code, ret.Msg)
		return ""
	}

	if len(ret.Result.Songs) > 0 {
		for _, v := range ret.Result.Songs {
			dt := v.Dt / 1000
			if dt == duration {
				return strconv.Itoa(v.Id)
			}
		}
	}

	return ""
}

func GetLyric(id string) (string, string) {

	var lyrc, tlyrc string
	m := make(map[string]interface{})

	m["id"] = id
	m["os"] = "pc"
	m["lv"] = -1
	m["kv"] = -1
	m["tv"] = -1
	m["os"] = "pc"
	m["csrf_token"] = ""

	req, _ := json.Marshal(m)
	params, encSecKey, _ := EncParams(string(req))

	resp, e := post(LyricApi, params, encSecKey)
	if e != nil {
		log.Println(e)
		return lyrc, tlyrc
	}

	result := &LyricInfo{}

	e2 := json.Unmarshal(resp, result)

	if e2 != nil {
		log.Println(e2)
		return lyrc, tlyrc
	}

	//
	code := result.Code

	if 200 != code {
		log.Printf("code: %v, msg: %v \n", code, result.Msg)
		return lyrc, tlyrc
	}
	lyrc, tlyrc = result.Lrc.Lyric, result.Tlyric.Lyric
	return lyrc, tlyrc

}

func GetHotComments(id string) ([]HotCommentInfo, []HotCommentInfo) {

	empty := make([]HotCommentInfo, 0)

	m := make(map[string]interface{})

	m["rid"] = ""
	m["limit"] = 50
	m["offset"] = 0
	m["total"] = true
	m["csrf_token"] = ""

	req, _ := json.Marshal(m)
	params, encSecKey, _ := EncParams(string(req))

	resp, e := post(fmt.Sprintf(CommentApi, id), params, encSecKey)
	if e != nil {
		log.Println(e)
		return empty, empty
	}

	ret := &CommentInfo{}

	e2 := json.Unmarshal(resp, ret)

	if e2 != nil {
		log.Println(e2)
		return empty, empty
	}
	code := ret.Code

	if 200 != code {
		log.Printf("code: %v, msg: %v \n", code, ret.Msg)
		return empty, empty
	}

	return ret.HotCommentInfo, ret.CommentInfo
}

func post(_url, params, encSecKey string) ([]byte, error) {
	client := &http.Client{}
	form := url.Values{}
	form.Set("params", params)
	form.Set("encSecKey", encSecKey)

	request, err := http.NewRequest("POST", _url, strings.NewReader(form.Encode()))

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Host", "music.163.com")
	request.Header.Set("Origin", "http://music.163.com")
	request.Header.Set("User-Agent", UserAgent)

	request.Header.Set("Cookie", Cookie)

	resp, err := client.Do(request)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	resBody, resErr := ioutil.ReadAll(resp.Body)
	if resErr != nil {
		log.Println(err)
		return nil, resErr
	}
	return resBody, nil
}

func save(path string, src io.Reader) {
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		log.Printf("Write eror: %v \n", err)
		return
	}
	n, err := io.Copy(out, src)
	if err != nil {
		log.Println(err)
	}
	log.Printf("<- %v, size: %v", path, n)
}
