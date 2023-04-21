package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
	openai "github.com/sashabaranov/go-openai"
	"image/png"
	"net/http"
	"net/url"
	"os"
)

func main() {
	key := "xxx"
	gptConfig := openai.DefaultConfig(key)

	{
		// 国内需要代理
		transport := &http.Transport{}
		proxyUrl, err := url.Parse("http://localhost:10809")
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
		// 创建一个 HTTP 客户端，并将 Transport 对象设置为其 Transport 字段
		gptConfig.HTTPClient = &http.Client{
			Transport: transport,
		}
	}

	c := openai.NewClientWithConfig(gptConfig)
	ctx := context.Background()

	// Sample image by link
	reqUrl := openai.ImageRequest{
		Prompt:         "车水马龙",
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              5,
	}

	respUrl, err := c.CreateImage(ctx, reqUrl)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return
	}
	fmt.Println(respUrl.Data[0].URL) // 输出图片的URL

	// Example image as base64
	reqBase64 := openai.ImageRequest{
		Prompt:         "车水马龙",
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              7,
	}

	respBase64, err := c.CreateImage(ctx, reqBase64)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		fmt.Printf("Base64 decode error: %v\n", err)
		return
	}

	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		fmt.Printf("PNG decode error: %v\n", err)
		return
	}

	file, err := os.Create(bfile.Join(bfile.SelfDir(), "example.png"))
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		fmt.Printf("PNG encode error: %v\n", err)
		return
	}

	fmt.Println("The image was saved as example.png")
}
