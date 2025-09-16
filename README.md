 # 短链接项目
 
## 搭建项目的骨架

## 1. 建库建表
## 2. 搭建go-zero框架的骨架
### 2.1. 生成api代码
   goctl api go -api shorturl.api -dir .
### 2.2. 生成model代码
```bash
   goctl model mysql datasource -url="root:root@tcp(127.0.0.1:3306)/shortener" -table="short_url_map" -dir="./model" -c
   goctl model mysql datasource -url="root:root@tcp(127.0.0.1:3306)/shortener" -table="sequence" -dir="./model" -c
```