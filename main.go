package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var (
	help   bool
	listen string
	tmpDir string
	cosDir string
)

type Response struct {
	Url string `json:"url"`
}

func init() {
	flag.BoolVar(&help, "help", false, "显示帮助")
	flag.StringVar(&listen, "http", "127.0.0.1:8016", "监听地址")
	flag.StringVar(&tmpDir, "tmp", "/tmp/qcloud-cos-tmpdir", "临时目录")
	flag.StringVar(&cosDir, "cosdir", "/blog/", "COS目录")
}

func cosUpload(filename string) (response Response, e error) {
	basename := filepath.Base(filename)

	// remove right '/' or space
	remoteDir := strings.TrimRightFunc(cosDir, func(r rune) bool {
		return r == '/' || unicode.IsSpace(r)
	})

	// add left '/' if not has
	if !strings.HasPrefix(remoteDir, "/") {
		remoteDir = "/" + remoteDir
	}
	cosfile := fmt.Sprintf(`%s/%s`, remoteDir, basename)

	// 执行命令
	var output bytes.Buffer
	cmd := exec.Command("coscmd", "-d", "upload", filename, cosfile)
	cmd.Stdout = &output
	cmd.Stderr = &output
	e = cmd.Run()
	if e != nil {
		return
	}

	// 使用正则解析
	re, e := regexp.Compile(`(https://.*\.myqcloud\.com):443 "PUT (/[^ ]*) HTTP/1.1" 200 0`)
	if e != nil {
		return
	}
	match := re.FindStringSubmatch(output.String())
	if len(match) < 3 {
		e = fmt.Errorf(`cannot find Url from output: %s`, output.String())
		return
	}
	url := fmt.Sprintf(`%s%s`, match[1], match[2])

	// 打印日志
	log.Println(`upload success:`, url)
	return Response{Url: url}, nil
}

func Upload(w http.ResponseWriter, r *http.Request) {

	// 获取上传的文件
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `cannot get upload file:`+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 写入文件
	uploadFile := fmt.Sprintf(`%s/%s`, tmpDir, header.Filename)
	writer, e := os.Create(uploadFile)
	if e != nil {
		http.Error(w, `cannot create file:`+uploadFile+e.Error(), http.StatusInternalServerError)
		return
	}
	defer writer.Close()
	io.Copy(writer, file)

	// 调用coscmd上传文件
	response, e := cosUpload(uploadFile)
	if e != nil {
		http.Error(w, `cosUpload error:`+e.Error(), http.StatusInternalServerError)
		return
	}

	// 返回JSON字符串
	marshal, e := json.Marshal(response)
	if e != nil {
		http.Error(w, `json encode error:`+e.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(marshal)
	return
}

func init() {
	e := os.MkdirAll(tmpDir, 0755)
	if e != nil {
		log.Fatalln(e)
	}
	http.HandleFunc("/upload", Upload)
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	server := http.Server{Addr: listen}
	go func() {
		e := server.ListenAndServe()
		if e != nil {
			log.Fatalln(`serve error:`, e)
		}
	}()
	log.Println(`Server already started...`)
	select {}
}
