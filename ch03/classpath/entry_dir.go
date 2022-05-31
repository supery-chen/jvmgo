package classpath

import (
	"io/ioutil"
	"path/filepath"
)

// 定义DirEntry结构体
type DirEntry struct {
	// 有一个属性,存储目录的绝对路径
	absDir string
}

// Go没有专门的构造函数, 所以这里用new开头的函数来创建结构体实例
func newDirEntry(path string) *DirEntry {
	// 先把path转换为绝对路径, 如果转换过程中出现错误, 则调用panic函数终止程序执行
	absDir, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	// 否则创建DirEntry实例并返回
	return &DirEntry{absDir}
}

// Go结构体不需要显式实现接口, 只要方法匹配即可. 下面的两个方法就是实现了Entry接口的方法

func (de *DirEntry) readClass(className string) ([]byte, Entry, error) {
	// 先把目录和class文件名拼成一个完整的路径
	fileName := filepath.Join(de.absDir, className)
	// 然后调用ioutil包提供的ReadFile函数来读取class文件内容
	data, err := ioutil.ReadFile(fileName)
	// 最后返回
	return data, de, err
}

func (de *DirEntry) String() string {
	// 直接返回目录的绝对路径
	return de.absDir
}
