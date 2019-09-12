package tests

import (
	"encoding/json"
	"fmt"
	"github.com/zhuzhiqiang18/go_orm"
	"github.com/zhuzhiqiang18/go_orm/model"
	"gopkg.in/guregu/null.v3"

	"time"
)


/**
插入
 */
func TestSave(){
	db, err := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	db.DBSetting().SetTableFormat(go_orm.HUMP_UNDERLINE)
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
	list := make([]model.Student,0)
	err := db.FindQuery(&list, "select * from student ")
	if err!=nil {
		fmt.Println(err)
	}else {
		for _,stu := range list {
			fmt.Println(stu)
		}

	}

}





/**
但表条件查询
 */
func TestFindQueryWhere()  {
	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	list := make([]model.Student,0)
	err := db.FindQuery(&list, "select * from student where name =?","张三")
	if err!=nil {
		fmt.Println(err)
	}else {
		for _,stu := range list {
			fmt.Println(stu.Name)
			fmt.Println(stu.No)
			fmt.Println(stu.Address)
		}
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
    tx:=db.Begin()//获取事务
	for i:=0;i<10;i++{
		re, lastInsertId := tx.Save(&student)
		fmt.Println("改变条数",re)
		fmt.Println("最后插入主键",lastInsertId)
	}
	defer func() {
		err:=recover()
		if err !=nil {
			tx.Rollback()//事务回滚
			return
		}
	}()

	panic("事务回滚")

	tx.Commit()//事务提交

}

func TestTx1()  {
	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	var student model.Student
	student.Name="张三"
	student.No="00000000"
	student.ClassId=1
	tx:=db.Begin()//获取事务

	tx.Update(&student,"name")

	var c model.Class
	c.Name="张三"
	c.Id=2
	tx.Save(&c)

	defer func() {
		err:=recover()
		if err !=nil {
			tx.Rollback()//事务回滚
			return
		}
	}()

	panic("事务回滚")

	tx.Commit()//事务提交

}

func TestGql()  {
	var gql go_orm.Gql
	students := make([]model.Student,0)


	//select * from Student where name ="张三" and class_id = 1
	gql.Where("name = ? ").Where("class_id = ?").Bind(&students).SetPara("张三",1)

	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	err := db.FindGql(&gql)

	if err!=nil {
		fmt.Println(err)
		return
	}

	for _,stu := range students {
		fmt.Println(stu)
	}
}
/**
测试null包
 */
func TestNull()  {
	var teacher model.Teacher
	teacher.Name = null.NewString("zzq",true)
	teacher.Create = null.NewTime(time.Now(),true)
	teacher.IsReading = null.NewBool(true,true)
	teacher.High = null.NewFloat(160.256,true)
	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	count,last:= db.Save(&teacher);
	fmt.Println("改变条数" , count)
	fmt.Println("最后插入" , last)
}

/**
支持null值
 */
func TestFindNull()  {
	var gql go_orm.Gql
	teachers := make([]model.Teacher,0)
	//select * from Student where name ="zzq" and class_id = 1
	gql.Where("name = ? ").Where("is_reading = ?").Bind(&teachers).SetPara("zzq",1)

	db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	db.DBSetting().SetFieldFormat(go_orm.HUMP_UNDERLINE)//驼峰下划线
	db.DBSetting().SetTableFormat(go_orm.HUMP_UNDERLINE)//驼峰下划线
	defer db.Close()
	err := db.FindGql(&gql)
	if err != nil{
		panic(err)
	}

	for _,th := range teachers {
		fmt.Println(th)
		jsonStr,_:=json.Marshal(th)
		fmt.Println(string(jsonStr))
//{"Id":1,"Name":"zzq","Address":null,"No":null,"ClassId":null,"Create":"2019-08-20T14:48:12Z","IsReading":true,"High":null,"Weight":null}
	}
}


