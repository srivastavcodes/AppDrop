# API Request Examples

This file contains example requests for all API endpoints using curl commands.

## Base URL
```
http://localhost:4000
```

---

## Page Endpoints

### 1. Create Page - POST /pages

**Request:**
```bash
curl -X POST http://localhost:4000/pages \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Home",
    "route": "/home",
    "is_home": true
  }'
```

**JSON Payload:**
```json
{
  "name": "Home",
  "route": "/home",
  "is_home": true
}
```

**Response (201 Created):**
```json
{
  "page": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "app_id": "00000000-0000-0000-0000-000000000001",
    "name": "Home",
    "route": "/home",
    "is_home": true,
    "created_at": "2024-02-10T10:30:00Z",
    "updated_at": "2024-02-10T10:30:00Z"
  }
}
```

**More Examples:**

Create a Collection page:
```bash
curl -X POST http://localhost:4000/pages \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Collections",
    "route": "/collections",
    "is_home": false
  }'
```

Create a Product Detail page:
```bash
curl -X POST http://localhost:4000/pages \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Product Detail",
    "route": "/product/:id",
    "is_home": false
  }'
```

**Error Cases:**

Empty name (400):
```bash
curl -X POST http://localhost:4000/pages \
  -H "Content-Type: application/json" \
  -d '{
    "name": "",
    "route": "/test"
  }'
```

Duplicate route (409):
```bash
curl -X POST http://localhost:4000/pages \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Another Home",
    "route": "/home"
  }'
```

---

### 2. List All Pages - GET /pages

**Request:**
```bash
curl http://localhost:4000/pages
```

**Response (200 OK):**
```json
{
  "pages": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "app_id": "00000000-0000-0000-0000-000000000001",
      "name": "Home",
      "route": "/home",
      "is_home": true,
      "created_at": "2024-02-10T10:30:00Z",
      "updated_at": "2024-02-10T10:30:00Z"
    },
    {
      "id": "223e4567-e89b-12d3-a456-426614174001",
      "app_id": "00000000-0000-0000-0000-000000000001",
      "name": "Collections",
      "route": "/collections",
      "is_home": false,
      "created_at": "2024-02-10T10:31:00Z",
      "updated_at": "2024-02-10T10:31:00Z"
    }
  ]
}
```

---

### 3. Get Single Page - GET /pages/:id

**Request:**
```bash
curl http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000
```

**Response (200 OK):**
```json
{
  "page": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "app_id": "00000000-0000-0000-0000-000000000001",
    "name": "Home",
    "route": "/home",
    "is_home": true,
    "created_at": "2024-02-10T10:30:00Z",
    "updated_at": "2024-02-10T10:30:00Z",
    "widgets": [
      {
        "id": "456e7890-e89b-12d3-a456-426614174002",
        "page_id": "123e4567-e89b-12d3-a456-426614174000",
        "type": "banner",
        "position": 0,
        "config": {
          "image_url": "https://example.com/banner.jpg",
          "link": "/collections/new"
        },
        "created_at": "2024-02-10T10:32:00Z",
        "updated_at": "2024-02-10T10:32:00Z"
      }
    ]
  }
}
```

**Error Case (404):**
```bash
curl http://localhost:4000/pages/00000000-0000-0000-0000-000000000000
```

Response:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "the requested resource could not be found"
  }
}
```

---

### 4. Update Page - PUT /pages/:id

**Request (Update name only):**
```bash
curl -X PUT http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Home Page"
  }'
```

**JSON Payload (partial update):**
```json
{
  "name": "Home Page"
}
```

**Request (Update multiple fields):**
```bash
curl -X PUT http://localhost:4000/pages/223e4567-e89b-12d3-a456-426614174001 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "All Collections",
    "route": "/collections-all",
    "is_home": false
  }'
```

**JSON Payload (full update):**
```json
{
  "name": "All Collections",
  "route": "/collections-all",
  "is_home": false
}
```

**Request (Set as home page):**
```bash
curl -X PUT http://localhost:4000/pages/223e4567-e89b-12d3-a456-426614174001 \
  -H "Content-Type: application/json" \
  -d '{
    "is_home": true
  }'
```

**Response (200 OK):**
```json
{
  "page": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "app_id": "00000000-0000-0000-0000-000000000001",
    "name": "Home Page",
    "route": "/home",
    "is_home": true,
    "created_at": "2024-02-10T10:30:00Z",
    "updated_at": "2024-02-10T10:35:00Z"
  }
}
```

**Error Cases:**

Empty name (400):
```bash
curl -X PUT http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": ""
  }'
```

Duplicate route (409):
```bash
curl -X PUT http://localhost:4000/pages/223e4567-e89b-12d3-a456-426614174001 \
  -H "Content-Type: application/json" \
  -d '{
    "route": "/home"
  }'
