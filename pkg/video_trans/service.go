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
	// s.clearInterval()
	s.storeInterval()
	return s
}

func (s *service) TransState(filepath string) (TransItem, error) {
	s.accessLock.RLock()
	itemInMemory, ok := s.TransList[filepath]
	s.accessLock.RUnlock()
	if ok {
		return *itemInMemory, nil
	}

	filename := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	filenameWithoutSuffix := strings.Split(filename, ".")[0]

	newTransItem := &TransItem{
		InFileName:  s.rootPath + filepath,
		OutFileName: s.storePath + filenameWithoutSuffix + ".m3u8",
		Filename:    filenameWithoutSuffix,
		Progress:    "0.00%",
		CreateAt:    time.Now(),
	}

	go func() {

		s.accessLock.Lock()
		s.TransList[filepath] = newTransItem
		s.accessLock.Unlock()

		TransWithProgress(newTransItem.InFileName, newTransItem.OutFileName,
			func(progress string) {
				s.TransList[filepath].Progress = progress
			},
			func(err error) {
				s.accessLock.Lock()
				delete(s.TransList, filepath)
				s.accessLock.Unlock()
			})
	}()

	return *newTransItem, nil
}

func (s *service) clearInterval() {

	ticker := time.NewTicker(DefaultM3U8TTL / 2)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.accessLock.Lock()
				for k, v := range s.TransList {
					if time.Since(v.CreateAt) > DefaultM3U8TTL {
						delete(s.TransList, k)
						_ = os.Remove(v.OutFileName)
						deleteFiles, _ := filepath.Glob(s.storePath + v.Filename + "*.ts")
						for _, f := range deleteFiles {
							_ = os.Remove(f)
						}
					}
				}
				s.accessLock.Unlock()
			}
		}
	}()
}

func (s *service) storeInterval() {

	ticker := time.NewTicker(DefaultStoreInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.accessLock.RUnlock()
				finishList := make(map[string]*TransItem)
				for k, v := range s.TransList {
					if v.Progress == "100.00%" {
						finishList[k] = v
					}
				}
				s.accessLock.RUnlock()
				cacheJson, err := json.Marshal(finishList)
				fmt.Println("storeInterval", string(cacheJson))
				if err == nil {
					err = os.WriteFile(s.cachePath, cacheJson, os.ModePerm)
					if err != nil {
						fmt.Println("storeInterval", err)
					}
				}
			}
		}
	}()
}
