package go_orm

import (
	"fmt"
)

/**
sql 生成器
 */
type Gql struct {
	andSql string
	andCount int64
	orSql string
	//orCount int64
	orderBySql string
	para []interface{}
	fields[] string
	isCount bool
	selectSql string
	t interface{}
	
}

func (gql *Gql) New() *Gql {
	return &Gql{}
}

func (gql *Gql) Where(andsql string) *Gql {
	if gql.andCount==0 {
		gql.andSql+=" where " +andsql+" "
	}else {
		gql.andSql+=" and " +andsql+" "
	}
	gql.andCount++
	return gql
}

func (gql *Gql) Or(orSql string) *Gql {

	gql.orSql+=" or " +orSql +" "

	return gql
}

func (gql *Gql) Order(orderSql string) *Gql {
	gql.orderBySql="order by "+orderSql +" "
	return gql
}


func (gql *Gql) Count() *Gql {
	gql.isCount=true
	return gql
}

func (gql *Gql)Fields(field ...string)  {
	gql.fields=field
}

func (gql *Gql) GetGql() string {
	if gql.t == nil {
		panic("需要绑定一个结构体")
	}
	ss:=""
	if 0 == len(gql.fields) {
		gql.selectSql = " * "
	}else {
		for _,fd := range gql.fields{
			ss+= fd +" ,"
		}
		gql.selectSql=ss[:len(ss)-1]
	}
	
	if gql.isCount  {
		return fmt.Sprintf("%s count(%s) %s %s %s","select",gql.selectSql,gql.andSql,gql.orSql,gql.orderBySql)
	}

	return fmt.Sprintf("%s %s %s %s %s","select",gql.selectSql,gql.andSql,gql.orSql,gql.orderBySql)
}

func (gql *Gql) setPara(para ...interface{}) *Gql {
	gql.para=para
	return gql
}

func (gql *Gql) GetPara() *[]interface{} {
	return &(gql.para)
}
func (gql *Gql) GetBind() interface{}  {
	return gql.t
}

func (gql *Gql) Bind(bind interface{}) *Gql {
	gql.t=bind
	return gql
}

func (gql *Gql) GetFields() ([] string) {
	return gql.fields
}





