package gui

import (
	"github.com/google/gxui"
	"github.com/google/gxui/gxfont"
	"github.com/google/gxui/samples/flags"
)

const (
	width  = 1024
	height = 320
)

func appMain(driver gxui.Driver) {
	theme := flags.CreateTheme(driver)
	window := theme.CreateWindow(width, height, title)
	window.SetScale(flags.DefaultScaleFactor)

	image := theme.CreateImage()

	layout := theme.CreateLinearLayout()
	layout.SetSizeMode(gxui.Fill)

	font32, _ := driver.CreateFont(gxfont.Default, 32)

	label := theme.CreateLabel()
	label.SetFont(font32)
	label.SetColor(gxui.White)

	window.OnClose(driver.Terminate)

	window.AddChild(layout)
	layout.AddChild(image)
	layout.AddChild(label)
}
