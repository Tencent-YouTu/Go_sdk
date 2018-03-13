/*
* File Name:	detectface.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2015-08-25
 */
package main

import (
	"fmt"
	"github.com/Tencent-YouTu/Go_sdk"
	"io/ioutil"
	"os"
)

func main() {
	//Register your app on http://open.youtu.qq.com
	//Get the following details
	appID := uint32(0)
	secretID := ""
	secretKey := ""
	userID := ""

	as, err := youtu.NewAppSign(appID, secretID, secretKey, userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewAppSign() failed: %s\n", err)
		return
	}
	imgData, err := ioutil.ReadFile("../../testdata/imageA.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ReadFile() failed: %s\n", err)
		return
	}
	//yt := youtu.Init(as, youtu.TencentYunHost)
	yt := youtu.Init(as, youtu.DefaultHost)
	df, err := yt.DetectFace(imgData, false, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DetectFace() failed: %s", err)
		return
	}
	fmt.Printf("df: %#v\n", df)
}
