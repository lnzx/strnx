package internal

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/lnzx/strnx/tools"
	"io"
	"log"
	"net/http"
	"os"
)

var lastVersion = 835
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
		err = sendSms()
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

func sendSms() error {
	body, err := json.Marshal(map[string]string{
		"mobile": mobile,
		"event":  "l1-node有新版本!",
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
