---
name: protobuf-contracts
description: Define, review, and regenerate Protobuf schemas and their Go gRPC stubs.
---

## protobuf-contracts

### Conventions

- Package: `fontis.<service>.v1`. File: snake_case matching service name.
- Every message and field has a comment explaining its purpose.
- Field numbers start at 1. Use `reserved` for deleted fields.
- Enums start at 0 with `UNSPECIFIED` value.
- Breaking changes (field removal, type change, number reassignment) require new package version.
- Non-breaking changes (new field, new message, new enum value) allowed within a version.

### Workflow

1. Edit `.proto` files in `contracts/protobuf/<service>/v1/`.
2. Run `make protoc-gen` to regenerate Go stubs into `runtime/<service>/proto/`.
3. CI verifies generated code matches source: `make protoc-gen && git diff --exit-code`.
4. Never manually edit generated `*.pb.go` files.

### Adding a new service contract

1. Create `contracts/protobuf/<service>/v1/<service>.proto`.
2. Define service, RPCs, messages. Follow the existing pattern in identity or auth.
3. Run `make protoc-gen`. Verify stubs appear in `runtime/<service>/proto/`.
4. Update `runtime/<service>/go.mod` if new dependencies are needed.
