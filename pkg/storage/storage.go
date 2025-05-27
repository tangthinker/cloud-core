package storage

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type FileStat struct {
	Filename  string
	Size      int64
	Path      string
	Dir       bool
	Timestamp int64
}

func (fs *FileStat) String() string {
	return fmt.Sprintf("Filename: %s, Size: %d, Path: %s, Dir: %t, Timestamp: %d",
		fs.Filename, fs.Size, fs.Path, fs.Dir, fs.Timestamp)
}

type Storage interface {
	// LS 列出目录下的文件
	LS(path string) ([]FileStat, error)

	// Stat 获取文件信息
	Stat(path string) (*FileStat, error)

	// Get 获取文件二进制
	Get(path string) ([]byte, error)

	// GetStream 获取文件流
	GetStream(path string) (io.ReadCloser, error)

	// Download 下载文件
	Download(path string, writer io.Writer) error

	// Upload 上传文件
	Upload(path string, filename string, reader io.Reader) error

	// Remove 删除文件
	Remove(path string) error

	// Mkdir 创建目录
	Mkdir(path string) error

	// Rename 重命名文件
	Rename(oldPath, newPath string) error

	// Copy 复制文件
	Copy(srcPath, dstPath string) error

	// Move 移动文件
	Move(srcPath, dstPath string) error

	RealPath(path string) string
}

type CommonStorage struct {
	RootPath string
}

func NewCommonStorage(rootPath string) Storage {
	return &CommonStorage{RootPath: rootPath}
}

func (c *CommonStorage) LS(path string) ([]FileStat, error) {
	realPath := c.realPath(path)
	dirStat, err := os.Stat(realPath)
	if err != nil {
		return nil, err
	}
	if !dirStat.IsDir() {
		return nil, fmt.Errorf("path %s is not a directory", path)
	}

	file, err := os.Open(realPath)
	if err != nil {
		return nil, err
	}

	files, err := file.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var fileStats []FileStat
	for _, f := range files {
		cur := FileStat{
			Filename:  f.Name(),
			Size:      f.Size(),
			Path:      path + f.Name(),
			Dir:       f.IsDir(),
			Timestamp: f.ModTime().Unix(),
		}
		if f.IsDir() {
			cur.Path += "/"
		}
		if strings.HasPrefix(cur.Filename, ".") {
			continue
		}
		fileStats = append(fileStats, cur)
	}

	return fileStats, nil
}

func (c *CommonStorage) Stat(path string) (*FileStat, error) {
	realPath := c.realPath(path)
	fileStat, err := os.Stat(realPath)
	if err != nil {
		return nil, err
	}
	return &FileStat{
		Filename:  fileStat.Name(),
		Size:      fileStat.Size(),
		Path:      path,
		Dir:       fileStat.IsDir(),
		Timestamp: fileStat.ModTime().Unix(),
	}, nil
}

func (c *CommonStorage) Get(path string) ([]byte, error) {
	realPath := c.realPath(path)
	isDir, err := isDirectory(realPath)
	if err != nil {
		return nil, err
	}
	if isDir {
		return nil, fmt.Errorf("path %s is a directory", path)
	}

	file, err := os.Open(realPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *CommonStorage) GetStream(path string) (io.ReadCloser, error) {
	realPath := c.realPath(path)
	isDir, err := isDirectory(realPath)
	if err != nil {
		return nil, err
	}
	if isDir {
		return nil, fmt.Errorf("path %s is a directory", path)
	}

	file, err := os.Open(realPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (c *CommonStorage) Download(path string, writer io.Writer) error {
	realPath := c.realPath(path)
	isDir, err := isDirectory(realPath)
	if err != nil {
		return err
	}
	if isDir {
		return fmt.Errorf("path %s is a directory", path)
	}

	file, err := os.Open(realPath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(writer, file)

	if err != nil {
		return err
	}

	return nil
}

func (c *CommonStorage) Upload(path string, filename string, reader io.Reader) error {
	realPath := c.realPath(path)
	isDir, err := isDirectory(realPath)
	if err != nil {
		return err
	}
	if !isDir {
		return fmt.Errorf("path %s is not a directory", path)
	}

	file, err := os.Create(realPath + filename)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, reader)

	if err != nil {
		return err
	}

	return nil
}

func (c *CommonStorage) Remove(path string) error {
	realPath := c.realPath(path)
	return os.RemoveAll(realPath)
}

func (c *CommonStorage) Mkdir(path string) error {
	realPath := c.realPath(path)
	return os.MkdirAll(realPath, 0755)
}

func (c *CommonStorage) Rename(oldPath, newPath string) error {
	realPath := c.realPath(oldPath)
	newRealPath := c.realPath(newPath)
	return os.Rename(realPath, newRealPath)
}

func (c *CommonStorage) Copy(srcPath, dstPath string) error {
	realFile := c.realPath(srcPath)
	dstRealPath := c.realPath(dstPath)

	srcFile, err := os.Open(realFile)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	dstFile, err := os.Create(dstRealPath)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommonStorage) Move(srcPath, dstPath string) error {
	realFile := c.realPath(srcPath)
	dstRealPath := c.realPath(dstPath)

	return os.Rename(realFile, dstRealPath)
}

func (c *CommonStorage) RealPath(path string) string {
	realPath := c.realPath(path)
	return realPath
}

func (c *CommonStorage) realPath(path string) string {
	return c.RootPath + path
}

func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
