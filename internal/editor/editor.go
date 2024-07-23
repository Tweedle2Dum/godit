package editor

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

var originaltermios *unix.Termios

func CTRLKEY (c rune) rune {
	return c & 0x1f
} 



func Die (s string, err error)  {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", s, err)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", s)
	}
	os.Exit(1)
}

func EnableRawMode()  {
	fd := int(os.Stdin.Fd())
	var err error
	originaltermios, err = unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		Die("TCGETS",err)
	}
	ct := *originaltermios

	ct.Lflag = ct.Lflag &^ (unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN )
	ct.Iflag = ct.Iflag &^(unix.IXON | unix.ICRNL | unix.BRKINT | unix.INPCK | unix.ISTRIP)
	ct.Oflag = ct.Oflag &^(unix.OPOST)
	ct.Cflag = ct.Cflag &^(unix.CS8)
	ct.Cc[unix.VMIN] = 0 
	ct.Cc[unix.VTIME] = 1
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &ct); err != nil {
		fmt.Println("Error setting terminal to raw mode:", err)
		Die("TCSETS",err)
	}

	

}

func disableRawMode() {
	fd := int(os.Stdin.Fd())
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, originaltermios); err != nil {
		Die("error restoring terminal state",err)
	}

	
}

func editorReadKey () rune {
	var buffer [1]byte
		buffer[0] = '\000'
		_, err := os.Stdin.Read(buffer[:])
		if err != nil && err!=io.EOF && err != syscall.EAGAIN  {
			Die("read",err)
		}
		return rune(buffer[0])
}


func ProcessKeyPress () {
	c:= editorReadKey()
	switch c {
	case CTRLKEY('q') :
		os.Stdout.Write([]byte("\x1b[2J"))
		os.Stdout.Write([]byte("\x1b[H"))
		disableRawMode()
		os.Exit(0)
	}
}




func RefreshScreen () {
	os.Stdout.Write([]byte("\x1b[2J"))
	os.Stdout.Write([]byte("\x1b[H"))
}