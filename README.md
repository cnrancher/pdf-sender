# 环境变量配置：

- SMTP_RANCHER_TO_DAY: 每日定期接收消息邮箱
- SMTP_RANCHER_TO_MON: 每月定期接收消息邮箱
- DAY_CRON: 每日定时时间
- MON_CRON: 每月定时时间

- SMTP_USER: 服务邮箱用户名
- SMTP_PWD: 服务邮箱用户密码
- SMTP_ENDPOINT: SMTP端点
- SMTP_PORT: 服务邮箱端口
- SENDER_EMAIL: 服务邮箱端口

- Rancher2_PDF_URL: Rancher2.x PDF文件地址
- Rancher2_PWD: Rancher2.x PDF文件密码
- RKE_PDF_URL: RKE PDF文件地址
- RKE_PWD: RKE PDF文件密码
- K3s_PDF_URL: K3s PDF文件地址
- K3s_PWD: K3s PDF文件密码
- Octopus_PDF_URL: Octopus PDF文件地址
- Octopus_PWD: Octopus PDF文件密码
- Harvester_PDF_URL: Harvester PDF文件地址
- Harvester_PWD: Harvester PDF文件密码

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
