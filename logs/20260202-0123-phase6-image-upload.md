# Phase 6: Image Upload

**Date**: 2026-02-02  
**Phase**: 6 - Image Upload

## Summary

Implemented profile image upload functionality using AWS SDK v2 with Garage S3-compatible storage. Members can now upload, complete, and delete their profile images through pre-signed URLs.

## Changes

### AWS SDK v2 Integration

Added AWS SDK for Go v2 with support for:
- S3-compatible endpoints (Garage)
- Pre-signed PUT URLs for uploads
- Pre-signed GET URLs for downloads
- Path-style addressing (required for Garage)

### S3 Service (`internal/storage/s3.go`)

- `NewS3Service()` - Creates S3 service with Garage configuration
- `GenerateUploadURL()` - Generates pre-signed PUT URL for image upload
- `GenerateDownloadURL()` - Generates pre-signed GET URL for image download
- `DeleteObject()` - Deletes objects from S3
- `ObjectExists()` - Checks if object exists in S3
- `GetPublicURL()` - Returns public URL format for objects

### Upload Handler (`internal/storage/handler.go`)

- `POST /upload/profile-image` - Generate pre-signed upload URL
  - Validates content type (JPEG, PNG, WebP)
  - Generates object key: `profile-images/{member_id}/{timestamp}.{ext}`
  - Returns upload URL with 15-minute expiration

- `POST /upload/profile-image/complete` - Confirm upload and update DB
  - Validates object exists in S3
  - Updates member's `image_url` in database

- `DELETE /upload/profile-image` - Delete profile image
  - Deletes from S3
  - Clears `image_url` in database

### Configuration (`internal/config/config.go`)

Added S3/Garage configuration:
- `S3_ENDPOINT` - Garage endpoint URL
- `S3_BUCKET` - Bucket name (default: lscs-core)
- `S3_ACCESS_KEY` - Access key ID
- `S3_SECRET_KEY` - Secret access key
- `S3_REGION` - Region (default: garage)

### Routes (`internal/server/routes.go`)

Added session-protected upload routes:
- `POST /upload/profile-image` → `uploadHandler.GenerateUploadURLHandler`
- `POST /upload/profile-image/complete` → `uploadHandler.CompleteUploadHandler`
- `DELETE /upload/profile-image` → `uploadHandler.DeleteImageHandler`

## Files Affected

**New Files:**
- `internal/storage/s3.go` - S3 service implementation
- `internal/storage/handler.go` - Upload handler

**Modified Files:**
- `go.mod` - Added AWS SDK v2 dependencies
- `internal/config/config.go` - Added S3 configuration
- `internal/server/server.go` - Added S3 service and upload handler initialization
- `internal/server/routes.go` - Added upload routes
- `PLAN.md` - Updated Phase 6 status

## Environment Variables

```env
# S3/Garage Storage (new)
S3_ENDPOINT=https://garage.example.com
S3_BUCKET=lscs-core
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key
S3_REGION=garage
```

## Notes

- Storage service is optional - API works without S3 configuration
- Pre-signed URLs expire after 15 minutes
- Supported image types: image/jpeg, image/png, image/webp
- Object keys follow pattern: `profile-images/{member_id}/{timestamp}.{ext}`
- Phase 7 will initialize the Next.js frontend
