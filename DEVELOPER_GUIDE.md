# BishOps Toolchain - Developer Guide

## Quick Start

### Prerequisites
- Go 1.23+ (`go version`)
- AWS credentials configured locally (`~/.aws/credentials`)
- Slack bot token (with `files:read`, `chat:write`, `files:write` scopes)

### Setup

1. **Clone & Install Dependencies**
   ```bash
   go mod download
   ```

2. **Generate App Key**
   ```bash
   go run . artisan key:generate
   ```

3. **Create `.env` File**
   ```bash
   cp .env.example .env
   ```

   Update with:
   ```bash
   APP_NAME=BishOps Tools
   APP_ENV=development
   APP_DEBUG=true
   AWS_PROFILE=<your-aws-profile>
   SLACK_BOT_TOKEN=xoxb-<your-token>
   MCP_PORT=8080
   ```

4. **Run the Server**
   ```bash
   go run .
   ```

   Server starts on `http://localhost:8080`

### Test MCP Integration

```bash
# In another terminal, test the MCP server
curl -X POST http://localhost:8080/health
```

---

## Development Commands

### Goravel Artisan Commands

```bash
# Generate controller
go run . artisan make:controller UserController

# Generate migration
go run . artisan make:migration create_users_table

# Run migrations
go run . artisan migrate

# Create seeder
go run . artisan make:seeder UserSeeder
```

### Building & Deployment

```bash
# Build binary
go build -o toolchain .

# Run binary
./toolchain

# Docker build
docker build -t bishops-toolchain:latest .

# Docker run
docker run -e AWS_PROFILE=default \
           -e SLACK_BOT_TOKEN=xoxb-... \
           -p 8080:8080 \
           -v ~/.aws:/root/.aws \
           bishops-toolchain:latest
```

---

## Key Files to Modify

### Adding a New MCP Tool

1. **Create Handler** (`app/handlers/my_handler.go`)
   ```go
   package handlers

   import (
       "context"
       "github.com/modelcontextprotocol/go-sdk/mcp"
   )

   type MyHandler struct{}

   func (h *MyHandler) Handle(ctx context.Context, input *MyInput) error {
       // Your logic here
       return nil
   }
   ```

2. **Create MCP Controller** (`app/mcp/controllers/my_tool.go`)
   ```go
   package controllers

   import (
       "context"
       "github.com/modelcontextprotocol/go-sdk/mcp"
   )

   type MyToolArguments struct {
       Param1 string `json:"param1" mcp:"description"`
   }

   func NewMyTool(handler *MyHandler) (*mcp.Tool, mcp.ToolHandlerFor[MyToolArguments, any]) {
       return &mcp.Tool{
           Name:        "my-tool",
           Description: "What this tool does",
       }, func(ctx context.Context, req *mcp.CallToolRequest, params MyToolArguments) (*mcp.CallToolResult, any, error) {
           // Call handler logic
           // Return results
       }
   }
   ```

3. **Register in Kernel** (`app/mcp/kernel.go`)
   ```go
   func Tools() {
       srv := facades.Mcp().Server()
       handler := &handlers.MyHandler{}
       tool, toolHandler := controllers.NewMyTool(handler)
       mcp.AddTool(srv, tool, toolHandler)
   }
   ```

### Adding AWS Service Access

The AWS package automatically exposes CloudFront and S3. To add a new AWS service:

1. **Update `packages/aws/client/aws.go`**
   ```go
   import "github.com/aws/aws-sdk-go-v2/service/myservice"

   func (a *AwsClient) MyService() *myservice.Client {
       if a.myserviceClient == nil {
           a.myserviceClient = myservice.NewFromConfig(a.cfg)
       }
       return a.myserviceClient
   }
   ```

2. **Update contracts** (`packages/aws/contracts/aws.go`)
   ```go
   type Aws interface {
       // ... existing methods
       MyService() *myservice.Client
   }
   ```

3. **Use in handlers**
   ```go
   import facades "github.com/desygner/desygner-toolchain/packages/aws/facades"

   awsCli := facades.Aws()
   myServiceClient := awsCli.MyService()
   ```

### Posting to Slack from Handlers

Tools can send Slack messages through the `extension.Input` type:

```go
input := &handlers.MyInput{
    Source: "...",
    Input: extension.Input{
        SlackChannel: "C123...",
        SlackTS: "1234567890.123456",
    },
}

// In handler
input.SendMessage(ctx, ":spin: Processing...")
input.SendMessage(ctx, ":greentick1: Done!")
```

---

## Debugging

### Enable Debug Logging

1. Set `APP_DEBUG=true` in `.env`
2. Set log level in `config/logging.go`
3. Logs appear in `storage/logs/`

### Common Issues

**Issue**: "distribution not found"
- **Cause**: CloudFront domain not in aliases
- **Fix**: Add domain to CloudFront distribution aliases
- **Debug**: Check `cache_key = cdn.{profile}.distribution.{domain}`

