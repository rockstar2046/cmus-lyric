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

	duration := 500 * time.Millisecond
	tick := time.NewTicker(duration)

	uiEvents := ui.PollEvents()
	for {
		select {
		case <-tick.C:
			pkg.Listen(curFile, curLyric, curPos, keys)
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "m":
				tick.Stop()
				pkg.DrawComments()
			case "y":
				tick = time.NewTicker(duration)
			case "?":
				tick.Stop()
				pkg.Help()
			}
		}
	}

}
