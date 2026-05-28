# Go vs Python CLI Performance

Measured on macOS ARM64 (Apple Silicon), warm cache. All API calls target the same `open-agent.liepin.com` endpoint.

| Operation | Go CLI | Python CLI | Improvement |
|-----------|--------|------------|-------------|
| Startup (`--help`) | 0.005s | 0.149s | **30x** |
| Cold startup | 0.020s | 0.338s | **17x** |
| `resume get` | 0.180s | 15.9s | **88x** |
| `job search` | 0.340s | 15.8s | **47x** |

Startup speed matters most for CLI tools invoked frequently by AI agents — every invocation pays the interpreter startup cost. A compiled Go binary has near-zero startup overhead.

The Python API call times include module import overhead (Typer, Pydantic, requests). The Go binary is a single statically-linked executable with no runtime resolution.
