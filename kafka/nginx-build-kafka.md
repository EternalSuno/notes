# linux系统安装kafka扩展 - php
### 1. linux使用kafka需要安装librdkafka
##### git clone https://github.com/edenhill/librdkafka.git
##### cd librdkafka/
##### make && make install

### 2. 安装php扩展
##### git clone https://github.com/arnaud-lb/php-rdkafka.git
##### cd php-rdkafka/
##### /php路径/bin/phpize
##### ./configure --with-php-config=/php路径/bin/php-config
##### make && make install
##### php.ini中添加: extension = rdkafka.so
##### nginx -s reload
##### kill -USR2 PHP-FPM的PID
      



