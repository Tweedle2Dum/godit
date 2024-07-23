package main

import "github.com/godit/internal/editor"


func main() {
	 editor.EnableRawMode()
	for{
		editor.RefreshScreen()
		editor.ProcessKeyPress()
	}

}
