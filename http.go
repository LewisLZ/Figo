package Figo

import (
	"github.com/quexer/utee"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func DownLoad(localFileName, remoteFileName string, retryTimes int) error {
	if retryTimes <= 0 {
		retryTimes = 1
	}
	var err error
	for i := 0; i < retryTimes; i++ {
		if err = download(localFileName, remoteFileName); err == nil {
			return err
		} else if i < retryTimes-1 {
			log.Println("Retry 5 Sec Later")
			time.Sleep(time.Second * time.Duration(5))
		}
	}
	return err
}

func download(localFileName, remoteFileName string) error {
	if fileInfo, _ := os.Stat(localFileName); fileInfo == nil {
		os.Create(localFileName)
	}
	f, err := os.OpenFile(localFileName, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	stat, err := f.Stat() //获取文件状态
	if err != nil {
		return err
	}
	f.Seek(stat.Size(), 0)
	var req http.Request
	req.Method = "GET"
	req.Close = true
	if req.URL, err = url.Parse(remoteFileName); err != nil {
		return err
	}
	header := http.Header{}
	header.Set("Range", "bytes="+strconv.FormatInt(stat.Size(), 10)+"-")
	req.Header = header
	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		return err
	}
	written, err := io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	println("written: ", written)
	return nil
}

func HttpForeignHelp(w http.ResponseWriter, r *http.Request) {
	api := getHeadValue("url", r.Header)
	buf := make([]byte, 10*1024)
	r.Body.Read(buf)
	body := ioutil.NopCloser(r.Body)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", api, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	resp, err := client.Do(req)
	utee.Chk(err)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	utee.Chk(err)
	log.Println(string(data), err)
	w.Write(data)
}
