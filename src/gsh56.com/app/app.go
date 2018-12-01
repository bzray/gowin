package app

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

func init() {
	info, err := os.Stat(filepath)
	if err == nil && info.Size() > 0 {
		reader()
	} else {
		//GetWordHttp()
		GetWordFile()
	}
}

type Chinese struct {
	Word   string
	Pinyin string
	Desc   string
	Index  int
}

var filepath, temppath = "./w.bin", "./x.bin"
var Words, Pinyin = make(map[int]string), make(map[int]string)
var Data []Chinese
var Quit []int

func reader() {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	idx := 0
	for {
		b, _, err := buf.ReadLine()
		//line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("加载文件文字 %v 个", len(Data))
				return
			}
			log.Fatal(err)
		}

		line := string(b)

		d := strings.Split(line, ";")

		var ch Chinese
		if len(d) > 0 {
			ch = Chinese{
				Word: d[0],
			}
		} else {
			continue
		}

		if len(d) > 1 {
			ch.Pinyin = d[1]
		}
		if len(d) > 2 {
			ch.Desc = d[2]
		}
		ch.Index = idx
		idx++
		Data = append(Data, ch)
	}

}

func writer(data []Chinese) {

	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, v := range data {
		f.WriteString(v.Word + ";" + v.Pinyin + ";" + v.Desc + "\n")
	}
	f.Sync()
	log.Printf("写入姓氏 %v 个", len(Data))
}

func clear() {
	Quit = append([]int{})
	Data = append([]Chinese{})
	Words, Pinyin = make(map[int]string), make(map[int]string)
}

func GetWordHttp() {

	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}

	resp, err := client.Get("http://www.wxwww.cn/index.php/home/MyArticle/baijiaxing")
	if err != nil {
		log.Printf("error:%v", err)
	}

	defer resp.Body.Close()

	dom, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Printf("error:%v", err)
	}

	st := dom.Find(".zy-content").First()
	if st != nil {
		clear()

		idx := 0
		st.Find(".zy-hanzi > td").Each(func(i int, selection *goquery.Selection) {
			Words[i] = selection.Text()
			idx = i
		})

		st.Find(".zy-pinyin > td").Each(func(i int, selection *goquery.Selection) {
			Pinyin[i] = selection.Text()
		})

		//for k,v := range Words{
		//	if v!="。" && v!="，"{
		//		Data = append(Data, Chinese{
		//			Word:v,
		//			Pinyin:Pinyin[k],
		//			Index:k,
		//		})
		//	}
		//}

		for n := 0; n < idx; n++ {
			v := Words[n]
			if v != "，" && v != "。" {
				Data = append(Data, Chinese{
					Word:   v,
					Pinyin: Pinyin[n],
					Index:  n,
				})
			}
		}

		log.Printf("请求姓氏 %v 个", len(Data))

		writer(Data)
	}
}

func GetWordFile() {
	f, err := os.Open(temppath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	clear()

	buf := bufio.NewReader(f)
	idx := 0
	for {
		b, _, err := buf.ReadLine()
		//line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		line := string(b)

		d := strings.Split(line, "：")

		var ch Chinese
		if len(d) > 0 {
			v := strings.Split(d[0], ",")
			ch = Chinese{
				Word:   v[0],
				Pinyin: v[1],
			}
			if len(d) > 1 {
				ch.Desc = d[1]
			}
		} else {
			continue
		}

		ch.Index = idx
		idx++
		Data = append(Data, ch)
	}
	log.Printf("加载姓氏 %v 个", len(Data))
	writer(Data)

}

func GetScroll(channel chan Chinese, closeChannel chan bool) {
	size := len(Data)
	rand.Seed(time.Now().UnixNano())

	for {
		t := rand.Intn(size)

		if t >= 0 && t < size {
			select {
				case channel <- Data[t]:
					String(Data[t])
				case <-closeChannel:
					close(channel)
					log.Println("发送结束")
					return
			}

		}
	}
}

func GetRandom() Chinese{
	size := len(Data)
	rand.Seed(time.Now().UnixNano())

	cnt := 0
loop:
	for {
		t := rand.Intn(size)

		for _, v := range Quit {
			if v == t {
				if cnt > 600 {
					log.Fatal("已无数据循环,请检查数据")
				}
				cnt++
				goto loop
			}
		}
		cnt = 0

		if t >= 0 && t < size {
			Quit = append(Quit, t)
			String(Data[t])
			return Data[t]
		}
	}
}

func String(c Chinese){
	fmt.Printf("姓氏: %v(%v) \n解释:%v\n", c.Word, c.Pinyin, c.Desc)
}
