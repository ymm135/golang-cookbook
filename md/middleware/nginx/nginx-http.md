- # nginx-http-core  

目录:  
- [基础知识](#基础知识)
  - [http 模块](#http-模块)
  - [server 模块](#server-模块)
  - [nginx配置测试`nginx -t`](#nginx配置测试nginx--t)
  - [项目一般配置](#项目一般配置)
- [疑问及拓展](#疑问及拓展)
  - [HTTP 严格传输安全 (HSTS) 和 NGINX](#http-严格传输安全-hsts-和-nginx)
  - [websocket配置](#websocket配置)


## 基础知识 

https://nginx.org/en/docs/http/ngx_http_core_module.html  

### http 模块  
https://nginx.org/en/docs/http/ngx_http_core_module.html#http

### server 模块  
https://nginx.org/en/docs/http/ngx_http_core_module.html#server  

Sets configuration for a virtual server. There is no clear separation between IP-based (based on the IP address) and name-based (based on the “Host” request header field) virtual servers. Instead, the listen directives describe all addresses and ports that should accept connections for the server, and the server_name directive lists all server names. Example configurations are provided in the “How nginx processes a request” document.  

- ### https://nginx.org/en/docs/http/request_processing.html  

Nginx 处理一个请求的流程大致如下：

1. **接收请求：** Nginx 首先接收到客户端的 HTTP 请求。

2. **寻找匹配的服务器块：** Nginx 通过请求中的主机名（Host 头）和端口号来寻找相匹配的 `server` 块。

3. **选择 location 块：** 在找到的 `server` 块中，Nginx 根据请求的 URI 选择最佳匹配的 `location` 块。

4. **执行请求：** 在 `location` 块中，Nginx 根据配置执行相关操作，比如代理传递、重定向或返回静态文件内容。

5. **生成响应：** Nginx 处理请求后，生成 HTTP 响应并发送给客户端。

Nginx 的配置极其灵活，可以根据需要进行高度定制。更详细的处理流程和配置选项可以在 [Nginx 官方文档](https://nginx.org/en/docs/http/request_processing.html)中找到。  



### nginx配置测试`nginx -t`  
```sh
nginx -t
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

### 项目一般配置  
```sh
user  root;
worker_processes  auto;

error_log  /var/log/nginx/error.log notice;
pid        /var/run/nginx.pid;

events {
    worker_connections  1024;
}


http {
    include       mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  65;

    ssl_protocols  TLSv1.2;
    ssl_ciphers 'ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!MD5:!PSK';

    #gzip  on;

    #include /etc/nginx/conf.d/*.conf;

    client_max_body_size 1024m;
    proxy_connect_timeout 3600; #单位秒
    proxy_send_timeout 3600; #单位秒
    proxy_read_timeout 3600; #单位秒

    server {
        listen       443  ssl;
        server_name  localhost;
        root  /opt/app/frontend;
	
        # ssl证书
        ssl_certificate /etc/nginx/crt/ssl.crt;
        ssl_certificate_key /etc/nginx/crt/ssl.key;
        ssl_ciphers 'ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!MD5:!PSK';

        location / {
            try_files $uri $uri/ /index.html;
            index  index.html;
            add_header Cache-Control no-cache;
        }

        location /api {
            proxy_set_header Host $http_host;
            proxy_set_header  X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

            proxy_pass http://127.0.0.1:8000;
            proxy_cookie_path / /api;
            proxy_redirect default;
            rewrite ^/api/(.*) /$1 break;
            client_max_body_size 500m;
        }

        location ^~ /ws/  {
            proxy_pass http://127.0.0.1:8000/ws/;
            proxy_set_header  X-Real-IP  $remote_addr;
            proxy_set_header Host $host:8000;
            proxy_http_version 1.1;
            proxy_set_header Connection keep-alive;
            proxy_set_header Keep-Alive 600;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_connect_timeout 60;
            proxy_read_timeout 600;
        }
    }
}
```

> 其中包含资源配置/后端配置/websocket配置   

## 疑问及拓展  
### HTTP 严格传输安全 (HSTS) 和 NGINX  

add_header指令的继承规则
NGINX 配置块add_header从其封闭块继承指令，因此您只需将该add_header指令放置在顶级server块中即可。有一个重要的例外：如果块add_header本身包含指令，它不会从封闭块继承标头，并且您需要重新声明所有add_header指令：  

> 意思就是顶级和子级的配置是独立,需要独立配置  

漏洞详情:
```sh
Attack Details
URLs where HSTS is not enabled:
https://10.25.30.126/
https://10.25.30.126/crossdomain.xml
https://10.25.30.126/sitemap.xml.gz
https://10.25.30.126/smooth
https://10.25.30.126/login
https://10.25.30.126/clientaccesspolicy.xml
https://10.25.30.126/sitemap.xml
```

> 导致问题的原因在于静态资源配置  

```sh
        location / {
            try_files $uri $uri/ /index.html;
            index  index.html;
            add_header Cache-Control no-cache;
        }
```

所以需要在`location /`块中增加`add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;`  
```sh
        location / {
            try_files $uri $uri/ /index.html;
            index  index.html;
            add_header Cache-Control no-cache;
            add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
        }
```

重启nginx后,查看请求头:  

```sh
$ curl -Ik https://10.25.30.126
HTTP/1.1 200 OK
Server: nginx/1.24.0
Date: Tue, 30 Jan 2024 15:35:20 GMT
Content-Type: text/html
Content-Length: 10961
Last-Modified: Fri, 26 Jan 2024 08:27:15 GMT
Connection: keep-alive
ETag: "65b36ce3-2ad1"
Cache-Control: no-cache
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
Accept-Ranges: bytes
```

> `Strict-Transport-Security: max-age=31536000; includeSubDomains; preload` 

### websocket配置  
```sh
  location ^~ /ws/  {
            proxy_pass http://127.0.0.1:8000/ws/;
            proxy_set_header  X-Real-IP  $remote_addr;
            proxy_set_header Host $host:8000;
            proxy_http_version 1.1;
            proxy_set_header Connection keep-alive;
            proxy_set_header Keep-Alive 600;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_connect_timeout 60;
            proxy_read_timeout 600;
        }
```

在 Nginx 配置中关于 WebSocket 的设置主要通过 `location` 块来实现。这个配置段的作用和原理如下：

1. **`proxy_pass` 指令：** 将请求转发到指定的内部服务器地址，这里是 `http://127.0.0.1:8000/ws/`。这意味着所有匹配 `/ws/` 的请求都将被转发到本地的 `8000` 端口上的 WebSocket 服务。

2. **设置头部信息：** 通过 `proxy_set_header` 指令设置头部信息，以确保 WebSocket 的正确功能。例如，`X-Real-IP` 设置为客户端的 IP 地址，`Upgrade` 和 `Connection` 头部被设置以支持 WebSocket 协议的升级机制。

3. **保持连接活跃：** `keep-alive` 和 `Keep-Alive` 头部用于维持长连接，这对于 WebSocket 连接至关重要。

4. **超时设置：** `proxy_connect_timeout` 和 `proxy_read_timeout` 指定了连接和读取的超时时间。

这种配置允许 Nginx 作为反向代理处理 WebSocket 连接，确保 WebSocket 请求被正确地路由和处理。
