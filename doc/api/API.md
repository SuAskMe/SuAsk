# SuAsk API 文档

## 概述

本文档详细描述了 SuAsk 匿名提问平台的所有 API 接口。SuAsk 是一个为学院教师和学生提供匿名提问与交流的平台。

## 基础信息

- **协议**: HTTP/HTTPS
- **数据格式**: JSON
- **字符编码**: UTF-8
- **认证方式**: JWT Token

## 状态码

| 状态码 | 描述 |
|--------|------|
| 0 | 成功 |
| 1 | 失败 |

## 认证机制

大部分 API 接口需要认证，用户需要在请求头中添加 `Authorization` 字段：

```
Authorization: Bearer <token>
```

## API 接口列表

### 1. 登录认证相关

#### 1.1 用户登录

**接口地址**: `POST /login`

**请求参数**:
```json
{
  "name": "用户名",
  "password": "密码"
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "type": "Token格式",
    "token": "用户的Token",
    "role": "用户角色",
    "id": 1
  }
}
```

#### 1.2 心跳检测

**接口地址**: `POST /user/heartbeat`

**请求参数**: 无

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1
  }
}
```

#### 1.3 用户登出

**接口地址**: `POST /user/logout`

**请求参数**: 无

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

#### 1.4 刷新 Token

**接口地址**: `/refresh-token`

**请求参数**:
```json
{
  "token": "原Token"
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "新Token"
  }
}
```

### 2. 用户注册相关

#### 2.1 用户注册

**接口地址**: `POST /register`

**请求参数**:
```json
{
  "name": "用户名",
  "password": "密码",
  "token": "注册邮箱成功时传递的Token",
  "email": "注册邮箱"
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1
  }
}
```

### 3. 用户信息相关

#### 3.1 获取当前用户信息

**接口地址**: `GET /user`

**请求参数**: 无

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "用户名",
    "role": "角色",
    "nickname": "昵称",
    "introduction": "简介",
    "avatar": "头像文件链接",
    "email": "邮箱",
    "themeId": 1,
    "question_box_perm": "提问箱权限"
  }
}
```

#### 3.2 根据ID获取用户信息

**接口地址**: `GET /info/user?id={id}`

**请求参数**:
```
id: 用户ID
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "用户名",
    "role": "角色",
    "nickname": "昵称",
    "introduction": "简介",
    "avatar": "头像文件链接"
  }
}
```

#### 3.3 更新用户信息

**接口地址**: `PUT /user`

**请求参数**:
```json
{
  "nickname": "昵称",
  "introduction": "简介",
  "avatar": "头像文件",
  "themeId": 1
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1
  }
}
```

#### 3.4 更新密码

**接口地址**: `PUT /user/password`

**请求参数**:
```json
{
  "email": "邮箱",
  "code": "验证码",
  "password": "新密码"
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1
  }
}
```

#### 3.5 发送验证码

**接口地址**: `POST /user/send-code`

**请求参数**:
```json
{
  "email": "邮箱",
  "type": "方式"
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "msg": "消息"
  }
}
```

#### 3.6 忘记密码

**接口地址**: `POST /user/forget-password`

**请求参数**:
```json
{
  "email": "邮箱",
  "code": "验证码",
  "password": "新密码"
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1
  }
}
```

### 4. 教师信息相关

#### 4.1 获取教师列表

**接口地址**: `GET /info/teacher`

