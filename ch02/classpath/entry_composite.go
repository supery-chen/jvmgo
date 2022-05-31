package classpath

import (
	"errors"
	"strings"
)

// CompositeEntry是由更小的Entry组成, 正好可以表示成[]Entry
// 在Go语言中, 数组属于比较底层的数据结构, 很少直接使用. 在大部分情况下, 使用更遍历的slice类型
type CompositeEntry []Entry

// 构造函数把参数(路径列表)按分隔符分成小路径, 然后把每个小路径都转换成具体的Entry实例
func newCompositeEntry(pathList string) CompositeEntry {
	// 初始化为空数组
	compositeEntry := []Entry{}
	// 遍历小路径
	for _, path := range strings.Split(pathList, pathListSeparator) {
		// 对每个小路径调用newEntry
		entry := newEntry(path)
		compositeEntry = append(compositeEntry, entry)
	}
	return compositeEntry
}

// 依次调用每一个子路径的readClass方法
func (ce CompositeEntry) readClass(className string) ([]byte, Entry, error) {
	for _, entry := range ce {
		// 如果成功读取到class数据, 返回数据, 否则继续
		data, from, err := entry.readClass(className)
		if err == nil {
			return data, from, nil
		}
	}
	// 如果遍历完所有子路径还没有找到class文件, 则返回错误
	return nil, nil, errors.New("class not found: " + className)
}

// 调用每一个子路径的String()方法, 然后把得到的字符串用路径分隔符拼接起来即可
func (ce CompositeEntry) String() string {
	// 创建一个slice, 类型为string, 长度为子entry的大小
	strs := make([]string, len(ce))
	// 遍历每个entry, 填充其string到slice中
	for i, entry := range ce {
		strs[i] = entry.String()
	}
	// 拼接并返回
	return strings.Join(strs, pathListSeparator)
}
