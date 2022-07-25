<span style="color: red">  go_firefly 开发脚手架，基于事件方式注册路由，使用xml 作为sql文件 </span>

# 基于golang的后台管理系统，适合新手学习，简单，清爽

## 功能特点

#### 前端基于layui,juicer前端模板引擎。

Layui: https://www.layui.com/demo/

juicer: https://github.com/PaulGuo/Juicer

#### 后端基于gin开发。

gin框架: https://github.com/gin-gonic/gin

#### MVC 设计模式,快速入门,方便上手。

#### goview 模板引擎，服务端html 渲染更简单，并且支持自定义 delimers。

goview: https://github.com/foolin/goview

#### GoMybatis 操作数据库，结构更简单清晰, 使用xml 方式，将sql 放在程序代码外

GoMybatis：https://github.com/zhuxiujia/GoMybatis

#### golang 基于事件方式注册路由，将路由文件分割到不文件中自动注册

eventbus：https://github.com/asaskevich/EventBus


## 二次开发 & 技术交流

#### 扫码备注: 'firefly',

![avatar](/static/img/qr.jpg)

## 环境要求

Mysql: 5.6+

## 目录说明

#### /resource 用于系统默认的配置文件

#### /resource/mybatis 用于存储 sql ymal 文件，支持yaml1.2 并增加yaml 点位标识

#### /src golang源代码

#### /static 用于存储前端css/js/img

#### /views 模板文件

## 界面载图

#### 登录界面

![avatar](/static/img/login.png)

#### 后台管理

![avatar](/static/img/home.png)

## 使用说明

#### 下载代码

```bash
git clone https://github.com/wzhsh90/go_firefly.git
cd go_firefly
go run main.go
```

#### 示例sql数据表

```sql
CREATE TABLE `sys_company_t` (
                                 `id` char(24) NOT NULL,
                                 `com_name` varchar(100) DEFAULT NULL,
                                 `com_desc` varchar(100) DEFAULT NULL,
                                 `flag` tinyint(1) DEFAULT '0',
                                 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

```

#### 示例xml 写sql

``` xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <resultMap id="BaseResultMap" tables="sys_company_t">
        <result column="id" property="Id" langType="string"/>
        <result column="com_desc" property="ComDesc" langType="string"/>
        <result column="com_name" property="ComName" langType="string"/>
    </resultMap>
    <sql id="base_col">
        id,com_name,com_desc
    </sql>
    <insertTemplete/>

    <select id="list">
        select <include refid="base_col"/> from sys_company_t
        <where>
            <if test="name!=''">
                and com_name like #{name}
            </if>
        </where>
        limit #{pageIndex},#{pageSize}
    </select>
    <select id="listCount">
        select count(*) from sys_company_t
        <where>
            <if test="name!=''">
                and com_name like #{name}
            </if>
        </where>
    </select>
    <select id="existName">
        select count(*) from sys_company_t where com_name=#{name}
    </select>
    <select id="Get">
        select <include refid="base_col"/> from sys_company_t where id=#{id}
    </select>
    <delete id="del">
        delete from sys_company_t where id=#{id}
    </delete>
    <update id="update">
        update sys_company_t set com_name=#{ComName},com_desc=#{ComDesc} where id=#{id}
    </update>
</mapper>


```

***** * 默认用户/名称: FireFly / firefly

#### 运行程序

```bash
默认运行开发模式
go run main.go= go run main.go --env dev
go run #开发模式: go run main.go --env dev
go run #生产模式: cargo run main.go --env prod
```

                  
#### rust 开发脚手架
rust 版本： git clone https://github.com/wzhsh90/rust_firefly.git
