[](https://img.shields.io/github/stars/milkomeda-org/doc-scanner)
[](https://img.shields.io/github/issues/milkomeda-org/doc-scanner)
[](https://img.shields.io/github/license/milkomeda-org/doc-scanner)
# Content Scanner
### Method Of Use
Write matching comments at the interface layer and execute the program

Js.go.java.php file is currently available for scanning. Please add other types to compile

### Script parameter description
> -type (the document type)
  -url (upload the address)
  -key (project api key)
  -token (project api token)
  -speed (number of concurrent upload requests)
  -cat (base 64 encoding of fixed directory, directory using / split)
  -searchPath (the directory or file in which the scan is performed)
  -templatePath (template file path)

### The annotation examples
```
/**
 * #doc
 * @class
 * @title
 * @description 
 * @method
 * @query id must int query_id
 * @url /api/update
 * @return {"code":0,"msg":"ok","data":"success","ts":"2020-10-30T03:59:45.591Z"}
 * @return_param id int query_id
 * @remark query request
 */
```

### Explanatory Note
> Comments that begin with /** and end with */ and contain #doc are then parsed by the document

| parameter name | use     |
|--------|----------|
|@class| classify
|@title| name
|@description| describe
|@method| method
|@url| request path
|@header| The First Parameter
|@query| the Query Parameter
|@body| body Parameters
|@return| return Value Example
|@return_param| returned Value Specification
|@remark| remark
|@number| serialNumber

### Document system supported
- ShowDoc

### To support
- Yapi
- Rap

### Hint
@header and @query and @body and @return_param writes the parameters in the order of the corresponding table in the template, separated by Spaces

The JetBrains IDE adds the following to generate templates quickly using a custom Live Template
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
