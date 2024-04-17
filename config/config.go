package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Unknwon/goconfig"
)

const configFile = "/conf/conf.ini"

var File *goconfig.ConfigFile

// 加载文件时, 先走init方法
func init() {
	//  文件系统读取当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + configFile
	// 允许自定义文件位置
	len := len(os.Args)
	if len > 1 {
		// 存在命令行参数
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + configFile
		}
	}
	if !fileExist(configPath) {
		// 如果这里就没有, 就没有必要往下执行了
		panic(errors.New("配置文件不存在"))
	}
	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("读取配置文件出错: ", err)
	}
	// fmt.Println(File)
}

func fileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func A() {
	fmt.Println("wocao")
}
