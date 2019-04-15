# gobatis

目前代码都是基于mysql编写测试的,其他数据库暂时还未做兼容处理


* mapper配置  
1. mapper可以配置namespace属性  
1. mapper可以包含: select, insert, update, delete标签  
1. mapper子标签id属性则为标签唯一标识, 必须配置属性
1. 其中select标签必须包含resultType属性，resultType可以是: map, maps, array, arrays, struct, structs, value
  
* 标签说明  
select: 用于查询操作   
insert: 用于插入sql操作  
update: 用于更新sql操作  
delete: 用于删除sql操作

* resultType说明  
map: 则数据库查询结果为map  
maps: 则数据库查询结果为map数组  
array: 则数据库查询结果为值数组  
arrays: 则数据库查询结果为多个值数组  
struct: 则数据库查询结果为单个结构体  
structs: 则数据库查询结果为结构体数组  
value: 则数据库查询结果为单个数值  
 
以下是mapper配置示例: mapper/userMapper.xml
```xml
<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE mapper PUBLIC "gobatis"
        "https://raw.githubusercontent.com/youkale/gobatis/master/gobatis.dtd">
<mapper namespace="userMapper">
    <select id="findMapById" resultType="struct">
        SELECT * FROM user where id=#{id} order by id
    </select>
    <select id="findMapByValue" resultType="map">
        SELECT * FROM user where id=1 order by id
    </select>
    <select id="findMapByValues" resultType="maps">
        SELECT * FROM user  order by id
    </select>
    <select id="findSliceByValue" resultType="slice">
        SELECT id FROM user where id=1  order  by id
    </select>
    <select id="findSlicesByValue" resultType="slices">
        SELECT * FROM user order by id
    </select>

    <select id="findArrayByValue" resultType="array">
        SELECT * FROM user where id=1  order by id
    </select>

    <select id="findValueByValue" resultType="value">
        SELECT id FROM user where id=1 order by id
    </select>

    <select id="findStructByStruct" resultType="struct">
        SELECT id, user_name, password FROM user where id=#{Id} order by id
    </select>
    <select id="findById" resultType="struct">
        SELECT * FROM user where id=#{id} order by id
    </select>
    <select id="queryStructs" resultType="structs">
        SELECT * FROM user order by id
    </select>
    <select id="queryStructsByCond" resultType="structs">
        SELECT id, user_name, password, pwd, email FROM user
        <where>
            <if test="Name != nil and Name != ''">and user_name = #{Name}</if>
        </where>
        order by id
    </select>
    <select id="queryStructsByCond2" resultType="structs">
        SELECT id, name, crtTm, pwd, email FROM user
        <trim prefixOverrides="and" prefix="where" suffixOverrides="," suffix="and 1=1">

            <if test="Name != nil and Name != ''">and user_name = #{Name}</if>
        </trim>
        order by id
    </select>
    <select id="queryStructsByCond3" resultType="structs">
        SELECT id, name, crtTm, pwd, email FROM user
        <trim prefixOverrides="and" prefix="where" suffixOverrides="," suffix="and 1=1">
            <choose>
                <when test="Age % 3 == 0">
                    and age = #{Age}
                </when>
                <when test="Age % 2 == 0 ">
                    and age = #{Age}
                </when>
                <otherwise>
                    and name = 'otherwise'
                </otherwise>
            </choose>
            <if test="Name != nil and Name != ''">and name = #{Name}</if>
            <if test="Password % 2 == 0 ">and pwd = #{Password} </if>

        </trim>
        order by id
    </select>

    <update id="updateByCond">
        update user
        <set>
            <if test="Name != nil and Name2 != ''">name = #{Name},</if>
            <if test="Password != nil and Password != ''">pwd = #{Password},</if>
        </set>
        where id = #{Id}
    </update>
    <insert id="saveUser">
        insert into user (user_name, age, addr)
        values (#{Name}, #{Email}, #{CrtTm})
    </insert>
    <delete id="deleteById">
        delete from user where id=#{id}
    </delete>
</mapper>
```

## 使用方法
example.go
```go
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
	Addr      gobatis.NullString
	Passwd    string
	IsDisable bool
	Money     gobatis.NullFloat64
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
	fmt.Printf("%v,%v\n",v,err)


	sean := User{}
	sean.Addr = gobatis.NullString{String:"火星"}
	sean.Age = 22
	sean.IsDisable = true
	sean.Money= gobatis.NullFloat64{
		Float64:10000.00,
	}
	sean.Passwd = "password"
	sean.Total = 1.22
	sean.UserName = "sean"

	runner, err := batis.Begin()
	i, i2, err := runner.Insert("userMapper.save", &sean)
	fmt.Printf("%v/%v/%v\n",i,i2,err)
	err = runner.Commit()

	begin, _ := batis.Begin()

	sean.Id = i
	sean.UserName = "sdfsdfs"

	update, _ := begin.Update("userMapper.updateById", &sean)

	err = begin.Commit()

	fmt.Printf("%v/%v",update,err)

	i3, err := batis.Delete("userMapper.deleteById", &sean)

	fmt.Printf("%v//%v",i3,err)

}

```