**请求参数**: 无

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "teachers": [
      {
        "id": 1,
        "responses": 0,
        "name": "老师名字",
        "avatarUrl": "老师头像链接",
        "introduction": "老师简介",
        "email": "老师邮箱",
        "perm": "提问箱权限"
      }
    ]
  }
}
```

#### 4.2 获取指定教师的问题列表

**接口地址**: `GET /info/teacher/pin?teacher_id={teacher_id}`

**请求参数**:
```
teacher_id: 教师ID
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "question_list": [
      {
        "id": 1,
        "title": "问题标题",
        "contents": "问题内容",
        "views": 0,
        "created_at": 1234567890,
        "image_urls": ["图片链接1", "图片链接2"]
      }
    ]
  }
}
```

#### 4.3 更新教师提问箱权限

**接口地址**: `PUT /teacher/perm`

**请求参数**:
```json
{
  "perm": "权限(public, private, protected)"
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1
  }
}
```

### 5. 问题相关

#### 5.1 添加问题

**接口地址**: `POST /questions/add`

**请求参数**:
```json
{
  "dst_user_id": 1,
  "title": "问题标题",
  "content": "问题内容",
  "is_private": 0,
  "files": ["文件1", "文件2"]
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1
  }
}
```

#### 5.2 获取公开问题列表

**接口地址**: `GET /questions/public?sort_type={sort_type}&page={page}`

**请求参数**:
```
sort_type: 排序类型(0-3)
page: 页码(从1开始)
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "question_list": [
      // 公开问题列表
    ],
    "remain_page": 5
  }
}
```

#### 5.3 搜索公开问题关键字

**接口地址**: `GET /questions/public/keywords?keyword={keyword}&sort_type={sort_type}`

**请求参数**:
```
keyword: 关键字
sort_type: 排序类型(0-3)
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "words": [
      {
        "value": "关键字"
      }
    ]
  }
}
```

#### 5.4 根据关键字获取公开问题列表

**接口地址**: `GET /questions/public/search?keyword={keyword}&sort_type={sort_type}&page={page}`

**请求参数**:
```
keyword: 关键字
sort_type: 排序类型(0-3)
page: 页码(从1开始)
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "question_list": [
      // 公开问题列表
    ],
    "remain_page": 5
  }
}
```

#### 5.5 获取教师问题列表

**接口地址**: `GET /questions/teacher?sort_type={sort_type}&page={page}&teacher_id={teacher_id}`

**请求参数**:
```
sort_type: 排序类型(0-3)
page: 页码(从1开始)
teacher_id: 教师ID
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "question_list": [
      // 教师问题列表
    ],
    "remain_page": 5
  }
}
```

#### 5.6 搜索教师问题关键字

**接口地址**: `GET /questions/teacher/keywords?keyword={keyword}&sort_type={sort_type}&teacher_id={teacher_id}`

**请求参数**:
```
keyword: 关键字
sort_type: 排序类型(0-3)
teacher_id: 教师ID
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "words": [
      {
        "value": "关键字"
      }
    ]
  }
}
```

### 6. 教师个人问题管理

#### 6.1 获取对我的提问

**接口地址**: `GET /teacher/question/all?page={page}&sort_type={sort_type}`

**请求参数**:
```
page: 页码(从1开始)
sort_type: 排序类型(0-3)
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "question_list": [
      // 问题列表
    ],
    "remain_page": 5
  }
}
```

#### 6.2 获取对我的提问的关键字

**接口地址**: `GET /teacher/question/keywords?keyword={keyword}&sort_type={sort_type}`

**请求参数**:
```
keyword: 关键字
sort_type: 排序类型(0-3)
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "words": [
      {
        "value": "关键字"
      }
    ]
  }
}
```

#### 6.3 搜索对我的提问

**接口地址**: `GET /teacher/question/search?keyword={keyword}&sort_type={sort_type}&page={page}`

**请求参数**:
```
keyword: 关键字
sort_type: 排序类型(0-3)
page: 页码(从1开始)
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "question_list": [
      // 问题列表
    ],
    "remain_page": 5
  }
}
```

#### 6.4 获取未回复提问

**接口地址**: `GET /teacher/question/unanswered?page={page}&sort_type={sort_type}`

**请求参数**:
```
page: 页码(从1开始)
sort_type: 排序类型(0-3)
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "question_list": [
      // 未回复问题列表
    ],
    "remain_page": 5
  }
}
```

#### 6.5 获取已回复提问

**接口地址**: `GET /teacher/question/answered?page={page}&sort_type={sort_type}`

**请求参数**:
```
page: 页码(从1开始)
sort_type: 排序类型(0-3)
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "question_list": [
      // 已回复问题列表
    ],
    "remain_page": 5
  }
}
```

#### 6.6 获取置顶提问

**接口地址**: `GET /teacher/question/top`

**请求参数**: 无

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "question_list": [
      // 置顶问题列表
    ],
    "remain_page": 5
  }
}
```

#### 6.7 置顶提问

**接口地址**: `POST /teacher/question/pin`

**请求参数**:
```json
{
  "question_id": 1
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_pinned": true
  }
}
```

### 7. 回答相关

#### 7.1 获取问题回复

**接口地址**: `GET /answer?question_id={question_id}`

**请求参数**:
```
question_id: 问题ID
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "question": {
      // 问题详情
    },
    "answer_list": [
      // 回答列表
    ],
    "can_reply": true
  }
}
```

#### 7.2 点赞回复

**接口地址**: `POST /answer/upvote`

**请求参数**:
```json
{
  "question_id": 1,
  "answer_id": 1
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_upvoted": true,
    "upvote_num": 10
  }
}
```

#### 7.3 添加回答

**接口地址**: `POST /answer/add`

**请求参数**:
```json
{
  "question_id": 1,
  "in_reply_to": 1,
  "content": "回答内容",
  "files": ["文件1", "文件2"]
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1
  }
}
```

### 8. 收藏相关

#### 8.1 获取收藏列表

**接口地址**: `GET /favorites?sort_type={sort_type}&page={page}`

**请求参数**:
```
sort_type: 排序类型(0-1)
page: 页码(从1开始)
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "favorite_list": [
      // 收藏列表
    ],
    "remain_page": 5
  }
}
```

### 9. 历史记录相关

#### 9.1 获取历史记录列表

**接口地址**: `GET /history?sort_type={sort_type}&page={page}`

**请求参数**:
```
sort_type: 排序类型(0-3)
page: 页码(从1开始)
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "favorite_list": [
      // 历史记录列表
    ],
    "remain_page": 5
  }
}
```

#### 9.2 搜索历史记录关键字

**接口地址**: `GET /history/keywords?keyword={keyword}&sort_type={sort_type}`

**请求参数**:
```
keyword: 关键字(长度2-100)
sort_type: 排序类型(0-3)
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "words": [
      {
        "value": "关键字"
      }
    ]
  }
}
```

### 10. 通知相关

#### 10.1 获取通知

**接口地址**: `GET /notification?user_id={user_id}`

**请求参数**:
```
user_id: 用户ID
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "new_question": [
      // 新问题通知
    ],
    "new_reply": [
      // 新回复通知
    ],
    "new_answer": [
      // 新回答通知
    ]
  }
}
```

#### 10.2 更新已读信息

**接口地址**: `PUT /notification`

**请求参数**:
```json
{
  "id": 1
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "is_read": true
  }
}
```

#### 10.3 删除提醒

**接口地址**: `DELETE /notification`

**请求参数**:
```json
{
  "id": 1
}
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

#### 10.4 获取提醒数目

**接口地址**: `GET /notification/count?user_id={user_id}`

**请求参数**:
```
user_id: 用户ID
```

**响应参数**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "new_question_count": 0,
    "new_reply_count": 0,
    "new_answer_count": 0
  }
}
```