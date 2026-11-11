# Third-party material and release notices

Inventory reviewed against the checked-in Go module versions and web lockfile on 2026-11-11. This is a release preparation record, not a substitute for checking the exact artifacts being shipped. Revisit it whenever dependencies or assets change.

## CLI

The locally built CLI reports Go 1.27.3 and two third-party modules in `go version -m bin/explain` (checked 2026-11-11):

| Module | Version | License and included notice |
| --- | --- | --- |
| Go runtime and standard library | Go 1.27.3 | BSD-style: [license](../third-party/go-runtime/LICENSE) and [patent grant](../third-party/go-runtime/PATENTS) |
| github.com/yuin/goldmark | v1.8.6 | MIT: [license and copyright](../third-party/goldmark/LICENSE) |
| go.yaml.in/yaml/v3 | v3.0.5 | MIT for the listed libyaml-derived files and Apache-2.0 for the remaining files: [license](../third-party/go-yaml/LICENSE), [NOTICE](../third-party/go-yaml/NOTICE), [Apache-2.0 text](../third-party/go-yaml/Apache-2.0.txt) |

For each binary distribution, include the project's MIT and CC BY 4.0 texts, the content attribution identified in the README, and the complete `third-party/` notices above. Apache-2.0 requires its license text and applicable NOTICE attribution when distributing the covered work; retain relevant upstream notices and mark changes to modified Apache-licensed files. Check the actual module graph, Go toolchain and build inputs again at release; update the runtime notices if the toolchain changes. No release packaging exists yet, so this requirement is not automatically enforced.

## Website and build tools

`web/package-lock.json` contains the pinned npm dependency graph, including optional platform packages and development-only tools. Its declared licenses include MIT, ISC, Apache-2.0, BSD, MPL-2.0 and LGPL-3.0-or-later. For example, `lightningcss` is MPL-2.0 and optional `@img/sharp-libvips-*` packages declare LGPL-3.0-or-later. Their presence in the lockfile does not by itself establish that their code is present in the generated site. License fields are inventory leads, not verified complete notices or a compatibility conclusion.

The local production output inspected on 2026-11-11 contained HTML, the authored favicon, one CSS file and four small JavaScript files. The scripts trace to the project’s Catalog component and copy-buttons source; the CSS traces to the project stylesheet. No third-party client code or assets were identified in that output. Astro, sharp, lightningcss and other npm packages are build or test tools for this output, not automatically bundled merely because they are installed. Inspect the exact production output again before publication, including inline scripts and generated assets, and provide any newly applicable copyright, license and NOTICE text and source or modification notices. If redistributing native build dependencies, assess their LGPL/MPL obligations for that distribution separately. Recheck the lockfile and installed platform packages for each release; do not claim the site audit is complete until the actual artifacts have been checked.

The tracked site assets include an SVG favicon. The site links to upstream manuals and related projects as references; the project documentation says no third-party explanation corpus was imported. This review did not independently establish authorship or copyright status of every original or AI-assisted passage, and the history and provenance audit remains a public-release gate.
