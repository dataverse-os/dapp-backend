# dapp-backend

为了简化开发者的理解门槛。有两个可选的做法：

1. 仓库下可以增设两个目录 1）链接到Ceramic节点的仓库   2）开发者部署 应用后端+Ceramic节点的 云脚本  

2. 单独开一个仓库，放置 应用后端+Ceramic节点的 云脚本  


完善readme

## Usage

### Use Docker Compose

```shell
#初始化ceramic目录
docker run -it --rm -v ~/.ceramic/:/root/.ceramic/ ceramicnetwork/js-ceramic:latest
```

```yaml
version: "3.9"
services:
  ceramic:
    image: ceramicnetwork/js-ceramic:latest
    volumes:
      - ~/.ceramic/:/root/.ceramic/
    healthcheck:
      test: ["CMD-SHELL", "curl -f http://localhost:7007/api/v0/node/healthcheck || exit 1"]
      interval: 1m30s
      timeout: 10s
      retries: 3
      start_period: 40s

  dapp-backend:
    image: dataverseos/dapp-backend:latest
    environment:
      - DID_PRIVATE_KEY={YOUR_PRIVATE_KEY_HERE}
      - CERAMIC_URL=http://ceramic:7007
    depends_on:
      - ceramic
```