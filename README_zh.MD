![](https://img.shields.io/github/stars/milkomeda-org/doc-scanner)
![](https://img.shields.io/github/issues/milkomeda-org/doc-scanner)
![](https://img.shields.io/github/license/milkomeda-org/doc-scanner)
# 内容扫描器
### 使用方法
在接口层编写匹配注释，执行程序即可
.js .go .java .php文件目前可用于扫描，其他类型请自行添加

### 脚本参数说明
| 参数名 | 用途     | 类型
|--------|----------|----------|
|type| 文档系统类型,参数值为枚举showdoc/yapi|string
|url| 上传地址|string
|key| 项目api_key|string
|token| 项目api_token|string
|speed| 上传请求并发速度|int
|cat| 前置固定目录的base64|string
|searchPath | 目录/文件|string
|templatePath| 模板文件路径|string
|ask| 是否先询问再继续执行|boolean

### 注释示例
```
/**
 * #doc
 * @class
 * @title
 * @description
 * @method post
 * @query id 必选 int queryID
 * @header id 必选 int queryID
 * @body id 必选 int queryID
 * @url /api/update
 * @return {"code":0,"msg":"ok","data":"success","ts":"2020-10-30T03:59:45.591Z"}
 * @return_param id int query_id
 * @remark query request
 */
```

### 注释说明
> 注释以/**开头 */结尾 并且注释中包含#doc就会被执行文档解析

| 参数名 | 用途     |
|--------|----------|
|@class| 分类
|@title| 名称
|@description| 描述
|@method| 方法
|@url| 地址
|@header| 头参数
|@query| query参数
|@body| body参数
|@return| 返回值示例
|@return_param| 返回值说明
|@remark| 备注
|@number| 编号

### 已支持文档系统
- ShowDoc
- Yapi

### 即将支持
- Rap


### 提示
@header @query @body @return_param 按照参数名 是否必选 类型 描述顺序编写参数，使用空格分割

JetBrains IDE可使用自定义Live Template添加如下内容，快捷生成模板
```
/**
 * #doc
 * @class $class$
 * @title $title$
 * @description $description$
 * @method $method$
 * @url /$path$
 * @return {"code":0,"msg":"ok","data":"success","ts":"2020-10-30T03:59:45.591Z"}
 * @remark $remark$
 */
```

Yapi文档系统的cat为项目分类数字id，从分类url里获取，非数字无法上传


### 在Intellij平台和Interactive Assistant一起使用
0.通过Release下载可执行文件或自行编译可执行文件
1.通过Settings->Plugin->安装Interactive Assistant
2.通过Settings->Other Settings->Interactive Assistant->
3.复制以下配置模板到Map选项卡的输入框:
```
{
  "$type": "exe",
  "$args_name": "searchPath",
  "$path": "\"C:\\Users\\Administrator\\Desktop\\doc-scanner\\execute.exe\"",
  "url": "http://172.16.2.101:4999/server/index.php?s\u003d/api/item/updateByApi",
  "key": "e9f0bdd396a768399c63ef86d70ccc322044412143",
  "token": "834e06eb69e21565d997cf15a1159da21794468976",
  "speed": "10",
  "cat": "NDE3NzY\u003d ",
  "templatePath": "\"C:\\Users\\Administrator\\Desktop\\doc-scanner\\template.txt\"",
  "ask": "false",
  "type": "showdoc"
}
```
3.然后在右键菜单和编辑菜单->Abilities便可使用
