---
applyTo: "**/*_test.go"
---

## 📋 General Testing Principles

### 1. **Package Naming Convention**
- Always use `<package>_test` package name (external testing)
- Example: `package seatrepo_test`, `package usecase_test`

### 2. **Required Imports**
```go
import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    // For database mocking
    "github.com/DATA-DOG/go-sqlmock" 
    "github.com/jmoiron/sqlx"
    
    // For Redis mocking  
    "github.com/go-redis/redismock/v9"
    
    // For gomock (when using mocks)
    "github.com/golang/mock/gomock"
    
    // Framework errors for error type assertions
    errsFramework "github.com/kittipat1413/go-common/framework/errors"
)
```

## 🏗️ Test Structure Patterns

### 1. **Database Repository Tests**

#### Main Test File (`main_test.go`)
```go
package <repo>_test

import (
    "testing"
    "<project>/internal/domain/repository"
    <repopackage> "<project>/internal/infra/db/repository/<repo>"
    "<project>/pkg/testhelper"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/jmoiron/sqlx"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func initTest(t *testing.T) *testhelper.RepoTestHelper[repository.<RepoInterface>] {
    return testhelper.NewRepoTestHelper(t, func(db *sqlx.DB) repository.<RepoInterface> {
        return <repopackage>.New<RepoName>Repository(db)
    })
}

func TestNew<RepoName>Repository(t *testing.T) {
    db, _, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    mockDB := sqlx.NewDb(db, "sqlmock")

    // Execute
    repo := <repopackage>.New<RepoName>Repository(mockDB)

    // Assert
    assert.NotNil(t, repo)
}

func Test<RepoName>RepositoryImpl_WithTx(t *testing.T) {
    h := initTest(t)
    defer h.Done()

    // Create a mock transaction
    txDB, _, err := sqlmock.New()
    require.NoError(t, err)
    defer txDB.Close()

    transactionDB := sqlx.NewDb(txDB, "sqlmock")

    // Execute
    txRepo := h.Repository.WithTx(transactionDB)

    // Assert
    assert.NotNil(t, txRepo)
    assert.NotEqual(t, h.Repository, txRepo, "WithTx should return a new repository instance")
}
```

#### Method Test Files
```go
func Test<RepoName>RepositoryImpl_<MethodName>(t *testing.T) {
    // Test data setup
    testID := uuid.New()
    expectedEntity := entity.<EntityName>{
        // ... entity fields
    }

    tests := []struct {
        name          string
        setupMock     func(mock sqlmock.Sqlmock)
        expected<Field> <Type>
        expectedError bool
        errorType     error
    }{
        {
            name: "successful <operation>",
            setupMock: func(mock sqlmock.Sqlmock) {
                rows := sqlmock.NewRows([]string{"id", "name", "..."}).
                    AddRow(testID, "test-name", ...)
                mock.ExpectQuery("SELECT ...").WillReturnRows(rows)
            },
            expected<Field>: &expectedEntity,
            expectedError:  false,
        },
        {
            name: "database connection error",
            setupMock: func(mock sqlmock.Sqlmock) {
                mock.ExpectQuery("SELECT ...").WillReturnError(sql.ErrConnDone)
            },
            expected<Field>: nil,
            expectedError:  true,
            errorType:      &errsFramework.DatabaseError{},
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            h := initTest(t)
            defer h.Done()

            tt.setupMock(h.Mock)

            // Execute
            result, err := h.Repository.<MethodName>(context.Background(), testID)

            // Assert
            if tt.expectedError {
                require.Error(t, err)
                assert.Contains(t, err.Error(), "[repository <repo>/<file> <MethodName>]")
                
                if tt.errorType != nil {
                    assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
                }
            } else {
                require.NoError(t, err)
                // Add specific assertions based on return type
            }

            // Verify all expectations were met
            h.AssertExpectationsMet(t)
        })
    }
}
```

### 2. **Redis Repository Tests**

