// Copyright 2018 ra <rockagen@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package pkg

import (
	"bytes"
	ui "github.com/gizak/termui"
	"io/ioutil"
	"log"
	"netease/pkg"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func Listen(curFile string, curLyric map[int][]string, curPos int, keys []int) {

	pos, file, dt := cmusRemote()
	if pos > 0 {
		if curFile != file {
			curFile = file
			curLyric = loadLyrics(file)
			if curLyric == nil {
				drawEmpty()
				pkg.FetchLyricCmus(file, dt)
				curLyric =loadLyrics(file)
				return
			} else {
				keys = make([]int, 0, len(curLyric))
				for k := range curLyric {
					keys = append(keys, k)
				}
				sort.Ints(keys)
			}
		}

		var tmpPos int
		for i, n := range keys {
			if pos < n {
				tmpPos = keys[i-1]
				break
			}
		}
		if curPos == tmpPos && curPos != 0 {
			return

		}

		curPos = tmpPos

		list := make([]string, 2*len(keys))

		idx, cline := 0, 0

		for _, v := range keys {

			data := curLyric[v]

			if curPos == v && v != 0 {

				text := data[0]

				if len(text) < 1 {
					text = "..."
				}
				data[0] = "[" + text + "](fg-cyan)"
				if len(data[1]) > 0 {
					data[1] = "[" + data[1] + "](fg-cyan)"
					cline++
				}
				cline = idx

			}
			list[idx] = data[0]
			idx++
			if len(data[1]) > 0 {
				list[idx] = data[1]
				idx++
			}
		}
		draw(list, cline)

	} else {
		drawEmpty()
	}
}

func drawEmpty() {
	draw([]string{"", "[no lyrics](fg-red)"}, 0)
}

func draw(list []string, cline int) {

	height := ui.TermHeight()
	ls := ui.NewList()
	ls.BorderLabel = "[" + list[0] + "](fg-white,bg-blue)"
	ls.PaddingTop = 2
	ls.Height = height
	ls.Width = 50
	ls.Border = false
	idx := 1
	if cline+2 > height {
		idx = cline - 1
	}
	ls.Items = list[idx:]

	ui.Render(ls)
}

func loadLyrics(path string) map[int][]string {

	pathIdx := strings.LastIndexAny(path, ".")

	lpath := path[:pathIdx] + ".lyric"
	tlpath := path[:pathIdx] + ".t.lyric"

	titleIdx := strings.LastIndexAny(path, "/")
	title := path[titleIdx+1 : pathIdx]

	content, e := ioutil.ReadFile(lpath)
	if e != nil {
		return nil
	}
	// translate lyric
	var tlines []string
	lines := strings.Split(string(content), "\n")
	t_content, e := ioutil.ReadFile(tlpath)
	if e == nil {
		tlines = strings.Split(string(t_content), "\n")
	}

	m := make(map[int][]string)

	lyricMap := buildLyricMap(lines)
	tlyricMap := buildLyricMap(tlines)

	for k, v := range lyricMap {
		t1 := v
		t2 := tlyricMap[k]
		m[k] = []string{t1, t2,}
	}
	m[0] = []string{title, ""}
	return m
}

func buildLyricMap(lyric []string) map[int]string {
	m := make(map[int]string)
	re := regexp.MustCompile("^\\[([0-9]+):([0-9]+).*](.*)")
	for _, v := range lyric {
		ar := re.FindStringSubmatch(v)
		if len(ar) > 3 {
			mi, _ := strconv.Atoi(ar[1])
			sec, _ := strconv.Atoi(ar[2])

			pos := 60*mi + sec

			m[pos] = ar[3]
		}

	}
	return m
}

//
// return [position,file,duration]
func cmusRemote() (int, string, int) {
	cmd := exec.Command("cmus-remote", "-Q")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal("cmus exited.")
		os.Exit(-1)
	}
	info := strings.Split(out.String(), "\n")

	if len(info) < 1{
		log.Fatal("cmus exited.")
		os.Exit(-1)
	}
	//status stopped
	status := strings.Split(info[0], " ")[1]
	if status != "playing" {
		return 0, "", 0
	}

	if status == "pause" {
		return 1, "", 0
	}
	idx := strings.Index(info[1], " ") + 1

	file := info[1][idx:]
	position := strings.Split(info[3], " ")[1]
	duration := strings.Split(info[2], " ")[1]
	pos, err := strconv.Atoi(position)
	dt, err := strconv.Atoi(duration)
	return pos, file, dt
}