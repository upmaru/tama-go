# Security Policy

## Supported Versions

We actively support the following versions of tama-go with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| 0.x.x   | :x:                |

## Reporting a Vulnerability

We take the security of tama-go seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to: security@upmaru.com

You should receive a response within 48 hours. If for some reason you do not, please follow up via email to ensure we received your original message.

### What to Include

Please include the following information in your report:

- Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit the issue

This information will help us triage your report more quickly.

### What to Expect

After submitting a report, you can expect:

1. **Acknowledgment**: We'll acknowledge receipt of your vulnerability report within 48 hours
2. **Initial Assessment**: We'll provide an initial assessment within 72 hours
3. **Regular Updates**: We'll keep you informed of our progress throughout the investigation
4. **Resolution Timeline**: We aim to resolve critical vulnerabilities within 7 days, high severity within 14 days, and others within 30 days

### Security Measures

The tama-go library implements several security measures:

#### Input Validation
- All API inputs are validated before processing
- Proper error handling prevents information leakage
- Rate limiting and timeout configurations help prevent abuse

#### Authentication & Authorization
- Secure API key management
- No hardcoded credentials in the codebase
- Proper handling of authentication tokens

#### Transport Security
- All API communications use HTTPS
- Certificate validation is enforced
- No sensitive data is logged

#### Dependency Management
- Regular dependency updates via Dependabot
- Automated vulnerability scanning with govulncheck
- Security-focused code review process

### Security Best Practices for Users

When using tama-go, please follow these security best practices:

#### API Key Management
```go
// ✅ Good - Use environment variables
apiKey := os.Getenv("TAMA_API_KEY")
if apiKey == "" {
    log.Fatal("TAMA_API_KEY environment variable is required")
}

// ❌ Bad - Never hardcode API keys
apiKey := "your-api-key-here" // Don't do this!
```

#### Configuration Security
```go
// ✅ Good - Use reasonable timeouts
config := tama.Config{
    BaseURL: "https://api.tama.io",
    APIKey:  os.Getenv("TAMA_API_KEY"),
    Timeout: 30 * time.Second, // Reasonable timeout
}

// ❌ Bad - Avoid extremely long timeouts
config := tama.Config{
    Timeout: 0, // No timeout can lead to hanging connections
}
```

#### Error Handling
```go
// ✅ Good - Handle errors without exposing sensitive info
result, err := client.Neural.GetSpace(spaceID)
if err != nil {
    log.Printf("Failed to get space: %v", err)
    return fmt.Errorf("operation failed")
}

// ❌ Bad - Don't expose internal errors to end users
if err != nil {
    return err // May contain sensitive information
}
```

### Vulnerability Disclosure Timeline

Our typical vulnerability disclosure timeline:

1. **Day 0**: Vulnerability reported
2. **Day 1-2**: Acknowledgment sent, initial triage
3. **Day 3-5**: Detailed investigation and fix development
4. **Day 5-7**: Testing and validation of fix
5. **Day 7-10**: Security release preparation
6. **Day 10-14**: Public disclosure (coordinated with reporter)

### Security Contact

For security-related questions or concerns, please contact:

- **Email**: security@upmaru.com
- **Response Time**: Within 48 hours
- **Encryption**: PGP key available upon request

### Acknowledgments

We appreciate the security research community's efforts to improve the security of open source software. Security researchers who responsibly disclose vulnerabilities will be acknowledged in our security advisories (unless they prefer to remain anonymous).

### Legal

This security policy is subject to our [Terms of Service](https://upmaru.com/terms) and [Privacy Policy](https://upmaru.com/privacy).

---

**Last Updated**: December 2024
**Version**: 1.0