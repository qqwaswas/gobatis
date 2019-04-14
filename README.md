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
        "https://raw.githubusercontent.com/wenj91/gobatis/master/gobatis.dtd">
<mapper namespace="userMapper">
    <select id="findMapById" resultType="map">
        SELECT id, name FROM user where id=#{id} order by id
    </select>
    <select id="findMapByValue" resultType="map">
            SELECT id, name FROM user where id=#{0} order by id
    </select>
    <select id="findStructByStruct" resultType="struct">
        SELECT id, name, crtTm FROM user where id=#{Id} order by id
    </select>
    <select id="queryStructs" resultType="structs">
        SELECT id, name, crtTm FROM user order by id
    </select>
    <select id="queryStructsByOrder" resultType="structs">
        SELECT id, name, crtTm FROM user order by ${id} desc
    </select>
    <insert id="insertStruct">
        insert into user (name, email, crtTm)
        values (#{Name}, #{Email}, #{CrtTm})
    </insert>
    <delete id="deleteById">
        delete from user where id=#{id}
    </delete>
    <select id="queryStructsByCond" resultType="structs">
         SELECT id, name, crtTm, pwd, email FROM user
         <where>
             <if test="Name != nil and Name != ''">and name = #{Name}</if>
         </where>
         order by id
    </select>
     <select id="queryStructsByCond2" resultType="structs">
         SELECT id, name, crtTm, pwd, email FROM user
         <trim prefixOverrides="and" prefix="where" suffixOverrides="," suffix="and 1=1">
              <if test="Name != nil and Name != ''">and name = #{Name}</if>
         </trim>
         order by id
    </select>
    <update id="updateByCond">
        update user
        <set>
            <if test="Name != nil and Name2 != ''">name = #{Name},</if>
        </set>
        where id = #{Id}
    </update>
</mapper>
```

## 使用方法
example.go
```go
package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 引入驱动
	"github.com/wenj91/gobatis"        // 引入gobatis
)

// 实体结构示例， tag：field为数据库对应字段名称
type User struct {
	Id    gobatis.NullInt64  `field:"id"`
	Name  gobatis.NullString `field:"name"`
	Email gobatis.NullString `field:"email"`
	CrtTm gobatis.NullTime   `field:"crtTm"`
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
    
    	u := User{}
    
    	err = runner.Select("userMapper.findMapById", map[string]interface{}{"id": 1})(&u)
    
    	fmt.Printf("%v, error%v",u,err)
}
```
