package storage

import (
	"bytes"
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {
	t.Log("TestStorage")

	s := NewCommonStorage("/Users/tal/code/GoProject/cloud-core")

	ls, err := s.LS("/")
	if err != nil {
		t.Error(err)
		return
	}

	for _, i := range ls {
		fmt.Println(i.String())
	}

	utils, err := s.Get("/internal/storage/storage.go")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(string(utils))

}

func TestStorageDownload(t *testing.T) {
	t.Log("TestStorageDownload")

	s := NewCommonStorage("/Users/tal/code/GoProject/cloud-core")

	buff := bytes.Buffer{}
	err := s.Download("/internal/storage/utils.go", &buff)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(buff.String())
}

func TestCommonStorage_Stat(t *testing.T) {
	t.Log("TestCommonStorage_Stat")

	s := NewCommonStorage("/Users/tal/code/GoProject/cloud-core")

	stat, err := s.Stat("/internal/storage/utils.go")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(stat.String())
}

func TestCommonStorage_Upload(t *testing.T) {
	t.Log("TestCommonStorage_Upload")

	s := NewCommonStorage("/Users/tal/code/GoProject/cloud-core")

	err := s.Upload("/internal/storage/", "test.txt", bytes.NewReader([]byte("Hello world!\n heiheihei!")))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCommonStorage_Remove(t *testing.T) {
	t.Log("TestCommonStorage_Remove")

	s := NewCommonStorage("/Users/tal/code/GoProject/cloud-core")

	err := s.Remove("/internal/storage/test.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCommonStorage_Mkdir(t *testing.T) {
	t.Log("TestCommonStorage_Mkdir")

	s := NewCommonStorage("/Users/tal/code/GoProject/cloud-core")

	err := s.Mkdir("/internal/storage/test")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCommonStorage_Rename(t *testing.T) {
	t.Log("TestCommonStorage_Rename")

	s := NewCommonStorage("/Users/tal/code/GoProject/cloud-core")

	err := s.Rename("/internal/storage/test.txt", "/internal/storage/test1.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCommonStorage_Copy(t *testing.T) {
	t.Log("TestCommonStorage_Copy")

	s := NewCommonStorage("/Users/tal/code/GoProject/cloud-core")

	err := s.Copy("/internal/storage/test1.txt", "/internal/storage/test2.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCommonStorage_Move(t *testing.T) {
	t.Log("TestCommonStorage_Move")

	s := NewCommonStorage("/Users/tal/code/GoProject/cloud-core")

	err := s.Move("/internal/storage/test2.txt", "/internal/storage/test1/test2.txt")
	if err != nil {
		t.Error(err)
		return
	}
}
