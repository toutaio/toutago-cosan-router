# Community Guidelines

Welcome to the Cosan Router community! This document outlines our expectations and guidelines for participation.

## Our Pledge

We are committed to providing a welcoming and inspiring community for all. We pledge to make participation in our project and community a harassment-free experience for everyone, regardless of:

- Age
- Body size
- Disability
- Ethnicity
- Gender identity and expression
- Level of experience
- Nationality
- Personal appearance
- Race
- Religion
- Sexual identity and orientation

## Communication Channels

### GitHub Issues
- **Bug Reports**: Use for reproducible bugs
- **Feature Requests**: Propose new features
- **Questions**: Ask for help (also see Discussions)

### GitHub Discussions
- **General Discussion**: Ideas, feedback, showcases
- **Q&A**: Community help and support
- **Announcements**: Project updates

### Pull Requests
- **Code Contributions**: Follow CONTRIBUTING.md
- **Documentation**: Improvements and fixes
- **Examples**: New use cases

## Getting Started

### New Contributors

1. **Read the Documentation**
   - README.md for overview
   - CONTRIBUTING.md for guidelines
   - CODE_OF_CONDUCT.md for standards

2. **Explore the Codebase**
   - Check out examples/
   - Read the ADRs in docs/adr/
   - Review test files for patterns

3. **Find an Issue**
   - Look for "good first issue" label
   - Check "help wanted" issues
   - Ask in Discussions if unsure

### Asking Questions

**Before asking:**
- Search existing issues and discussions
- Check documentation and examples
- Review troubleshooting guide

**When asking:**
- Provide context (what you're trying to do)
- Include minimal code sample
- Specify environment (Go version, OS, etc.)
- Show what you've tried

### Reporting Bugs

Use the bug report template. Include:
- Clear description
- Steps to reproduce
- Expected vs actual behavior
- Minimal code sample
- Environment details
- Error messages/stack traces

### Suggesting Features

Use the feature request template. Include:
- Problem statement (what needs solving)
- Proposed solution
- Alternative approaches considered
- Example usage
- Willingness to help implement

## Code Contributions

### Before Starting

1. **Check existing issues/PRs** to avoid duplicate work
2. **Discuss major changes** in an issue first
3. **Follow SOLID principles** - this is a reference implementation
4. **Maintain test coverage** - aim for >90%

### Contribution Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests (maintain >90% coverage)
5. Update documentation
6. Run tests (`go test ./...`)
7. Run race detector (`go test -race ./...`)
8. Commit with clear messages
9. Push to your fork
10. Open a Pull Request

### Code Standards

- **Follow Go conventions**: golangci-lint, go fmt
- **Write tests**: Unit, integration, and examples
- **Document code**: Godoc comments on public APIs
- **Keep it SOLID**: Follow the project's architectural principles
- **No external dependencies**: Use stdlib only (core router)

### Pull Request Guidelines

- Use the PR template
- Link related issues
- Describe changes clearly
- Include tests
- Update CHANGELOG.md
- Keep PRs focused (one feature/fix per PR)
- Respond to review feedback promptly

## Code Review

### For Contributors

- Be patient - maintainers review when available
- Be responsive to feedback
- Ask for clarification if needed
- Be open to suggestions
- Learn from the process

### For Reviewers

- Be constructive and respectful
- Explain reasoning behind suggestions
- Recognize good work
- Focus on code, not the person
- Help mentormentoring mindset

## Recognition

We value all contributions:

- **Code**: Features, fixes, refactoring
- **Documentation**: Guides, examples, typo fixes
- **Testing**: New tests, test improvements
- **Bug Reports**: Detailed, reproducible reports
- **Discussions**: Helping others, sharing knowledge
- **Reviews**: Constructive feedback on PRs

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes (for significant contributions)
- Thanked in announcements

## Project Maintenance

### Triaging Issues

Maintainers will:
- Label issues appropriately
- Ask for more information if needed
- Close duplicates/invalid issues
- Prioritize based on impact and effort

### Release Schedule

- **Patch releases**: As needed for bugs
- **Minor releases**: Monthly for features
- **Major releases**: When breaking changes needed

### Response Times

We aim for:
- **Critical bugs**: 24-48 hours
- **Other issues**: Within 1 week
- **Pull requests**: Initial review within 1 week

## Getting Help

### Documentation
1. README.md
2. examples/
3. docs/guides/
4. API documentation (pkg.go.dev)

### Community
1. GitHub Discussions (preferred for questions)
2. GitHub Issues (for bugs/features)
3. Stack Overflow (tag: cosan-router)

### Security Issues

**Do not** open public issues for security vulnerabilities.

Use GitHub Security Advisories:
https://github.com/toutaio/toutago-cosan-router/security/advisories/new

## Governance

### Decision Making

- **Minor changes**: Maintainer discretion
- **Features**: Discussion in issues
- **Breaking changes**: Broad community input
- **Architecture**: ADR process

### Maintainers

Current maintainers are listed in MAINTAINERS.md

Maintainers have:
- Commit access
- Issue/PR triage
- Release authority
- Direction setting

## Growing the Community

### Ways to Help

- **Answer questions** in Discussions
- **Review pull requests** 
- **Write blog posts** about Cosan
- **Give talks** at meetups/conferences
- **Create tutorials** and guides
- **Report bugs** you find
- **Suggest improvements**

### Community Goals

- Foster learning about SOLID principles
- Build production-ready tools
- Support each other
- Share knowledge
- Have fun!

## Updates

These guidelines may evolve. Check back periodically for updates.

Last updated: 2025-12-30

---

**Thank you for being part of the Cosan community!** 🎉

Questions? Start a discussion: https://github.com/toutaio/toutago-cosan-router/discussions
