# plugin-builder-aur

The `plugin-builder-aur` plugin candy of the [opencharly/charly](https://github.com/opencharly/charly)
candy library, as a standalone repo (the candy de-submodule cutover, plugin
kind). The Go module lives at `candy/plugin-builder-aur/` with module path
`github.com/opencharly/plugin-builder-aur/candy/plugin-builder-aur`; the charly resolver fetches this repo at the pinned tag and
the compiled-in wiring imports the module at that path.
