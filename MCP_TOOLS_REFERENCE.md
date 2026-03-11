# MCP Tools Reference

This document describes all available MCP tools exposed by the BishOps Toolchain server.

---

## Tool 1: `cloudfront-invalidate`

**Purpose**: Invalidate CloudFront cache for a specific URL path

**Handler**: `app/handlers/invalidate.go` → `Invalidate.Handle()`

### Request Parameters

```json
{
  "url": "https://cdn.example.com/assets/image.png",
  "slack": {
    "channel": "C123ABC456DEF",
    "ts": "1234567890.123456"
  }
}
```

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `url` | string | Yes | Full URL to invalidate. Can be any CloudFront domain alias. Path is extracted for invalidation. |
| `slack.channel` | string | No | Slack channel ID where result should be posted |
| `slack.ts` | string | No | Slack thread timestamp for threaded replies |

### How It Works

1. **Distribution Resolution**
   - Extracts hostname from URL
   - Queries CloudFront API for distributions matching the domain
   - Uses 24-hour cache to reduce API calls
   - Fails if domain not found in any distribution's aliases

2. **Path Extraction**
   - Parses URL path (e.g., `/assets/image.png`)
   - URL-decodes for special characters
   - Passes exact path to CloudFront invalidation

3. **Invalidation Request**
   - Creates `InvalidationBatch` with single path
   - Sets `CallerReference` to unique ID: `{app-name}/{env}/{timestamp}`
   - Submits to CloudFront API

4. **Progress Tracking**
   - Notifies Slack immediately after submission with AWS console link
   - Waits up to 5 minutes for invalidation to complete
   - Sends final confirmation to Slack thread

5. **Error Handling**
   - Fails if distribution not found
   - Fails if invalidation request rejected
   - Fails if wait timeout exceeded
   - Sends error message to Slack thread

### Response

**Success Response**:
```json
{
  "type": "text",
  "text": "invalidation done. EABCDEF123456G7H"
}
```

**Slack Notifications**:
1. Initial: `:spin: invalidation submitted https://console.aws.amazon.com/cloudfront/v4/home?region=us-east-1#/distributions/E123ABC/invalidations/details/EABCDEF123456G7H`
2. Final: `:greentick1: invalidation completed https://console.aws.amazon.com/cloudfront/v4/home?region=us-east-1#/distributions/E123ABC/invalidations/details/EABCDEF123456G7H`

### Example Usage

```python
# From N8N or Claude context
await mcp_call("cloudfront-invalidate", {
    "url": "https://assets.example.com/images/logo.png",
    "slack": {
        "channel": "C00000000001",
        "ts": "1609459200.000100"
    }
})
```

### Edge Cases

- **URL with query params**: Query params are ignored, only path is invalidated
  - Input: `https://cdn.example.com/image.png?v=2`
  - Invalidated path: `/image.png`

- **Root path**: Invalidating domain root
  - Input: `https://cdn.example.com/`
  - Invalidated path: `/`

- **Wildcard patterns**: If you need to invalidate multiple paths, call tool multiple times
  - Single call cannot do: `/*` (would need separate tool)

---

## Tool 2: `file-replace`

**Purpose**: Replace an S3 object with a file from Slack, keeping automatic backup

**Handler**: `app/handlers/replace.go` → `Replace.Handle()`

### Request Parameters

```json
{
  "source": "https://files.slack.com/files-pri/T123456/F123456/image.png?t=xoxp-token",
  "target": "https://cdn.example.com/images/logo.png",
  "slack": {
    "channel": "C123ABC456DEF",
    "ts": "1234567890.123456"
  }
}
```

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `source` | string | Yes | Slack private file URL (from Slack file picker) |
| `target` | string | Yes | Target location: S3 URL, CloudFront URL, or virtual-hosted S3 URL |
| `slack.channel` | string | No | Slack channel ID where Slack thread will receive backup file |
| `slack.ts` | string | No | Slack thread timestamp for threaded backup upload |

### Supported Target URL Formats

All of these are valid:
```
s3://bucket/path/to/file.png
s3.amazonaws.com/bucket/path/to/file.png
bucket.s3.amazonaws.com/path/to/file.png
bucket.s3.us-east-1.amazonaws.com/path/to/file.png
https://cdn.example.com/path/to/file.png (CloudFront)
```

### How It Works

1. **Bucket Resolution**
   - Attempts direct parsing from S3 URL format
   - If CloudFront URL or unknown format, queries CloudFront → S3 bucket lookup
   - Caches bucket resolution for 24 hours

2. **Parallel Downloads** (concurrent)
   - **Slack File**: Downloaded from private URL to temp file
   - **S3 Current**: Downloaded from S3 to temp backup file
   - Both happen simultaneously
   - MIME type detected from Slack file

3. **Upload Phase** (concurrent)
   - **Backup to Slack**: Original S3 file uploaded to Slack thread with timestamp
   - **New to S3**: Slack file uploaded to S3 with:
     - Detected MIME type as `Content-Type`
     - `ACL: public-read` (for CDN access)
     - Original S3 metadata preserved where possible