#### For redismock
```go
func TestNew<RepoName>Repository(t *testing.T) {
    client, _ := redismock.NewClientMock()

    // Execute
    repo := <repopackage>.New<RepoName>Repository(client)

    // Assert
    assert.NotNil(t, repo)
}

func Test<RepoName>RepositoryImpl_<MethodName>(t *testing.T) {
    tests := []struct {
        name          string
        setupMock     func(mock redismock.ClientMock)
        expectedOK    bool
        expectedError bool
        errorType     error
    }{
        {
            name: "successful operation",
            setupMock: func(mock redismock.ClientMock) {
                mock.ExpectHGet("key", "field").SetVal("value")
            },
            expectedOK:    true,
            expectedError: false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client, mock := redismock.NewClientMock()
            repository := <repopackage>.New<RepoName>Repository(client)

            tt.setupMock(mock)

            // Execute
            result, err := repository.<MethodName>(context.Background(), args...)

            // Assert
            if tt.expectedError {
                require.Error(t, err)
                assert.Contains(t, err.Error(), "[repository <repo>/<file> <MethodName>]")
                
                if tt.errorType != nil {
                    assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
                }
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

#### For gomock (with external dependencies)
```go
func TestNew<RepoName>Repository(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDependency := <deps>_mocks.NewMock<Dependency>(ctrl)

    // Execute
    repo := <repopackage>.New<RepoName>Repository(mockDependency)

    // Assert
    assert.NotNil(t, repo)
}

