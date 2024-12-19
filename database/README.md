## ER图

* **使用ddb文件**

  1. 开[drawDB | Online database diagram editor and SQL generator](https://drawdb.vercel.app/)
  2. 新建一个Mysql表
  3. 导入图表，选择suask.ddb

* **使用Natvicat**: navicat可以自动生成ER图，但需要双击才能查看表中的细节，可以根据个人喜好使用


## 说明 2024-12-15
* 删除数据库中的trigger, 增加可编程性, 方便后续扩展

## 说明

* 对于点赞数，数据库有两个trigger来记录，不需要手动实现
* 对于通知，目前没有trigger，需要手动实现
* `user`表和`notification`表都实现了软删除，在`gorm`中对程序员透明，并不需要额外操作

## 数据库文件使用

`suask.sql`文件是只包含数据库结构的文件

`suask_data.sql`文件内包含基础导入数据

我对user表进行了初始填充，共有三个用户 `root` `teacher` 和 `student`

**密码**均为`123456`，可以进行登录操作