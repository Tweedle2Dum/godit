package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unicode"

	"github.com/godit/internal/editor"
)

func CTRLKEY (c rune) rune {
	return c & 0x1f
} 

func main() {
	t:= editor.EnableRawMode()
	defer editor.DisableRawMode(t)
	for {
		var buffer [1]byte
		buffer[0] = '\000'
		_, err := os.Stdin.Read(buffer[:])
		if err != nil && err!=io.EOF && err != syscall.EAGAIN  {
			editor.DisableRawMode(t)
			editor.Die("read",err)
		}
		c:= rune(buffer[0])
		
		if unicode.IsControl(c) {
			fmt.Printf("%d\r\n", c)

		}else{
			fmt.Printf("%d ('%c')\r\n",c,c)
		}
		if(c == CTRLKEY('q')) {break}
	}

}
