# dagger-server
websocket代理服务器

### 自动安装
```
curl -fsSL  https://raw.githubusercontent.com/midoks/dagger/main/dagger-server/scripts/install.sh | sh
```

### DEV

```
curl -fsSL  https://raw.githubusercontent.com/midoks/dagger/main/dagger-server/scripts/install_dev.sh | sh
```

# Nginx配置

```
location / {
   proxy_redirect off;
   proxy_pass http://127.0.0.1:12345;
   proxy_http_version 1.1;
   proxy_set_header Upgrade $http_upgrade;
   proxy_set_header Connection "upgrade";
   proxy_set_header Host $http_host;
}
```

# 用户模式
```
./dagger-server user -m add -u test -p test
./dagger-server user -m delete -u test
./dagger-server user -m mod -u test -p test
./dagger-server user -m list
./dagger-server user -m query -u test
```