package main

import (
	"fmt"
	"gsh56.com/app"
	"log"
	"os"
	"time"
)

func main() {

	if len(os.Args) >0 && os.Args[0]=="d"{
		exec(DOS)
	}else{
		exec(GUI)
	}

	
	defer func() {
		if err := recover(); err != nil {
			name := "error.txt"
			_, err := os.Stat(name)

			if err == nil {
				file, _ := os.Create(name)
				file.WriteString(errLog(err))
				file.Close()
			} else {
				file, _ := os.Open(name)
				file.WriteString(errLog(err))
				file.Close()
			}
		}
	}()
}

const (
	DOS = iota
	GUI
)

func exec(tp int) {
	if tp == DOS {
		app.Dos()
	} else if tp == GUI {
		app.Gui()
	} else {
		log.Fatal("参数错误")
	}
}

func errLog(err error) string {
	log.Printf("%v\n", err)
	return fmt.Sprintf("%v:%v\n", time.Now().Format("2006-01-02 15:04:05"), err)
}
