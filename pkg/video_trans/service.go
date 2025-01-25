package video_trans

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type TransItem struct {
	InFileName  string    `json:"in_file_name"`
	OutFileName string    `json:"out_file_name"`
	Filename    string    `json:"filename"`
	Progress    string    `json:"progress"`
	CreateAt    time.Time `json:"create_at"`
}

const (
	DefaultM3U8Path      = "/tmp/m3u8/"
	DefaultM3U8TTL       = time.Hour * 24
	DefaultStoreInterval = time.Minute * 3
)

type Service interface {
	TransState(filepath string) (TransItem, error)
}

type service struct {
	TransList  map[string]*TransItem
	accessLock sync.RWMutex
	storePath  string
	rootPath   string
	cachePath  string
	ticker     *time.Ticker
}

func NewService(root string) Service {
	storePath := root + DefaultM3U8Path
	transList := make(map[string]*TransItem)
	_ = os.MkdirAll(storePath, os.ModePerm)

	cachePath := storePath + "cache.txt"
	cacheFile, err := os.ReadFile(cachePath)
	if err == nil {
		_ = json.Unmarshal(cacheFile, &transList)
		fmt.Println("from file:", transList)
	}

	s := &service{
		TransList: transList,
		storePath: storePath,
		rootPath:  root,
		cachePath: cachePath,
	}
	s.storeInterval()
	return s
}

func (s *service) TransState(filepath string) (TransItem, error) {
	fmt.Println("time: ", time.Now(), "filepath: ", filepath)

	s.accessLock.RLock()
	itemInMemory, ok := s.TransList[filepath]
	s.accessLock.RUnlock()
	if ok {
		return *itemInMemory, nil
	}

	// 使用filepath包安全地处理文件路径
	filename := filepath.Base(filepath)
	filenameWithoutSuffix := strings.TrimSuffix(filename, filepath.Ext(filename))

	newTransItem := &TransItem{
		InFileName:  s.rootPath + filepath,
		OutFileName: s.storePath + filenameWithoutSuffix + ".m3u8",
		Filename:    filenameWithoutSuffix,
		Progress:    "0.00%",
		CreateAt:    time.Now(),
	}

	s.accessLock.Lock()
	s.TransList[filepath] = newTransItem
	for k, v := range s.TransList {
		fmt.Println("trans:", k, v)
	}
	s.accessLock.Unlock()

	go func(filepath string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("TransState goroutine panic: %v, filepath: %s\n", err, filepath)
			}
		}()

		TransWithProgress(newTransItem.InFileName, newTransItem.OutFileName,
			func(progress string) {
				s.accessLock.Lock()
				s.TransList[filepath].Progress = progress
				s.accessLock.Unlock()
			},
			func(err error) {
				s.accessLock.Lock()
				delete(s.TransList, filepath)
				_ = os.Remove(newTransItem.OutFileName)
				// 使用filepath.Join和filepath.Glob正确处理文件路径
				pattern := filepath.Join(s.storePath, newTransItem.Filename+"*.ts")
				deleteFiles, err := filepath.Glob(pattern)
				if err == nil {
					for _, f := range deleteFiles {
						_ = os.Remove(f)
					}
				}
				s.accessLock.Unlock()
			})
	}(filepath)

	return *newTransItem, nil
}

func (s *service) storeInterval() {
	s.ticker = time.NewTicker(DefaultStoreInterval)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("storeInterval panic: %v\n", err)
			}
			// 确保ticker被关闭
			if s.ticker != nil {
				s.ticker.Stop()
			}
		}()

		for {
			select {
			case <-s.ticker.C:
				s.accessLock.Lock()
				finishList := make(map[string]*TransItem)
				for k, v := range s.TransList {
					if v.Progress == "100.00%" {
						finishList[k] = v
						fmt.Println("store:", k, v)
					}
				}
				cacheJson, err := json.Marshal(finishList)
				fmt.Println("storeInterval", string(cacheJson))
				if err == nil {
					err = os.WriteFile(s.cachePath, cacheJson, os.ModePerm)
					if err != nil {
						fmt.Println("storeInterval", err)
					}
				}
				s.accessLock.Unlock()
			}
		}
	}()
}
