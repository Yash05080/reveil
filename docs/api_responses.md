# API Response Structures

## 1. Safe Post (201 Created / 200 OK)
When a post is created or listed and has **no flags**:

```json
{
  "success": true,
  "data": {
    "id": "984b04bd-83cf-418a-aecc-986d22f09692",
    "community_id": "aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "user_id": "11111111-1111-1111-1111-111111111111",
    "content": "I really think we should discuss the impact of climate change on local agriculture.",
    "content_type": "text",
    "like_count": 0,
    "comment_count": 0,
    "created_at": "2025-12-11T14:14:02.95Z",
    "updated_at": "2025-12-11T14:14:02.95Z",
    "is_edited": false,
    "is_removed": false
  }
}
```

## 2. Flagged Post (201 Created / 200 OK)
When a post is created or listed and **has been flagged** (by either Light or Heavy moderation):
The response includes a `moderation` object. You can use this to show a warning to the user on the frontend.

```json
{
  "success": true,
  "data": {
    "id": "5e2a56a4-e82e-4971-a44c-21092850426b",
    "community_id": "aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "user_id": "11111111-1111-1111-1111-111111111111",
    "content": "I will doxx you and your family.", 
    "content_type": "text",
    ...
    "moderation": {
      "is_flagged": true,
      "flag_reason": "Contains blocked phrase: 'doxx you'",
      "severity_level": 5
    }
  }
}
```

### Key Fields
- `moderation.is_flagged`: (bool) True if any flag exists.
- `moderation.flag_reason`: (string) The reason for the highest severity flag (e.g., "Hate Speech (ML)" or "Contains blocked phrase...").
- `moderation.severity_level`: (int) 1-5. 
  - 5: Blocked Phrase / Severe.
  - 3: ML Flagged / Medium.

## 3. Decryption
The `content` field in the response is always the **Decrypted Plain Text**. The backend handles encryption/decryption transparently.
