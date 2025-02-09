<h1>DataDog Trace Debug</h1>

```bash
docker run -p 8866:8866 sirhawell/ddtd
```

in your instrumented code, start tracer with address = localhost:8866

e.g
```go
tracer.Start(tracer.WithAgentAddr("localhost:8866"))
```

open a browser window at localhost:8866

click refresh to see traces

<h1>Build/Run from Source</h1>

for agent:
```
export HTTP_ADDRESS=:8866
export HTTP_PANEL_ROOT=./panel/.output/public
go run ./cmd
```

for panel:
```
cd ./panel
npm install
npm run dev
```
