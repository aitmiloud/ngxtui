# NgxTUI Architecture

This document describes the architecture and organization of the NgxTUI project.

## Project Structure

```
ngxtui/
├── cmd/
│   └── ngxtui/           # Application entry point
│       └── main.go       # Main function, bootstraps the app
├── internal/             # Private application code
│   ├── app/              # Bubble Tea application logic
│   │   ├── init.go       # Model initialization
│   │   ├── update.go     # Update logic (state transitions)
│   │   └── view.go       # View rendering orchestration
│   ├── model/            # Data models and types
│   │   └── types.go      # Model structs, messages, keybindings
│   ├── nginx/            # NGINX service layer
│   │   └── service.go    # NGINX operations (list, enable, disable, etc.)
│   ├── styles/           # UI styling
│   │   └── styles.go     # Lipgloss styles and color palette
│   └── ui/               # UI components and rendering
│       └── views.go      # View rendering functions
├── go.mod                # Go module definition
├── go.sum                # Go dependencies checksums
└── README.md             # Project documentation
```

## Architecture Principles

This project follows the **Model-View-Update (MVU)** pattern, which is the foundation of Bubble Tea applications:

### 1. **Model** (`internal/model/`)
- Defines the application state
- Contains all data structures (Site, Model, etc.)
- Defines message types for state transitions
- Defines keybindings and their behavior

### 2. **Update** (`internal/app/update.go`)
- Handles all state transitions
- Processes messages (keyboard input, ticks, status updates)
- Executes commands and side effects
- Pure function: `(Model, Msg) -> (Model, Cmd)`

### 3. **View** (`internal/app/view.go` + `internal/ui/`)
- Renders the current state to the terminal
- Pure function: `Model -> String`
- Delegates to specialized rendering functions in `ui` package

## Package Responsibilities

### `cmd/ngxtui`
- **Purpose**: Application entry point
- **Responsibilities**:
  - Parse command-line arguments (if any)
  - Check permissions
  - Initialize and run the Bubble Tea program
  - Handle top-level errors

### `internal/app`
- **Purpose**: Bubble Tea application logic
- **Responsibilities**:
  - Initialize the model (`init.go`)
  - Handle state updates (`update.go`)
  - Orchestrate view rendering (`view.go`)
  - Implement the `tea.Model` interface

### `internal/model`
- **Purpose**: Data models and types
- **Responsibilities**:
  - Define application state structure
  - Define message types for communication
  - Define keybindings
  - No business logic, just data structures

### `internal/nginx`
- **Purpose**: NGINX service layer
- **Responsibilities**:
  - Interface with NGINX system
  - List, enable, disable sites
  - Test configuration
  - Reload NGINX
  - Read logs
  - Abstract away system-specific details

### `internal/styles`
- **Purpose**: UI styling
- **Responsibilities**:
  - Define color palette
  - Define reusable lipgloss styles
  - Maintain consistent visual design
  - No rendering logic, just style definitions

### `internal/ui`
- **Purpose**: UI rendering components
- **Responsibilities**:
  - Render specific views (tabs, tables, charts, etc.)
  - Apply styles to content
  - Format data for display
  - No state management, just rendering

## Design Patterns

### 1. **Separation of Concerns**
Each package has a single, well-defined responsibility. This makes the code:
- Easier to understand
- Easier to test
- Easier to modify
- More reusable

### 2. **Dependency Injection**
Services (like `nginx.Service`) are created and passed where needed, making the code:
- More testable (can mock services)
- More flexible (can swap implementations)
- Less coupled

### 3. **Message Passing**
State changes happen through messages, not direct mutation:
- Makes state transitions explicit
- Easier to debug (can log all messages)
- Easier to test (can replay messages)
- Follows Bubble Tea conventions

### 4. **Pure Functions**
View and update functions are pure (no side effects in the function body):
- Predictable behavior
- Easier to test
- Easier to reason about
- Commands handle side effects separately

## Data Flow

```
User Input
    ↓
Update (processes input, returns new model + commands)
    ↓
Commands (side effects: API calls, file I/O, etc.)
    ↓
Messages (results of commands)
    ↓
Update (processes messages, updates state)
    ↓
View (renders current state)
    ↓
Terminal Display
```

## Adding New Features

### Adding a New Tab
1. Add tab type to `model.TabType` in `internal/model/types.go`
2. Add rendering logic in `internal/ui/views.go`
3. Add tab handling in `internal/app/update.go`
4. Update tab bar rendering in `internal/ui/views.go`

### Adding a New NGINX Operation
1. Add method to `nginx.Service` in `internal/nginx/service.go`
2. Add action to menu in `internal/ui/views.go`
3. Add case in `executeAction` in `internal/app/update.go`

### Adding a New Message Type
1. Define message type in `internal/model/types.go`
2. Add handler in `Update` function in `internal/app/update.go`
3. Create command that sends the message (if needed)

## Testing Strategy

### Unit Tests
- Test individual functions in isolation
- Mock dependencies (e.g., `nginx.Service`)
- Test edge cases and error conditions

### Integration Tests
- Test interactions between packages
- Test full update cycles
- Test command execution

### Manual Testing
- Run the application
- Test all user interactions
- Verify visual appearance

## Best Practices

1. **Keep functions small**: Each function should do one thing well
2. **Use meaningful names**: Names should describe what the code does
3. **Document public APIs**: Add comments to exported types and functions
4. **Handle errors**: Always check and handle errors appropriately
5. **Follow Go conventions**: Use `gofmt`, follow Go idioms
6. **Keep state immutable**: Don't modify the model in place, return a new one
7. **Use constants**: Define magic strings and numbers as constants

## References

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://golang.org/doc/effective_go.html)
