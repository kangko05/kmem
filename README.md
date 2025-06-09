# KMEM - Personal Cloud Storage Solution

## Problem Statement

My parents constantly struggled with limited phone storage space.

## Tech Stack

**Backend**: Go, Gin Framework, PostgreSQL, FFmpeg  
**Frontend**: React, TypeScript, TailwindCSS  
**Infrastructure**: Docker

## Architecture

<div align="center">
  <img src="./architecture.svg" alt="System Architecture" width="700" style="max-width: 100%; height: auto;">
</div>

## Key Features

### File Management

- Duplicate detection using SHA256 hashing
- Support for images (JPEG, PNG, GIF) and videos (MP4, AVI, MOV)
- Search and filter functionality with infinite scroll

### Performance

- Background thumbnail generation using Go routines
- In-memory caching with TTL and LRU eviction
- Multiple thumbnail sizes for responsive loading
- Queue-based processing to prevent UI blocking

### Security & Reliability

- JWT authentication with automatic token refresh
- File validation and type checking
- Soft delete with scheduled cleanup jobs
- ZFS filesystem for data integrity and snapshots

## Technical Implementation

### Async Processing

- Implemented worker pool for better performance

### Storage Strategy

- ZFS filesystem for data integrity and compression
- Automatic snapshots for backup and recovery
- Direct filesystem integration for optimal performance

## Results

- **Performance**: Background processing prevents upload blocking
- **Storage**: Duplicate detection reduces unnecessary storage usage
- **Reliability**: ZFS snapshots provide automatic backup protection

## Technical Challenges Solved

1. **Thumbnail Generation** → Background processing with worker queues
2. **User Experience** → Responsive design
3. **Data Protection** → ZFS snapshots and duplicate detection

## Quick Start

Start with docker:

```bash
POSTGRES_PASSWORD=yourpass JWT_SECRET_KEY=jwtsecret docker-compose up
```

Run development script:

```bash
dev.sh
```

Add account:

```shell
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"username": "testuser", "password":"testpass"}' \ # username and password here
    localhost:8000/auth/signup
```

- Access application at 'http://localhost:5173'
- Currently using vite dev server for frontend for now

## Future Enhancements

### Core Features

- [ ] Album functionality for photo organization
- [ ] Tag system for better searching and filtering performance
- [ ] Automatic tagging system using AI
- [ ] Bulk file operations (delete, tag, rename)

### Advanced Features

- [ ] Metadata extraction (EXIF data utilization for search)
- [ ] Progressive Web App (PWA) for mobile app-like experience
- [ ] Bulk upload with drag & drop (entire folder upload)

### System & Monitoring

- [ ] System logging implementation
- [ ] System monitoring and metrics
- [ ] Storage quota management (per-user limits)
- [ ] API rate limiting

### Performance & Infrastructure

- [ ] CDN integration for global content delivery
- [ ] Image optimization (WebP conversion, etc.)
- [ ] Redis clustering for enhanced caching
- [ ] Lazy loading improvements
- [ ] Real-time sync across devices

### Backup & Storage

- [ ] ZFS snapshot-based backup automation
- [ ] Cloud backup integration (S3, Google Drive)
- [ ] Advanced compression algorithms

---

**A practical solution for family photo storage and sharing**
