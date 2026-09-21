[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# Open Cluster Management API

## Communication

- Be concise. Skip preambles and postambles.
- Comments explain why, not what.
- Errors should be actionable and specific.

## Constraints

- DCO is required. Sign commits with `git commit -s`.
- Generated deepcopy, client, informer, lister, and CRD files are read-only. Change API types, then run `make update`.
- Preserve compatibility for served API versions and follow Kubernetes API conventions.
- Source files require the project copyright header. Use `make update` to add or verify headers.
- Follow `CONTRIBUTING.md` for AI-assisted contribution and disclosure requirements.

## Commands

Run `make help` for available targets. Common workflows:

```text
make test       # Unit tests
make update     # Regenerate code, CRDs, and copyright headers
make verify     # Run code generation, CRD, copyright, and lint checks
```

## Style

- Keep package names lowercase and aligned with their directory names.
- Keep JSON and `+kubebuilder` tags compatible with the existing API schema.
