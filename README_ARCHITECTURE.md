# BishOps Toolchain - Complete Documentation Index

This repository contains a production-ready **MCP (Model Context Protocol) Server** that bridges AI systems with AWS infrastructure operations through Slack.

## 📚 Documentation Overview

This guide helps you understand, develop, and operate the BishOps Toolchain.

### 1. **[ARCHITECTURE.md](./ARCHITECTURE.md)** - System Design & Components
**Start here if you want to understand:**
- What this system is and why it exists
- High-level architecture and data flow
- All major components (Goravel, MCP, AWS, Slack packages)
- How the three MCP tools work
- Service provider pattern and dependency injection
- Deployment options (Docker, Kubernetes)

**Key Sections:**
- System architecture diagram
- Component breakdown
- Data flow examples
- Configuration reference
- Security considerations

**Best for:** New developers, architects, code reviewers

---

### 2. **[DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md)** - Setup & Development
**Start here if you want to:**
- Get the project running locally
- Understand common development tasks
- Add new MCP tools
- Integrate new AWS services
- Post messages to Slack from handlers
- Debug issues

**Key Sections:**
- Quick start guide (5 minutes)
- Development commands
- Step-by-step: Add a new MCP tool
- Step-by-step: Add a new AWS service
- Testing & debugging
- N8N integration guide
- Troubleshooting checklist

**Best for:** Developers, DevOps engineers, integrators

---

### 3. **[MCP_TOOLS_REFERENCE.md](./MCP_TOOLS_REFERENCE.md)** - Tool Specifications
**Start here if you want to:**
- Understand what each MCP tool does
- See detailed parameters and responses
- Know how tools work under the hood
- Learn about edge cases and performance
- Use tools from N8N/Claude

**Key Sections:**
- Detailed reference for each of 3 tools:
  - `cloudfront-invalidate` - Cache invalidation
  - `file-replace` - S3 file updates with backups
  - `bucket-resolver` - S3 bucket discovery
- Parameter specifications
- Response formats
- Technical implementation details
- Performance characteristics
- Rate limits and quotas

**Best for:** AI integration engineers, N8N workflow builders, API consumers

---

## 🚀 Quick Links

### For Different Roles

**I'm a Product Manager/Architect**
→ Read [ARCHITECTURE.md](./ARCHITECTURE.md) - System Design section

