package com

import (
	"encoding/base64"
	"eve-client/cli"
	"eve-client/inp"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var HTTP_HOST = "http://localhost:8000"

func GetReq(url string) []byte {
	fmt.Println(url)
	req := fiber.Get(url)
	_, byte, err := req.Bytes()
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second * 2)
	}

	return byte
}

func PostReq(url string) []byte {
	req := fiber.Post(url)
	_, byte, err := req.Bytes()
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second * 2)
	}

	return byte
}

func DeleteReq(url string) []byte {
	req := fiber.Delete(url)
	_, byte, err := req.Bytes()
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second * 2)
	}

	return byte
}

var file = cli.CURRENT_MENU.Header

func ExecFile() {
	file = cli.CURRENT_MENU.Header
	data := PostReq(HTTP_HOST + "/exec/" + file)
	fmt.Println(string(data))
	for {
		in := inp.Inp()
		if in == inp.Back || in == inp.Left {
			return
		}
	}
}

var fileFuncs = []cli.Opt{
	{
		Name: "Execute",
		Func: ExecFile,
	},
	{
		Name: "Read File",
		Func: ReadFile,
	},
	{
		Name: "Delete File",
		Next: cli.Menu{
			Header: "Are you sure?",
			Opts: []cli.Opt{
				{
					Name: "Yes",
					Func: DeleteFile,
				},
			},
		},
	},
}

func GetFileList() []cli.Opt {
	file = cli.CURRENT_MENU.Header
	data := GetReq(HTTP_HOST + "/files")
	rawOpts := strings.Split(strings.TrimSpace(string(data)), " ")
	var opts []cli.Opt
	for _, opt := range rawOpts {
		opts = append(opts, cli.Opt{Name: opt, Next: cli.Menu{Header: opt, Opts: fileFuncs}})
	}
	return opts
}

func SendFile() {
	file = cli.CURRENT_MENU.Header
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	str := base64.URLEncoding.EncodeToString(data)
	out := PostReq(HTTP_HOST + "/files/" + file + "/" + str)
	fmt.Println(string(out))
}

func DeleteFile() {
	file = cli.CURRENT_MENU.Header
	out := DeleteReq(HTTP_HOST + "/files/" + file)
	fmt.Println(string(out))
}

func ReadFile() {
	file = cli.CURRENT_MENU.Header
	out := GetReq(HTTP_HOST + "/files/" + file)
	fmt.Println(string(out))
	for {
		in := inp.Inp()
		if in == inp.Back || in == inp.Left {
			return
		}
	}
}
