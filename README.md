# oauth2-server
Go OAuth2 Web Server

```bash
http://localhost:8080/oauth2/authorize?client_id=oauth2_client&redirect_uri=http://localhost/callback&response_type=code&state=somestate&scope=read_write
```

```bash
http://localhost:8080/oauth2/authorize?client_id=oauth2_client&response_type=code&state=somestate&scope=read_write&redirect_uri=http%3a%2f%2flocalhost%2fcallback
```