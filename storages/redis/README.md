To launch a container with Redis, run the command from the `storages/redis` directory:
```bash
docker-compose up -d
```

To access Redis from a container, use the following commands:
```
docker exec -it <container id> sh
redis-cli
auth <password>
```