package main

import (
	"fmt"

	"github.com/funbinary/go_example/pkg/errors"
)

func ExampleErrorNew() {
	newerr := errors.New("cuowu")
	fmt.Println("===============%s===================")
	fmt.Printf("new %s\n", newerr) //new cuowu

	fmt.Println("===============%v===================")
	fmt.Printf("new %v\n", newerr) //new cuowu

	fmt.Println("===============%+v===================")
	fmt.Printf("new %+v\n", newerr)
	//new cuowu
	//main.main
	//		D:/workspace/kds/base_package/errors/example/example.go:9
	//runtime.main
	//		D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//		D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581

	fmt.Println("===============%#+v===================")
	fmt.Printf("new %#+v\n", newerr)

	//new cuowu
	//main.main
	//		D:/workspace/kds/base_package/errors/example/example.go:9
	//runtime.main
	//		D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//		D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581

	fmt.Println("===============%-v===================")
	fmt.Printf("new %-v\n", newerr) //new cuowu

	fmt.Println("===============%#-v===================")
	fmt.Printf("new %#-v\n", newerr) //new cuowu

}

func ExampleErrorErrorf() {
	artile := "斗破苍穹"
	id := 10086
	errorferr := errors.Errorf("ID:%d 文章:%s", id, artile)
	fmt.Println("===============%s===================")
	fmt.Printf("errorf %s\n", errorferr) //errorf ID:10086 文章:斗破苍穹

	fmt.Println("===============%v===================")
	fmt.Printf("errorf %v\n", errorferr) //errorf ID:10086 文章:斗破苍穹

	fmt.Println("===============%+v===================")
	fmt.Printf("errorf %+v\n", errorferr) //errorf ID:10086 文章:斗破苍穹
	//main.ExampleErrorErrorf
	//	D:/workspace/kds/base_package/errors/example/example.go:48
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:65
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581

	fmt.Println("===============%#+v===================")
	fmt.Printf("errorf %#+v\n", errorferr) //errorf ID:10086 文章:斗破苍穹
	//main.ExampleErrorErrorf
	//	D:/workspace/kds/base_package/errors/example/example.go:48
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:65
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581

	fmt.Println("===============%-v===================")
	fmt.Printf("errorf %-v\n", errorferr) //errorf ID:10086 文章:斗破苍穹

	fmt.Println("===============%#-v===================")
	fmt.Printf("errorf %#-v\n", errorferr) //errorf ID:10086 文章:斗破苍穹
}

func ExampleErrorWrap() {
	err := errors.New("first error")
	wraperr := errors.Wrap(err, "wrap error")
	fmt.Println("===============%s===================")
	fmt.Printf("Wrap %s\n", wraperr) //Wrap wrap error

	fmt.Println("===============%v===================")
	fmt.Printf("Wrap %v\n", wraperr) //Wrap wrap error

	fmt.Println("===============%+v===================")
	fmt.Printf("Wrap %+v\n", wraperr)
	//Wrap first error
	//main.ExampleErrorWrap
	//	D:/workspace/kds/base_package/errors/example/example.go:85
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:126
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581
	//wrap error
	//main.ExampleErrorWrap
	//	D:/workspace/kds/base_package/errors/example/example.go:86
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:126
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581

	fmt.Println("===============%#+v===================")
	fmt.Printf("Wrap %#+v\n", wraperr)
	//Wrap first error
	//main.ExampleErrorWrap
	//	D:/workspace/kds/base_package/errors/example/example.go:85
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:126
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581
	//wrap error
	//main.ExampleErrorWrap
	//	D:/workspace/kds/base_package/errors/example/example.go:86
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:126
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581

	fmt.Println("===============%-v===================")
	fmt.Printf("Wrap %-v\n", wraperr) //Wrap wrap error

	fmt.Println("===============%#-v===================")
	fmt.Printf("Wrap %#-v\n", wraperr) //Wrap wrap error

}

func ExampleErrorWrapf() {
	err := errors.New("first error")
	wrapferr := errors.Wrapf(err, "wrapf error #%d", 2)
	fmt.Println("===============%s===================")
	fmt.Printf("Wrapf %s\n", wrapferr) //Wrapf wrapf error #2

	fmt.Println("===============%v===================")
	fmt.Printf("Wrapf %v\n", wrapferr) //Wrapf wrapf error #2

	fmt.Println("===============%+v===================")
	fmt.Printf("Wrapf %+v\n", wrapferr)
	//Wrapf first error
	//main.ExampleErrorWrapf
	//	D:/workspace/kds/base_package/errors/example/example.go:144
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:207
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581
	//wrapf error #2
	//main.ExampleErrorWrapf
	//	D:/workspace/kds/base_package/errors/example/example.go:145
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:207
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581

	fmt.Println("===============%#+v===================")
	fmt.Printf("Wrapf %#+v\n", wrapferr)
	//Wrapf first error
	//main.ExampleErrorWrapf
	//	D:/workspace/kds/base_package/errors/example/example.go:144
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:207
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581
	//wrapf error #2
	//main.ExampleErrorWrapf
	//	D:/workspace/kds/base_package/errors/example/example.go:145
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:207
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581

	fmt.Println("===============%-v===================")
	fmt.Printf("Wrapf %-v\n", wrapferr) //Wrapf wrapf error #2

	fmt.Println("===============%#-v===================")
	fmt.Printf("Wrapf %#-v\n", wrapferr) //Wrapf wrapf error #2

}

