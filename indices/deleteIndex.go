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

package indices

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
)

// The delete APi allows you to delete one or more indices through an API. This operation may fail
// if the elasitsearch configuration has been set to forbid deleting indexes.
func Delete(index string) (api.BaseResponse, error) {
	var url string
	var retval api.BaseResponse

	if len(index) > 0 {
		url = fmt.Sprintf("/%s", index)
	} else {
		return retval, fmt.Errorf("You must specify at least one index to delete")
	}

	body, err := api.DoCommand("DELETE", url, nil)
	if err != nil {
		return retval, err
	}

	jsonErr := json.Unmarshal(body, &retval)
	if jsonErr != nil {
		return retval, jsonErr
	}

	return retval, err
}
