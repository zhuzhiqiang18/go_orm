package tests

import (
	"fmt"
	"github.com/zhuzhiqiang18/go_orm"
	"github.com/zhuzhiqiang18/go_orm/model"

	"time"
)


/**
插入
 */
func TestSave(){
	db, err := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
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
	res, lastInsertId := db.Save(&student)
	if err !=nil {
		fmt.Println("err")
		return
	}
	fmt.Println("改变行数",res)
	fmt.Println("最后插入的id",lastInsertId)

}

/**
删除
 */
func TestDelete(){
	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	var student model.Student
	student.ClassId=1
	res := db.Delete(&student,"class_id")
	fmt.Println("改变行数",res)
}

/**
更新
 */
func TestUpdate(){
	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	var student model.Student
	student.Name="张三"
	student.No="00000000"
	res := db.Update(&student,"name")
	fmt.Println("改变行数",res)
}

/**
单表全查询
 */
func TestFindQuery()  {
	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	//传类型地址
	list := db.FindQuery(&model.Student{}, nil)
	for _,stu := range *list {
		fmt.Println(stu.(model.Student))
	}
}

/**
指定返回字段
 */
func TestFindQueryField()  {
	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	list := db.FindQuery(&model.Student{}, nil,"Name","No","Address")
	for _,stu := range *list {
		fmt.Println(stu.(model.Student))
		/*fmt.Println(stu.(model.Student).Name)
		fmt.Println(stu.(model.Student).No)
		fmt.Println(stu.(model.Student).Address)*/
	}

}

/**
但表条件查询
 */
func TestFindQueryWhere()  {
	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	list := db.FindQuery(&model.Student{}, map[string]interface{}{"name": "张三"},"Name","No","Address")
	for _,stu := range *list {
		fmt.Println(stu.(model.Student).Name)
		fmt.Println(stu.(model.Student).No)
		fmt.Println(stu.(model.Student).Address)
	}
}


/**
测试sql直接执行
 */
func TestNativeSql() {
	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	re := db.NativeSql("delete  from student")
	fmt.Println("改变条数",re)

}

/**
测试非自增主键 返回最后插入值
 */
func TestAutoInsertId()  {
	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	var class model.Class
	class.Name="tset"
	class.Id=1

	re, lastInsertId := db.Save(&class)
	fmt.Println("改变条数",re)
	fmt.Println("最后插入主键",lastInsertId)
}
/**
测试事务回滚
 */
func TestTx()  {
	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	var student model.Student
	student.Name="张三"
	student.No="00000000"
    tx:=db.Begin()
	for i:=0;i<10;i++{
		re, lastInsertId := tx.Save(&student)
		fmt.Println("改变条数",re)
		fmt.Println("最后插入主键",lastInsertId)
	}
	defer func() {
		err:=recover()
		if err !=nil {
			tx.Rollback()
		}
	}()

	panic("事务回滚")

	tx.Commit()

}