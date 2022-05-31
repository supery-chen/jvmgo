package main

import (
	"fmt"
	"jvmgo/ch03/classpath"
	"strings"
)

func main() {
	// 手写解析命令行参数
	cmd := parseCmd()
	if cmd.versionFlag {
		// 如果命令行参数指定了-version,则打印版本信息并退出
		fmt.Println("version 0.0.1")
	} else if cmd.helpFlag || cmd.class == "" {
		// 如果命令行参数指定了-help,或未指定class,则打印帮助信息并退出
		printUsage()
	} else {
		// 否则,启动JVM
		startJVM(cmd)
	}
}

func startJVM(cmd *Cmd) {
	cp := classpath.Parse(cmd.XjreOption, cmd.cpOption)
	// 先打印出命令行参数
	fmt.Printf("classpath:%v class:%v args:%v\n", cp, cmd.class, cmd.args)
	// 处理className
	className := strings.Replace(cmd.class, ".", "/", -1)
	// 读取主类数据
	classData, _, err := cp.ReadClass(className)
	if err != nil {
		fmt.Printf("Could not find or load main class %s\n", cmd.class)
	}
	// 打印到控制台
	fmt.Printf("class data:%v\n", classData)
}
