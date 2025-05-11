# 📈 Tracing with Uptrace via OpenTelemetry

This project uses [Uptrace](https://uptrace.dev) as a tracing backend with OpenTelemetry. Traces are automatically collected across the entire stack — including:

- HTTP API via middleware
- Usecase and business logic layers
- Database queries
- External calls (if instrumented)

Powered by [go-common](https://github.com/kittipat1413/go-common), tracing is already integrated and only requires basic setup to start collecting spans.

---

## 🚀 Getting Started with Uptrace

### 1. 📝 Sign up at [https://uptrace.dev](https://uptrace.dev)

- Create a project
- Copy your **DSN** from the project settings
  > Example: `'https://<token>@api.uptrace.dev?grpc=4317'`

---

### 2. 📄 Configure OpenTelemetry Collector

Create a file named `otel-collector-config.yaml` in the project root:

```yaml
exporters:
  debug:
  otlp/uptrace:
    endpoint: https://otlp.uptrace.dev:4317
    tls: { insecure: false }
    headers:
      uptrace-dsn: 'YOUR_UPTRACE_DSN'

receivers:
  otlp:
    protocols:
      grpc:
      http:

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [debug, otlp/uptrace]
```
> Replace `'YOUR_UPTRACE_DSN'` with your actual Uptrace DSN.

### 3. 🌐 Set OTEL endpoint in your environment

Make sure your application sets the following env variable:
```bash
export OTEL_COLLECTOR_ENDPOINT=127.0.0.1:4317
```
> You can set this in `env.yaml`, or via your deployment config.

### 4. 🐳 Start OpenTelemetry Collector

Add this to your `docker-compose.yaml`:
```yaml
services:
  otel-collector:
    image: otel/opentelemetry-collector-contrib:0.123.0
    volumes:
      - ./otel-collector-config.yaml:/etc/otelcol-contrib/config.yaml
    ports:
      - "4317:4317"
      - "4318:4318"
```
Start the collector:
```bash
docker compose up -d otel-collector
```

## ✅ You’re Done!

Now run your application and interact with the API.

Open your Uptrace dashboard — you should start seeing traces!

Below is an example of what a successful trace looks like:
![Uptrace Example1](/docs/images/uptrace1.png)
![Uptrace Example2](/docs/images/uptrace2.png)
![Uptrace Example3](/docs/images/uptrace3.png)
