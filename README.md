# Blog
This is a Work In Progress Blog engine to host my own personal blog.

The current stage is to migrate my existing work over to this new Docker + Go architecture

## Current features
- Basic Docker setup
- HTML is minified
- Compiled TS and SCSS files
- Go backend
- Docker Build caching

## Planned Features (Roughly in planned order)
- User accounts
- SQL page management
- Comment threads
- AB testing system
- Bilingual English 日本語 page translations
- A half decent looking fronted
- Selfhost with a domain and HTTPS

## Account Privilege levels
0. Standard User Account. 
1. First payment tier. Early Access
2. Second payment tier. Add additional perks. Might go unused
3. Third payment tier. Add additional perks
4. A level for friends and family. Special options if I add a dead man's switch
5. Owner

## Nice to haves
- Early access subscriptions
- SEO optimizations
- Lynx support (No JS or CSS needed)
- Data analysis panel
- Dead Man's switch to ensure travel safety
- A really good looking fronted
- .onion site available. A Vanity onion would be nice
- Host on a proper server

### Page AB testing system plan
- A cookie should be set by the backend indicating the page version
- One cookie corresponds to multiple pages, with positional arguments indicating which
- The backend needs to share a page index with nginx
- This page index will be used to determine the page cache.
- a blank 2 chars indicates that a test was not decided.
- On first page access, the cookie will update to a random active test
- There should also be a global test cookie, which will also influence the cache
