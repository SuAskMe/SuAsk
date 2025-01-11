# SuAsk • 速问
![image](https://github.com/user-attachments/assets/68cbe1a7-1d9a-4487-a639-c86c15fbf6e9)

一个匿名提问箱，为了解决~~苏老师~~匿名提问箱经常跑路的问题

软件工程实训项目

---

## 它可以实现什么
![image](https://github.com/user-attachments/assets/a87d6177-c133-469e-8de1-f8d574bc923f)

- 完全匿名，用户可以畅所欲言，后台不会对用户信息进行记录
- 拥有公开论坛，用户可以自由发表提问
- 教师匿名提问箱，用户可以匿名提问教师问题，内置学院教师信息，开箱即用
- 提供提醒功能，你不会错过任何一条有关信息

---

## 如何使用

### 环境依赖
go 1.22
goframe v2.7.4


### 启动命令
```shell
git clone https://github.com/DreamingLri/SuAsk.git
cd SuAsk
gf run main.go
```

除此之外，您还需要配置此项目的[前端](https://github.com/DreamingLri/SuAsk-Web)

---

## Q&A

**Q: 为什么要创建这个项目？**
A: 学院的苏老师深受商用匿名提问箱跑路的迫害，~~据他所说已经用跑路了3个~~，便希望我们开发一个属于学院的匿名提问箱

**Q: 为什么使用go-frame？**
A: 这个项目本质上还是本专业的中级实训大作业，当时指明要求使用go语言来编写后端，而go-frame是一个较成熟的框架，便采用了

**Q: 我发现了bug怎么办？**
A: 欢迎提交[Issue](https://github.com/DreamingLri/SuAsk/issues/new)



