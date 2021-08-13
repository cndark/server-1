# 基础说明

* 语言主要使用 golang, nodejs
* 游戏进程采用长连接
* 通信协议采用 protobuf
* 数据库使用 mongodb 4.2

# 环境变量说明

变量名        |值                 |说明
---           |---               |---
PERF_MON      |true              |进程日志中显示监控条目
LOG_LEVEL_INFO|true              |进程日志等级设置为Info
DEV_MODE      |true              |设置为开发模式
PROJ_NAME     |比如: p1          |
WORK_DIR      |/data/${PROJ_NAME}|设置部署工作路径

# 部署说明

* 系统配置
    - alias, vimrc, selinux, limits, sudoers, sshd
    - wget, curl, nodejs, mongodb
    - update openssl
    - diskmount
    - useradd, sshkey
    - net.ipv4.ip_local_port_range
* 在 /etc/bashrc 中设置好以上环境变量
* 做镜像时记住修改 sshd_config 中的 MaxStartups 为 100:30:100
* 做镜像时, 不要把私钥保留在镜像中
