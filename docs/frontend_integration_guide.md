# Reveil API - Frontend Integration Guide

## 1. Authentication
The API relies on **Supabase Auth**.

### Headers
Every request to the backend must include the `Authorization` header with a valid JWT from Supabase.
```
Authorization: Bearer <SUPABASE_JWT>
```

---

## 2. Communities
Community management is minimal for now.

### **List Communities**
`GET /api/communities`
- **Response**: Array of communities.
- **Usage**: Use `id` for creating posts.

---

## 3. Posts
Posts are encrypted at rest. The API handles decryption before sending response to frontend (for MVP).
**New**: All posts require a `title` and `content`.

### **Create Post**
`POST /api/communities/{community_id}/posts`

- **Payload**:
```json
{
  "title": "My Post Title",         // [REQUIRED] Post Title
  "content": "My post body text...", // [REQUIRED] Post Body
  "content_type": "text",            // [REQUIRED] "text", "image", "link"
  "image_url": "https://..."         // [OPTIONAL] Unencrypted URL for images
}
```

- **Response**:
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "title": "My Post Title",
    "content": "My post body text...",
    "image_url": "https://...",
    "created_at": "...",
    "moderation": { ... } // If flagged immediately (rare)
  }
}
```

### **List Posts**
`GET /api/communities/{community_id}/posts`
- **QueryParams**: `limit` (default 20, max 50), `before` (timestamp for pagination).
- **Response**: Array of Post Objects. Note: `title` and `content` are decrypted strings.

### **Edit Post**
`PUT /api/posts/{post_id}`
- **Payload**:
```json
{
  "title": "Updated Title", // [OPTIONAL]
  "content": "Updated content..." // [REQUIRED]
}
```
- **Behavior**: Updates trigger re-moderation. If toxic, response includes `moderation` flags.

### **Delete Post**
`DELETE /api/posts/{post_id}`

### **Report Post**
`POST /api/posts/{post_id}/report`
- **Payload**: `{"reason": "Spam"}`

---

## 4. Comments
Comments support infinite nesting via `parent_id`.

### **Create Comment**
`POST /api/communities/{community_id}/posts/{post_id}/comments`
- **Payload**:
```json
{
  "content": "This is a comment",
  "parent_id": "uuid" // [OPTIONAL] For nested replies
}
```

### **List Comments**
`GET /api/communities/{community_id}/posts/{post_id}/comments`
- **Response**: Flat list of comments.
- **Frontend Logic**: You must reconstruct the tree using `id` and `parent_id`.

---

## 5. Likes
`POST /api/posts/{post_id}/like`
`POST /api/comments/{comment_id}/like`
- Toggles like status. Returns updated count.

---

## 6. Real-time (SSE)
Connect to: `GET /api/events?community_id={id}`
- Events: `post_created`, `comment_created`.

---

## 7. Moderation
If a post is flagged (on create or edit), the response `data` will contain a `moderation` object:
```json
"moderation": {
  "is_flagged": true,
  "flag_reason": "Contains blocked phrase",
  "severity_level": 5
}
```
**UI Handling**: You should handle these flags (e.g., show warning, blur content).

---

## 8. Images
- Use `image_url` field for images.
- This field is **never encrypted**.
- The `content` field is for the text body and **is encrypted**.
