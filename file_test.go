package polevpnmobile

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestFileReadFromTail(t *testing.T) {

	f, err := os.Open("./test.log")

	if err != nil {
		t.Fatal("open file fail,", err)
	}

	i, err := f.Stat()

	if err != nil {
		t.Fatal("open file fail,", err)
	}

	lines := make([][]byte, 0)
	lineBuf := make([]byte, 4096)
	var offset int64 = 0

	buf := make([]byte, 1)
	index := 0

	count := 0

	size := 0

	for {

		offset -= 1

		f.Seek(int64(offset), io.SeekEnd)
		_, err := f.Read(buf)
		if err != nil {
			break
		}
		size++
		lineBuf[index] = buf[0]
		index++
		if buf[0] == '\n' {
			lines = append(lines, lineBuf[:index])
			index = 0
			lineBuf = make([]byte, 4096)
			count++
			if count == 100 {
				break
			}
		}

		if -offset >= i.Size() {
			break
		}
	}

	if index > 0 {
		lines = append(lines, lineBuf[:index])
	}

	out := make([]byte, size)
	index = 0
	for i := len(lines) - 1; i > -1; i-- {
		line := lines[i]
		for j := len(line) - 1; j > -1; j-- {
			out[index] = line[j]
			index++
		}
	}

	fmt.Print(string(out))
}
