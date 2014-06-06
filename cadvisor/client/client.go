// Copyright 2014 Google Inc. All Rights Reserved.
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

package cadvisor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/lmctfy/cadvisor/info"
)

type Client struct {
	baseUrl string
}

func NewClient(URL string) (*Client, error) {
	c := &Client{
		baseUrl: strings.Join([]string{
			URL,
			"api/v1.0",
		}, "/"),
	}
	_, err := c.MachineInfo()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (self *Client) machineInfoUrl() string {
	return strings.Join([]string{self.baseUrl, "machine"}, "/")
}

func (self *Client) MachineInfo() (minfo *info.MachineInfo, err error) {
	u := self.machineInfoUrl()
	ret := new(info.MachineInfo)
	err = self.httpGetJsonData(ret, u, "machine info")
	if err != nil {
		return
	}
	minfo = ret
	return
}

func (self *Client) containerInfoUrl(name string) string {
	if name[0] == '/' {
		name = name[1:]
	}
	return strings.Join([]string{self.baseUrl, "container", name}, "/")
}

func (self *Client) httpGetJsonData(data interface{}, url, infoName string) error {
	resp, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("unable to get %v: %v", infoName, err)
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		err = fmt.Errorf("unable to decode %v: %v", infoName, err)
		return err
	}
	return nil
}

func (self *Client) ContainerInfo(name string) (cinfo *info.ContainerInfo, err error) {
	u := self.containerInfoUrl(name)
	ret := new(info.ContainerInfo)
	err = self.httpGetJsonData(ret, u, fmt.Sprintf("container info for %v", name))
	if err != nil {
		return
	}
	cinfo = ret
	return
}
