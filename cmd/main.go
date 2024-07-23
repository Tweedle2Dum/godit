package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unicode"

	"github.com/godit/internal/editor"
)

func main() {
	t:= editor.EnableRawMode()
	defer editor.DisableRawMode(t)
	for {
		var buffer [1]byte
		buffer[0] = '\000'
		n, err := os.Stdin.Read(buffer[:])
		if err != nil && err!=io.EOF && err != syscall.EAGAIN  {
			editor.DisableRawMode(t)
			editor.Die("read",err)
		}
		if n>0 && buffer[0] == 'q' {
			fmt.Println("\nExiting...")
			break
		}
		if unicode.IsControl(rune(buffer[0])) {
			fmt.Printf("%d\r\n", rune(buffer[0]))

		}else{
			fmt.Printf("%d ('%c')\r\n",rune(buffer[0]),rune(buffer[0]))
		}
		
	}

}
