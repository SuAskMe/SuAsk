# SuAsk
an anonymous Question Box for Dr. Su



# 环境搭建

本项目使用goframe框架进行开发，且最新版goframe要求go语言版本大于1.22，请查看自己的go语言版本时候符合条件

goframe官方文档：[GoFrame官网 - 类似PHP-Laravel,Java-SpringBoot的Go企业级开发框架](https://goframe.org/)

**本项目推荐使用linux环境进行开发**

### Linux & Mac 安装

```shell
wget -O gf https://github.com/gogf/gf/releases/latest/download/gf_$(go env GOOS)_$(go env GOARCH) && chmod +x gf && ./gf install -y && rm ./gf
```

### Windows 安装

请到官方文档下自行获取相关链接教程：[工具安装-install | GoFrame官网 - 类似PHP-Laravel,Java-SpringBoot的Go企业级开发框架](https://goframe.org/docs/cli/install)

### 通用安装方式 （推荐）

```shell
go install github.com/gogf/gf/cmd/gf/v2@latest
```



# 使用

**本项目已经初始化好仓库并且搭好脚手架**

查看goframe版本

```shell
gf -v
```

在`SuAsk`目录下启动项目

```shell
gf run main.go
```

构建项目 (目前尚未配置

```
gf build
```



# 项目目录

```
/
├── api                 请求接口输入/输出数据结构定义
├── hack                项目开发工具、脚本
├── internal            业务逻辑存放目录，核心代码
│   ├── cmd             入口指令与其他命令工具目录
│   ├── consts          常量定义目录
│   ├── controller      控制器目录，接收/解析用户请求
│   ├── dao             数据访问对象目录，用于和底层数据库交互
│   ├── logic           核心业务逻辑代码目录
│   ├── model           数据结构管理模块，管理数据实体对象，以及输入与输出数据结构定义
│   |   ├── do          数据操作中业务模型与实例模型转换，由工具维护，不能手动修改
│   │   └── entity      数据模型是模型与数据集合的一对一关系，由工具维护，不用手动修改。
│   └── service         业务接口定义层。具体的接口实现在logic中进行注入。
├── manifest            包含程序编译、部署、运行、配置的文件
├── resource            静态资源文件
├── utility
├── go.mod
└── main.go             程序入口文件
```



# 学习

* **3小时快速入门课程：**[Go语言Web开发|GoFrame框架入门_哔哩哔哩_bilibili](https://www.bilibili.com/video/BV1Uu4y1u7kX?spm_id_from=333.788.videopod.episodes&vd_source=03dcb05712764a790687a977fee70f9d)

* **课程笔记：**[gitee链接](https://gitee.com/unlimited13/code/blob/master/GO/GoFrame.md)
