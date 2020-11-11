// Created by vinson on 2020/11/4.

package utils

import (
	"fmt"
	term "github.com/nsf/termbox-go"
	"os"
	"strings"
)

func reset() {
	_ = term.Sync()
}

// Selection the function receive an array of fileInfo
// return the array index from user selected, -1 means selected all
func Selection(files *[]os.FileInfo) int {
	var isExit = false
	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer func() {
		term.Close()
		if isExit {
			os.Exit(0)
		}
	}()
	// 打印到控制到等待输入
	selectIndex := -1
loop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyArrowUp:
				reset()
				selectIndex--
			case term.KeyArrowDown:
				reset()
				selectIndex++
			case term.KeyEnter:
				reset()
				break loop
			case term.KeyEsc:
				reset()
				isExit = true
				break loop
			default:
				reset()
				fmt.Println("Invalid input")
			}
		case term.EventError:
			panic(ev.Err)
		}
		switch {
		case selectIndex < -1:
			selectIndex = -1
		case selectIndex > len(*files)-1:
			selectIndex = len(*files) - 1
		default:
			CallClear()
			var fsBuilder strings.Builder
			fsBuilder.WriteString("\033[1;32;40mPlease select and continue or press ESC to exit\033[0m\n\n")
			if selectIndex == -1 {
				fsBuilder.WriteString("\033[4;33;40mAll\033[0m\n")
			} else {
				fsBuilder.WriteString("All\n")
			}
			for i, f := range *files {
				if i == selectIndex {
					fsBuilder.WriteString("\033[4;33;40m")
					fsBuilder.WriteString(f.Name())
					fsBuilder.WriteString("\033[0m\n")
				} else {
					fsBuilder.WriteString(f.Name())
					fsBuilder.WriteString("\n")
				}
			}
			fmt.Print(fsBuilder.String())
		}
	}
	CallClear()
	return selectIndex
}
