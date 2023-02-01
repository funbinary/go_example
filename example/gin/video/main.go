package main

/*
   gin框架实现文件下载功能
*/

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
	"github.com/gin-gonic/gin"
)

// 主函数
func main() {
	r := gin.Default()

	//Get路由，动态路由
	r.GET("/GetRecord/:path", DowFile)

	//监听端口
	err := r.Run(":18888")
	if err != nil {
		fmt.Println("error")
	}
}

// 文件下载功能实现
func DowFile(c *gin.Context) {
	//通过动态路由方式获取文件名，以实现下载不同文件的功能
	path := c.Param("path")
	//拼接路径,如果没有这一步，则默认在当前路径下寻找
	filename := bfile.Join("/data/video/record", path, path+".mp4")
	//响应一个文件
	c.File(filename)
	return
}
