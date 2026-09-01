# Blog
This is a Work In Progress Blog engine to host my own personal blog.

## Tech Stack
#### Infrastructure
- Docker Compose: to simplify cross-platform development & deployment
- Nginx: routes requests to appropriate API endpoints
#### Frontend
- HTML: Mostly generated from backend
    - html-minifier-next: Minifies HTML files for more efficient net usage
- SCSS: Page styling
- TS: Improvements to Frontend behavior
    - esbuild: TS transpiler and bundler
#### Backend
- Go: Generates HTML pages from templates, and used as the main API
    - jackc/pgx: go bindings to postgreSQL
    - valyala/fasttemplate: Makes HTML template files more efficient to use
    - crypto: Used for secure password generation and validation
    - google/uuid: Generates random file names without the risk of name collision
- Python: Isolated worker for safe image re-encoding
    - FastAPI: Used to simplify communication between backend and workers (Should remove eventually)
    - Pillow: Used for profile picture re-encoding
#### Database
- PostgreSQL: Manages user accounts, translations, comments, etc

## Current features
- Simple user accounts
- Safely sanitized user profile pictures
- SQL page management
- Simple Bilingual page translation system
- Side by side page editor

## Planned Features (Roughly in planned order)
- Comment threads
- AB testing system
- A half decent looking fronted
- Selfhost with a domain and HTTPS

## Nice to haves
- Early access subscriptions
- SEO optimizations
- Lynx support (No JS or CSS needed)
- Bot prevention with Anubis
- Data analysis panel
- Dead Man's switch to ensure travel safety
- A really good looking fronted
- .onion site available. A Vanity onion would be nice
- Host on a proper server

### Account Privilege levels
It would probably be better to switch to something more granular eventually
0. ReadOnly    : Not signed in
1. User        : Standard User Account.
2. EarlyAccess : First payment tier. Early Access
3. Bonus       : Second payment tier. Add additional perks. Might go unused
4. VIP         : Third payment tier. Add additional perks
5. Trusted     : Friends and Family. Special options if I add a dead man's switch.
6. Owner       : Allows for creating pages

### Page AB testing system plan
- A cookie should be set by the backend indicating the page version
- One cookie corresponds to multiple pages, with positional arguments indicating which
- The backend needs to share a page index with nginx
- This page index will be used to determine the page cache.
- a blank 2 chars indicates that a test was not decided.
- On first page access, the cookie will update to a random active test
- There should also be a global test cookie, which will also influence the cache