func ExampleErrorWithStack() {
	err := errors.New("first error")
	withStackErr := errors.WithStack(err)
	fmt.Println("===============%s===================")
	fmt.Printf("WithStack %s\n", withStackErr) //WithStack first error

	fmt.Println("===============%v===================")
	fmt.Printf("WithStack %v\n", withStackErr) //WithStack first error

	fmt.Println("===============%+v===================")
	fmt.Printf("WithStack %+v\n", withStackErr)
	//WithStack first error
	//main.ExampleErrorWithStack
	//	D:/workspace/kds/base_package/errors/example/example.go:203
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:267
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581
	//main.ExampleErrorWithStack
	//	D:/workspace/kds/base_package/errors/example/example.go:204
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:267
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581

	fmt.Println("===============%#+v===================")
	fmt.Printf("WithStack %#+v\n", withStackErr)
	//WithStack first error
	//main.ExampleErrorWithStack
	//	D:/workspace/kds/base_package/errors/example/example.go:203
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:267
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581
	//main.ExampleErrorWithStack
	//	D:/workspace/kds/base_package/errors/example/example.go:204
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:267
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581

	fmt.Println("===============%-v===================")
	fmt.Printf("WithStack %-v\n", withStackErr) //WithStack first error

	fmt.Println("===============%#-v===================")
	fmt.Printf("WithStack %#-v\n", withStackErr) //WithStack first error
}

func ExampleErrorWithMessage() {
	err := errors.New("first error")
	withMessageErr := errors.WithMessage(err, "with message")
	fmt.Println("===============%s===================")
	fmt.Printf("WithMessage %s\n", withMessageErr) //WithMessage with message

	fmt.Println("===============%v===================")
	fmt.Printf("WithMessage %v\n", withMessageErr) //WithMessage with message

	fmt.Println("===============%+v===================")
	fmt.Printf("WithMessage %+v\n", withMessageErr)
	//WithMessage first error
	//main.ExampleErrorWithMessage
	//	D:/workspace/kds/base_package/errors/example/example.go:259
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:322
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581
	//with message

	fmt.Println("===============%#+v===================")
	fmt.Printf("WithMessage %#+v\n", withMessageErr)
	//WithMessage first error
	//main.ExampleErrorWithMessage
	//	D:/workspace/kds/base_package/errors/example/example.go:259
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:322
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581
	//with message

	fmt.Println("===============%-v===================")
	fmt.Printf("WithMessage %-v\n", withMessageErr) //WithMessage with message

	fmt.Println("===============%#-v===================")
	fmt.Printf("WithMessage %#-v\n", withMessageErr) //WithMessage with message

}

func ExampleErrorWithMessagef() {
	err := errors.New("first error")
	withMessagefErr := errors.WithMessagef(err, "with message,#%d", 2)
	fmt.Println("===============%s===================")
	fmt.Printf("withMessagefErr %s\n", withMessagefErr) //withMessagefErr with message,#2

	fmt.Println("===============%v===================")
	fmt.Printf("withMessagefErr %v\n", withMessagefErr) //withMessagefErr with message,#2

	fmt.Println("===============%+v===================")
	fmt.Printf("withMessagefErr %+v\n", withMessagefErr)
	//withMessagefErr first error
	//main.ExampleErrorWithMessagef
	//	D:/workspace/kds/base_package/errors/example/example.go:301
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:351
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.goexit
	//	D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581
	//with message,#2

	fmt.Println("===============%#+v===================")
	fmt.Printf("withMessagefErr %#+v\n", withMessagefErr)
	//withMessagefErr first error
	//main.ExampleErrorWithMessagef
	//	D:/workspace/kds/base_package/errors/example/example.go:301
	//main.main
	//	D:/workspace/kds/base_package/errors/example/example.go:351
	//runtime.main
	//	D:/dev/go1.17.5/src/runtime/proc.go:255
	//runtime.D:/dev/go1.17.5/src/runtime/asm_amd64.s:1581
	//with message,#2

	fmt.Println("===============%-v===================")
	fmt.Printf("withMessagefErr %-v\n", withMessagefErr) //withMessagefErr with message,#2

	fmt.Println("===============%#-v===================")
	fmt.Printf("withMessagefErr %#-v\n", withMessagefErr) //withMessagefErr with message,#2

}

const (
	// Error codes below 1000 are reserved future use by the
	// "github.com/bdlm/errors" package.
	ConfigurationNotValid int = iota + 1000
	ErrInvalidJSON
	ErrEOF
	ErrLoadConfigFailed
)

