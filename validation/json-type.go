package validation

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func ImperfectJsonPatch(resp *http.Response) *http.Response {
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp
	}
	resp.Body.Close()
	m := map[string]interface{}{}
	err = json.Unmarshal(rawBody, &m)
	if err != nil {
		rawBody = changeJson(rawBody)
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	return resp
}

func changeJson(body []byte) []byte {
	past := '\a'

deleteCnt:
	for num := 0; num < 5; num++ {
		for i, v := range body {
			if past == ':' && v == ',' {
				body = deleteJsonRow(body, i-1)
				continue deleteCnt
			}
			past = rune(v)
		}
		break
	}
	return body
}

func deleteJsonRow(body []byte, deleteIndex int) []byte {
	if deleteIndex < 0 || len(body) < deleteIndex {
		return body
	}

	startIndex := -1
	endIndex := -1

	for i := deleteIndex; 0 <= i; i-- {
		if body[i] == ',' {
			startIndex = i
			break
		}
		if i == 0 {
			startIndex = 0
			break
		}
	}
	for i := deleteIndex; i < len(body); i++ {
		if body[i] == ',' {
			endIndex = i
			break
		}
	}

	if startIndex < 0 {
		return body
	}
	if endIndex < 0 {
		return body
	}

	body = append(body[:startIndex], body[endIndex:]...)
	return body
}
