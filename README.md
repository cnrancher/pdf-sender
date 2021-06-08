# PDF-SENDER

一个用来发送 Rancher 中文文档的服务

## 环境变量配置

- SMTP_RANCHER_TO_DAY: 每日定期接收消息邮箱
- SMTP_RANCHER_TO_MON: 每月定期接收消息邮箱
- DAY_CRON: 每日定时时间
- MON_CRON: 每月定时时间

- SMTP_USER: 服务邮箱用户名
- SMTP_PWD: 服务邮箱用户密码
- SMTP_ENDPOINT: SMTP 端点
- SMTP_PORT: 服务邮箱端口
- SENDER_EMAIL: 服务邮箱端口

- Rancher2_PDF_URL: Rancher2.x PDF 文件地址
- Rancher2_PWD: Rancher2.x PDF 文件密码
- RKE_PDF_URL: RKE PDF 文件地址
- RKE_PWD: RKE PDF 文件密码
- RKE2_PDF_URL: RKE2 PDF 文件地址
- RKE2_PWD: RKE2 PDF 文件密码
- K3s_PDF_URL: K3s PDF 文件地址
- K3s_PWD: K3s PDF 文件密码
- Octopus_PDF_URL: Octopus PDF 文件地址
- Octopus_PWD: Octopus PDF 文件密码
- Harvester_PDF_URL: Harvester PDF 文件地址
- Harvester_PWD: Harvester PDF 文件密码

- DB_HOST_IP: 数据库地址
- DB_USERNAME: 数据库用户名
- DB_PASSWORD: 数据库密码
- DB_NAME: 数据库名

- ALI_REGION_ID: 阿里云区域
- ALI_ACCESS_KEYID: 阿里云访问令牌
- ALI_ACCESS_SECRET: 阿里云访问令牌密钥
- ALI_SIGN_NAME: 阿里云短信签名
- ALI_TEMPLATE_CODE: 阿里云短信模板

## 建表 SQL 语句

```sql
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

## 后续功能规划

1. 支持从配置文件中读取发送配置，发送配置包括
2. 构建镜像时将默认配置复制到容器中使用
3. 精简启动参数

支持配置文件设置所有从环境变量中读取的参数，包括但不限于

- 文档下载配置
- 阿里云 API 访问配置
- oss 配置
- 邮件发送配置
