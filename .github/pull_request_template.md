## Description

Please provide a brief description of the changes in this pull request.

## Type of Change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Code refactoring
- [ ] Test improvement
- [ ] CI/CD improvement
- [ ] Dependency update

## Related Issues

Closes #(issue number)
Related to #(issue number)

## Changes Made

Please describe the changes made in this PR:

- [ ] Added/modified functionality in neural package
- [ ] Added/modified functionality in sensory package
- [ ] Updated client configuration
- [ ] Added/updated tests
- [ ] Updated documentation
- [ ] Updated examples
- [ ] Other: ___________

## Breaking Changes

If this PR introduces breaking changes, please describe them here and provide migration instructions:

```go
// Before
oldMethod()

// After
newMethod()
```

## Testing

- [ ] All existing tests pass
- [ ] New tests have been added for new functionality
- [ ] Tests cover edge cases
- [ ] Examples have been tested
- [ ] Manual testing performed

### Test Coverage

Please describe what you tested:

- [ ] Unit tests
- [ ] Integration tests
- [ ] Example code
- [ ] Error handling
- [ ] Edge cases

## Documentation

- [ ] Code comments updated
- [ ] README.md updated (if applicable)
- [ ] API_REFERENCE.md updated (if applicable)
- [ ] Examples updated (if applicable)
- [ ] Godoc comments added/updated

## Code Quality

- [ ] Code follows project style guidelines
- [ ] No new lint warnings introduced
- [ ] Code is properly formatted (gofmt)
- [ ] Dependencies are up to date
- [ ] No security vulnerabilities introduced

## Checklist

- [ ] I have read the contributing guidelines
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

## Screenshots/Examples (if applicable)

If your changes include UI changes or new functionality, please provide examples:

```go
// Example usage of new feature
config := tama.Config{
    BaseURL: "https://api.tama.io",
    APIKey:  "your-api-key",
}
client := tama.NewClient(config)

// Your new feature usage
result, err := client.NewFeature()
```

## Performance Impact

- [ ] No performance impact
- [ ] Performance improved
- [ ] Performance may be affected (please explain)

If performance is affected, please provide before/after metrics or explain the impact.

## Security Considerations

- [ ] No security impact
- [ ] Security improved
- [ ] Potential security impact (please explain)

Please describe any security considerations or improvements.

## Additional Notes

Any additional information that reviewers should know:

## Reviewer Notes

@reviewers please pay special attention to:

- [ ] API design changes
- [ ] Breaking changes
- [ ] Security implications
- [ ] Performance impact
- [ ] Test coverage
- [ ] Documentation accuracy