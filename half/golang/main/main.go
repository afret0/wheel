package main

import (
	"log"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	// 创建消息列表
	messageList := widgets.NewList()
	messageList.Title = "聊天记录"
	messageList.Rows = []string{}

	// 创建输入框
	input := widgets.NewParagraph()
	input.Title = "输入消息"
	input.Text = ""

	// 设置布局
	grid := termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		termui.NewRow(0.9, messageList),
		termui.NewRow(0.1, input),
	)

	// 渲染UI
	termui.Render(grid)

	// 处理用户输入
	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Enter>":
			// 发送消息
			messageList.Rows = append(messageList.Rows, "You: "+input.Text)
			input.Text = ""
		default:
			// 处理输入
			input.Text += e.ID
		}

		termui.Render(grid)
	}
}
