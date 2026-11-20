# Show this command guide (also the default for plain just).
@help:
    printf '\n%s\n\n' 'Shell Explain — development commands'
    printf '%s\n' 'Commands:'
    '{{ just_executable() }}' --list --unsorted --list-heading '' --list-prefix '  just ' --color never
    printf '\n%s\n' 'Prerequisites: just, Go 1.26+, Node.js 22.23.4+, and npm.'
    printf '%s\n' 'Web setup: npm ci --prefix web'
    printf '%s\n' 'Browser-test setup: cd web && npx playwright install chromium'
    printf '\n%s\n\n' 'Details: docs/development.md. No command publishes or deploys; infra-preview only reads AWS.'

# Build the offline explain CLI with embedded Markdown.
@build:
    printf '\n%s\n%s\n' 'Building explain CLI...' 'Go may download pinned dependencies on the first build.'
    go build -trimpath -o bin/explain ./cmd/explain
    printf '\n%s\n%s\n' 'Built: bin/explain' 'Try: ./bin/explain grep'

# Build release archives for a version into dist/ (local dry run; publishes nothing).
@release-build version='v0.0.0-dev':
    printf '\n%s\n' 'Building release archives (needs GNU tar; nothing is published)...'
    bash scripts/build-release.sh '{{ version }}' dist
    printf '\n%s\n' 'Built: dist/ (archives and checksums.txt)'

# Validate Markdown and generate Astro's shared content data.
@content:
    printf '\n%s\n' 'Validating content/ and exporting the shared catalog...'
    go run ./cmd/content-export --source content --output web/src/data/catalog.json
    printf '\n%s\n' 'Generated: web/src/data/catalog.json (do not edit; Markdown is the source).'

# Check Go analysis, website formatting, and Astro types.
@check:
    printf '\n%s\n' '[1/3] Checking Go static analysis...'
    go vet ./...
    printf '\n%s\n' '[2/3] Checking website source formatting...'
    npm --prefix web run format:check
    printf '\n%s\n' '[3/3] Validating shared content and checking Astro types...'
    npm --prefix web run check
    printf '\n%s\n' 'All analysis, formatting, and Astro checks passed.'

# Run Go tests, temporary-fixture examples, and available host-tool checks.
@test:
    printf '\n%s\n' 'Running Go tests (examples use temporary fixture directories)...'
    go test ./...
    printf '\n%s\n' 'Go tests passed. Host-tool checks may skip unavailable implementations.'

# Validate content and build the static website in web/dist.
@web-build:
    printf '\n%s\n' 'Building the static Astro website from shared Markdown...'
    npm --prefix web run build
    printf '\n%s\n%s\n' 'Built: web/dist/' 'Local output only; nothing was published.'

# Start the local Astro development server in the foreground.
@web-dev:
    printf '\n%s\n%s\n' 'Starting local Astro development (validates/exports content first)...' 'Astro prints the local URL below. Press Ctrl+C to stop.'
    npm --prefix web run dev -- --ignore-lock

# Test live Markdown updates, then CLI/site parity, usability, and accessibility.
@web-test:
    printf '\n%s\n%s\n' 'Testing live Markdown updates and browser behavior against fresh builds...' 'Uses a free port, so a running just web-dev is not affected.'
    npm --prefix web test
    printf '\n%s\n' 'Browser tests passed; the test preview server has stopped.'

# Install the Pulumi program's dependencies (infra/pulumi).
@infra-setup:
    npm ci --prefix infra/pulumi

# Typecheck the Pulumi program and test the CloudFront function.
@infra-check:
    npm run typecheck --prefix infra/pulumi
    npm test --prefix infra/pulumi

# Preview infrastructure changes (read-only; needs AWS sign-in, see infra/README.md).
@infra-preview:
    cd infra/pulumi && pulumi preview --stack prod
