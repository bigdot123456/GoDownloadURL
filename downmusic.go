package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	//ref:="http://f2.htqyy.com/play7/597/mp3/2"
	headstr := "http://f2.htqyy.com/play7/"
	tailstr := "/mp3/2"

	Mydir := "Music"
	TrueDir := CreateDateDir(Mydir)
	var URLlist []string
	var musicURL string
	start := 1000
	end := 2000
	for ; start <= end; start++ {
		musicURL = headstr + strconv.Itoa(start) + tailstr
		download1(musicURL, start, TrueDir)
		URLlist = append(URLlist, musicURL)
	}
	fmt.Println(URLlist)
}
func download1(musicURL string, index int, foldername string) {

	//var durl = "https://dl.google.com/go/go1.10.3.darwin-amd64.pkg";

	_, err := url.ParseRequestURI(musicURL)
	if err != nil {
		panic("网址错误")
	}

	//filename := path.Base(uri.Path)
	filename := foldername + "/" + strconv.Itoa(index) + ".mp3"
	log.Println("[*] Filename " + filename)

	nt := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s]To download %s\n", nt, filename)

	newFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer newFile.Close()

	//client := http.Client{Timeout: 900 * time.Second}
	//resp, err := client.Get(musicURL)
	client := http.DefaultClient
	client.Timeout = time.Second * 60 //设置超时时间

	var resp, err1 = client.Get(musicURL)
	if err1 != nil {
		panic(err1)
	}
	defer resp.Body.Close()

	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
}



func download(musicURL string, index int, foldername string) {

	//var durl = "https://dl.google.com/go/go1.10.3.darwin-amd64.pkg";

	_, err := url.ParseRequestURI(musicURL)
	if err != nil {
		panic("网址错误")
	}

	//filename := path.Base(uri.Path)
	filename := foldername + "/" + strconv.Itoa(index) + ".mp3"
	log.Println("[*] Filename " + filename)

	client := http.DefaultClient
	client.Timeout = time.Second * 60 //设置超时时间
	resp, err := client.Get(musicURL)
	if err != nil {
		panic(err)
	}
	if resp.ContentLength <= 0 {
		log.Println("[*] Destination server does not support breakpoint download.")
	}
	raw := resp.Body
	defer raw.Close()
	reader := bufio.NewReaderSize(raw, 1024*32)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)

	buff := make([]byte, 32*1024)
	written := 0
	go func() {
		for {
			nr, er := reader.Read(buff)

			if nr > 0 {
				nw, ew := writer.Write(buff[0:nr])
				if nw > 0 {
					written += nw
				}
				if ew != nil {
					err = ew
					break
				}
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			// directly ignore the following message!
			if er != nil && er.Error() == "use of closed network connection" {
				break
			}

			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
		}

		if err != nil {
			panic(err)
		}
	}()

	spaceTime := time.Second * 1
	ticker := time.NewTicker(spaceTime)
	lastWtn := 0
	stop := false

	for {
		select {
		case <-ticker.C:
			speed := written - lastWtn
			log.Printf("[*] Speed %s / %s \n", bytesToSize(speed), spaceTime.String())
			if written-lastWtn == 0 {
				ticker.Stop()
				stop = true
				break
			}
			lastWtn = written
		}
		if stop {
			break
		}
	}
}

func bytesToSize(length int) string {
	var k = 1024 // or 1024
	var sizes = []string{"Bytes", "KB", "MB", "GB", "TB"}
	if length == 0 {
		return "0 Bytes"
	}
	i := math.Floor(math.Log(float64(length)) / math.Log(float64(k)))
	r := float64(length) / math.Pow(float64(k), i)
	return strconv.FormatFloat(r, 'f', 3, 64) + " " + sizes[int(i)]
}

// CreateDateDir 根据当前日期来创建文件夹
func CreateDateDir(Path string) string {
	folderName := time.Now().Format("20060102")
	folderPath := filepath.Join(Path, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		err = os.MkdirAll(folderPath, 0777) //0777也可以os.ModePerm
		if err != nil {
			panic("can't mkdir" + folderPath)
		}
		_ = os.Chmod(folderPath, 0777)
	}
	return folderPath
}
