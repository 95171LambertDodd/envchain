# envchain

> A tool for composing and validating layered environment configs across dev, staging, and prod.

---

## Installation

```bash
go install github.com/yourorg/envchain@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/envchain.git && cd envchain && go build ./...
```

---

## Usage

Define your environment layers in a `envchain.yaml` file:

```yaml
base: .env.base
layers:
  - .env.staging
  - .env.local
validate:
  - DATABASE_URL
  - API_KEY
  - PORT
```

Then run:

```bash
envchain resolve --env staging
```

This merges each layer in order, with later layers taking precedence, and validates that all required variables are present before outputting the final config.

Export directly to your shell:

```bash
eval $(envchain resolve --env prod --export)
```

Check for missing or conflicting variables without applying:

```bash
envchain validate --env prod
```

---

## How It Works

1. **Base layer** provides shared defaults across all environments.
2. **Named layers** (dev, staging, prod) override base values as needed.
3. **Local overrides** allow per-machine customization without affecting version control.
4. **Validation** ensures required variables are defined before your app starts.

---

## License

MIT © yourorg