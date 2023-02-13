package ui

import (
	"fmt"
	"github.com/nathan-fiscaletti/consolesize-go"
	"strings"
	"time"
)

var syms = [6]string{"   ", "=  ", "== ", "===", " ==", "  ="}

type ProgressBar struct {
	cols int // the max number of characters per line in the terminal
	c    chan string
}

func (p *ProgressBar) getEmptyLineLength(s string) int {
	emptyLength := p.cols - len(s)
	if emptyLength < 0 {
		emptyLength = 0
	}
	return emptyLength
}

func (p *ProgressBar) handle() {
	var i uint8
	var emptyLength int
	var data string
	for {
		if len(p.c) > 0 {
			data = <-p.c
			if data == "" {
				close(p.c)
				break
			}
			emptyLength = p.getEmptyLineLength(data)
			if emptyLength > 8 {
				emptyLength -= 8
			}
		}
		i++
		fmt.Printf("[%s] %s..%s\r", syms[i%6], data, strings.Repeat(" ", emptyLength))
		time.Sleep(time.Second / 9)
	}
}

func (p *ProgressBar) println(text string) {
	emptyLength := p.getEmptyLineLength(text)
	fmt.Println(text, strings.Repeat(" ", emptyLength))
}

func (p *ProgressBar) UpdateStatus(text ...string) {
	p.c <- strings.Join(text, " ")
}

func (p *ProgressBar) Println(text string) {
	p.println(text)
}

func (p *ProgressBar) Error(err error) {
	p.println(err.Error())
}

func (p *ProgressBar) Start() {
	go p.handle()
}

func (p *ProgressBar) Finish(text string) {
	p.c <- ""
	p.println(text)
}

func CreateProgressBar() *ProgressBar {
	cols, _ := consolesize.GetConsoleSize()
	p := &ProgressBar{
		cols: cols - 1,
		c:    make(chan string, 1),
	}

	return p
}
