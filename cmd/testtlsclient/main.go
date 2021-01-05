/*
 * Copyright (c) 2021 wuqingtao
 * testtlsclient is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 *
 * Author: wuqingtao (wqt_1110@qq.com)
 * CreateTime: 2021-01-05 22:26:12+0800
 */

package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	c := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	u := "https://127.0.0.1:8123/abc"
	resp, err := c.Get(u)
	if err != nil {
		log.Printf("Get(%s) %s", u, err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll %s", err)
	}
	log.Printf("response data %s", b)
}
