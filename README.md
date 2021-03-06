# go_orm
初学GO 自己封装的orm 旨在技术学习

# 目录
1. [ORM](#orm)
2. [Model约定](#model约定)
3. [如何使用](#引入包)
4. [链接数据库](#链接数据库)
5. [CURD使用方法](#curd使用方法)
6. [事务](#事务)
7. [GQL](#gql)
8. [NULL](#null)
9. map
10. 分页
11.  [命名规则（驼峰转为下划线）](#命名规则)


# ORM 
使用反射
# Model约定
1. tag sql 代表数据库中的字段 不注明 则使用结构体字段 
2. 数据库表 不区分大小写
```go
type Student struct {
	Id int `sql:"id"` 
	Name string `sql:"name"`
	Address string `sql:"address"`
	No string
	ClassId int `sql:"class_id"`
}
```
# 引入包
> import "github.com/zhuzhiqiang18/go_orm"
# 链接数据库
只支持MYSQL
```go
 db, err := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
   	defer db.Close()
```
# CURD使用方法
## 插入
```go
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
  	fmt.Println("改变行数",res)
  	fmt.Println("最后插入的id",lastInsertId)
```
## 删除
```go
  db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
  	defer db.Close()
  	var student model.Student
  	student.ClassId=1
  	res:= db.Delete(&student,"class_id")
  	fmt.Println("改变行数",res)
```
## 更改
```go
db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	var student model.Student
	student.Name="张三"
	student.No="00000000"
	res:= db.Update(&student,"name")
	fmt.Println("改变行数",res)
```
## 查询
### 单表全查询
```go
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
```
### 单表指定字段查询
```go
db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	//传类型地址
	list := make([]model.Student,0)
	err := db.FindQuery(&list, "select name,no from student ")
	if err!=nil {
		fmt.Println(err)
	}else {
		for _,stu := range list {
			fmt.Println(stu)
		}

	}
```
### 条件查询
>条件查询使用tag sql字段 


```go
db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	defer db.Close()
	//传类型地址
	list := make([]model.Student,0)
	err := db.FindQuery(&list, "select name,no from student where name=? ","张三")
	if err!=nil {
		fmt.Println(err)
	}else {
		for _,stu := range list {
			fmt.Println(stu)
		}

	}
```
### 联合查询
待更新……
# 事务
```go
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
```
# GQL
> gql 是一个 sql 生成器

```go
    //select * from Student where name ="张三" and class_id = 1
	var gql go_orm.Gql
	gql.Where("name = ? ").Where("class_id = ?").Bind(&model.Student{}).SetPara("张三",1)

   //select name,class_id from Student where name ="张三" and class_id = 1
   	new(go_orm.Gql).Where("name = ? ").Where("class_id = ?").Bind(&model.Student{}).SetPara("张三",1).Fields("name","class_id")
   
   	//select name,class_id from Student where name ="张三" and class_id = 1 or class_id = 2
   	new(go_orm.Gql).Where("name = ? ").Where("class_id = ?").Or("class_id = ?").Bind(&model.Student{}).SetPara("张三",1,2).Fields("name","class_id")
   
   	//select name,class_id from Student where name ="张三" and class_id = 1 or class_id = 2 order by id desc
   	new(go_orm.Gql).Where("name = ? ").Where("class_id = ?").Or("class_id = ?").Order("id desc").Bind(&model.Student{}).SetPara("张三",1,2).Fields("name","class_id")
   
```
## gql查询
```go
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
```
# NULL
>GO中没有NULL 为适配数据库中的NULL 以及JSON的NULL 引入第三方包 [gopkg.in/guregu/null.v3](https://github.com/guregu/null) 方便数据适配

## 示例
```go
type Teacher struct {
	Id null.Int `sql:"id"`
	Name null.String `sql:"name"`
	Address null.String `sql:"address"`
	No null.String `sql:"no"`
	ClassId null.Int `sql:"class_id"`
	Create null.Time `sql:"create_date"`
	IsReading null.Bool `sql:"is_reading"`
	High null.Float `sql:"high"`
	Weight null.Float `sql:"weight"`
}



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



```
# 命名规则
```go
db, _ := go_orm.Open("root","123456","127.0.0.1",3306,"go_test")
	db.DBSetting().SetFieldFormat(go_orm.HUMP_UNDERLINE)//字段驼峰转下划线
	db.DBSetting().SetTableFormat(go_orm.HUMP_UNDERLINE)//表名驼峰转下划线
```
字段设置驼峰转下划线 如字段tag为空则转下划线方式 如设置tag 则使用tag

