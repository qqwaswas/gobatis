package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wenj91/gobatis"
)

type User struct {
	Id       int64    `field:"id"`
	UserName string `field:"user_name"`
	Age      int8    `field:"age"`
	Addr     string `field:"addr"`
	Passwd string `field:"passwd"`
	IsDisable bool `field:"is_disable"`
	Money float32 `field:"money"`
	Total float64 `field:"total"`
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
	}

	batis, err := gobatis.NewGoBatis(context.Background(), config)
	if nil != err {
		panic(err)
	}

	runner, err := batis.Begin()

	//u := User{}
	//
	//err = runner.Select("userMapper.findMapById", map[string]interface{}{"id": 1})(&u)
	//
	//fmt.Printf("%v, error%v",u,err)

	var us  []User

	runner.Select("userMapper.queryStructs",map[string]interface{}{})(&us)

	fmt.Printf("%v",us)

}
