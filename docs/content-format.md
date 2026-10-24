# Shared Markdown content format

This is a proposed format to validate before implementation. No content parser or command pages exist yet.

## Authoring source

Keep shell explanations in `.md` files, with YAML front matter and a restricted Markdown body. Start with English. Use one canonical file per command entry, for example `content/en/grep.md`. Front matter holds structured metadata, not a second copy of prose. The body holds the explanation, examples, and portability details shown by both consumers.

Proposed metadata example, shown here only as documentation:

```yaml
---
schema_version: 1
id: grep
title: grep
summary: Search text for lines that match a pattern.
aliases: [ggrep]
tags: [text, search, logs]
implementations: [gnu, apple, freebsd, busybox]
references:
  - https://www.gnu.org/software/grep/manual/
---
```

An implementation list declares authored coverage; it does not prove every implementation supports every option. Only list a variant once its relevant sections have been reviewed. Aliases represent documentation lookup names, not permission to assume the installed implementation; `ggrep` still needs appropriate selection evidence.

Use a small prescribed section structure: overview, use cases, examples, options, common mistakes, portability, and references. Repeated content belongs in common sections. Variant-specific subsections use stable machine-readable identifiers defined by the schema, not free-form prose heuristics.

## Variant sections

Proposed convention: level-two headings identify named views, such as `## Common`, `## GNU`, `## Apple`, and `## BusyBox`; level-three headings divide examples/options/mistakes within a view. The validator maps only documented headings to variant IDs and rejects ambiguous duplicates. An alternative directive format can be chosen during the rendering spike, but do not invent an untested syntax in production content.

Compose the common section with the chosen implementation section in both CLI and web. Common examples must actually work across the implementations the page claims. If a command differs too much, keep the common overview brief and author explicit per-implementation examples.

Model version-sensitive options with structured metadata once the first real incompatibility requires it. Do not silently show newest-version options to an older detected binary. When no applicable version data exists, show that limitation rather than infer complete compatibility.

## Markdown subset

Support headings, paragraphs, ordered/unordered lists, inline code, fenced code, and links. Prefer short sections over wide tables for command explanations; project planning documents may still use tables. Preserve code exactly while wrapping prose for terminal width. Do not depend on raw HTML, MDX, JavaScript, JSX, or Astro-only components.

Every example needs a purpose and a short explanation of important options/arguments. Include expected results where deterministic. Explain placeholders such as `FILE` consistently. Use generic fixtures and quote patterns/paths appropriately. Avoid shell prompts in code that users are expected to copy.

Option lookup needs a structured association between an option spelling and its Markdown subsection. Do not implement it by searching for arbitrary substrings; `-n` inside unrelated prose must not become a false option definition. Support aliases and value-taking options explicitly.

## Rendering contract

The CLI and website must agree on selected content IDs, order, code text, links, warnings, and variant applicability. HTML and terminal styling can differ. Browser pages can expose all variants, while CLI defaults to one selected view.

Build derived search records and JSON from the same source. Include schema version, content revision, command ID, implementation, version applicability, and sections in machine output. Do not include terminal escape codes. Missing content is a structured result, not generated prose pretending to be a manual.

Validate schema, unknown keys, duplicate IDs/aliases/headings, unsupported Markdown, unsafe URLs, broken links, and missing variant coverage. Golden fixtures should compare web and terminal selections. Examples need implementation-specific tests in disposable environments; document untested coverage honestly.

## Content ownership

Write original examples and explanations. Link upstream manuals for detailed authority. Borrowing project ideas does not permit copying their text. If reuse is approved, preserve provenance, attribution, licensing, and modification notices as required. Choose code and content licenses before importing a corpus or publishing reusable releases.
