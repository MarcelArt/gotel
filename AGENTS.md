# AI Agent Guidelines & Rules

This project enforces strict guidelines for AI coding agents to protect the workspace, preserve development state, and ensure code cleanliness.

## Critical Rules (Strictly Enforced)

1. **No Server Execution**
   - **Rule:** Do not start or run the backend web server (e.g., `go run main.go serve` or similar).
   - **Rationale:** Running the server in the background can conflict with user-managed sessions, ports, or background states. Only the USER is allowed to start, stop, or manage the application server.

2. **No Database Migrations**
   - **Rule:** Do not run migration CLI commands (e.g., `go run main.go migrate` or `--drop`).
   - **Rationale:** Wiping or modifying database tables disrupts user-seeded testing data. Only the USER is allowed to execute database migrations.

3. **Strict Separation of Concerns (Frontend vs. Backend)**
   - **Rule:** Do not touch the backend code (`internal/v1/features/` or other backend repository/service directories) if the chat session or active task is focused on frontend features (views, HTML templates, CSS, JS, handlers, routers).
   - **Rule:** If a backend modification is needed to support a frontend feature, do not make the change yourself. Instead, document the exact issue, code path, and recommendation, and request that the USER apply the fix.
   - **Rule:** Conversely, if the task is explicitly focused on backend refactoring/features, do not modify frontend templates or stylesheets.

4. **UI Layout and State Rules (Collapsible Sidebar, Breadcrumbs, & File Uploads)**
   - **Collapsible Sidebar:** The sidebar collapse/minimize state is managed via Alpine.js and stored in `localStorage` under `sidebarCollapsed` (persisting the state across navigations). Always use `:class="{ 'collapsed': sidebarCollapsed }"` for toggling responsive styles on sidebar/main content.
   - **Dynamic Breadcrumbs:** Page breadcrumbs are located in the sticky top header. Main navigation tabs should clear custom sub-breadcrumbs, while nested/detail pages (such as asset instances, stock logs, or transaction lists) must dispatch their breadcrumb hierarchy on load via Alpine's `$dispatch('update-breadcrumbs', [...])`. Clickable breadcrumbs must route back via HTMX (`hx-get`) and trigger the correct active `currentTab` state in the sidebar.
   - **Modernized File Uploads:** All file upload fields must be styled as modern drag-and-drop dropzone components instead of plain default browser input controls. They must utilize Alpine.js to show dynamic image preview thumbnails and allow clearing selected files via a reset button.

## General Coding Standards

- **Maintain Documentation Integrity:** Retain all pre-existing comments, docstrings, and struct annotations unless explicitly requested otherwise.
- **Link Reference Formatting:** When mentioning files or functions, use clickable links in Markdown format referencing absolute file URIs or code lines.
- **CSS Cache Busting:** Whenever making changes to CSS files, you must update/increment the version query parameter (e.g., `?v=X.Y.Z`) in the template's `<link>` stylesheet tags to bypass browser caching and ensure the updates are immediately active.