```

---

### 5. Delete Page - DELETE /pages/:id

**Request:**
```bash
curl -X DELETE http://localhost:4000/pages/223e4567-e89b-12d3-a456-426614174001
```

**Response (200 OK):**
```json
{
  "message": "page successfully deleted"
}
```

**Error Cases:**

Delete home page (409):
```bash
curl -X DELETE http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000
```

Response:
```json
{
  "error": {
    "code": "CONFLICT",
    "message": "cannot delete home page"
  }
}
```

Non-existent page (404):
```bash
curl -X DELETE http://localhost:4000/pages/00000000-0000-0000-0000-000000000000
```

---

## Widget Endpoints

### 6. Create Widget - POST /pages/:id/widgets

**Request (Banner widget):**
```bash
curl -X POST http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "banner",
    "config": {
      "image_url": "https://example.com/banner.jpg",
      "link": "/collections/new",
      "alt_text": "New Collection Banner"
    }
  }'
```

**JSON Payload:**
```json
{
  "type": "banner",
  "config": {
    "image_url": "https://example.com/banner.jpg",
    "link": "/collections/new",
    "alt_text": "New Collection Banner"
  }
}
```

**Request (Product Grid widget):**
```bash
curl -X POST http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "product_grid",
    "config": {
      "collection_id": "summer-collection",
      "columns": 2,
      "show_price": true,
      "show_add_to_cart": true
    }
  }'
```

**Request (Text widget):**
```bash
curl -X POST http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "text",
    "config": {
      "content": "Welcome to our store!",
      "text_align": "center",
      "font_size": 24,
      "color": "#333333"
    }
  }'
```

**Request (Image widget):**
```bash
curl -X POST http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "image",
    "config": {
      "src": "https://example.com/hero.jpg",
      "alt": "Hero Image",
      "aspect_ratio": "16:9"
    }
  }'
```

**Request (Spacer widget - no config):**
```bash
curl -X POST http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "spacer"
  }'
```

**Response (201 Created):**
```json
{
  "widget": {
    "id": "456e7890-e89b-12d3-a456-426614174002",
    "page_id": "123e4567-e89b-12d3-a456-426614174000",
    "type": "banner",
    "position": 0,
    "config": {
      "image_url": "https://example.com/banner.jpg",
      "link": "/collections/new",
      "alt_text": "New Collection Banner"
    },
    "created_at": "2024-02-10T10:32:00Z",
    "updated_at": "2024-02-10T10:32:00Z"
  }
}
```

**Error Cases:**

Invalid widget type (400):
```bash
curl -X POST http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "carousel"
  }'
```

Response:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "widget type must be one of: banner, product_grid, text, image, spacer"
  }
}
```

Page not found (404):
```bash
curl -X POST http://localhost:4000/pages/00000000-0000-0000-0000-000000000000/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "banner"
  }'
```

---

### 7. Update Widget - PUT /widgets/:id

**Request (Update type and config):**
```bash
curl -X PUT http://localhost:4000/widgets/456e7890-e89b-12d3-a456-426614174002 \
  -H "Content-Type: application/json" \
  -d '{
    "type": "image",
    "config": {
      "src": "https://example.com/updated-banner.jpg",
      "alt": "Updated Banner"
    }
  }'
```

**JSON Payload:**
```json
{
  "type": "image",
  "config": {
    "src": "https://example.com/updated-banner.jpg",
    "alt": "Updated Banner"
  }
}
```

**Request (Update position only):**
```bash
curl -X PUT http://localhost:4000/widgets/456e7890-e89b-12d3-a456-426614174002 \
  -H "Content-Type: application/json" \
  -d '{
    "position": 5
  }'
```

**Request (Update config only):**
```bash
curl -X PUT http://localhost:4000/widgets/456e7890-e89b-12d3-a456-426614174002 \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "image_url": "https://example.com/new-banner.jpg",
      "link": "/sale"
    }
  }'
```

**Response (200 OK):**
```json
{
  "widget": {
    "id": "456e7890-e89b-12d3-a456-426614174002",
    "page_id": "123e4567-e89b-12d3-a456-426614174000",
    "type": "image",
    "position": 0,
    "config": {
      "src": "https://example.com/updated-banner.jpg",
      "alt": "Updated Banner"
    },
    "created_at": "2024-02-10T10:32:00Z",
    "updated_at": "2024-02-10T10:40:00Z"
  }
}
```

**Error Cases:**

Invalid widget type (400):
```bash
curl -X PUT http://localhost:4000/widgets/456e7890-e89b-12d3-a456-426614174002 \
  -H "Content-Type: application/json" \
  -d '{
    "type": "video"
  }'
```

Widget not found (404):
```bash
curl -X PUT http://localhost:4000/widgets/00000000-0000-0000-0000-000000000000 \
  -H "Content-Type: application/json" \
  -d '{
    "type": "text"
  }'
```

---

### 8. Delete Widget - DELETE /widgets/:id

**Request:**
```bash
curl -X DELETE http://localhost:4000/widgets/456e7890-e89b-12d3-a456-426614174002
```

**Response (200 OK):**
```json
{
  "message": "widget successfully deleted"
}
```

