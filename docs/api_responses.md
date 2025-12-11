# Révéil API Response Standard

All API responses follow the standard `JSend` inspired format.

## Success Response
```json
{
  "success": true,
  "data": { ... } // Object or Array
}
```

## Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE", // e.g., VALIDATION_ERROR, UNAUTHORIZED
    "message": "Human readable message"
  }
}
```

## Data Objects

### Post Object
```json
{
  "id": "uuid",
  "community_id": "uuid",
  "user_id": "uuid",
  "title": "Post Title",       // Decrypted
  "content": "Post Content",   // Decrypted
  "content_type": "text",      // text, image, link
  "image_url": "https://...",  // Nullable
  "like_count": 0,
  "comment_count": 0,
  "created_at": "ISO8601",
  "updated_at": "ISO8601",
  "is_edited": false,
  "is_removed": false,
  "moderation": {              // Only present if flagged
    "is_flagged": true,
    "flag_reason": "Toxic",
    "severity_level": 5
  }
}
```

### Comment Object
```json
{
  "id": "uuid",
  "post_id": "uuid",
  "parent_id": "uuid", // or null
  "user_id": "uuid",
  "content": "Comment text", // Decrypted
  "like_count": 0,
  "created_at": "ISO8601",
  "is_removed": false
}
```
