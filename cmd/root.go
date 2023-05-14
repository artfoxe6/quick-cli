/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var name = ""
var old = "github.com/artfoxe6/quick-gin"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "quick-cli",
	Short: "快速初始化基于 quick-gin 的项目",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: setup,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&name, "name", "", "项目名称")
	rootCmd.PersistentFlags().StringVar(&old, "from", "github.com/artfoxe6/quick-gin", "原始项目名，默认 github.com/artfoxe6/quick-gin")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setup(cmd *cobra.Command, args []string) {
	if name == "" {
		if len(args) != 1 {
			fmt.Println("请提供待初始化的项目名, 请使用 ./quick-cli blog 或 ./quick-cli --name=blog")
			return
		}
		name = args[0]
	}

	// 删除git文件夹
	bash := exec.Command("rm", "-rf", ".git")
	_, err := bash.Output()
	if err != nil {
		fmt.Printf("err: rm -rf .git: %v\n", err)
		return
	}

	// 删除git文件夹
	bash = exec.Command("cp", "config/config.ini.example", "config/config.ini")
	_, err = bash.Output()
	if err != nil {
		fmt.Printf("err: copy ini: %v\n", err)
		return
	}

	// 获取当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current directory: %v\n", err)
		return
	}

	// 替换项目名字
	replaceHelloWithWord(currentDir)

	// 删除git文件夹
	bash = exec.Command("rm", "./quick-cli")
	_, err = bash.Output()
	if err != nil {
		fmt.Printf("err: rm quick-cli: %v\n", err)
		return
	}

	// 删除quick-cli
	_ = os.Remove(currentDir + "quick-cli")
	_ = os.Remove(currentDir + "quick-cli-mac")
	_ = os.Remove(currentDir + "quick-cli.exe")

	// 修改当前目录的名称
	newSuffix := strings.Split(name, "/")
	oldSuffix := strings.Split(currentDir, "/")
	newDir := strings.ReplaceAll(currentDir, oldSuffix[len(oldSuffix)-1], newSuffix[len(newSuffix)-1])
	err = os.Rename(currentDir, newDir)
	if err != nil {
		fmt.Println("无法修改目录名称:", err)
		return
	}

	fmt.Println("如果发现当前目录名没有变化，请执行 cd ../" + newSuffix[len(newSuffix)-1])

	fmt.Println("最后一步，请手动修改 config/config.ini 中的数据库参数")

}

func replaceHelloWithWord(rootDir string) {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		err = replaceStringInFile(path)
		if err != nil {
			//fmt.Printf("Failed to replace strings in file %s: %v\n", path, err)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", rootDir, err)
	}
}

func replaceStringInFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 替换字符串
	newData := strings.ReplaceAll(string(data), old, name)

	oldSuffix := strings.Split(old, "/")
	newSuffix := strings.Split(name, "/")
	newData = strings.ReplaceAll(string(newData), oldSuffix[len(oldSuffix)-1], newSuffix[len(newSuffix)-1])

	// 将替换后的数据写回文件
	err = os.WriteFile(filePath, []byte(newData), 0)
	if err != nil {
		return err
	}

	return nil
}
