package internal

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/lnzx/strnx/tools"
	"io"
	"log"
	"net/http"
	"os"
)

var lastVersion = 883
var smsApiKey = os.Getenv("SMS_API_KEY")
var mobile = os.Getenv("MOBILE")

const (
	versionUrl = "https://orchestrator.strn.pl/requirements"
	notifyUrl  = "https://ntfyx.fly.dev/sms"
)

func CheckVersionJob() {
	rsp, err := tools.Get(versionUrl)
	if err != nil {
		log.Println(err)
		if rsp != nil {
			err = rsp.Body.Close()
			if err != nil {
				return
			}
		}
		return
	}
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var version Version
	if err = json.Unmarshal(body, &version); err != nil {
		log.Println(err)
		return
	}
	log.Println("cron found lastVersion:", version.LastVersion)
	if version.LastVersion > lastVersion {
		// 内容不能超过12个字符 Node更新:999
		msg := fmt.Sprintf("Node更新:%d", version.LastVersion)
		err = sendSms(msg)
		if err != nil {
			log.Println(err)
			return
		}
		lastVersion = version.LastVersion
		log.Println("send sms ok")
	}
}

type Version struct {
	LastVersion int `json:"lastVersion"`
	MinVersion  int `json:"minVersion"`
}

func sendSms(content string) error {
	body, err := json.Marshal(map[string]string{
		"mobile": mobile,
		"event":  content,
	})

	req, err := http.NewRequest("POST", notifyUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", smsApiKey)
	req.Header.Set("Content-Type", "application/json")

	rsp, err := tools.Do(req)
	if err != nil {
		if rsp != nil {
			err = rsp.Body.Close()
			if err != nil {
				return err
			}
		}
		return err
	}
	return nil
}
