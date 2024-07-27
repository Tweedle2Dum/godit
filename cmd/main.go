package main

import "github.com/godit/internal/editor"


func main() {
	 editor.EnableRawMode()
	 editor.InitEditor()

	for{
		editor.RefreshScreen()
		editor.ProcessKeyPress()
	}

}
