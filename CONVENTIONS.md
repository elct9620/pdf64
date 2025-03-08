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
