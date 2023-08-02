# Setup Instructions
1. Export ENV or Add to .env
```
export APP_ADDR=0.0.0.0:8080
RPC_API_URL=https://rpc.canxium.org
```

2. Build & start API
`go build cmd/api/main`
`./main`

3. Get supply
`curl --location --request GET '127.0.0.1:8080/info/cau'`


3. Get total supply
`curl --location --request GET '127.0.0.1:8080/info/cau?q=totalSupply'`


3. Get circulating supply
`curl --location --request GET '127.0.0.1:8080/info/cau?q=circulatingSupply'`
