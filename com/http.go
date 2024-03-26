package com

import (
	"encoding/base64"
	"eve-client/cli"
	"eve-client/inp"
	"eve-client/log"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var LOG log.Logger = log.Logger{
	Emoji: "📶",
}

var (
	HTTP_HOST = "http://10.0.0.111:8000"
	baseDir   = "./files/"
)

func GetReq(url string) []byte {
	req := fiber.Get(url)
	_, byte, err := req.Bytes()
	if err != nil {
		LOG.Error(err)
		time.Sleep(time.Second * 2)
	}

	return byte
}

func PostReq(url string) []byte {
	req := fiber.Post(url)
	_, byte, err := req.Bytes()
	if err != nil {
		LOG.Error(err)
		time.Sleep(time.Second * 2)
	}

	return byte
}

func DeleteReq(url string) []byte {
	req := fiber.Delete(url)
	_, byte, err := req.Bytes()
	if err != nil {
		LOG.Error(err)
		time.Sleep(time.Second * 2)
	}

	return byte
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
		Func: DeleteFile,
	},
}

var localFileFuncs = []cli.Opt{
	{
		Name: "Send file",
		Func: SendFile,
	},
	{
		Name: "Execute file remotely",
		Func: func() {
			SendFile()
			ExecFile()
		},
	},
}

var FileList []cli.Opt

func SelectUploadFile() []cli.Opt {
	dir, err := os.ReadDir("./files")
	if err != nil {
		LOG.Error(err)
	}

	var opts []cli.Opt
	for _, file := range dir {
		opts = append(opts, cli.Opt{Name: file.Name(), Next: cli.Menu{Header: file.Name(), Opts: localFileFuncs}})
	}
	return opts
}

var file = cli.CURRENT_MENU.Header

func ExecFile() {
	file = cli.CURRENT_MENU.Header
	data := PostReq(HTTP_HOST + "/exec/" + file)
	LOG.Message(string(data))
	inp.WaitForAny()
}

func GetFileList() {
	file = cli.CURRENT_MENU.Header
	data := GetReq(HTTP_HOST + "/files")
	rawOpts := strings.Split(strings.TrimSpace(string(data)), " ")
	var opts []cli.Opt
	for _, opt := range rawOpts {
		opts = append(opts, cli.Opt{Name: opt, Next: cli.Menu{Header: opt, Opts: fileFuncs}})
	}
	FileList = opts
}

func SendFile() {
	LOG.Info("Sending file")
	file = cli.CURRENT_MENU.Header
	data, err := os.ReadFile(baseDir + file)
	if err != nil {
		LOG.Error(err)
		return
	}
	str := base64.URLEncoding.EncodeToString(data)
	LOG.Info("File sent")
	out := PostReq(HTTP_HOST + "/files/" + file + "/" + str)
	LOG.Message("SERVER:", string(out))
	inp.WaitForAny()
}

func DeleteFile() {
	file = cli.CURRENT_MENU.Header
	out := DeleteReq(HTTP_HOST + "/files/" + file)
	LOG.Message(string(out))
	inp.WaitForAny()
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
