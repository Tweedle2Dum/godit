package main

import (
	"fmt"
	"os"
	"github.com/godit/internal/editor"
)

func main() {
	var buffer [1]byte
	t,err := editor.EnableRawMode()
	if(err!=nil){
		return
	}
	defer editor.DisableRawMode(t)
	for {
		n, err := os.Stdin.Read(buffer[:])
		if err != nil {
			fmt.Println("Error reading input:", err)
		}
		if n>0 && buffer[0] == 'q' {
			fmt.Println("\nExiting...")
			break
		}
	}

}