**Error Case (404):**
```bash
curl -X DELETE http://localhost:4000/widgets/00000000-0000-0000-0000-000000000000
```

Response:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "the requested resource could not be found"
  }
}
```

---

### 9. Reorder Widgets - POST /pages/:id/widgets/reorder

**Request:**
```bash
curl -X POST http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000/widgets/reorder \
  -H "Content-Type: application/json" \
  -d '{
    "widget_ids": [
      "789e1234-e89b-12d3-a456-426614174005",
      "456e7890-e89b-12d3-a456-426614174002",
      "567e8901-e89b-12d3-a456-426614174003",
      "678e9012-e89b-12d3-a456-426614174004"
    ]
  }'
```

**JSON Payload:**
```json
{
  "widget_ids": [
    "789e1234-e89b-12d3-a456-426614174005",
    "456e7890-e89b-12d3-a456-426614174002",
    "567e8901-e89b-12d3-a456-426614174003",
    "678e9012-e89b-12d3-a456-426614174004"
  ]
}
```

**Description:**
The array order determines the new widget positions:
- First ID → position 0
- Second ID → position 1
- Third ID → position 2
- etc.

**Response (200 OK):**
```json
{
  "message": "widgets successfully reordered"
}
```

**Error Cases:**

Empty array (400):
```bash
curl -X POST http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000/widgets/reorder \
  -H "Content-Type: application/json" \
  -d '{
    "widget_ids": []
  }'
```

Response:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "widget_ids array cannot be empty"
  }
}
```

Widget doesn't belong to page (400):
```bash
curl -X POST http://localhost:4000/pages/123e4567-e89b-12d3-a456-426614174000/widgets/reorder \
  -H "Content-Type: application/json" \
  -d '{
    "widget_ids": [
      "456e7890-e89b-12d3-a456-426614174002",
      "00000000-0000-0000-0000-999999999999"
    ]
  }'
```

Response:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "some widgets do not belong to this page"
  }
}
```

Page not found (404):
```bash
curl -X POST http://localhost:4000/pages/00000000-0000-0000-0000-000000000000/widgets/reorder \
  -H "Content-Type: application/json" \
  -d '{
    "widget_ids": ["456e7890-e89b-12d3-a456-426614174002"]
  }'
```

---

## Complete Workflow Example

Here's a complete workflow to build a home page with widgets:

```bash
# 1. Create home page
curl -X POST http://localhost:4000/pages \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Home",
    "route": "/home",
    "is_home": true
  }'
# Note the returned page ID, e.g., PAGE_ID=123e4567-e89b-12d3-a456-426614174000

# 2. Add a banner widget
curl -X POST http://localhost:4000/pages/PAGE_ID/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "banner",
    "config": {
      "image_url": "https://example.com/hero.jpg"
    }
  }'
# Note: WIDGET1_ID

# 3. Add a product grid
curl -X POST http://localhost:4000/pages/PAGE_ID/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "product_grid",
    "config": {
      "collection_id": "featured"
    }
  }'
# Note: WIDGET2_ID

# 4. Add a text widget
curl -X POST http://localhost:4000/pages/PAGE_ID/widgets \
  -H "Content-Type: application/json" \
  -d '{
    "type": "text",
    "config": {
      "content": "Shop our latest collection"
    }
  }'
# Note: WIDGET3_ID

# 5. Reorder widgets (text first, then banner, then grid)
curl -X POST http://localhost:4000/pages/PAGE_ID/widgets/reorder \
  -H "Content-Type: application/json" \
  -d '{
    "widget_ids": ["WIDGET3_ID", "WIDGET1_ID", "WIDGET2_ID"]
  }'

# 6. Get the complete page with widgets
curl http://localhost:4000/pages/PAGE_ID

# 7. Update the banner
curl -X PUT http://localhost:4000/widgets/WIDGET1_ID \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "image_url": "https://example.com/new-hero.jpg",
      "link": "/new-arrivals"
    }
  }'

# 8. List all pages
curl http://localhost:4000/pages
```

---

## Testing Invalid Requests

### Malformed JSON
```bash
curl -X POST http://localhost:4000/pages \
  -H "Content-Type: application/json" \
  -d '{invalid json'
```

### Unknown Fields
```bash
curl -X POST http://localhost:4000/pages \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "route": "/test",
    "unknown_field": "value"
  }'
```

### Invalid UUID Format
```bash
curl http://localhost:4000/pages/not-a-uuid
```

### Missing Required Fields
```bash
curl -X POST http://localhost:4000/pages \
  -H "Content-Type: application/json" \
  -d '{}'
```

---

## Notes

1. **UUIDs**: Replace example UUIDs with actual IDs from your responses
2. **Content-Type**: Always include `Content-Type: application/json` header for POST/PUT
3. **Response Headers**: POST requests return `Location` header with resource URL
4. **Position**: Widget position is auto-assigned on creation (starts at 0)
5. **Config**: Widget config is optional and can be empty object or omitted
6. **Partial Updates**: PUT endpoints support partial updates (only send fields to change)