4. **Progress Tracking**
   - Posts `:spin:` message at key milestones
   - Includes bucket location info
   - Final message: `:greentick1: File has been replaced`
   - Backup file attached to Slack thread for recovery

5. **Error Handling**
   - Fails if Slack file cannot be downloaded
   - Fails if S3 object cannot be accessed
   - Fails if S3 upload is rejected
   - Cleans up temp files on error
   - Sends error details to Slack thread

### Response

**Success Response**:
```json
{
  "type": "text",
  "text": "replacement done"
}
```

**Slack Notifications**:
1. Initial: `:spin: Looking up the location for https://cdn.example.com/images/logo.png`
2. Found: `:spin: Location found at my-bucket. Download the original item as backup`
3. Downloaded: `:spin: Download complete. Replace the original item with the new one`
4. Final: `:greentick1: File has been replaced`
5. (Backup file in thread): "Backup copy of: https://cdn.example.com/images/logo.png"

### Example Usage

```python
# From N8N or Claude context
# User uploads file in Slack, gets a private URL
await mcp_call("file-replace", {
    "source": "https://files.slack.com/files-pri/T123/F456/image.png?t=xoxp",
    "target": "https://assets.example.com/logo.png",
    "slack": {
        "channel": "C789",
        "ts": "1609459200.000100"
    }
})
```

### Technical Details

**Temp Files**:
- Created in system temp directory (`/tmp` on Unix)
- Names: `{filename}-backup`, `{filename}-replace`
- Automatically cleaned up after operation

**MIME Type Detection**:
```
image/png, image/jpeg, image/svg+xml
application/pdf, text/css, application/javascript
... (covers 200+ types via Go's mime package)
```

**Concurrency Pattern**:
```go
errsync := syncs.NewErrGroup()
errsync.Go(func() error { return downloadSlackFile(...) })
errsync.Go(func() error { return downloadS3Object(...) })
err := errsync.Wait()  // Waits for both
```

### Edge Cases

- **File already exists at target**: Gets replaced, original backed up
- **Different file size**: Handled correctly, S3 ContentLength updated
- **File with same name already in Slack thread**: New backup uploaded with timestamp
- **Network interruption**: Temp files cleaned up, operation fails gracefully
- **Large files**: Streamed to avoid memory overload
  - Slack max: 1 GB
  - S3 max: 5 TB (for multi-part upload)

---

## Tool 3: `bucket-resolver`

**Purpose**: Resolve the source S3 bucket for a given CloudFront URL

**Handler**: `app/handlers/cdn.go` → `CDN.ResolveBucket()`

### Request Parameters

```json
{
  "url": "https://cdn.example.com/assets/styles/theme.css",
  "slack": {
    "channel": "C123ABC456DEF",
    "ts": "1234567890.123456"
  }
}
```

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `url` | string | Yes | CloudFront URL or S3 URL to resolve |
| `slack.channel` | string | No | Currently unused, but available for future notifications |
| `slack.ts` | string | No | Currently unused, but available for future notifications |

### How It Works

1. **URL Parsing**
   - Extracts hostname and path from URL

2. **Distribution Lookup**
   - Queries CloudFront API for distribution matching hostname
   - Uses 24-hour cache
   - Fails if not found in any distribution's aliases

3. **Path Pattern Matching**
   - Iterates through cache behaviors in distribution config
   - Converts CloudFront path patterns (`*` wildcards) to regex
   - Matches URL path against each pattern
   - Stops at first match

4. **Origin Resolution**
   - Extracts `TargetOriginId` from matching cache behavior
   - Finds corresponding origin in distribution config
   - Extracts origin domain name and path

5. **Bucket Extraction**
   - Uses regex to parse S3 bucket from origin domain:
     - Pattern: `^(?P<bucket>.*)\.s3\.([^.]+)\.amazonaws\.com$`
   - Returns bucket name
   - If pattern doesn't match, returns domain name as fallback

### Response

**Success Response**:
```json
{
  "type": "text",
  "text": "{\n  \"DistributionID\": \"E123ABC456DEF\",\n  \"Bucket\": \"my-website-assets\"\n}"
}
```

### Example Usage

```python
await mcp_call("bucket-resolver", {
    "url": "https://cdn.example.com/assets/styles/theme.css",
    "slack": {
        "channel": "C789",
        "ts": "1609459200.000100"
    }
})

# Response: { "DistributionID": "E123...", "Bucket": "assets-bucket" }
```

### Technical Details

**CloudFront Path Pattern Examples**:

| Pattern | Regex | Matches |
|---------|-------|---------|
| `/images/*` | `^/images/.*$` | `/images/logo.png`, `/images/test.jpg` |
| `/*.js` | `^/.*\.js$` | `/app.js`, `/vendor/lib.js` |
| `/api/*` | `^/api/.*$` | `/api/users`, `/api/users/123` |

