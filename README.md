# oauth2-server
Go OAuth2 Web Server

# OAuth2 Library

[OAuth2](https://github.com/nilorg/oauth2)

# Code 授权模式
```bash
http://localhost:8080/oauth2/authorize?client_id=oauth2_client&redirect_uri=http://localhost/callback&response_type=code&state=somestate&scope=read_write
```

# Token 授权模式
```bash
http://localhost:8080/oauth2/authorize?client_id=oauth2_client&redirect_uri=http://localhost/callback&response_type=token&state=somestate&scope=read_write
```