type defaultCoder struct {
	// C refers to the integer code of the ErrCode.
	C int

	// External (user) facing error text.
	Ext string
}

// Code returns the integer code of the coder.
func (coder defaultCoder) Code() int {
	return coder.C

}

// String implements stringer. String returns the external error message,
// if any.
func (coder defaultCoder) String() string {
	return coder.Ext
}

func init() {
	errors.Register(defaultCoder{ConfigurationNotValid, "ConfigurationNotValid error"})
	errors.Register(defaultCoder{ErrInvalidJSON, "Data is not valid JSON"})
	errors.Register(defaultCoder{ErrEOF, "End of input"})
	errors.Register(defaultCoder{ErrLoadConfigFailed, "Load configuration file failed"})
}

func ExampleErrorWithCode() {
	withCodeErr := errors.WithCode(ConfigurationNotValid, "error with code, %s", "none")
	fmt.Println("===============%s===================")
	fmt.Printf("WithCode %s\n", withCodeErr) //WithCode ConfigurationNotValid error

	fmt.Println("===============%v===================")
	fmt.Printf("WithCode %v\n", withCodeErr) //WithCode ConfigurationNotValid error

	fmt.Println("===============%+v===================")
	fmt.Printf("WithCode %+v\n", withCodeErr)
	//WithCode error with code, none - #0 [D:/workspace/kds/base_package/errors/example/example.go:379 (main.ExampleErrorWithCode)] (1000) ConfigurationNotValid error

	fmt.Println("===============%#+v===================")
	fmt.Printf("WithCode %#+v\n", withCodeErr)
	//WithCode [{"caller":"#0 D:/workspace/kds/base_package/errors/example/example.go:379 (main.ExampleErrorWithCode)","code":1000,"error":"error with code, none","message":"ConfigurationNotValid error"}]

	fmt.Println("===============%-v===================")
	fmt.Printf("WithCode %-v\n", withCodeErr)
	//WithCode error with code, none - #0 [D:/workspace/kds/base_package/errors/example/example.go:379 (main.ExampleErrorWithCode)] (1000) ConfigurationNotValid error

	fmt.Println("===============%#-v===================")
	fmt.Printf("WithCode %#-v\n", withCodeErr)
	//WithCode [{"caller":"#0 D:/workspace/kds/base_package/errors/example/example.go:379 (main.ExampleErrorWithCode)","code":1000,"error":"error with code, none","message":"ConfigurationNotValid error"}]

}

func ExampleErrorWithC() {
	err := errors.New("NO valid")
	wrapCErr := errors.WrapC(err, ConfigurationNotValid, "配置无效")
	fmt.Println("===============%s===================")
	fmt.Printf("WrapC %s\n", wrapCErr) //WrapC ConfigurationNotValid error

	fmt.Println("===============%v===================")
	fmt.Printf("WrapC %v\n", wrapCErr) //WrapC ConfigurationNotValid error

	fmt.Println("===============%+v===================")
	fmt.Printf("WrapC %+v\n", wrapCErr)
	//WrapC 配置无效 - #1 [D:/workspace/kds/base_package/errors/example/example.go:406 (main.ExampleErrorWithC)] (1000) ConfigurationNotValid error; NO valid - #0 [D:/workspace/kds/base_package/errors/example/example.go:405 (main.ExeErrorWithC)] (1) NO valid
	fmt.Println("===============%#+v===================")
	fmt.Printf("WrapC %#+v\n", wrapCErr)
	//WrapC [{"caller":"#1 D:/workspace/kds/base_package/errors/example/example.go:406 (main.ExampleErrorWithC)","code":1000,"error":"配置无效","message":"ConfigurationNotValid error"},{"caller":"#0 D:/workspace/kds/base_package/errexample/example.go:405 (main.ExampleErrorWithC)","code":1,"error":"NO valid","message":"NO valid"}]

	fmt.Println("===============%-v===================")
	fmt.Printf("WrapC %-v\n", wrapCErr)
	//WrapC 配置无效 - #1 [D:/workspace/kds/base_package/errors/example/example.go:406 (main.ExampleErrorWithC)] (1000) ConfigurationNotValid error

	fmt.Println("===============%#-v===================")
	fmt.Printf("WrapC %#-v\n", wrapCErr)
	//WrapC [{"caller":"#1 D:/workspace/kds/base_package/errors/example/example.go:406 (main.ExampleErrorWithC)","code":1000,"error":"配置无效","message":"ConfigurationNotValid error"}]
}

func main() {
	//err := fmt.Errorf("aasd")
	//fmt.Printf("%+v", err)
	//ExampleErrorNew()
	//ExampleErrorErrorf()
	//ExampleErrorWrap()
	//ExampleErrorWrapf()
	//ExampleErrorWithStack()
	//ExampleErrorWithMessage()
	//ExampleErrorWithMessagef()
	//ExampleErrorWithCode()
	ExampleErrorWithC()
}
