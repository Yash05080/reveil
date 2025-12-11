# Reveil API: Frontend Integration Guide

## 1. Connection & Authentication

### Base URL
```
http://localhost:8080/api
```

### Authentication
All API requests (except `/health`) require a **JWT Token**.
- **Header**: `Authorization: Bearer <YOUR_JWT_TOKEN>`
- **Token Source**: Currently generated via `cmd/generate_token`. In production, this will come from Supabase Auth login.

### Error Handling
Common HTTP Status Codes:
- `200 OK` / `201 Created`: Success.
- `400 Bad Request`: Validation error (check `details` in JSON).
- `401 Unauthorized`: Missing or invalid token.
- `403 Forbidden`: Valid token but actions not allowed (e.g., deleting someone else's post).
- `500 Internal Server Error`: Server-side issue.

---

## 2. Key Data Models (TypeScript Interfaces)

Use these interfaces to type your frontend responses.

```typescript
// The main Post object returned by the API
interface Post {
  id: string;
  community_id: string;
  user_id: string;
  
  // Decrypted content - display directly
  content: string; 
  content_type: 'text' | 'image';
  image_url?: string;
  
  like_count: number;
  comment_count: number;
  
  created_at: string; // ISO 8601
  updated_at: string;
  
  is_edited: boolean;
  is_removed: boolean;
  
  // Important: Present ONLY if the post was flagged
  moderation?: ModerationStatus; 
}

interface ModerationStatus {
  is_flagged: boolean;      // If true, show warning
  flag_reason?: string;     // e.g., "Contains blocked phrase..."
  severity_level?: number;  // 1-5 scale. 5 is highest.
}

interface APIResponse<T> {
  success: boolean;
  data: T;
  error?: string;   // Present on failure
  code?: string;    // App-specific error code
}
```

---

## 3. Workflows & Endpoints

### A. Viewing the Feed (List Posts)
**Endpoint**: `GET /communities/{community_id}/posts`

**Query Params**:
- `limit`: (Optional) Number of posts (Default: 20, Max: 50).
- `before`: (Optional) Timestamp for pagination (load older posts).
- `user_id`: (Optional) Filter by specific user.
- `content_type`: (Optional) Filter by 'text' or 'image'.

**Usage**:
1. Encryption is transparent. The `content` field contains plain text.
2. Check `post.moderation`. If present, overlay a warning (e.g., "This post may contain harmful content").
3. **Pagination**: Take the `created_at` of the last post in the list and pass it as `before` in the next request.

---

### B. Creating a Post
**Endpoint**: `POST /communities/{community_id}/posts`

**Payload**:
```json
{
  "content": "Hello World",
  "content_type": "text",
  "image_url": "https://example.com/image.png" // Optional
}
```

**Response (201 Created)**: Returns the created `Post` object.
**Moderation Handling**:
- **Light Model (Instant)**: If the response includes `moderation` (e.g., specific keywords), show the user an immediate warning: *"Your post was flagged for: [reason]"*.
- **Heavy Model (Async)**: The post might be flagged *after* creation. Listen to SSE (see below) or refresh the feed to see updated flags.

---

### C. Real-time Updates (SSE)
**Endpoint**: `GET /communities/{community_id}/posts/stream`

**Connection**:
Use `EventSource` to listen for real-time updates (new posts).

```javascript
const sse = new EventSource(`http://localhost:8080/api/communities/${communityId}/posts/stream`);

sse.onmessage = (event) => {
  const payload = JSON.parse(event.data);
  
  if (payload.event_type === "post_created") {
    const newPost = payload.payload; // This is a Post object
    // Prepend 'newPost' to your feed list
  }
};
```

---

## 4. Moderation UI Guidelines

Since the API supports dual-layer moderation, the frontend should reflect this:

1.  **Severity 5 (Blocked Phrases)**:
    *   Consider blurring the content by default.
    *   Show a red warning badge: "Violates Community Guidelines".
    *   Display `flag_reason`.

2.  **Severity 3 (AI Flagged)**:
    *   Show a yellow/orange warning icon.
    *   Text: "Flagged by AI as potentially inappropriate".
    *   Allow user to click "Show Content".

3.  **No Flag**:
    *   Display normally.

## 5. Summary of Endpoints

| Method | URL | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Check API status (No Auth required). |
| `GET` | `/communities/{id}/posts` | List posts (Feed). Pagination supported. |
| `POST` | `/communities/{id}/posts` | Create new post. Triggers moderation. |
| `PUT` | `/posts/{id}` | Update post content. |
| `DELETE` | `/posts/{id}` | Soft delete post. |
| `GET` | `/communities/{id}/posts/stream` | SSE stream for real-time posts. |
