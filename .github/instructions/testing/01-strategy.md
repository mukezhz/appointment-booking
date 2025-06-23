# Testing Strategy

## Tools
- Ginkgo: BDD testing
- TestContainers: Integration testing
- Testify: Assertions & mocks

## Test Types
1. Unit Tests (business logic)
2. Integration Tests (DB, API)
3. Concurrency Tests

## Running Tests
```bash
# All tests
go test ./... -v

# Ginkgo tests with coverage
ginkgo -v --cover -r ./domain/... ./pkg/...
```

## Best Practices
1. Use test suites
2. Mock external dependencies
3. Use isolated test databases
4. Test edge cases
