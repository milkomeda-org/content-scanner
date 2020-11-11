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
|type| 文档系统类型|string
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
 * @query id 必选 int 曲谱ID
 * @header id 必选 int 曲谱ID
 * @body id 必选 int 曲谱ID
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

### 即将支持
- Yapi
- Rap


### 提示
@header @query @body @return_param 按照模板对应表格顺序编写参数，使用空格分割

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
