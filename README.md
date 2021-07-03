# Database
```sql
CREATE ROLE yoel LOGIN PASSWORD 'cascadesheet'
CREATE DATABASE shopping OWNER yoel
```

```bash
sudo vim /etc/systemd/system/shopping.service
```

```
[Unit]
Description="API service de GO para el sistema de venta shopping"

[Service]
ExecStart=/home/ubuntu/shopping/api/shopping
WorkingDirectory=/home/ubuntu/shopping/api/
User=ubuntu
Restart=always

[Install]
WantedBy=multi-user.target
```


activar el servicio
```bash
sudo systemctl enable shopping.service
```

iniciar el servicio
```bash
sudo systemctl start shopping.service
```

verificar el servicio
```bash
sudo systemctl status shopping.service
```

PROXY
```bash
server {
        listen 80;
        server_name api.yarbys.com;

        location / {
                proxy_pass http://127.0.0.1:1323;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection 'upgrade';
                proxy_set_header Host $host;
                proxy_cache_bypass $http_upgrade;
        }
}
```

Check
```bash
nginx -t
```

```bash
ln -s /etc/nginx/sites-available/api.yarbys.com /etc/nginx/sites-enabled/api.yarbys.com
```

```bash
sudo service nginx restart
systemctl status nginx
```