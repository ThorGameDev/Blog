# Blog
This is a Work In Progress Blog engine to host my own personal blog.

## Current features
- Basic Docker setup
- HTML is minified
- Compiled TS and SCSS files
- Go backend
- Docker Build caching
- Simple user accounts
- SQL page management
- Simple Bilingual page translation system

## Planned Features (Roughly in planned order)
- User accounts
- Comment threads
- AB testing system
- A half decent looking fronted
- Selfhost with a domain and HTTPS

## Nice to haves
- Early access subscriptions
- SEO optimizations
- Lynx support (No JS or CSS needed)
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
