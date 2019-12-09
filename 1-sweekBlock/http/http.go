package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)



func PostMsg(msg, ip string) []byte {
	cli := &http.Client{
		Timeout: 200 * time.Second,
	}

	body := bytes.NewBuffer([]byte(msg))
	resq, err := http.NewRequest("POST", ip, body)
	if err != nil {
		fmt.Println("resq faild")
	}
	resq.Close = true

	rep, err := cli.Do(resq)
	if err != nil {
		fmt.Println("get response faild!")
	}

	defer rep.Body.Close()
	rcv, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		fmt.Println("get body faild")
	}

	return rcv
}
