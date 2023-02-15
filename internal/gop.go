package internal

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	//go:embed tpl/tpl.go
	preCodes string
)

type LineError struct {
	Msg    string
	Number int
}

func (e *LineError) Error() string {
	return fmt.Sprintf("line %d: %s", e.Number, e.Msg)
}

// TODO: change to igop, not request play.goplus.org
func gopExecute(code string, preLines int) (string, error) {
	log.Println("gop execute code")
	inputData := url.Values{}
	inputData.Add("body", code)
	resp, err := http.DefaultClient.Post(
		"https://play.goplus.org/compile",
		"application/x-www-form-urlencoded; charset=UTF-8",
		strings.NewReader(inputData.Encode()),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	gopResp := &struct {
		Errors string                     `json:"Errors"`
		Events []struct{ Message string } `json:"Events"`
	}{}
	err = decoder.Decode(gopResp)
	if err != nil {
		return "", err
	}
	if gopResp.Errors != "" {
		return "", parseGopPlayError(gopResp.Errors, preLines)
	}
	if len(gopResp.Events) == 0 || len(gopResp.Events[0].Message) == 0 {
		return "", errors.New("empty response from https://play.goplus.org")
	}
	return gopResp.Events[0].Message, nil
}

func parseGopPlayError(ori string, preLines int) *LineError {
	log.Println("begin to parse goplay error:", ori)
	log.Println("pre lines:", preLines)
	res := &LineError{Number: -1, Msg: ori}
	sep := ".gop:"
	i := strings.Index(ori, sep)
	if i == -1 {
		sep = ".go:"
		i = strings.Index(ori, sep)
	}
	if i == -1 {
		return res
	}
	msg := ori[i+len(sep):]
	i = strings.Index(msg, "\n")
	if i == -1 {
		i = strings.Index(msg, `\n`)
	}
	if i != -1 {
		msg = msg[:i]
	}
	res.Msg = msg
	sep = ": "
	i = strings.Index(msg, sep)
	if i == -1 {
		return res
	}
	res.Msg = msg[i+len(sep):]
	msg = msg[:i]
	i = strings.LastIndex(msg, ":")
	if i != -1 {
		msg = msg[:i]
	}
	n, err := strconv.Atoi(msg)
	if err != nil {
		return res
	}
	res.Number = n - preLines
	return res
}