**Caching Strategy**:
- Distribution list: 24 hours (rarely changes)
- Bucket resolution: 24 hours (immutable after calculation)
- Cache key: `cdn.{profile}.distribution.{hostname}`
- Cache key: `cdn.{profile}.bucket.{dist-id}.{path}`

### Edge Cases

- **Multiple origins**: Uses first matching path pattern
- **No matching pattern**: Uses default cache behavior's origin
- **CloudFront origin path**: Included in returned bucket path
  - If origin is `my-bucket/static`, returns `my-bucket/static`
- **Non-S3 origins**: Returns domain name if pattern doesn't match

---

## Common Parameters

### Slack Parameter Block

All tools accept an optional `slack` object:

```json
{
  "channel": "C123ABC456DEF",  // Slack channel ID (not name)
  "ts": "1234567890.123456"     // Thread timestamp
}
```

**How to get channel ID**:
- Right-click channel in Slack → "View channel details" → Copy ID
- Or use Slack API: `/list` command

**How to get thread timestamp**:
- Right-click message → "Copy link"
- `?thread_ts=1234567890.123456` is your timestamp
- Or via Slack API: `message.ts` field

### Error Response Format

All tools return errors in this format:

```json
{
  "type": "text",
  "text": "request failed. error: distribution not found"
}
```

Error messages appear both in:
1. MCP response to N8N
2. Slack thread (if channel/ts provided)

---

## Performance Characteristics

| Tool | Typical Duration | Max Duration | API Calls |
|------|------------------|--------------|-----------|
| cloudfront-invalidate | 10-30 seconds | 5 minutes | 2-3 (resolve + create + wait) |
| file-replace | 5-30 seconds | N/A (depends on file size) | 4-5 (2 downloads + 2 uploads) |
| bucket-resolver | <1 second | 2 seconds | 1-2 (cached after first call) |

---

## Slack Emoji Reference

Used in notifications:

| Emoji | Meaning | CSS |
|-------|---------|-----|
| `:spin:` | Processing/In progress | animation |
| `:greentick1:` | Success/Complete | ✓ |
| `:x:` | Error/Failed | ✗ |

---

## N8N Integration Examples

### Example 1: Simple Invalidation
```
Input Node (URL from Slack message)
  → Claude (with MCP)
    → cloudfront-invalidate
  → Response Node (post result)
```

### Example 2: File Update Workflow
```
Slack Webhook (file uploaded)
  → Parse filename
  → Claude (with MCP)
    → bucket-resolver (find current location)
    → file-replace (upload new version)
  → Post confirmation
```

### Example 3: Batch Operations
```
Input (array of URLs)
  → Loop: For each URL
    → cloudfront-invalidate
  → Aggregate results
  → Post summary
```

---

## Debugging Tools

### Test a Tool Directly

Using curl (if you have stdio access):
```bash
# Start server
go run .

# In another terminal, create test.json
cat > test.json << 'EOF'
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "cloudfront-invalidate",
    "arguments": {
      "url": "https://example.com/test.png"
    }
  }
}
EOF

# Send to server
cat test.json | go run .
```

### Inspect Cache

Add debug logging:
```go
// In cdn.go
logger.Infof("cache key: %s", key)
logger.Infof("cache hit: %v", cacher.Has(key))
```

### Monitor API Calls

Enable AWS SDK logging:
```go
// In packages/aws/client/aws.go
cfg, _ := config.LoadDefaultConfig(ctx,
    config.WithClientLogMode(aws.LogRequestWithBody),
)
```

---

## Rate Limits & Quotas

**CloudFront API**:
- 100 requests/sec
- Burst: 200/sec
- Invalidation: 3000/month (free tier)

**S3 API**:
- Effective request rate: 3500 PUT/COPY/POST/DELETE per second
- 5500 GET/HEAD per second

**Slack API**:
- File upload: 1 per 1 second per bot
- Message posting: 1 per 1 second per channel

These tools should not hit limits in normal usage.

---

## Security

### Input Validation

- **URL validation**: Parsed with Go's `net/url`
- **Path validation**: URL-decoded to prevent directory traversal
- **Slack channel**: Must match CloudFront distribution (no direct S3 access from arbitrary channels)

### Output Sanitization

- **Error messages**: Don't expose full paths
- **Console links**: Include distribution ID only (anonymized to user)
- **Logs**: Slack tokens truncated in logs

### Permissions Required

**AWS IAM**:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "cloudfront:GetDistribution",
        "cloudfront:GetDistributionConfig",
        "cloudfront:ListDistributions",
        "cloudfront:CreateInvalidation",
        "cloudfront:GetInvalidation"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": "arn:aws:s3:::my-bucket/*"
    }
  ]
}
```

**Slack**:
- `files:read` - Download files
- `files:write` - Upload files
- `chat:write` - Post messages
- `channels:read` - Read channel info (optional)