**Issue**: "failed to download slack file"
- **Cause**: Invalid file URL or Slack token permissions
- **Fix**: Verify `files:read` scope in Slack bot token
- **Debug**: Log the file URL and token (sanitized)

**Issue**: "S3 access denied"
- **Cause**: AWS credentials don't have S3:PutObject permission
- **Fix**: Update IAM policy:
   ```json
   {
       "Effect": "Allow",
       "Action": ["s3:GetObject", "s3:PutObject"],
       "Resource": "arn:aws:s3:::my-bucket/*"
   }
   ```

---

## Testing

### Unit Tests
```bash
go test ./app/handlers -v
```

### Integration Tests
Requires AWS credentials + Slack token:
```bash
go test ./tests/feature -v
```

### Test Example
```go
// tests/feature/file_replace_test.go
func TestFileReplace(t *testing.T) {
    handler := &handlers.Replace{}
    input := &handlers.ReplaceInput{
        Source: "slack-url",
        Target: "s3://bucket/key",
        Input: extension.Input{
            SlackChannel: "C123",
            SlackTS: "123.456",
        },
    }

    err := handler.Handle(context.Background(), input)
    assert.NoError(t, err)
}
```

---

## Performance Tuning

### Cache Configuration
Update `packages/aws/facades/aws.go`:
```go
// Cache TTL for distributions (24 hours default)
// Reduce for more fresh data, increase for fewer API calls
ttl := 24 * time.Hour
```

### Concurrency
File replacement uses `syncs.ErrGroup` for parallel I/O:
```go
// Both downloads happen concurrently
// Both uploads happen concurrently
// 2 Slack API calls + 2 S3 API calls in parallel
```

### Database Migrations
If adding database features:
```bash
go run . artisan migrate:refresh --seed
```

---

## N8N Integration

The MCP server is designed to be called from N8N via the MCP node:

### N8N MCP Node Configuration
```
MCP Server: /path/to/toolchain (or Docker image)
Server Args: go run .
Transport: stdio
```

### Example N8N Workflow
```
1. Slack Webhook (receives message)
   ↓
2. Claude with MCP (interprets intent + calls tool)
   ├─ cloudfront-invalidate
   ├─ file-replace
   └─ bucket-resolver
   ↓
3. Response Handler (formats result)
   ↓
4. Slack Message (posts result)
```

### Environment in N8N
Set as environment variables in N8N:
```
AWS_PROFILE=default
AWS_REGION=us-east-1
SLACK_BOT_TOKEN=xoxb-...
MCP_PORT=8080
```

---

## Code Organization

### Handler Responsibilities
- **Input Validation**: Check parameters
- **CloudFront/S3 Operations**: Call AWS SDK
- **Slack Notifications**: Progress updates
- **Error Handling**: User-friendly messages
- **Resource Cleanup**: Temp files, connections

### Controller Responsibilities
- **MCP Protocol**: Handle CallToolRequest/Response
- **Parameter Parsing**: Extract arguments
- **Progress Notifications**: Send to Slack via MCP
- **Logging**: Structured logs with context

### Package Responsibilities
- **AWS Package**: CloudFront + S3 client management
- **Slack Package**: Slack API wrapper + message helpers
- **MCP Package**: MCP server setup + service binding

---

## Monitoring

### Logs Location
```
storage/logs/
├── stack-{date}.log      # Application logs
├── error-{date}.log      # Error logs
└── single.log           # Combined logs
```

### Slack Integration Logging
All operations log to Slack thread:
- `:spin:` = processing
- `:greentick1:` = success
- `:x:` = failure
- Includes AWS console links for tracking

### CloudWatch Integration
Cloudwatch client available:
```go
cw := facades.Aws().Cloudwatch()
// Can emit custom metrics for monitoring
```

---

## Helpful Commands

```bash
# View recent logs
tail -f storage/logs/stack-*.log

# Check AWS credentials
aws sts get-caller-identity --profile default

# Verify Slack token
curl -H "Authorization: Bearer $SLACK_BOT_TOKEN" https://slack.com/api/auth.test

# Load test MCP server
# (Use MCP client SDK or N8N test node)

# Format code
go fmt ./...

# Lint
golangci-lint run ./...

# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o toolchain .
```

---

## Contributing

1. Create feature branch: `git checkout -b feature/my-feature`
2. Make changes following existing patterns
3. Add tests: `app/handlers/my_handler.go` → `tests/feature/my_handler_test.go`
4. Run tests: `go test ./...`
5. Format: `go fmt ./...`
6. Commit with meaningful message
7. Push and create PR

---

## Troubleshooting Checklist

- [ ] AWS credentials working? (`aws sts get-caller-identity`)
- [ ] Slack token valid? (test in Slack API tester)
- [ ] MCP server running? (check port 8080)
- [ ] CloudFront domain in aliases? (check AWS console)
- [ ] S3 bucket permissions? (check IAM policy)
- [ ] Logs showing errors? (check `storage/logs/`)
- [ ] Rate limits hit? (check CloudFront API metrics)
- [ ] File too large? (check Slack file size limits)
