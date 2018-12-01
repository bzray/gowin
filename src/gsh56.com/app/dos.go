package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Dos() {
	print()

	buf := bufio.NewScanner(os.Stdin)
	for buf.Scan() {
		t := strings.TrimSpace(buf.Text())
		if t == "" || t == "1" {
			GetRandom()
			//}else if t=="2" {
			//GetWordHttp()
			//GetWordFile()
		} else if t == "2" {
			Quit = append([]int{})
		} else if t == "3" {
			os.Exit(0)
		}
		print()
	}
}

func print() {
	fmt.Print("\n")
	fmt.Println("-------------------------------------")
	fmt.Println("------------姓氏筛选 V1.0--------------")
	fmt.Println("---1、随机选姓。2、清空已选。3、退出系统---")
	fmt.Println("-------------------------------------")
}
