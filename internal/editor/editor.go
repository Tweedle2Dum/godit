package editor

import (
	"fmt"
	"os"
	"golang.org/x/sys/unix"
)


func Die (s string, err error)  {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", s, err)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", s)
	}
	os.Exit(1)
}

func EnableRawMode() (*unix.Termios) {
	fd := int(os.Stdin.Fd())
	t, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		Die("TCGETS",err)
	}
	ct := *t

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

	return t

}

func DisableRawMode(s *unix.Termios) error {
	fd := int(os.Stdin.Fd())
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, s); err != nil {
		return fmt.Errorf("error restoring terminal state: %v", err)
	}

	return nil
}
