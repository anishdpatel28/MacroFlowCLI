# MacroFlow Testing Guide

Comprehensive test suite for MacroFlow CLI.

## Running Tests

### Run All Tests

```bash
cd MacroFlowCLI
go test -v
```

### Run Specific Tests

```bash
# Run tests matching a pattern
go test -v -run TestProjectCreation

# Run tests for a specific function
go test -v -run TestMacro
```

### Run with Coverage

```bash
# Generate coverage report
go test -cover

# Generate detailed coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Benchmarks

```bash
# Run all benchmarks
go test -bench=.

# Run specific benchmark
go test -bench=BenchmarkProjectCreation
```

## Test Coverage

The test suite covers:

### ✅ Project Management (8 tests)
- `TestProjectCreation` - Creating new projects
- `TestDuplicateProject` - Preventing duplicate projects
- `TestGetProjects` - Retrieving all projects
- `TestProjectByPath` - Finding projects by directory path
- `TestDeleteProject` - Deleting projects and associated macros
- `TestInvalidProjectOperations` - Error handling for invalid operations
- `TestEmptyStorage` - Operations on empty storage
- `TestStoragePersistence` - Data persistence across instances

### ✅ Macro Management (8 tests)
- `TestMacroCreation` - Creating new macros
- `TestDuplicateMacro` - Preventing duplicate macros
- `TestGetMacrosByProject` - Retrieving macros for a project
- `TestGetMacroByName` - Finding macros by name
- `TestGetAllMacros` - Retrieving all macros
- `TestDeleteMacro` - Deleting single macros
- `TestDeleteMacros` - Deleting multiple macros
- `TestInvalidMacroOperations` - Error handling for invalid operations

### ✅ Import/Export (3 tests)
- `TestExportData` - Exporting all data to JSON
- `TestImportData` - Importing data (replace mode)
- `TestImportDataMerge` - Importing data (merge mode)

### ✅ Data Integrity (2 tests)
- `TestJSONMarshaling` - JSON serialization/deserialization
- `TestStoragePersistence` - Data persistence verification

### ✅ Performance (3 benchmarks)
- `BenchmarkProjectCreation` - Project creation performance
- `BenchmarkMacroCreation` - Macro creation performance
- `BenchmarkProjectLookup` - Path lookup performance

## Test Statistics

```bash
# Count tests
grep -c "^func Test" macroflow_test.go
# Output: 21 tests

# Count benchmarks
grep -c "^func Benchmark" macroflow_test.go
# Output: 3 benchmarks

# Count lines of test code
wc -l macroflow_test.go
# Output: ~700+ lines
```

## Expected Output

When running all tests, you should see:

```
=== RUN   TestProjectCreation
--- PASS: TestProjectCreation (0.00s)
=== RUN   TestDuplicateProject
--- PASS: TestDuplicateProject (0.00s)
=== RUN   TestGetProjects
--- PASS: TestGetProjects (0.00s)
...
PASS
ok      github.com/anishdpatel28/macroflow    0.XXXs
```

## Coverage Report

Target: **>80% coverage**

To see what's covered:

```bash
go test -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Continuous Testing

### Watch Mode (with external tool)

```bash
# Install gotestsum
go install gotest.tools/gotestsum@latest

# Run tests on file change
gotestsum --watch
```

### Pre-commit Testing

Add to `.git/hooks/pre-commit`:

```bash
#!/bin/sh
cd MacroFlowCLI
go test ./... || exit 1
```

## Test Organization

Tests are organized by functionality:

```
macroflow_test.go
├── Setup/Cleanup (setupTest)
├── Project Tests
│   ├── Creation & Validation
│   ├── Retrieval & Lookup
│   └── Deletion
├── Macro Tests
│   ├── Creation & Validation
│   ├── Retrieval & Lookup
│   └── Deletion
├── Import/Export Tests
├── Data Integrity Tests
└── Performance Benchmarks
```

## Adding New Tests

Template for new tests:

```go
func TestNewFeature(t *testing.T) {
    cleanup := setupTest(t)
    defer cleanup()

    // Test setup
    project, _ := testStorage.CreateProject("test", "/tmp/test")

    // Test execution
    result := yourNewFeature(project)

    // Assertions
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

## Debugging Tests

### Verbose Output

```bash
go test -v -run TestSpecificTest
```

### Print Debug Info

Add to tests:

```go
t.Logf("Debug info: %+v", variable)
```

### Run Single Test

```bash
go test -v -run "^TestProjectCreation$"
```

## Best Practices

1. **Isolation**: Each test is isolated with `setupTest()`
2. **Cleanup**: Always defer cleanup functions
3. **Clear Names**: Test names describe what they test
4. **Fast**: Tests run in milliseconds
5. **Independent**: Tests don't depend on each other
6. **Comprehensive**: Test both success and failure cases

## Troubleshooting

### Tests Fail Locally

```bash
# Clean and retry
go clean -testcache
go test -v
```

### Permission Issues

```bash
# Ensure temp directory is writable
ls -la /tmp
```

### Import Errors

```bash
# Update dependencies
go mod tidy
go test -v
```

## Future Enhancements

Potential additions:
- [ ] Integration tests for CLI commands
- [ ] End-to-end tests with actual shell execution
- [ ] Parallel test execution
- [ ] Stress tests for concurrent operations
- [ ] Fuzz testing for input validation
- [ ] Mock filesystem tests

---

**Coverage Goal**: 80%+ code coverage  
**Test Speed**: All tests < 1 second  
**Maintenance**: Update tests with new features
