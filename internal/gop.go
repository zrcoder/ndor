package internal

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type InputLine struct {
	Line int
	Code string
}

type Input []InputLine

func (i Input) String() string {
	buf := strings.Builder{}
	for _, v := range i {
		buf.WriteString(v.Code)
		buf.WriteByte('\n')
	}
	return buf.String()
}

func regularInput(code string) Input {
	arr := strings.Split(code, "\n")
	res := make([]InputLine, 0, len(arr))
	for i, v := range arr {
		vt := strings.TrimSpace(v)
		if vt == "" {
			continue
		}
		line := InputLine{Line: i}
		if isOper(vt) {
			line.Code = addPrintPrefix(v)
		} else {
			line.Code = v
		}
		res = append(res, line)
	}
	return res
}

func isOper(s string) bool {
	i := strings.Index(s, " ")
	if i == -1 {
		i = strings.Index(s, "(")
	}
	if i == -1 {
		i = len(s)
	}
	_, ok := operations[s[:i]]
	return ok
}

func addPrintPrefix(s string) string {
	i := 0
	for ; i < len(s) && s[i] == ' ' || s[i] == '\t'; i++ {
	}
	return s[:i] + "print " + s[i:]
}

func gopExecute(input Input) (string, error) {
	inputData := url.Values{}
	inputData.Add("body", input.String())
	resp, err := http.DefaultClient.Post(
		"https://play.goplus.org/compile",
		"application/x-www-form-urlencoded; charset=UTF-8",
		strings.NewReader(inputData.Encode()),
	)
	if err != nil {
		log.Println("[test]", err)
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	gopResp := &struct{ Events []struct{ Message string } }{}
	err = decoder.Decode(gopResp)
	if err != nil {
		log.Println("[test]", err)
		return "", err
	}
	if len(gopResp.Events) == 0 || len(gopResp.Events[0].Message) == 0 {
		log.Println("[test]")
		return "", errors.New("empty response from https://play.goplus.org")
	}
	return gopResp.Events[0].Message, nil
}
