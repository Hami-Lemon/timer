package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os"
	"time"
)

const (
	dur         = time.Duration(1000 / 24) //24帧
	defaultText = "00:00:00.00"
	width       = 280
	height      = 220
)

type Counter struct {
	ticker  *time.Ticker
	start   time.Time
	pass    time.Duration
	running bool
	stop    chan struct{}
}

func main() {
	c := Counter{
		running: false,
		stop:    make(chan struct{}, 1),
	}
	var textLabel *walk.TextLabel
	var resetBtn, startBtn *walk.PushButton
	update := func() {
		for {
			select {
			case <-c.stop:
				c.ticker.Stop()
				return
			case now := <-c.ticker.C:
				c.pass = now.Sub(c.start)
				millisecond := c.pass / time.Millisecond
				ms := millisecond % 1000 / 10
				millisecond /= 1000
				s := millisecond % 60
				millisecond /= 60
				m := millisecond % 60
				h := millisecond / 60
				_ = textLabel.SetText(fmt.Sprintf("%02d:%02d:%02d.%02d", h, m, s, ms))
			}
		}
	}
	window := MainWindow{
		Title: "计时器",
		Bounds: Rectangle{
			X: 1200,
			Y: 300,
		},
		Size:            Size{Width: width, Height: height},
		MinSize:         Size{Width: width, Height: height},
		MaxSize:         Size{Width: width, Height: height},
		Layout:          VBox{},
		DoubleBuffering: true,
		Children: []Widget{
			TextLabel{
				MinSize:       Size{Width: width, Height: 120},
				Text:          defaultText,
				Font:          Font{Family: "msyh", PointSize: 22},
				TextAlignment: AlignHCenterVCenter,
				AssignTo:      &textLabel,
			},
			HSplitter{
				MaxSize: Size{Width: width, Height: 40},
				Children: []Widget{
					PushButton{
						Text:     "重置",
						AssignTo: &resetBtn,
						OnClicked: func() {
							c.pass = 0
							_ = textLabel.SetText(defaultText)
							_ = startBtn.SetText("开始")
						},
					},
					PushButton{
						Text:     "开始",
						AssignTo: &startBtn,
						OnClicked: func() {
							c.running = !c.running
							if c.running {
								resetBtn.SetEnabled(false)
								c.start = time.Now().Add(-c.pass)
								_ = startBtn.SetText("暂停")
								c.ticker = time.NewTicker(dur * time.Millisecond)
								go update()
							} else {
								_ = startBtn.SetText("继续")
								resetBtn.SetEnabled(true)
								c.stop <- struct{}{}
								c.ticker.Stop()
							}
						},
					},
				},
			},
		},
	}
	code, err := window.Run()
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(code)
}
