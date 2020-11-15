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