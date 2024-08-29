package editor

import (
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"syscall"
	"github.com/godit/internal/buffer"
)

type Editor struct {
	t     *unix.Termios
	srows int
	scol  int
}

var e Editor

func CTRLKEY(c rune) rune {
	return c & 0x1f
}

func safeExit() {
	disableRawMode()
	os.Exit(0)
}

func editorDrawRows(b *buffer.Buffer) {
	var r int
	for r = 0; r < e.srows; r++ {
		b.BAppend([]byte("~"))
		b.BAppend([]byte("\x1b[K"));
		if(r < e.srows-1) {
			b.BAppend([]byte("\r\n"))

		}
	}
}

func Die(s string, err error) {

	os.Stderr.Write([]byte(s))
	os.Stderr.Write([]byte("\t -------> "))
	os.Stderr.Write([]byte(err.Error()))
	os.Stderr.Write([]byte("\n"))
	os.Exit(1)
}

func EnableRawMode() {
	fd := int(os.Stdin.Fd())
	var err error
	e.t, err = unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		Die("TCGETS", err)
	}
	ct := *e.t

	ct.Lflag = ct.Lflag &^ (unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN)
	ct.Iflag = ct.Iflag &^ (unix.IXON | unix.ICRNL | unix.BRKINT | unix.INPCK | unix.ISTRIP)
	ct.Oflag = ct.Oflag &^ (unix.OPOST)
	ct.Cflag = ct.Cflag &^ (unix.CS8)
	ct.Cc[unix.VMIN] = 0
	ct.Cc[unix.VTIME] = 1
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &ct); err != nil {
		Die("TCSETS", err)
	}

}

func disableRawMode() {
	fd := int(os.Stdin.Fd())
	os.Stdout.Write([]byte("\x1b[2J"))
	os.Stdout.Write([]byte("\x1b[H"))
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, e.t); err != nil {
		Die("error restoring terminal state", err)
	}
}

func editorReadKey() rune {
	var buffer [1]byte
	buffer[0] = '\000'
	_, err := os.Stdin.Read(buffer[:])
	if err != nil && err != io.EOF && err != syscall.EAGAIN {
		Die("read", err)
	}
	return rune(buffer[0])
}

func ProcessKeyPress() {
	c := editorReadKey()
	switch c {
	case CTRLKEY('q'):
		os.Stdout.Write([]byte("\x1b[2J"))
		os.Stdout.Write([]byte("\x1b[H"))
		safeExit()
	}
}

func RefreshScreen() {
	b := buffer.InitBuffer()
	b.BAppend([]byte("\x1b[?25l"));
	b.BAppend([]byte("\x1b[H"))
	editorDrawRows(b)
	b.BAppend([]byte("\x1b[H"))
	b.BAppend([]byte("\x1b[?25H"));

	os.Stdout.Write(b.Buf)
	b.BFree()

}

func getWindowSize(rows *int, cols *int) (int, error) {
	fd := int(os.Stdin.Fd())
	s, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil && s.Col == 0 {
		w, err := os.Stdout.Write([]byte("\x1b[999C\x1b[999B"))
		if w != 12 {
			return -1, err
		}
		c, err := getCursorPosition(&e.srows, &e.scol)
		return c, err
	}
	*rows = int(s.Row)
	*cols = int(s.Col)
	return 0, nil
}

func InitEditor() {
	c, err := getWindowSize(&e.srows, &e.scol)
	if c == -1 {
		Die("getWindowSize", err)
	}

}

func getCursorPosition(rows *int, cols *int) (int, error) {
	if _, err := os.Stdout.Write([]byte("\x1b[6n")); err != nil {
		return -1, err
	}

	if _, err := fmt.Fscanf(os.Stdin, "\x1b[%d;%d", rows, cols); err != nil {
		return -1, err
	}

	return 0, nil
}
