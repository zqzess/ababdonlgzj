package console

import (
	"abandonlgzj/config"
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Console struct {
	wg sync.WaitGroup
}

var (
	instance *Console
	once     sync.Once
)

func GetInstance() *Console {
	once.Do(func() {
		instance = &Console{}
	})
	return instance
}

var sp = config.GetInstance()

func (c *Console) Console() {
	// 创建一个新的扫描器来读取控制台输入
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("请输入命令('exit' 退出)： \n1. 重载配置文件")

	go func() {
		// 循环监听控制台输入
		for {
			// 读取一行输入
			if scanner.Scan() {
				input := scanner.Text()

				// 处理输入
				switch strings.TrimSpace(input) {
				case "1":
					sp.ReloadConfig()
				case "exit":
					fmt.Println("bye~")
					return
				default:
					fmt.Printf("未知命令: %s\n", input)
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "读取输入时出错:", err)
				return
			}
		}
	}()
}

func (c *Console) ConsoleInput(commands chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("请输入您的选择 （输入 exit 退出）")
	for {
		fmt.Print("1.重载\n")
		if scanner.Scan() {
			input := scanner.Text()
			if input == "exit" {
				fmt.Println("bye~~")
				close(commands) // 关闭命令通道
				return
			} else if input == "1" {
				commands <- "reload" // 发送 reload 命令到通道
			} else {
				fmt.Println("未知输入，请重试。")
			}
		}
	}
}

// CommandHandler 处理来自命令行的命令
func (c *Console) CommandHandler(commands chan string) {
	defer c.wg.Done() // 在这里确保 WaitGroup 的计数减一

	c.wg.Add(1)

	for command := range commands {
		switch command {
		case "reload":
			sp.ReloadConfig()
		default:
			fmt.Println("未知命令:", command)
		}
	}
}
