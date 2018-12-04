// Copyright 2018 ra <rockagen@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package main

import (
	"cmus-lyric/pkg"
	ui "github.com/gizak/termui"
)
import (
	"log"
	"time"
)

func main() {
	log.SetFlags(0)
	err := ui.Init()
	if err != nil {
		panic(err)
	}

	defer ui.Close()

	var curFile string

	var curLyric map[int][]string
	var curPos int

	var keys []int

	tick := time.Tick(500 * time.Millisecond)

	uiEvents := ui.PollEvents()
	for {
		select {
		case <-tick:
			pkg.Listen(curFile, curLyric, curPos, keys)
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		}

	}

}
