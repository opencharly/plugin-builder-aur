# plugin-builder-aur

The `aur` builder for OpenCharly — the Arch User Repository (AUR) package
builder as a plugin `builder:` word.

A candy that ships an `aur:` package section is built by this builder: the host
selects it by **detection** (the candy's `aur:` section), never by an authored
`external_builder:`. The provider serves the build-time multi-stage
(`FROM <builder> AS aur-build` + the `/tmp/aur-pkgs` COPY artifact) and the
deploy-time teardown IR (package removal).

## What it provides

| Capability | Surface |
|---|---|
| `builder:aur` | the `aur` builder word — the AUR/makepkg multi-stage builder |

The builder authors no `plugin_input` (it is triggered by detection, not an
authored field); its self-contained `#AurBuilderInput` (`schema/aur.cue`) is
empty by design and ships so the schema travels with the plugin.

## How to use it

Compose the plugin candy in a box's `candy:` list:

```yaml
- '@github.com/opencharly/plugin-builder-aur/candy/plugin-builder-aur:<tag>'
```

Then a candy whose packages come from the AUR declares them in an `aur:`
section; the host resolves them through this builder.

## Layout

- `candy/plugin-builder-aur/` — the plugin module: `plugin.go`,
  `schema/aur.cue`, `cmd/serve/main.go`.
- `charly.yml` — the root project manifest (`discover: candy`).
- `.github/workflows/tag-on-merge.yml` — CalVer tag + `CHANGELOG/` on merge.

## Related

- Owning skill: `/charly-image:image` — box/builder configuration and the box
  dependency graph. This candy carries no `skill:` entity of its own; the gap is
  tracked in
  [opencharly/opencharly#291](https://github.com/opencharly/opencharly/issues/291).
- `/charly-internals:plugin` — the plugin/provider model, including the
  `builder` provider class.
- [`opencharly/charly`](https://github.com/opencharly/charly) — the charly CLI.
