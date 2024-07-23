package editor

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func EnableRawMode() (*unix.Termios, error) {
	fd := int(os.Stdin.Fd())
	t, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		fmt.Println("Error getting terminal state for raw mode", err)
		return nil, err
	}
	ct := *t

	ct.Lflag = ct.Lflag &^ (unix.ECHO | unix.ICANON)

	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &ct); err != nil {
		fmt.Println("Error setting terminal to raw mode:", err)
		return nil, err
	}

	return t, nil

}

func DisableRawMode(s *unix.Termios) error {
	fd := int(os.Stdin.Fd())
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, s); err != nil {
		return fmt.Errorf("error restoring terminal state: %v", err)
	}

	return nil
}
