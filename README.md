# go_web_curd
go web 增删改查 封装了自己的orm

# ORM 
使用反射
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