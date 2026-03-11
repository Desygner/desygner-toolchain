# BishOps Toolchain Architecture

## Overview

**BishOps Toolchain** is an MCP (Model Context Protocol) server that exposes AWS CloudFront and S3 operations through a standardized interface. It's designed to be integrated with N8N automation workflows and Slack bots, enabling AI-assisted CDN cache invalidation and file replacement operations directly from Slack conversations.

**Core Purpose**: Bridge cloud infrastructure operations (AWS CloudFront/S3) with conversational AI platforms (Claude via N8N) through Slack, allowing users to manage CDN cache and file updates through natural language commands in Slack.

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        SLACK BOT (BishOps)                      │
│              (N8N Workflow with Claude Integration)              │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                    MCP Protocol (JSON-RPC over stdio)
                                  │
┌─────────────────────────────────▼───────────────────────────────┐
│                   DESYGNER TOOLCHAIN (Go App)                   │
│                    (MCP Server Implementation)                   │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────┐   │
│  │           MCP Server & Tool Registration                 │   │
│  │  ├─ cloudfront-invalidate                                │   │
│  │  ├─ file-replace                                         │   │
│  │  └─ bucket-resolver                                      │   │
│  └──────────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              Business Logic (Handlers)                   │   │
│  │  ├─ Invalidate: CloudFront cache invalidation            │   │
│  │  ├─ Replace: S3 file replacement with backup             │   │
│  │  └─ CDN: Bucket/Distribution resolution                 │   │
│  └──────────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────┐   │
│  │           External Service Integrations                  │   │
│  │  ├─ AWS SDK (CloudFront, S3, STS, CloudWatch)            │   │
│  │  └─ Slack API (File upload, Thread replies)              │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
        ▼                   ▼                   ▼
    ┌────────┐        ┌────────┐        ┌────────┐
    │   AWS  │        │ Cache  │        │  Logs  │
    │CloudFront│      │(Redis) │        │(CloudWatch)│
    │  S3    │        └────────┘        └────────┘
    └────────┘
