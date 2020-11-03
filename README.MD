# 自动化文档生成
### 使用方法
在接口层编写相关注释，运行脚本即可
当前提供扫描.js .go文件

### 脚本参数说明
> -url 上传地址
  -key 项目api_key
  -token 项目api_token
  -speed 上传请求并发速度
  -cat 固定目录的base64
  -searchPath 目录/文件
  -templatePath 模板文件路径

### 注释示例
```
/**
 * #doc
 * @class 曲谱
 * @title 更新
 * @description 更新曲谱信息
 * @method post
 * @query id 必选 int 曲谱ID
 * @header id 必选 int 曲谱ID
 * @body id 必选 int 曲谱ID
 * @url https://www.showdoc.com.cn/api/music-score/update
 * @return {"code":0,"msg":"ok","data":"success","ts":"2020-10-30T03:59:45.591Z"}
 * @return_param id int 曲谱ID
 * @remark 更新曲谱
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


### 提示
@header @query @body @return_param 按照模板对应表格顺序编写参数，使用空格分割
