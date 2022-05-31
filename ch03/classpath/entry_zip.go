package classpath

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"path/filepath"
)

type ZipEntry struct {
	// 存放ZIP或JAR文件的绝对路径
	absPath string
}

func newZipEntry(path string) *ZipEntry {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &ZipEntry{absPath}
}

// TODO 作者在书中说, readClass方法每次都要打开或关闭ZIP文件, 因此效率不是很高. 作者进行了优化, 在ch03\classpath\entry_zip2.go
func (ze *ZipEntry) readClass(className string) ([]byte, Entry, error) {
	// 首先打开ZIP包
	r, err := zip.OpenReader(ze.absPath)
	// 如果打开ZIP包出错, 直接返回
	if err != nil {
		return nil, nil, err
	}
	// defer 会在当前函数返回前执行传入的函数, 我的理解是类似java中的finally, 常用于关闭文件描述符, 关闭数据库连接以及解锁资源, 类似java中实现Closeable接口时自动关闭资源?
	// 确保r自动关闭
	defer r.Close()
	// 遍历ZIP包中的文件
	for _, f := range r.File {
		// 找到文件名与className相同的文件
		if f.Name == className {
			// 尝试打开文件, rc为ReadCloser, ReadCloser is the interface that groups the basic Read and Close methods.
			rc, err := f.Open()
			// 如果打开失败, 直接返回
			if err != nil {
				return nil, nil, err
			}
			// 确保rc自动关闭
			defer rc.Close()
			// 调用ioutil的ReadAll方法, 读取rc的全部内容. ReadAll reads from r until an error or EOF and returns the data it read.
			data, err := ioutil.ReadAll(rc)
			// 如果读取失败, 直接返回
			if err != nil {
				return nil, nil, err
			}
			// 返回正确结果
			return data, ze, nil
		}
	}
	// 如果ZIP包中未找到与className名称一致的文件, 构造一个错误提示并返回
	return nil, nil, errors.New("class not found: " + className)
}

func (ze *ZipEntry) String() string {
	return ze.absPath
}
