# exporter

## mongodb-exporter
dokcer启动命令
```
docker run -d -e MONGODB_URI="mongodb://账号:密码@ip:端口/?authSource=admin"  mongodb-exporter:latest 
```
kubectl 启动命令
```
修改对应的字段
kubectl apply -f mongo-exporter.yaml
```
## 页面
暴露的端口：8080

# 采集的指标

### 慢操作记录
只会获取system.profile7天的记录，并不是开源的exporter增日志记录，有所区别。

### 连接数

### 健康状态
存活：1 死亡：0

### IOPS
IOPS使⽤量=data_iops+log_iop

### 内存
物理内存, 虚拟内存

### 磁盘
MongoDB节点的空间使用量由data_size和log_size组成，即ins_size=data_size+log_size
