# CLAUDE.md - AI Agent Instructions for AUI Project

## 🚀 Session Management

### Starting a Session

When given the instruction **"start session"**, execute these steps in order:

1. **Read Core Documentation**
   ```bash
   # Read these files in sequence:
   1. CLAUDE.md (this file)
   2. IMPLEMENTATION_PLAN.md
   3. All files in sessions/ directory (if exists)
   ```

2. **Analyze Current State**
   - Run tests: `go test ./...`
   - Check git status: `git status`
   - Review current branch: `git branch --show-current`
   - Identify which phase from IMPLEMENTATION_PLAN.md we're in

3. **Report Status**
   Present a concise status report:
   ```
   📊 Project Status:
   - Current Phase: [Phase X from IMPLEMENTATION_PLAN.md]
   - Branch: [current branch name]
   - Tests: [passing/failing status]
   - Last Session: [summary of previous session if available]
   
   📝 Suggested Next Steps:
   1. [Most logical next task from plan]
   2. [Alternative task if blocked]
   
   ⏸️ Awaiting instructions...
   ```

4. **Wait for Instructions**
   Do not proceed until receiving explicit instructions from the user.

### Ending a Session

When given the instruction **"end session"**, execute these steps:

1. **Ensure Clean State**
   ```bash
   go test ./...           # Run all tests
   go fmt ./...           # Format code
   go mod tidy            # Clean dependencies
   git status             # Check for uncommitted changes
   ```

2. **Create Session Summary**
   Create a file `sessions/YYYY-MM-DD-HH-MM-summary.md` with:
   ```markdown
   # Session Summary: [Date Time]
   
   ## Overview
   [Brief description of session goals and outcomes]
   
   ## Completed Tasks
   - [ ] Task 1 with commit hash
   - [ ] Task 2 with commit hash
   
   ## Changes Made
   ### Files Created
   - file1.go: [purpose]
   
   ### Files Modified  
   - file2.go: [what changed]
   
   ## Tests Added/Modified
   - TestName: [what it tests]
   
   ## Current State
   - Phase: [current phase from plan]
   - Branch: [branch name]
   - Tests: [pass/fail count]
   - Coverage: [percentage]
   
   ## Next Session Recommendations
   1. [Specific next task]
   2. [Any blockers or considerations]
   
   ## Commit History
   ```
   [List of commits from this session]
   ```
   ```

3. **Commit the Summary**
   ```bash
   git add sessions/
   git commit -m "docs: add session summary for [date]"
   ```

## 🧪 Development Workflow

### Test-Driven Development (TDD)

**MANDATORY**: Always follow TDD strictly:

1. **Red Phase** - Write failing test first
   ```go
   // 1. Write test that describes desired behavior
   func TestNewFeature(t *testing.T) {
       // Test implementation
   }
   
   // 2. Run test to see it fail
   go test ./package/... -v
   ```

2. **Green Phase** - Write minimal code to pass
   ```go
   // Write ONLY enough code to make test pass
   // No extra features or optimizations
   ```

3. **Refactor Phase** - Improve code quality
   ```go
   // Only after test passes, refactor for:
   // - Clarity
   // - Performance  
   // - Maintainability
   ```

4. **Verify** - Ensure all tests still pass
   ```bash
   go test ./... -v
   ```

### Git Workflow

#### Branch Strategy

```bash
# Create feature branch from main
git checkout main
git pull origin main
git checkout -b feature/short-description

# Keep branches small and focused
# One feature/fix per branch
# Merge frequently to avoid conflicts
```

#### Commit Standards

**USE CONVENTIONAL COMMITS**: [type]([scope]): [description]

Types:
- `feat`: New feature
- `fix`: Bug fix
- `test`: Adding tests
- `refactor`: Code change that neither fixes nor adds feature
- `docs`: Documentation only
- `style`: Formatting, missing semicolons, etc
- `perf`: Performance improvement
- `chore`: Maintenance tasks

Examples:
```bash
git commit -m "feat(agent): add streaming response handler"
git commit -m "test(context): add file deduplication tests"
git commit -m "fix(ui): correct tab switching wraparound"
git commit -m "docs: update implementation plan phase 2"
```

#### Pre-Commit Checklist

**ALWAYS** run before committing:

