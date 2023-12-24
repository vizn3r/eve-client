package com

import "github.com/gofiber/fiber/v2"

func GetReq(url string) []byte {
	req := fiber.Get(url)
	_, byte, err := req.Bytes()
	if err != nil {
		panic(err)
	}

	return byte
}

func PostReq(url string) []byte {
	req := fiber.Post(url)
	_, byte, err := req.Bytes()
	if err != nil {
		panic(err)
	}

	return byte
}