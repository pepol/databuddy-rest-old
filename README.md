# databuddy
DataBuddy is a globally replicated datastore with on-demand cross-region replication features

## Contents

- `main.go` - contains the applications startup code, including cobra command parsing
  and viper configuration parsing
- `internal/` - contains all internal (server-specific) code/packages
- `schemas/` - contains [JSON Schema](https://json-schema.org/) definitions for APIs
- `server/` - contains implementation of server, including the API (depends on schemas)
- `tools/` - contains CI-specific tooling (separate module)
