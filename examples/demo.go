package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/youkale/gobatis"
)

type User struct {
	Id        int64
	UserName  string
	Age       int8
	Addr      string
	Passwd    string
	IsDisable bool
	Money     float32
	Total     float64
}

func main() {
	db, err := sql.Open("mysql",
		"root:toor@tcp(127.0.0.1:3306)/gobatis")

	if nil != err {
		panic(err)
	}

	err = db.Ping()
	if nil != err {
		panic(err)
	}

	config := &gobatis.Config{
		Db:          db,
		MapperPaths: []string{"./examples/mapper"},
		Debug:true,
		ColumnStyle:[]int{gobatis.StyleSnake},
	}

	batis, err := gobatis.NewGoBatis(context.Background(), config)
	if nil != err {
		panic(err)
	}


	u := User{}

	err = batis.Select("userMapper.findMapById", map[string]interface{}{"id": 1})(&u)

	fmt.Printf("%v, error%v\n", u, err)

	var us []User

	_ = batis.Select("userMapper.queryStructs", map[string]interface{}{})(&us)

	fmt.Printf("%v\n", us)

	m := make(map[string]interface{})

	err = batis.Select("userMapper.findMapByValue",map[string]interface{}{})(&m)

	fmt.Printf("%v, err %v", m,err)

	var ms []map[string]interface{}

	err = batis.Select("userMapper.findMapByValues",map[string]interface{}{})(&ms)

	fmt.Printf("%v",ms)

	var ss []interface{}

	err = batis.Select("userMapper.findSliceByValue",map[string]interface{}{})(&ss)

	fmt.Printf("%v,%v",ss,err)

	var sss [][]interface{}

	err = batis.Select("userMapper.findSlicesByValue",map[string]interface{}{})(&sss)

	fmt.Printf("%v,%v",sss,err)

	var ar []interface{}
	err = batis.Select("userMapper.findArrayByValue",map[string]interface{}{})(&ar)
	fmt.Printf("%v,%v",ar,err)

	var ars []interface{}
	err = batis.Select("userMapper.findArraysByValue",map[string]interface{}{})(&ars)
	fmt.Printf("%v,%v",ars,err)

	var v int
	err = batis.Select("userMapper.findValueByValue",map[string]interface{}{})(&v)
	fmt.Printf("%v,%v",v,err)

}

