# go_web_curd
go web 增删改查 封装了自己的orm

# ORM 
使用反射
# Model约定
1. tag sql 代表数据库中的字段 不注明 则使用结构体字段 
2. 数据库 不区分大小写
```go
type Student struct {
	Id int `sql:"id"` 
	Name string `sql:"name"`
	Address string `sql:"address"`
	No string
	ClassId int `sql:"class_id"`
}
```
# CURD使用方法
## 插入
```go
    var student model.Student
	student.Name="张三"
	student.Address="中国"
	student.No="123456"
	student.ClassId=1
	res:= db.Save(student)
	fmt.Println("改变行数",res)
```
## 删除
```go
    var student model.Student
	student.ClassId=1
	res:= db.Delete(student,"class_id")//删除条件 需要使用tag sql字段
	fmt.Println("改变行数",res)
```
## 更改
```go
var student model.Student
	student.Name="张三"
	student.No="00000000"
	res:= db.Update(student,"name")//更改条件 需要使用tag sql字段
	fmt.Println("改变行数",res)
```
## 查询
待更新……