# Architecture Decision Records (ADR) Index

This directory contains Architecture Decision Records documenting the key architectural decisions made during the design and implementation of Cosan router.

## What is an ADR?

An Architecture Decision Record (ADR) is a document that captures an important architectural decision made along with its context and consequences.

## ADR Format

Each ADR includes:
- **Status:** Accepted, Proposed, Deprecated, Superseded
- **Date:** When the decision was made
- **Context:** The problem or situation requiring a decision
- **Decision:** The choice made
- **Rationale:** Why this choice was made
- **Trade-offs:** Pros and cons
- **Consequences:** Impact of the decision
- **References:** Related documentation

## Current ADRs

### [ADR 001: Interface-First Design](./001-interface-first-design.md)
**Status:** Accepted  
**Date:** 2025-12-29

Establishes interface-first design as the foundational architectural principle for Cosan. Every major component is defined by an interface rather than a concrete implementation.

**Key Points:**
- All components are mockable and testable
- Supports SOLID principles (especially DIP and ISP)
- Enables pluggable implementations
- Small performance overhead (~2%) for significant testability gains

### [ADR 002: Radix Tree for Route Matching](./002-radix-tree-routing.md)
**Status:** Accepted  
**Date:** 2025-12-29

Chooses radix tree (trie) as the route matching strategy over alternatives like regex or simple maps.

**Key Points:**
- O(k) lookup time (k = path length)
- Supports path parameters (`:id`) and wildcards (`*`)
- Industry-proven algorithm used by Chi, Gin, Echo
- Balanced performance and features
- Memory efficient through prefix sharing

### [ADR 003: Error Handling Strategy](./003-error-handling-strategy.md)
**Status:** Accepted  
**Date:** 2025-12-29

Defines handlers as returning `error` with centralized error processing via middleware.

**Key Points:**
- Idiomatic Go error handling
- Explicit errors in function signatures
- Centralized error processing
- Excellent testability
- Composable middleware for error transformation

### [ADR 004: Optional Ecosystem Integrations](./004-optional-ecosystem-integrations.md)
**Status:** Accepted  
**Date:** 2025-12-29

Establishes how Cosan integrates with Toutā ecosystem components (datamapper, fith-renderer, nasc-dependency-injector) through optional adapter interfaces.

**Key Points:**
- Cosan works standalone with zero Toutā dependencies
- Optional adapters for ecosystem integration
- Users choose what to integrate
- Functional options pattern for configuration
- Clean separation via interfaces

## Decision Process

Architectural decisions should be documented when:

1. The decision has significant impact on the codebase
2. Multiple alternatives were considered
3. The choice has important trade-offs
4. Future maintainers need to understand the rationale
5. The decision might be challenged or reconsidered

## Reading Order

For new contributors, read ADRs in this order:

1. **ADR 001 (Interface-First Design)** - Foundation of everything
2. **ADR 003 (Error Handling)** - Affects all handlers
3. **ADR 002 (Radix Tree Routing)** - Core routing algorithm
4. **ADR 004 (Ecosystem Integrations)** - How Cosan fits in Toutā

## Future ADRs

Potential topics for future ADRs:

- **ADR 005: Context Implementation** - DefaultContext design choices
- **ADR 006: Middleware Chain Execution** - Forward vs reverse chain
- **ADR 007: Route Group Implementation** - How groups share middleware
- **ADR 008: Static File Serving** - Strategy for serving static content
- **ADR 009: Testing Strategy** - Unit vs integration test approach
- **ADR 010: Versioning Strategy** - Semantic versioning and compatibility

## Contributing ADRs

When proposing a new ADR:

1. Create a new file: `XXX-short-title.md`
2. Number sequentially (next available number)
3. Use the template structure (see existing ADRs)
4. Submit as a pull request
5. Update this index

## References

- [ADR GitHub Org](https://adr.github.io/)
- [Documenting Architecture Decisions](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions)
- [Lightweight ADRs](https://www.thoughtworks.com/radar/techniques/lightweight-architecture-decision-records)

## Status

Last Updated: 2025-12-30  
Total ADRs: 4  
Accepted: 4  
Proposed: 0  
Deprecated: 0
