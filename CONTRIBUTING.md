
<h1 align="center">Contributing To DockMate 🐳</h1>
<h3 align="center">Thanks for your interest in contributing to <b>DockMate</b></h3>

## Table of Contents

- [Ways to Contribute](#ways-to-contribute)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
  - [Running Locally](#running-locally)
  - [Building](#building)
- [Making Changes](#making-changes)
  - [1. Create a Branch](#1-create-a-branch)
  - [2. Make Your Changes](#2-make-your-changes)
  - [3. Test](#3-test)
- [Submitting Changes](#submitting-changes)
- [Code Guidelines](#code-guidelines)
- [Commit Messages](#commit-messages)
- [Getting Help](#getting-help)

## Ways to Contribute
- 🐞 Report bugs
- ⭐ Suggest features
- 📝 Improve documentation
- ❗ Submit bug fixes
- ✨ Add new features

## Getting Started

### Prerequisites
- Go 1.24 or higher
- Docker installed and running
- Git

### Setup
1. Fork the repository
2. Clone your fork:
   ```
   git clone https://github.com/YOUR_USERNAME/dockmate.git
   cd dockmate
   ```
3. Add upstream remote:
   ```
   git remote add upstream https://github.com/shubh-io/dockmate.git
   ```
4. Install dependencies:
   ```
   go mod download
   ```

### Running Locally
```
go run .
```

### Building
```
go build -o dockmate
./dockmate
```

## Making Changes

### 1. Create a Branch

First, sync with the latest changes:
```
git checkout main
git pull upstream main
```

Then create your feature branch:
```
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-compose-support`
- `fix/memory-leak`
- `docs/update-readme`

### 2. Make Your Changes
- Write clean, readable code
- Follow Go conventions
- Add comments where needed
- Test your changes

### 3. Test
Make sure everything still works:
```
go test ./...
go run . # Manual testing
```

## Submitting Changes

1. Commit your changes:
   ```
   git add .
   git commit -m "Add: brief description of changes"
   ```
2. Push to your fork:
   ```
   git push origin feature/your-feature-name
   ```
3. Open a Pull Request on GitHub
4. Describe what you changed and why
5. Link any related issues

## Code Guidelines

- Run `go fmt` before committing
- Keep functions small and focused
- Use meaningful variable names
- Add comments for complex logic
- Update documentation for new features

## Commit Messages

Use clear, descriptive commit messages:
- `Add: feature name` for new features
- `Fix: bug description` for bug fixes
- `Docs: what you updated` for documentation
- `Refactor: what you improved` for code improvements

## Getting Help

- 💬 Open an issue for questions
- 🐞 Report bugs via GitHub Issues
- 🎯 Suggest features in Discussions

**Thanks for contributing :)**
