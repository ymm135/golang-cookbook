# nginx基础
## [官网](https://nginx.org/en/)  

## [快速入门](https://nginx.org/en/docs/beginners_guide.html)    

## 基础知识  

### 查看nginx版本及编译参数  
```sh
nginx -V
nginx version: nginx/1.24.0
built by gcc 9.4.0 (Ubuntu 9.4.0-1ubuntu1~20.04.1) 
built with OpenSSL 1.1.1f  31 Mar 2020
TLS SNI support enabled
configure arguments: --prefix=/usr/share/nginx --sbin-path=/usr/sbin/nginx --conf-path=/etc/nginx/nginx.conf --error-log-path=/var/log/nginx/error.log --http-log-path=/var/log/nginx/access.log --with-http_ssl_module --with-http_stub_status_module --with-http_realip_module --with-http_auth_request_module --with-http_v2_module
```

### nginx配置测试`nginx -t`  
```sh
nginx -t
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

### 常用指令  

Nginx 常用的命令包括：

1. **启动 Nginx**:
   ```
   sudo nginx
   ```

2. **停止 Nginx**:
   - 快速停止：
     ```
     sudo nginx -s stop
     ```
   - 优雅停止：
     ```
     sudo nginx -s quit
     ```

3. **重新加载配置文件**:
   ```
   sudo nginx -s reload
   ```

4. **重新打开日志文件**:
   ```
   sudo nginx -s reopen
   ```

5. **检查配置文件**:
   ```
   sudo nginx -t
   ```

6. **查看 Nginx 版本和编译配置**:
   ```
   nginx -v
   nginx -V
   ```

这些命令涵盖了 Nginx 的基本操作，包括启动、停止、配置管理等。在使用这些命令时，确保你有足够的权限（如使用 `sudo`）。  

