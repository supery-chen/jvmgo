package classpath

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// WildcardEntry实际上也是CompositeEntry, 所以就不再定义新的类型了

func newWildcardEntry(path string) CompositeEntry {
	// 去除path中的*
	baseDir := path[:len(path)-1]
	compositeEntry := []Entry{}
	// 定义walk方法
	walkFn := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//通配符类路径不能递归匹配子目录下的JAR文件
		if info.IsDir() && path != baseDir {
			return filepath.SkipDir
		}
		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") {
			jarEntry := newZipEntry(path)
			compositeEntry = append(compositeEntry, jarEntry)
		}
		return nil
	}
	filepath.Walk(baseDir, walkFn)
	return compositeEntry
}
