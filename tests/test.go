package tests

import (
	"fmt"
	"go_web_curd/model"
	"go_web_curd/persistent"
	"time"
)

func TestSave(){
	var student model.Student
	student.Name="张三"
	student.Address="中国"
	student.No="123456"
	student.ClassId=1
	student.Create = time.Now()
	student.IsReading =true
	res:= persistent.Save(&student)
	fmt.Println("改变行数",res)
}

func TestDelete(){
	var student model.Student
	student.ClassId=1
	res:= persistent.Delete(&student,"class_id")
	fmt.Println("改变行数",res)
}


func TestUpdate(){
	var student model.Student
	student.Name="张三"
	student.No="00000000"
	res:= persistent.Update(&student,"name")
	fmt.Println("改变行数",res)
}

func TestFindQuery()  {
	//传类型地址
	list := persistent.FindQuery(&model.Student{}, nil)
	for _,stu := range *list {
		fmt.Println(stu.(model.Student))
	}
}

func TestFindQueryField()  {
	list := persistent.FindQuery(&model.Student{}, nil,"Name","No","Address")
	for _,stu := range *list {
		fmt.Println(stu.(model.Student).Name)
		fmt.Println(stu.(model.Student).No)
		fmt.Println(stu.(model.Student).Address)
	}

}

func TestFindQueryWhere()  {
	list := persistent.FindQuery(&model.Student{}, map[string]interface{}{"name": "张三"},"Name","No","Address")
	for _,stu := range *list {
		fmt.Println(stu.(model.Student).Name)
		fmt.Println(stu.(model.Student).No)
		fmt.Println(stu.(model.Student).Address)
	}

}