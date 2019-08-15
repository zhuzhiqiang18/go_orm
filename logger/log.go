package logger

import "fmt"

func Info(sqlStr string,para interface{})  {
	fmt.Println("[INFO]:",sqlStr,para)
}

func Err(sqlStr string,para interface{})  {
	fmt.Println("[ERR]:",sqlStr,para)
}

func Debug(sqlStr string,para interface{})  {
	fmt.Println("[DEBUG]:",sqlStr,para)

}
