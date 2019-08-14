package tests

import (
	"fmt"
	"go_web_curd/model"
	"go_web_curd/orm"
	"time"
)

var db orm.Db

func TestSave(){
	db, err := orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	if err!=nil {
		fmt.Println(err)
		return
	}
	var student model.Student
	student.Name="张三"
	student.Address="中国"
	student.No="123456"
	student.ClassId=1
	student.Create = time.Now()
	student.IsReading =true
	res:= db.Save(&student)
	fmt.Println("改变行数",res)
}

func TestDelete(){
	db, _ := orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	var student model.Student
	student.ClassId=1
	res:= db.Delete(&student,"class_id")
	fmt.Println("改变行数",res)
}


func TestUpdate(){
	db, _ := orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	var student model.Student
	student.Name="张三"
	student.No="00000000"
	res:= db.Update(&student,"name")
	fmt.Println("改变行数",res)
}

func TestFindQuery()  {
	db, _ := orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	//传类型地址
	list := db.FindQuery(&model.Student{}, nil)
	for _,stu := range *list {
		fmt.Println(stu.(model.Student))
	}
}

func TestFindQueryField()  {
	db, _ := orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	list := db.FindQuery(&model.Student{}, nil,"Name","No","Address")
	for _,stu := range *list {
		fmt.Println(stu.(model.Student).Name)
		fmt.Println(stu.(model.Student).No)
		fmt.Println(stu.(model.Student).Address)
	}

}

func TestFindQueryWhere()  {
	db, _ := orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	list := db.FindQuery(&model.Student{}, map[string]interface{}{"name": "张三"},"Name","No","Address")
	for _,stu := range *list {
		fmt.Println(stu.(model.Student).Name)
		fmt.Println(stu.(model.Student).No)
		fmt.Println(stu.(model.Student).Address)
	}

}