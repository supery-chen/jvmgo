package main

import "fmt"

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
	fmt.Printf("classpath:%s class:%s args:%v\n", cmd.cpOption, cmd.class, cmd.args)
}