```

---

## Key Components

### 1. **Goravel Framework (Foundation)**
- **File**: `bootstrap/app.go`, `main.go`
- **Purpose**: Provides the application foundation and service provider system
- **Responsibilities**:
  - HTTP server setup (Gin-based)
  - Service container and dependency injection
  - Configuration management
  - Logging infrastructure

### 2. **MCP Package** (`packages/mcp/`)
- **Files**: `service_provider.go`, `service/mcp.go`, `contracts/mcp.go`
- **Purpose**: Wraps the MCP SDK and provides IoC container bindings
- **Key Responsibilities**:
  - Initialize MCP Server with app metadata
  - Register tools with the MCP server
  - Expose MCP interface through service container

**Service Provider Registration**:
```go
// Creates MCP Server with Implementation metadata
srv := mcp.NewServer(&mcp.Implementation{
    Name:    config.GetString("app.name"),      // "BishOps Tools"
    Version: config.GetString("app.version"),
}, nil)
```

### 3. **MCP Kernel** (`app/mcp/kernel.go`)
- **Purpose**: Central registration of MCP tools
- **Exposes 3 Tools**:

  #### Tool 1: `cloudfront-invalidate`
  - **Purpose**: Invalidate CloudFront cache for specific paths
  - **Input**: URL to invalidate + optional Slack context
  - **Process**:
    1. Resolve CloudFront distribution from URL hostname
    2. Extract path from URL
    3. Create invalidation request
    4. Wait for completion (up to 5 minutes)
    5. Post result to Slack thread
  - **Output**: Invalidation ID + AWS Console link

  #### Tool 2: `file-replace`
  - **Purpose**: Replace S3 objects with files from Slack, keeping backups
  - **Input**: Slack file URL + Target S3/CloudFront URL
  - **Process**:
    1. Resolve S3 bucket from target URL
    2. Download original S3 object as backup (parallel)
    3. Download new file from Slack (parallel)
    4. Upload backup to Slack thread
    5. Upload new file to S3
    6. Post completion to Slack
  - **Output**: Success message + backup link in Slack

  #### Tool 3: `bucket-resolver`
  - **Purpose**: Resolve the source S3 bucket for a CloudFront URL
  - **Input**: CloudFront or S3 URL
  - **Process**:
    1. Query CloudFront distribution config
    2. Match path patterns against cache behaviors
    3. Extract origin information
    4. Parse S3 bucket name
  - **Output**: JSON with bucket name and distribution ID

### 4. **Handlers** (`app/handlers/`)

#### `cdn.go` (CDN Base Class)
- **Methods**:
  - `ResolveBucket(ctx, url)`: Query CloudFront + extract bucket
  - `ResolveDistribution(ctx, url)`: Find distribution ID from domain
  - `getDistribution()`: Cached lookup with 24-hour TTL
  - `loadDistributions()`: Paginated CloudFront API call
  - `getBucket()`: Parse distribution config, match path patterns

#### `invalidate.go` (Invalidate Handler)
- **Entry Point**: `Handle(ctx, input)`
- **Workflow**:
  1. Creates CloudFront invalidation batch
  2. Notifies Slack with progress and AWS console link
  3. Waits for invalidation completion via CloudFront waiter
  4. Posts final success message to Slack thread

#### `replace.go` (File Replace Handler)
- **Entry Point**: `Handle(ctx, input)`
- **Workflow**:
  1. Creates temporary files for upload/backup
  2. Resolves bucket (tries direct S3 URL parsing first, falls back to CloudFront query)
  3. Runs parallel downloads:
     - Slack file → temp file
     - S3 object → temp backup
  4. Runs parallel uploads:
     - Backup → Slack thread as file
     - New file → S3 with public-read ACL
  5. Error handling: notifies Slack on failure
  6. Resource cleanup: removes temp files
- **MIME Type Detection**: Uses file inspection to set correct S3 content-type

### 5. **AWS Package** (`packages/aws/`)
- **Files**: `aws.go`, `client/aws.go`, `config/aws.go`, `facades/aws.go`, `service_provider.go`
- **Purpose**: AWS SDK wrapper and service provider
- **Exposes**:
  - CloudFront client
  - S3 client
  - Current AWS profile name
- **Configuration**:
  - Reads from environment: `AWS_PROFILE`, region
  - Caches distribution list (configurable TTL)
  - Supports multi-account operations

### 6. **Slack Package** (`packages/slack/`)
- **Files**: `slack.go`, `config/slack.go`, `facades/slack.go`, `service_provider.go`, `extension/input.go`
- **Purpose**: Slack API integration
- **Key Capabilities**:
  - Post messages to Slack threads
  - Download files from Slack (private URLs)
  - Upload files to Slack threads
- **Extension Input Type**:
  - Carries Slack channel + thread timestamp
  - Provides `SendMessage(ctx, text)` helper

---

## Data Flow: File Replacement Example

```
User in Slack:
"Replace image.png from this Slack file to https://cdn.example.com/assets/image.png"
        │
        ▼
N8N Workflow (with Claude):
1. Claude interprets user intent
2. Calls MCP tool: file-replace
   {
     "source": "https://files.slack.com/files-pri/...",
     "target": "https://cdn.example.com/assets/image.png",
     "slack": { "channel": "C123...", "ts": "1234567890.123456" }
   }
        │
        ▼
MCP Server receives CallToolRequest
        │
        ▼
FileReplace.Handle():
1. Parse target URL → extract bucket name
   Could be: "s3://bucket/key" or "cloudfront-url" or "s3.amazonaws.com/bucket/key"
        │
        ▼
2. Parallel Operations:
   ├─ S3GetObject(bucket, key) → temp backup
   │  (with ContentLength tracking)
   └─ SlackDownloadFile(source) → temp file
      (with MIME type detection)
        │
        ▼
3. Parallel Uploads:
   ├─ S3PutObject(bucket, key, new_file)
   │  ACL: public-read
   │  ContentType: detected from file
   └─ SlackUploadFile(backup_file)
      ThreadTimestamp: parent message
      InitialComment: "Backup copy of..."
        │
        ▼
4. Send Success Message to Slack Thread
        │
        ▼
MCP Server returns CallToolResult
        │
        ▼
N8N → Slack: "✅ File replacement complete"
```

---

## Configuration

### Environment Variables (`.env`)

```bash
# Application
APP_NAME=BishOps Tools
APP_ENV=production
APP_DEBUG=false
APP_KEY=<32-char encryption key>

# MCP Server
MCP_PORT=8080

# AWS
AWS_PROFILE=default
AWS_REGION=us-east-1
AWS_TTL=86400  # Cache TTL in seconds (24 hours)

