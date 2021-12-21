# windows系统搭建kafka
### 1. 首先需要安装JDK
##### https://www.oracle.com/java/technologies/downloads/
##### 下载安装完成后 添加下面两个系统环境变量
##### JAVA_HOME: D:\Java\jdk1.8.0_291 (jdk的安装路径)
##### Path: 增加 %JAVA_HOME%\BIN;
##### 添加完成后在CMD中执行 java -version检测配置是否成功

### 2. 安装zookeeper
##### 运行kafka之前需要先运行zookeeper
##### http://zookeeper.apache.org/releases.html
##### zookeeper安装好后 把zoo_sample.cfg重命名成zoo.cfg 
##### zoo.cfg  中dataDir的值改成“./zookeeper-3.4.13/data”
##### 添加系统环境变量 ZOOKEEPER_HOME: E:\apache-zookeeper-3.7.0 (zookeeper目录)
##### Path: 在现有的值后面添加 %ZOOKEEPER_HOME%\bin
##### 添加完成后再CMD中执行 zkserver

### 3. 安装kafka
##### http://kafka.apache.org/downloads.html
##### 建议安装2.12-2.8.0版本
##### 2.12-3.0及以上 在win10无法正常启动 日志无法写入
##### 安装完成后 config/server.properties 中 log.dirs的值改成 “./logs”
##### 运行kafka .\bin\windows\kafka-server-start.bat .\config\server.properties
 
 
### 4. 相关命令
##### 运行kafka之前需要先启动zookeeper
##### 创建topic
#####  .\kafka-topics.bat --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test
##### 创建producer
#####  .\kafka-console-producer.bat --broker-list localhost:9092 --topic test
#####  创建 consumer
#####  .\kafka-console-consumer.bat --bootstrap-server localhost:9092 --topic test --from-beginning
