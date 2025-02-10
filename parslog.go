package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	url         string = os.Args[1]
	ChatID      string = os.Args[2]
	GetMe       string = "/getMe"
	SendMessage string = "/sendMessage"
)

const (
	file     string = "tets.txt"
	ErrorLog string = "ERROR"
)

func main() {
	ReadLog()
}

func TimeNow() string {
	var tm string = time.Now().Format("15:04:05")
	time.Sleep(1 * time.Second)
	return tm
}

func Hostname() string {
	hostnm := exec.Command("hostname")
	stdout, _ := hostnm.Output()
	return string(stdout)
}

func HealthcheckBot() string {
	resp, err := http.Get(url + GetMe)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}
	return resp.Status
}

func PostRequest(text string) *http.Response {
	hostname := Hostname()
	var TextAndTime string = TimeNow() + ":" + hostname + " - " + text
	var PostUrl string = url + SendMessage
	var TextJson string = fmt.Sprintf(`{"chat_id": "%s", "text": "%s"}`, ChatID, TextAndTime)
	txt := []byte(TextJson)
	resp, err := http.Post(PostUrl, "application/json", bytes.NewBuffer(txt))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(resp)
	}
	return resp
}

func ReadLog() {
	file, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data := make([]byte, 500)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			continue
		}
		var strsplit []string = strings.Split(string(data[:n]), "\n")
		for i := 0; i < len(strsplit); i++ {
			if strings.Contains(strsplit[i], ErrorLog) {
				rr := PostRequest(strsplit[i])
				if rr.StatusCode != 200 {
					panic(rr.StatusCode)
				}
			}
		}
	}
}
