# 环境变量配置：

- SMTP_RANCHER_TO: 定期接收消息邮箱
- SMTP_USER: 服务邮箱用户名
- SMTP_PWD: 服务邮箱用户密码
- SMTP_ENDPOINT: SMTP端点
- SMTP_PORT: 服务邮箱端口
- PDF_URL: PDF文件地址
- PDF_PWD: PDF文件密码
- DB_HOST_IP: 数据库地址
- DB_USERNAME: 数据库用户名
- DB_PASSWORD:  数据库密码
- DB_NAME: 数据库名
- ALI_REGION_ID: 阿里云区域
- ALI_ACCESS_KEYID: 阿里云访问令牌
- ALI_ACCESS_SECRET: 阿里云访问令牌密钥
- ALI_SIGN_NAME: 阿里云短信签名
- ALI_TEMPLATE_CODE: 阿里云短信模板

# 建表SQL语句：

```
CREATE TABLE `user` (
  `uid` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `company` varchar(255) DEFAULT NULL,
  `position` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `savetime` datetime DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=77 DEFAULT CHARSET=utf8mb4;
```
