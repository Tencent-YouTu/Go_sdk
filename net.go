/*
* File Name:	net.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2015-08-25
 */
package youtu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
    "errors"
)

func (y *Youtu) interfaceURL(ifname string, urltype int) string {
    if urltype == 0 {
	    return fmt.Sprintf("%s/youtu/api/%s", y.host, ifname)
    } else if urltype == 1{
        return fmt.Sprintf("%s/youtu/imageapi/%s", y.host, ifname)
    } else {
        return fmt.Sprintf("%s/youtu/ocrapi/%s", y.host, ifname)
    }
}

func (y *Youtu) interfaceRequest(ifname string, req, rsp interface{}, urltype int) (err error) {
	url := y.interfaceURL(ifname, urltype)
	if y.debug {
		fmt.Printf("req: %#v\n", req)
	}
	data, err := json.Marshal(req)
	if err != nil {
		return
	}
    body, err := y.get(url, string(data))
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		if y.debug {
			fmt.Fprintf(os.Stderr, "body:%s\n", string(body))
		}
		return fmt.Errorf("json.Unmarshal() rsp: %#v failed: %s\n", rsp, err)
	}
	return
}

func (y *Youtu) get(addr string, req string) (rsp []byte, err error) {
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	httpreq, err := http.NewRequest("POST", addr, strings.NewReader(req))
	if err != nil {
		return
	}
	auth := y.sign()
	if y.debug {
		fmt.Fprintf(os.Stderr, "Authorization: %s\n", auth)
	}
	httpreq.Header.Add("Authorization", auth)
	httpreq.Header.Add("Content-Type", "text/json")
	httpreq.Header.Add("User-Agent", "")
	httpreq.Header.Add("Accept", "*/*")
	httpreq.Header.Add("Expect", "100-continue")
	resp, err := client.Do(httpreq)
	
    if err != nil {
		return
	}
    
    if resp.StatusCode != 200 {
        errStr := fmt.Sprintf("httperrorcode: %d \n", resp.StatusCode)    
        err = errors.New(errStr)
        return 
    }
    
	defer resp.Body.Close()
	rsp, err = ioutil.ReadAll(resp.Body)
	return
}
