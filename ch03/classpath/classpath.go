package classpath

import (
	"fmt"
	"os"
	"path/filepath"
)

// [为什么go中的receiver name不推荐使用this或者self](https://segmentfault.com/a/1190000025188649)

// 三个字段分别存放三种类路径
type Classpath struct {
	bootClasspath Entry
	extClasspath  Entry
	userClasspath Entry
}

// 使用-Xjre选项解析启动类和扩展类路径, 使用-classpath/-cp选项解析用户类路径
func Parse(jreOption, cpOption string) *Classpath {
	cp := &Classpath{}
	cp.parseBootAndExtClasspath(jreOption)
	cp.parseUserClasspath(cpOption)
	return cp
}

func (cp *Classpath) parseBootAndExtClasspath(jreOption string) {
	jreDir := getJreDir(jreOption)
	fmt.Printf("jreDir %s\n", jreDir)
	// jre/lib/*
	jreLibPath := filepath.Join(jreDir, "lib", "*")
	cp.bootClasspath = newWildcardEntry(jreLibPath)
}

// 优先使用用户输入的-Xjre选项作为jre目录
// 如果没有输入该选项, 则在当前目录下寻找jre目录
// 如果找不到, 尝试使用JAVA_HOME环境变量
// 如果都找不到, 则使用panic停止
func getJreDir(jreOption string) string {
	if jreOption != "" && exists(jreOption) {
		return jreOption
	}
	if exists("./jre") {
		return "./jre"
	}
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		return filepath.Join(jh, "jre")
	}
	panic("Can not find jre folder!")
}

// 用于判断目录是否存在
func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 如果用户没有提供-classpath/-cp选项, 则使用当前目录作为用户类路径
func (cp *Classpath) parseUserClasspath(cpOption string) {
	if cpOption == "" {
		cpOption = "."
	}
	cp.userClasspath = newEntry(cpOption)
}

// 从启动类路径、扩展类路径和用户类路径中搜索class文件
// 注意, 传递给ReadClass方法的类名不包含.class后缀,所以需要手动填充
func (cp *Classpath) ReadClass(className string) ([]byte, Entry, error) {
	className = className + ".class"
	if data, entry, err := cp.bootClasspath.readClass(className); err == nil {
		return data, entry, nil
	}
	if data, entry, err := cp.extClasspath.readClass(className); err == nil {
		return data, entry, nil
	}
	return cp.userClasspath.readClass(className)
}

// 返回用户类路径的字符串表示
func (cp *Classpath) String() string {
	return cp.userClasspath.String()
}
