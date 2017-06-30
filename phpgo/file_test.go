package phpgo

import (
	"testing"
)

func BenchmarkFileSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FileSize("file.go")
	}
}

func BenchmarkFileExist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FileExist("file.go")
	}
}

func TestFileSize(t *testing.T) {
	file := `string.go`
	size, err := FileSize(file)
	if err != nil {
		t.Error("FileSize Error", err)
		return
	}

	t.Log("FileSize", size)
}

func TestFileExist(t *testing.T) {
	file := `rand.go`
	if exist := FileExist(file); !exist {
		t.Error("FileExist no exist")
		return
	}

	t.Log("FileExist exist")
}

func TestFileMd5File(t *testing.T) {
	file := "rand.go"
	h, err := FileMd5File(file)
	if err != nil {
		t.Error("FileMd5File Error", err)
	}
	t.Log("FileMd5File:", h)
}
