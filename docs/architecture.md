# Product and architecture

## Confirmed direction

Develop shell-explain in a private repository through stages 1–4, with a CLI executable named explain and a planned static Astro website at shell-explain.ninegene.com. Stage 5 selects code/content licenses, opens the repository to the public, establishes external contribution guidance, and publishes the website. A shared collection of Markdown files supplies both outputs. The CLI must eventually be downloadable as a binary and installable using curl, Homebrew, and mise.

The initial experience is `explain grep`: what it does, common uses/options, examples, mistakes, and portability notes. It should work without network access after installation. Match documentation to the installed tool when feasible and expose manual selection.

## Language choice

Go is the implemented CLI language. This workload is primarily content loading, command selection, search, and formatting, rather than performance-intensive computation. Go supports target-specific builds and embedded files without requiring an interpreter on the user's machine. The parser uses the pinned pure-Go goldmark and YAML libraries. Rust was a viable alternative but is not needed to meet these goals.

References: [Go build/environment documentation](https://pkg.go.dev/cmd/go#hdr-Environment_variables), [Go embed](https://pkg.go.dev/embed), and [Rust target support](https://doc.rust-lang.org/rustc/platform-support.html).

## Shared-source pipeline

Implemented layout: `content/en/` contains authored command Markdown, `content/embed.go` embeds those exact files, `internal/catalog` validates and renders them, `internal/detect` selects implementations, `internal/cli` implements the interface, and `cmd/` provides the CLI and content exporter. `web/` contains Astro and browser tests; the root justfile exposes build/check commands.

The Go parser reads versioned Markdown into one normalized content model. Terminal rendering, JSON output, Astro pages, and static search data derive from it. The content exporter writes ignored web/src/data/catalog.json before Astro dev/build/check, and a Vite development plugin watches the original Markdown and regenerates it on saves/additions/deletions. Exports are serialized and atomically replaced; invalid edits retain the last valid catalog. Astro's file loader validates that generated model through a content collection; it does not independently parse the Markdown. Do not manually copy explanations into Go strings, Astro components, or a separate JSON database.

Schema version 1 and a deterministic SHA-256 content revision identify the bundle. Go tests reject malformed metadata/unsupported Markdown and preserve code; browser tests compare every documented implementation page with CLI JSON, including exact code and HTML. This single-parser design resolves the earlier choice between independent parsers and a normalized artifact.

## CLI scope

Initial commands: show an explanation, list entries, search entries, show a particular option, choose an implementation/platform, and emit structured JSON. Include version/help output, noninteractive plain text, terminal-width wrapping, and color controls. Do not require a pager or manipulate the user's shell configuration.

Normal lookup must not contact a hosted API or an AI service. Do not execute sample commands. Unknown commands/options should get an honest unavailable result and useful nearby matches, not invented help.

Pipeline/script parsing is deferred. If introduced, use a proper parser such as [mvdan/sh](https://github.com/mvdan/sh) and keep parsing separate from interpretation/execution. Explaining a syntactically valid command does not certify it as safe.

## Astro website

Use Astro's content collections/build-time loaders for shared Markdown metadata and static routes. Astro documents both Markdown rendering and collection schemas. See [content collections](https://docs.astro.build/en/guides/content-collections/) and [Markdown](https://docs.astro.build/en/guides/markdown-content/).

Build command pages, a catalog, tag pages, and a browser search index as static artifacts. The public reference needs no accounts, database, server renderer, or lab provisioning. Host the built files on GitHub Pages, with GitHub Actions building and deploying the Astro site. No Azure resources or paid cloud hosting are needed for this website.

Keep stages 1–4 private and use local Astro previews. Publish on GitHub Pages in stage 5 after the repository becomes public. GitHub Free supports Pages from public repositories; private-source Pages requires a qualifying paid GitHub plan and is not part of this plan. See [GitHub Pages availability](https://docs.github.com/en/pages/getting-started-with-github-pages/what-is-github-pages) and [Astro's Pages deployment guide](https://docs.astro.build/en/guides/deploy/github/).

Configure shell-explain.ninegene.com as the Pages custom domain, verify domain ownership, and configure DNS and HTTPS during authorized publication. Keep Astro's canonical site URL and asset paths consistent with the custom domain rather than assuming a repository subpath. Deployment and DNS changes remain outside this documentation-only update.

Provide explicit platform/implementation filters with shareable URLs. Default to a common portable view when possible. Browser OS detection must not claim to know which executable the visitor uses. Display version applicability, references, copyable examples, and links to related commands. Search/filter state should work with keyboard navigation and readable mobile layouts.

Allow ordinary Markdown only in explanation sources. Do not require MDX, Astro components, raw HTML, or executable template expressions, since those would not translate cleanly to CLI output. Sanitize rendered content and URLs. Avoid automatic execution, third-party tracking, and remote font/script dependencies unless justified.

## Integration contract

Other applications should pin a binary/content release and invoke a documented JSON interface with bounded input/output, or consume a versioned static export later. Keep JSON schema version, content revision, selected implementation, selection evidence, and warnings explicit. Terminal ANSI sequences must not appear in JSON.

Reference text belongs here; lesson-specific hints, quizzes, grading, and user progress belong to the consuming application. Keep the public project independent of any private consumer or infrastructure.
