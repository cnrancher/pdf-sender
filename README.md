# PDF-SENDER

一个用来发送 Rancher 中文文档的服务

## 启动配置

```text
NAME:
   pdf-sender - Send pdf documents to our lovely users.

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   v0.0.0-dev

COMMANDS:
   config   Print out the merged configuration
   run      Run pdf sender server
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug                        Enable Debug log for pdf sender
   --config-file value, -f value  (default: "/etc/pdf-sender.yml") [$CONFIG_FILE]
   --help, -h                     show help
   --version, -v                  print the version
```

配置文件示例可以参考程序默认配置 `package/init.yml`，配置文件内容如下：

```yaml
documents:
  - name: k3s
    title: K3s
    filename: K3s.pdf
    kind:
      - pdf
  - name: harvester
    title: Harvester
    filename: Harvester.pdf
    kind:
      - pdf
  - name: octopus
    title: Octopus
    filename: Octopus_CN_Doc.pdf
    kind:
      - pdf
  - name: rancher2
    title: Rancher2.x
    filename: Rancher2.x_CN_Doc.pdf
    kind:
      - pdf
  - name: rancher2.5
    title: Rancher2.5
    filename: Rancher2.5_CN_Doc.pdf
    kind:
      - pdf
  - name: rancher1
    title: Rancher1.6
    filename: rancher1.6.pdf
    kind:
      - pdf
  - name: rancher-ent
    title: 白皮书
    filename: Rancher企业版拓展功能概述及购买方式(2021).pdf
    ossPathPrefixOverride: RancherEntPDF
    kind:
      - ent
kinds:
  pdf:
    header: |
      您好，

      您可以通过下面的链接下载 Rancher 中文文档。

    footer: |
      Best Regards,
      SUSE/Rancher 中国团队
    subject: Rancher 2.x 中文文档
    senderName: Rancher Labs 中国
    description: "中文文档"
  ent:
    subject: 「Rancher 企业版拓展功能概述及购买方式」白皮书下载
    header: |
      您好，

      您可以通过下面的链接下载 「Rancher企业版拓展功能概述及购买方式」
    footer: |
      Best Regards,
      SUSE/Rancher 中国团队
    description: "企业版白皮书"
    senderName: Rancher Labs 中国
```

## 建表 SQL 语句

当前版本建表语句如下：

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
  `kind` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=77 DEFAULT CHARSET=utf8mb4;
```

如从 v0.2.x 升级到 v0.3.0，则需要对表进行更变

```sql
ALTER TABLE pdf.`user` ADD kind varchar(20) NULL;
```

## 后续版本规划

- 上传周期生成的文件到 oss
- 如果需要的情况下，做高可用支持
- 根据配置的结构自动生成命令行参数
