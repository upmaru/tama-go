# GitHub CI/CD Configuration for Tama-Go

This directory contains the GitHub Actions workflows and configuration for the Tama-Go client library. Our CI/CD pipeline ensures code quality, security, and reliable releases.

## 🚀 Workflows

### Main CI Workflow (`.github/workflows/ci.yml`)

The primary CI workflow runs on every push and pull request to `main` and `develop` branches.

**Jobs included:**
- **Test**: Runs tests across Go versions 1.21, 1.22, and 1.23
- **Lint**: Code quality checks using golangci-lint
- **Security**: Security scanning with Gosec
- **Mod Tidy**: Ensures go.mod and go.sum are up to date
- **Vulnerability Check**: Scans for known vulnerabilities using govulncheck
- **Formatting**: Verifies code formatting with gofmt
- **Examples**: Ensures example code builds successfully

### Release Workflow (`.github/workflows/release.yml`)

Automatically creates releases when version tags are pushed.

**Features:**
- Validates semantic versioning
- Generates changelog from commit history
- Creates GitHub releases
- Tests installation process
- Supports pre-release tags (alpha, beta, rc)

**Usage:**
```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### Coverage Workflow (`.github/workflows/coverage.yml`)

Comprehensive code coverage reporting and analysis.

**Features:**
- Generates detailed coverage reports
- Uploads to Codecov and Codacy
- Comments coverage on pull requests
- Compares coverage between branches
- Enforces minimum coverage threshold (70%)

## 🔧 Configuration Files

### `.golangci.yml`

Comprehensive linting configuration with 50+ enabled linters:

**Key features:**
- Performance optimization checks
- Security vulnerability detection
- Code style enforcement
- Import organization
- Variable naming conventions
- Error handling best practices

**Disabled for tests:**
- `gocyclo` - Complex test logic is acceptable
- `funlen` - Long test functions are common
- `gomnd` - Magic numbers in tests are OK

### `.github/dependabot.yml`

Automated dependency management:

**Schedules:**
- Go modules: Weekly (Mondays 09:00 UTC)
- GitHub Actions: Weekly (Mondays 10:00 UTC)
- Docker: Weekly (Tuesdays 09:00 UTC)

**Features:**
- Groups related dependencies
- Targets `develop` branch
- Automatic labeling
- Security update prioritization

### `.github/CODEOWNERS`

Defines code review requirements:

**Teams:**
- `@upmaru/maintainers` - Global owners
- `@upmaru/core-team` - Core library files
- `@upmaru/neural-team` - Neural package
- `@upmaru/sensory-team` - Sensory package
- `@upmaru/docs-team` - Documentation
- `@upmaru/devops-team` - CI/CD configuration

## 📋 Issue Templates

### Bug Report (`.github/ISSUE_TEMPLATE/bug_report.md`)
- Environment details
- Reproduction steps
- Expected vs actual behavior
- Code samples

### Feature Request (`.github/ISSUE_TEMPLATE/feature_request.md`)
- Use case description
- API design suggestions
- Priority assessment
- Acceptance criteria

### Documentation (`.github/ISSUE_TEMPLATE/documentation.md`)
- Documentation type
- Location specification
- Improvement suggestions
- Code examples

## 📝 Pull Request Template

Comprehensive PR template covering:
- Change description and type
- Testing requirements
- Documentation updates
- Security considerations
- Performance impact
- Breaking changes

## 🔒 Security

### Security Policy (`SECURITY.md`)
- Vulnerability reporting process
- Supported versions
- Security best practices
- Response timeline

### Security Measures
- Gosec security scanning
- Dependency vulnerability checks
- Secure coding guidelines
- No hardcoded secrets

## 🛠️ Local Development

### Prerequisites
```bash
# Install development tools
make install-tools
```

### Common Commands
```bash
# Run all checks locally (mimics CI)
make ci-check

# Individual checks
make test              # Run tests
make lint              # Run linter
make security-scan     # Security analysis
make vulnerability-check # Check for vulnerabilities
make test-coverage     # Generate coverage report

# Quick checks
make check             # Format, lint, test
make release-check     # All checks for release
```

### Pre-commit Checks

Before pushing code, run:
```bash
make ci-check
```

This runs the same checks as CI and catches issues early.

## 📊 Quality Gates

### Coverage Requirements
- **Minimum**: 70% (enforced)
- **Target**: 80%+ (recommended)
- **Trend**: Coverage should not decrease

### Security Requirements
- Zero high/critical security issues
- No known vulnerabilities
- All dependencies up to date

### Code Quality
- All linting rules pass
- Proper error handling
- No code duplication
- Clear documentation

## 🔄 Workflow Triggers

| Workflow | Trigger | Branches |
|----------|---------|----------|
| CI | Push, PR | main, develop |
| Release | Tag push | v*.*.* |
| Coverage | Push, PR | main, develop |

## 🏷️ Labels

Automatically applied labels:
- `dependencies` - Dependency updates
- `bug` - Bug reports
- `enhancement` - Feature requests
- `documentation` - Doc changes
- `security` - Security issues
- `ci` - CI/CD changes

## 📈 Monitoring

### Metrics Tracked
- Test execution time
- Coverage percentage
- Security scan results
- Dependency freshness
- Build success rate

### Notifications
- Slack/email on build failures
- Security alerts for vulnerabilities
- Coverage reports on PRs

## 🚨 Troubleshooting

### Common Issues

**Lint failures:**
```bash
make lint
golangci-lint run --fix
```

**Test failures:**
```bash
make test ARGS="-v -run TestSpecific"
```

**Coverage too low:**
```bash
make test-coverage
# Check coverage.html for details
```

**Security issues:**
```bash
make security-scan
# Review gosec-results.sarif
```

### Getting Help

1. Check the logs in GitHub Actions
2. Run the same command locally
3. Review the specific tool's documentation
4. Ask in the team Slack channel

## 🔮 Future Enhancements

Planned improvements:
- [ ] Integration test environment
- [ ] Performance benchmarking
- [ ] Multi-architecture builds
- [ ] Automated changelog generation
- [ ] Release notes automation
- [ ] SonarQube integration
- [ ] Docker container scanning

---

**Maintained by**: @upmaru/devops-team  
**Last Updated**: December 2024  
**Questions?**: Open an issue or contact the team