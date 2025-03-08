Coding Guidelines
===

Indentation
---

Use tabs for indentation.

Naming Conventions
---

- Use PascalCase for type names and public members.
- Use camelCase for local variables and method arguments.
- Use whold words in names when possible.

Comments
---

Only use comments for GoDoc documentation. Do not use comments for code explanations, try to write self-explanatory code.

Source Code Organization
---

- `pkg` for public packages.
- `internal` for private packages.
- `cmd` for main applications.

### Public Packages

- `api` defines the API interface.

### Internal Packages

- `controller` to implement the API interface.
- `usecase` to implement the business logic.
- `entity` to define the domain entities.

Testing
---

- Use `testing` package for unit and integration tests.
- Use `httptest` package for HTTP tests.
- Use `jmespath` to verify partial JSON responses in the integration tests.
- Use `_test` as package suffix for test packages excluding the main package.
- Write integration tests for the API.

### Example

```go
package main

import (
    "testing"
    "httptest"
)

func TestApiV1Convert(t *testing.T) {
    // Use table-driven tests
    tests := []struct {
        name string
        input string
        expected string
        error error
    }{
        {
            name: "Test 1",
            input: "Hello, World!",
            expected: "Hello, World!",
            error: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            apiV1 := v1.NewService()
            app := app.NewServer(apiV1)

            srv := httptest.NewServer(app)
            defer srv.Close()

            client := srv.Client()
            req, err := http.NewRequest("POST", srv.URL, bytes.NewBufferString(tt.input))
            if err != nil {
                t.Fatal(err)
            }

            resp, err := client.Do(req)
            if err != nil {
                t.Fatal(err)
            }

            if resp.StatusCode != http.StatusOK {
                t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
            }

            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                t.Fatal(err)
            }

            if string(body) != tt.expected {
                t.Errorf("expected body %q, got %q", tt.expected, string(body))
            }
        })
    }
}
```
