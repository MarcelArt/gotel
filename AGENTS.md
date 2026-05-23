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

## General Coding Standards

- **Maintain Documentation Integrity:** Retain all pre-existing comments, docstrings, and struct annotations unless explicitly requested otherwise.
- **Link Reference Formatting:** When mentioning files or functions, use clickable links in Markdown format referencing absolute file URIs or code lines.