**I'm a Backend Developer**
→ Start with [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - Quick Start

**I'm integrating with N8N/Claude**
→ Start with [MCP_TOOLS_REFERENCE.md](./MCP_TOOLS_REFERENCE.md) - Tool 1

**I'm debugging an issue**
→ Go to [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - Troubleshooting Checklist

**I need to add a new feature**
→ [ARCHITECTURE.md](./ARCHITECTURE.md) - Components, then [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - Add new tool section

---

## 🎯 What This System Does

In a nutshell:

```
User in Slack writes: "Invalidate /assets/image.png on CDN"
    ↓
N8N runs Claude with MCP integration
    ↓
Claude interprets intent and calls "cloudfront-invalidate" tool
    ↓
MCP Server:
  1. Resolves CloudFront distribution from URL
  2. Submits invalidation request
  3. Waits for completion (5 min max)
  4. Posts result with AWS console link to Slack thread
    ↓
User sees: "✓ Cache invalidated" in original Slack thread
```

---

## 📦 Key Features

✅ **Three Core MCP Tools:**
- Invalidate CloudFront cache
- Replace S3 files (with automatic backups)
- Resolve S3 buckets from CloudFront URLs

✅ **Slack Integration:**
- All operations post progress to Slack threads
- Automatic backup files attached to Slack
- Error messages sent to user's conversation
- Full audit trail in Slack

✅ **AWS Support:**
- CloudFront API for distribution management
- S3 for file operations
- Caching for performance (24-hour TTL)
- Multi-account support (via AWS_PROFILE)

✅ **Developer Friendly:**
- Goravel framework (Laravel-like patterns)
- Service provider pattern for extensibility
- Easy to add new tools or AWS services
- Comprehensive logging

✅ **Production Ready:**
- Docker & Kubernetes support
- Error handling & cleanup
- Structured logging
- Configurable via environment variables

---

## 🏗️ Project Structure

```
desygner-toolchain/
├── ARCHITECTURE.md                 ← System design (read first!)
├── DEVELOPER_GUIDE.md             ← Development setup & tasks
├── MCP_TOOLS_REFERENCE.md         ← Detailed tool specifications
├── README_ARCHITECTURE.md          ← This file (navigation guide)
│
├── app/
│   ├── handlers/                  ← Business logic (CDN, file replace)
│   ├── mcp/                       ← MCP tool registration & controllers
│   ├── http/                      ← HTTP endpoints (health, webhooks)
│   ├── models/                    ← Data models
│   └── providers/                 ← Service providers
│
├── packages/
│   ├── aws/                       ← AWS SDK wrapper
│   ├── mcp/                       ← MCP server setup
│   └── slack/                     ← Slack API wrapper
│
├── config/                        ← Application configuration
├── bootstrap/                     ← App initialization
├── infra/                         ← Kubernetes Helm chart
├── tests/                         ← Test files
├── main.go                        ← Entry point
├── go.mod                         ← Go dependencies
├── Dockerfile                     ← Container image
└── docker-compose.yml            ← Local dev setup
```

---

## 🔧 Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Framework** | Goravel (Go) | Application foundation, routing, DI |
| **Protocol** | MCP (Model Context Protocol) | Standardized tool-calling interface |
| **Cloud** | AWS SDK v2 (Go) | CloudFront, S3 operations |
| **Chat** | Slack API | Notifications and file handling |
| **HTTP** | Gin | Web framework |
| **Logging** | Zap/Goravel Log | Structured logging |
| **Caching** | Redis (via Goravel) | Distribution/bucket caching |
| **Container** | Docker | Deployment |
| **Orchestration** | Kubernetes (Helm) | Production deployment |

---

## 📋 Common Tasks

### Setup & Running
```bash
# 1. Setup
cp .env.example .env
# Edit .env with your AWS profile and Slack token

# 2. Run locally
go run .

# 3. Run in Docker
docker-compose up
```
→ See [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md)

### Development
```bash
# Add a new MCP tool
# → See [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) section "Adding a New MCP Tool"

# Add AWS service support
# → See [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) section "Adding AWS Service Access"

# Test tools
go test ./app/handlers -v
```

### Deployment
```bash
# Build Docker image
docker build -t bishops-toolchain .

# Deploy to Kubernetes
helm install toolchain ./infra/toolchain -f values.yaml
```

### Debugging
```bash
# View logs
tail -f storage/logs/stack-*.log

# Enable debug mode
APP_DEBUG=true go run .

# Test in isolation
# → See [MCP_TOOLS_REFERENCE.md](./MCP_TOOLS_REFERENCE.md) section "Debugging Tools"
```

---

## 🤝 Integration Points

### With N8N
The server is designed as an MCP server that N8N can call:
```
N8N MCP Node → starts: go run .
             → feeds: JSON-RPC over stdio
             → calls: cloudfront-invalidate, file-replace, bucket-resolver
```

### With Claude/AI
Tools are registered with MCP, so any AI system supporting MCP can call them:
```
Claude + N8N MCP Integration → reads tool definitions
                             → calls based on user intent
                             → processes results
```

### With Slack
All operations post progress to Slack threads:
```
Handler sends message → Slack API
                     → appears in user's thread
                     → with AWS console links
```

### With AWS
Operations use AWS SDK v2 with cached distribution queries:
```
Handler queries CloudFront → cached for 24 hours
                          → reduced API quota consumption
         queries S3 → standard rate limits apply
```

---

## 🔐 Security Model

**Authentication:**
- AWS: Via local credentials (`~/.aws/credentials` or assume role)
- Slack: Via bot token in env var

**Authorization:**
- AWS: IAM policy grants specific S3 and CloudFront permissions
- Slack: Bot token scopes limit what bot can do
- Internal: Slack channel/timestamp verify user context

**Data Protection:**
- No sensitive data in logs (tokens, paths truncated)
- Temporary files cleaned up after operations
- HTTPS for external API calls
- Public-read ACL for S3 files (CDN-appropriate)

See [ARCHITECTURE.md](./ARCHITECTURE.md) - Security Considerations

---

## 📊 Performance

| Operation | Typical Time | Max Time | Notes |
|-----------|--------------|----------|-------|
| cloudfront-invalidate | 10-30s | 5m | Waits for CloudFront completion |
| file-replace | 5-30s | N/A | Depends on file size |
| bucket-resolver | <1s | 2s | Usually cached |

**Optimization:**
- Distribution list cached 24 hours
- Bucket resolution cached 24 hours
- Parallel I/O for file uploads/downloads
- CloudFront waiter with exponential backoff

---

## 🐛 Troubleshooting Quick Links

**Distribution not found**
→ [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - Common Issues

**S3 access denied**
→ [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - Common Issues

**Slack file download failed**
→ [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - Common Issues

**Full troubleshooting checklist**
→ [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - Troubleshooting Checklist

---

## 📖 Learning Path

**For completely new developers:**
1. Start: [README.md](./README.md) (basic overview)
2. Understand: [ARCHITECTURE.md](./ARCHITECTURE.md) (system design)
3. Setup: [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) (get running)
4. Reference: [MCP_TOOLS_REFERENCE.md](./MCP_TOOLS_REFERENCE.md) (tools in detail)

**For integrators (N8N/Claude):**
1. Quick overview: This file (README_ARCHITECTURE.md)
2. Tool reference: [MCP_TOOLS_REFERENCE.md](./MCP_TOOLS_REFERENCE.md)
3. Integration guide: [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - N8N Integration section

**For operations/DevOps:**
1. Architecture: [ARCHITECTURE.md](./ARCHITECTURE.md) - Deployment section
2. Setup: [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - Quick Start
3. Monitoring: [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - Monitoring section

---

## 🤔 FAQ

**Q: How do I add a new MCP tool?**
A: See [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - "Adding a New MCP Tool"

**Q: Can this work without Slack?**
A: Yes, Slack integration is optional. Errors/results just won't post to Slack threads.

**Q: Does this work with other CloudFront distributions?**
A: Yes! Any distribution whose domain alias is used in the URL will be found.

**Q: Can I add more S3 buckets?**
A: Yes, as long as your AWS credentials have access to those buckets.

**Q: What happens if the MCP server crashes?**
A: N8N will retry according to its settings. Temp files are cleaned up automatically.

**Q: Is the tool thread-safe?**
A: Yes, all external operations use proper context handling and cleanup.

**Q: Can I extend this for other AWS services?**
A: Yes! See [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - "Adding AWS Service Access"

---

## 🎓 Key Concepts

**MCP (Model Context Protocol):**
A standard for defining tools that AI systems can call. Defined by Anthropic.

**Service Providers:**
Goravel pattern for registering services in a dependency injection container.

**Facades:**
Goravel pattern for accessing services statically (e.g., `facades.Aws()`)

**Handler:**
Business logic layer that performs actual AWS/Slack operations

**MCP Controller:**
Adapts Handler to MCP protocol (handles CallToolRequest/Response)

---

## 📞 Getting Help

1. Check the relevant documentation above
2. Search [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) troubleshooting section
3. Review logs in `storage/logs/`
4. Check AWS CloudWatch for API errors
5. Test Slack token in [Slack API tester](https://api.slack.com/tools/explorer)

---

## 📝 Related Resources

- **Model Context Protocol**: https://modelcontextprotocol.io/
- **Goravel Documentation**: https://www.goravel.dev/
- **AWS SDK for Go**: https://aws.github.io/aws-sdk-go-v2/
- **Slack API**: https://api.slack.com/

---

## 📄 Document Information

- **Last Updated**: 2024
- **Version**: 1.0
- **Scope**: Complete system documentation
- **Audience**: Developers, DevOps, Integrators

---

**Navigation:**
- 📖 [Full Architecture Details](./ARCHITECTURE.md)
- 🛠️ [Development Guide](./DEVELOPER_GUIDE.md)
- 📚 [Tool Specifications](./MCP_TOOLS_REFERENCE.md)
