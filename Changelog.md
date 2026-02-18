# Changelog

## internal/data/models.go
- Added `Stores StoreModel` to `Models` struct
- Added `StoreModel{Db: db}` in `NewModels`

---

## internal/data/pages.go
- `Page.AppId uuid.UUID` renamed to `Page.StoreId uuid.UUID`, JSON tag updated to `store_id`
- All SQL queries updated: `app_id` column references changed to `store_id`
- `GetAll()` removed
- `GetAllForApp(appID uuid.UUID)` replaced by `GetAllForStore(storeId uuid.UUID)`
- `changeIsHomePage`: SQL updated from `app_id` to `store_id`

---

## internal/data/store.go (new file)
- Added `Store` struct: `Id`, `Name`, `Slug`, `CreatedAt`, `UpdatedAt`
- Added `StoreModel` with methods: `Insert`, `Get`, `GetAll`, `Update`, `Delete`
- `Insert` and `Update` use `errors.As` for pq unique violation check

---

## internal/data/widgets.go
- Added `Get(id uuid.UUID) (*Widget, error)`
- `Update` query: added explicit `updated_at = NOW()`

---

## cmd/api/helpers.go
- `readIdParam` signature changed from `readIdParam(r *http.Request)` to `readIdParam(r *http.Request, name string)` — now reads any named URL param instead of hardcoding `"id"`

---

## cmd/api/store.go (new file)
- Added handlers: `listStoresHandler`, `showStoreHandler`, `createStoreHandler`, `updateStoreHandler`, `deleteStoreHandler`

---

## cmd/api/pages.go
- `listPagesForAppHandler` removed
- All handlers extract `store_id` from URL via `readIdParam(r, "store_id")`
- `listPagesHandler`: calls `GetAllForStore(id)` instead of removed `GetAll()`
- `createPageHandler`: removed `AppId` from input struct, added store existence check, `page.StoreId` set from URL param, Location header updated to `/stores/:store_id/pages/:page_id`
- `showPageHandler`, `updatePageHandler`, `deletePageHandler`: read both `store_id` and `page_id` from URL, added ownership check `if page.StoreId != storeId { notFoundResponse }`

---

## cmd/api/widgets.go
- `createWidgetHandler`: reads `store_id` and `page_id` from URL, added page ownership check, Location header updated to `/stores/:store_id/widgets/:id`
- `updateWidgetHandler`: reads `store_id` and `page_id` from URL, added page ownership check
- `reorderWidgetsHandler`: reads `store_id` and `page_id` from URL, added page ownership check

---

## cmd/api/routes.go
- Added store routes: `GET /stores`, `POST /stores`, `GET /stores/:store_id`, `PUT /stores/:store_id`, `DELETE /stores/:store_id`
- Page routes moved from `/pages` to `/stores/:store_id/pages`, param renamed from `:id` to `:page_id`
- Widget create and reorder routes moved from `/pages/:id/widgets` to `/stores/:store_id/pages/:page_id/widgets`
- Widget update and delete routes changed from `/widgets/:id` to `/stores/:store_id/widgets/:id`
- Removed `listPagesForAppHandler` route
