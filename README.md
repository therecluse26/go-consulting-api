# REST API for Consulting/Training Database

A RESTful API for a consulting/training application. Built for a microservice-oriented architecture.

Routes are defined within the `routes/routes.go` file and router methods are stored within subdirectories corresponding to database schemas.

### Current Features
- Auth handled through frontend application via Azure AD with Oauth2 (JWT)
- SQL Server backend
- Negroni for auth middleware
- ACLs controlled with Casbin
- Errors passed to Sentry.io
- Deployed with Docker
- Various caching methods for ACL and Auth keys
    - Memcached
    - Database
    - Local file
    - Local environment

### In Development
- All CRUD operations
- Fully implement access control

### Planned Features
- Additional DB engines
    - MySQL
    - PostgreSQL
    - Sqlite
- User-configurable endpoints
    - YAML or JSON files
    - Ability to bind a route to a sql query or new function
    - User-defined functions 
- Other auth providers
- Stateless configuration from external store

Configuration
-------------

Configuration is currently handled through JSON environment variables (looking into more scalable, central solutions, such as Hashicorp's Vault)

Example: 
```json
{
  "sqlHost": "localhost",
  "sqlPort": 1433,
  "sqlUser": "sa",
  "sqlDB": "Consulting_DB",
  "apiHost": "http://localhost:9988",
  "apiPort": 9988,
  "sqlPass": "Asdfasdf123!@#",
  "allowedOrigins": ["http://{{frontend_host}}", "http://localhost:9988"],
  "SentryHost": "https://{sentry_endpoint}",
  "cacheMethod": "local_file",
  "AuthHost": "https://login.microsoftonline.com/{{azure_id}}/oauth2",   
  "AuthClientId": "{{azure_client_id}}",   
  "AuthSecret": "{{azure_auth_secret}}"
}
```

Config JSON should be encoded into base64 and passed into the Docker container as `CFGJSON`

The a database seed file is included in `/build/SeedTestData.sql`