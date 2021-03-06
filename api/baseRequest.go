// Copyright 2012 Matthew Baird
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	//"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
)

func DoCommand(method string, url string, data interface{}) ([]byte, error) {
	var response map[string]interface{}
	var body []byte
	var httpStatusCode int
	req, err := ElasticSearchRequest(method, url)
	//log.Println(req.URL)
	if err != nil {
		return body, err
	}

	if data != nil {
		switch v := data.(type) {
		case string:
			req.SetBodyString(v)
		case io.Reader:
			req.SetBody(v)
		//case *bytes.Buffer:
		//	req.SetBody(v)
		default:
			err = req.SetBodyJson(v)
			if err != nil {
				return body, err
			}
		}

	}
	httpStatusCode, body, err = req.Do(&response)

	if err != nil {
		return body, err
	}
	if httpStatusCode > 304 {

		jsonErr := json.Unmarshal(body, &response)
		if jsonErr == nil {
			if error, ok := response["error"]; ok {
				status, _ := response["status"]
				return body, errors.New(fmt.Sprintf("Error [%s] Status [%v]", error, status))
			}
		}
		return body, jsonErr
	}
	return body, nil
}

// The API also allows to check for the existance of a document using HEAD
// This appears to be broken in the current version of elasticsearch 0.19.10, currently
// returning nothing
func Exists(pretty bool, index string, _type string, id string) (BaseResponse, error) {
	var response map[string]interface{}
	var body []byte
	var url string
	var retval BaseResponse
	var httpStatusCode int

	if len(_type) > 0 {
		url = fmt.Sprintf("/%s/%s/%s?%s", index, _type, id, Pretty(pretty))
	} else {
		url = fmt.Sprintf("/%s/%s?%s", index, id, Pretty(pretty))
	}
	req, err := ElasticSearchRequest("HEAD", url)
	if err != nil {
		// some sort of generic error handler
	}
	httpStatusCode, body, err = req.Do(&response)
	if httpStatusCode > 304 {
		if error, ok := response["error"]; ok {
			status, _ := response["status"]
			log.Println("Error: %v (%v)\n", error, status)
		}
	} else {
		// marshall into json
		jsonErr := json.Unmarshal(body, &retval)
		if jsonErr != nil {
			log.Println(jsonErr)
		}
	}
	//fmt.Println(string(body))
	return retval, err
}