# Slack
SLACK_BOT_TOKEN=xoxb-...
```

### Service Providers (loaded in `config/app.go`)
1. **Framework Providers** (Goravel): Log, Console, Database, Cache, HTTP, Route, etc.
2. **Custom Providers**: AppService, Auth, Route, GrpcService, Console, Queue, Event, Validation, Database
3. **Package Providers**: MCP, AWS, Slack
4. **HTTP**: Gin service provider

---

## Why This Architecture?

### 1. **MCP (Model Context Protocol) for Standardization**
- **Problem**: Different AI platforms (Claude, ChatGPT, etc.) need different interfaces
- **Solution**: MCP provides a standard tool-calling interface
- **Benefit**: N8N can swap Claude with another AI without changing tool definitions

### 2. **Slack Thread Context Passing**
- **Why**: Operations should be tied to the conversation
- **How**: Channel ID + timestamp embedded in tool parameters
- **Benefit**: Progress updates and results appear in the original thread, providing audit trail

### 3. **Parallel I/O Operations**
- **Why**: File replacement involves 4 I/O operations (download 2 files, upload 2 files)
- **How**: `syncs.ErrGroup` for concurrent operations
- **Benefit**: Faster operation completion, fewer timeouts

### 4. **Caching Distributions & Buckets**
- **Why**: CloudFront API calls are expensive and distributions rarely change
- **How**: 24-hour Redis cache with pattern-based key naming
- **Benefit**: Faster resolution, reduced AWS API quota consumption

### 5. **Service Provider Pattern (Goravel)**
- **Why**: Multiple packages need access to AWS, Slack, MCP
- **How**: IoC container with `facades.Aws()`, `facades.Slack()`, etc.
- **Benefit**: Loose coupling, testable, centralized configuration

### 6. **Separation of Concerns**
- **MCP Controllers**: Handle MCP protocol, parameter validation
- **Handlers**: Business logic (CDN operations)
- **Packages**: External integrations (AWS, Slack)
- **Benefit**: Easy to test, modify, and extend

---

## Deployment

### Docker
- **File**: `Dockerfile`
- Builds Go binary, runs on port 8080
- Expects environment variables at runtime

### Kubernetes (Helm)
- **Chart**: `infra/toolchain/`
- Includes: Deployment, Service, Ingress, HPA, Secrets
- Can be integrated with N8N K8s cluster

### N8N Integration
The MCP server is launched as a subprocess by N8N's MCP integration:
```
N8N MCP Node → starts process: go run . serve
             → connects via stdio
             → sends/receives JSON-RPC messages
```

---

## Security Considerations

1. **AWS Credentials**: Passed via `AWS_PROFILE` → uses local AWS config/credentials
2. **Slack Token**: Stored as `SLACK_BOT_TOKEN` environment variable
3. **File ACL**: Files uploaded to S3 are marked `public-read` (CDN-appropriate)
4. **Thread Context**: Slack channel/timestamp tied to operations for audit trail
5. **Error Messages**: Don't expose sensitive information (S3 paths kept private)

---

## Future Enhancements

1. **Additional CDN Providers**: CloudflareWorkers, AWS API Gateway invalidation
2. **Batch Operations**: Invalidate multiple paths in single request
3. **Cost Tracking**: Log operation costs to CloudWatch for chargeback
4. **Approval Workflows**: Require human approval for sensitive operations
5. **Audit Database**: Store all operations in persistent database (migrations setup exists)
6. **Document Search**: Implement `doc.Search()` handler for knowledge base integration
7. **Custom Validation**: Request validation service provider for input sanitization

---

## Testing

- **Framework**: `testify` (assertions)
- **Test Location**: `tests/` directory
- **Example**: `tests/feature/example_test.go`

To run tests:
```bash
go test ./...
```

---

## Project Structure Summary

```
.
├── app/
│   ├── console/          # CLI commands (artisan)
│   ├── grpc/             # gRPC implementations
│   ├── handlers/         # Business logic
│   ├── http/             # HTTP controllers & middleware
│   ├── mcp/              # MCP tool registration
│   ├── models/           # Data models (CDNMap, User, Doc)
│   └── providers/        # Service providers
├── bootstrap/            # App initialization
├── config/               # Configuration files
├── database/             # Migrations & seeders
├── packages/             # Reusable packages
│   ├── aws/              # AWS integration
│   ├── mcp/              # MCP wrapper
│   └── slack/            # Slack integration
├── infra/                # Infrastructure
│   └── toolchain/        # Helm chart
├── routes/               # Route definitions (HTTP, gRPC, API)
├── storage/              # Logs & temp files
├── tests/                # Test files
├── main.go               # Entry point
├── go.mod/.sum           # Dependency management
└── docker-compose.yml    # Local development setup
```