```bash
# 1. Format code
go fmt ./...

# 2. Clean dependencies
go mod tidy

# 3. Run linter (when configured)
# golangci-lint run

# 4. Run all tests
go test ./... -v

# 5. Check test coverage
go test ./... -cover

# 6. Verify no sensitive data
git diff --staged  # Review changes
```

#### Commit Frequency

- Commit after each test passes (TDD cycle)
- Commit when switching context
- Commit before attempting risky changes
- Small, atomic commits > large commits

## 📁 Project Structure Guidelines

### File Organization

```
aui/
├── cmd/           # Application entry points
├── internal/      # Private application code
│   ├── [domain]/  # Domain models and business logic
│   └── ui/        # TUI components
├── pkg/           # Public packages (APIs)
├── sessions/      # Session summaries
├── docs/          # User documentation
└── testdata/      # Test fixtures
```

### Testing Files

- Test files alongside implementation: `file.go` → `file_test.go`
- Test package same as implementation package
- Use table-driven tests when appropriate
- Mock external dependencies

## 🎯 Phase Tracking

Track progress against IMPLEMENTATION_PLAN.md:

### Current Phase Checklist

When working on a phase:

1. **Start Phase**
   - [ ] Create phase branch: `git checkout -b phase-N-description`
   - [ ] Review phase requirements in IMPLEMENTATION_PLAN.md
   - [ ] Create tracking issue/comment with phase checklist

2. **During Phase**
   - [ ] Follow TDD for each feature
   - [ ] Commit frequently with conventional commits
   - [ ] Update tests for new functionality
   - [ ] Document any deviations from plan

3. **Complete Phase**
   - [ ] All "Definition of Done" items checked
   - [ ] All tests passing
   - [ ] Documentation updated
   - [ ] Create PR to main branch
   - [ ] Update session summary with phase completion

## 🔧 Common Commands Reference

```bash
# Testing
go test ./...                     # Run all tests
go test ./internal/agent -v       # Run specific package tests verbose
go test -run TestName ./...        # Run specific test
go test ./... -cover              # Check coverage
go test -race ./...               # Race condition detection

# Building
go build -o aui cmd/aui/main.go   # Build binary
go install ./cmd/aui              # Install to GOPATH

# Dependencies
go mod tidy                       # Clean up dependencies
go mod download                   # Download dependencies
go get package@version            # Add specific version

# Git
git status                        # Check status
git diff                          # View changes
git log --oneline -10            # Recent commits
git stash                        # Temporarily store changes
git stash pop                    # Restore stashed changes

# Code Quality
go fmt ./...                      # Format code
go vet ./...                      # Examine code
golangci-lint run                # Comprehensive linting (when installed)
```

## ⚠️ Important Reminders

1. **ALWAYS use TDD** - No production code without failing test first
2. **Small, frequent commits** - Easier to review and revert if needed
3. **Run tests before committing** - Never break the build
4. **Document decisions** - Add comments for non-obvious choices
5. **Update session summaries** - Future agents need context
6. **Follow the plan** - IMPLEMENTATION_PLAN.md is the roadmap
7. **Ask if uncertain** - Better to clarify than assume

## 🎨 Code Style Guidelines

- Use descriptive variable names
- Keep functions small and focused
- Return early from functions when possible
- Handle errors explicitly
- Use interfaces for flexibility
- Prefer composition over inheritance
- Document exported functions
- Keep line length reasonable (~100 chars)

## 📋 Session Task Patterns

### Adding a New Feature

1. Read requirements from IMPLEMENTATION_PLAN.md
2. Write acceptance test (high-level)
3. Break into smaller unit tests
4. Implement feature using TDD
5. Refactor for clarity
6. Update documentation
7. Commit with descriptive message

### Fixing a Bug

1. Write test that reproduces bug
2. Verify test fails
3. Fix the bug
4. Verify test passes
5. Check for similar issues
6. Commit with "fix:" prefix

### Refactoring

1. Ensure comprehensive test coverage exists
2. Run tests (green)
3. Make incremental changes
4. Run tests after each change
5. Commit working states frequently
6. Final commit with "refactor:" prefix

---

*This document is the primary reference for AI agents working on the AUI project. It should be read at the start of every session and updated when workflows change.*