func Test<RepoName>RepositoryImpl_<MethodName>(t *testing.T) {
    tests := []struct {
        name              string
        setupMock         func(mock *<deps>_mocks.Mock<Dependency>)
        expectedError     bool
        expectedErrorMsg  string
        expectedErrorType error
    }{
        {
            name: "successful operation",
            setupMock: func(mock *<deps>_mocks.Mock<Dependency>) {
                mock.EXPECT().
                    <Method>(gomock.Any(), gomock.Any()).
                    Return("result", nil)
            },
            expectedError: false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockDependency := <deps>_mocks.NewMock<Dependency>(ctrl)
            tt.setupMock(mockDependency)

            repository := <repopackage>.New<RepoName>Repository(mockDependency)

            // Execute
            err := repository.<MethodName>(context.Background(), args...)

            // Assert
            if tt.expectedError {
                require.Error(t, err)
                assert.Contains(t, err.Error(), "[repository <repo>/<file> <MethodName>]")
                assert.Contains(t, err.Error(), tt.expectedErrorMsg)

                if tt.expectedErrorType != nil {
                    assert.ErrorAs(t, err, &tt.expectedErrorType, "Expected error to be of type %T", tt.expectedErrorType)
                }
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

### 3. **Usecase Tests**

#### Main Test File (`main_test.go`)
```go
package usecase_test

import (
    "testing"
    "time"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"

    <domain>_mocks "<project>/internal/domain/<domain>/mocks"
    <usecasepackage> "<project>/internal/usecase/<usecase>"

    "github.com/kittipat1413/go-common/framework/retry"
)

type testHelper struct {
    ctrl         *gomock.Controller
    retrier      retry.Retrier
    mockRepo     *<domain>_mocks.Mock<RepoInterface>
    usecase      <usecasepackage>.<UsecaseInterface>
}

func initTest(t *testing.T) *testHelper {
    ctrl := gomock.NewController(t)

    // Use real retrier with short retry configuration for tests
    queryBackoff, _ := retry.NewExponentialBackoffStrategy(10*time.Millisecond, 2.0, 100*time.Millisecond)
    retrier, _ := retry.NewRetrier(retry.Config{
        MaxAttempts: 3,
        Backoff:     queryBackoff,
    })

    mockRepo := <domain>_mocks.NewMock<RepoInterface>(ctrl)

    usecase := <usecasepackage>.New<UsecaseName>Usecase(
        retrier,
        mockRepo,
        // ... other dependencies
    )

    return &testHelper{
        ctrl:     ctrl,
        retrier:  retrier,
        mockRepo: mockRepo,
        usecase:  usecase,
    }
}

func (h *testHelper) Done() {
    h.ctrl.Finish()
}

func TestNew<UsecaseName>Usecase(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    queryBackoff, _ := retry.NewExponentialBackoffStrategy(10*time.Millisecond, 2.0, 100*time.Millisecond)
    retrier, _ := retry.NewRetrier(retry.Config{
        MaxAttempts: 3,
        Backoff:     queryBackoff,
    })
    mockRepo := <domain>_mocks.NewMock<RepoInterface>(ctrl)

    // Execute
    usecase := <usecasepackage>.New<UsecaseName>Usecase(retrier, mockRepo)

    // Assert
    assert.NotNil(t, usecase)
}
```

#### Method Test Files
```go
func Test<UsecaseName>Usecase_<MethodName>(t *testing.T) {
    tests := []struct {
        name           string
        setupMocks     func(h *testHelper)
        expectedResult *entity.<EntityName> // or other expected return type
        expectedError  bool
        errorType      error
        errorContains  string
    }{
        {
            name: "successful operation",
            setupMocks: func(h *testHelper) {
                h.mockRepo.EXPECT().
                    <Method>(gomock.Any(), gomock.Any()).
                    Return(expectedResult, nil)
            },
            expectedResult: expectedResult,
            expectedError:  false,
        },
        {
            name: "repository fails - retryable error eventually succeeds",
            setupMocks: func(h *testHelper) {
                // First attempt: database error (retryable)
                h.mockRepo.EXPECT().
                    <Method>(gomock.Any(), gomock.Any()).
                    Return(nil, errsFramework.NewDatabaseError("connection timeout", "timeout"))

                // Second attempt: succeeds
                h.mockRepo.EXPECT().
                    <Method>(gomock.Any(), gomock.Any()).
                    Return(expectedResult, nil)
            },
            expectedResult: expectedResult,
            expectedError:  false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            h := initTest(t)
            defer h.Done()

            // Setup mocks
            tt.setupMocks(h)

            // Execute
            ctx := context.Background()
            result, err := h.usecase.<MethodName>(ctx, args...)

            // Assert
            assert.Equal(t, tt.expectedOk, result)

            if tt.expectedError {
                require.Error(t, err)
                assert.Contains(t, err.Error(), "[usecase <usecase>/<file> <MethodName>]")
                
                if tt.errorContains != "" {
                    assert.Contains(t, err.Error(), tt.errorContains)
                }

                if tt.errorType != nil {
                    assert.ErrorAs(t, err, &tt.errorType, "Expected error to be of type %T", tt.errorType)
                }
            } else {
                assert.NoError(t, err)
                assert.Equal(t, pointer.GetValue(tt.expectedResult), pointer.GetValue(result))
                // ... more assertions
            }
        })
    }
}
```

## 🧪 Test Case Categories

### 1. **Happy Path Tests**
- Successful operations with valid data
- Edge cases with valid boundaries
- Empty but valid results

### 2. **Error Handling Tests**
- Database connection errors (`sql.ErrConnDone`, `context.DeadlineExceeded`)
- Redis connection errors (`connection refused`, `timeout`)
- Domain-specific errors (`NotFoundError`, `ConflictError`)
- Invalid input validation

### 3. **Retry Logic Tests (for usecases)**
- Retryable errors that eventually succeed
- Retryable errors that exhaust retries
- Non-retryable errors that don't retry

### 4. **Special Redis Tests**
- For `HEXPIRE` commands, use `ExpectDo` due to redismock limitations:
```go
mock.ExpectDo("HEXPIRE", key, ttlSeconds, "FIELDS", 1, field).SetVal([]int64{1})
```
- Use `mock.MatchExpectationsInOrder(false)` for unordered expectations
- Use `mock.Regexp().ExpectGet(\`pattern\`)` for dynamic keys

## ✅ Assertion Patterns

### 1. **Error Assertions**
```go
if expectedError {
    require.Error(t, err)
    assert.Contains(t, err.Error(), "[<location> <method>]") // Error prefix check
    
    if errorType != nil {
        assert.ErrorAs(t, err, &errorType, "Expected error to be of type %T", errorType)
    }
    
    if errorContains != "" {
        assert.Contains(t, err.Error(), errorContains)
    }
} else {
    require.NoError(t, err)
}
```

### 2. **Entity Comparisons**
```go
// For time fields, always compare in UTC
assert.Equal(t, expected.CreatedAt.UTC(), actual.CreatedAt.UTC())
assert.Equal(t, expected.UpdatedAt.UTC(), actual.UpdatedAt.UTC())

// For other fields, direct comparison
assert.Equal(t, expected.ID, actual.ID)
assert.Equal(t, expected.Name, actual.Name)
```

### 3. **Mock Expectations**
```go
// Always verify all mock expectations were met (for database tests)
h.AssertExpectationsMet(t)
```

## 🎯 Test Data Best Practices

### 1. **Create Realistic Test Scenarios**
- Use meaningful test names that describe the scenario
- Test both success and failure paths
- Include boundary conditions and edge cases
- Test validation errors with various invalid inputs
- Test different error types (NotFoundError, DatabaseError, etc.)

### 2. **Maintain Test Independence**
- Each test should be able to run independently
- Use fresh mocks for each test case (via `initTest()`)
- Don't rely on test execution order
- Clean up resources with `defer h.Done()`

## 📝 Naming Conventions

- Test files: `<method_name>_test.go`
- Test functions: `Test<StructName>_<MethodName>`
- Table-driven test names: Use descriptive scenario names
- Mock setup functions: `setupMock`, `setupMocks`
- Test helper functions: `initTest`, `testHelper.Done()`