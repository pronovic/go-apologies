# Apologies Go Library

_Note: I developed this code while working at a new job where GoLang was prevalent. I no longer work there, and don't anticipate working in GoLang much going forward.  So, this repository is archived and isn't really maintained._

This is a Go library that implements a game similar to the [Sorry](https://en.wikipedia.org/wiki/Sorry!_(game)) board game.

The implementation is based on [apologies](https://github.com/pronovic/apologies), which I wrote in 2020 as a Python techology demonstration project.  I wrote go-apologies to help learn Go and the Go ecosystem.  The apologies problem was large enough to force me to struggle with the new language, while still being approachable.  While the code is functional and the unit tests pass, please keep in mind that this was one of my first attempts at Go &mdash; my inexperience probably shows.

While the go-apologies code is functionally similar to apologies, it's organized differently, to reflect the differences in the languages.  As I first started writing the code, a given Python module (`source.py`) was usually mapped into an equivalent Go package (i.e. the subdirectory `source`).  However, I eventually refactored a lot of the code to work better with Go's naming conventions.  For instance, `game.py` was mostly moved to the `model` package, and some of the functionality in `rules.py` was moved into `model` and `reward`.  I wanted this to look like Go code, not Python code.

This isn't a complete duplicate of the original Python implementation.  For instance, I did not implement the simulation functionality, which was used mostly while I developed the reward scoring algorithm. The command line interface for the demo is different, and the demo looks a little different because Go's ncurses doesn't support unicode characters like ▶ or ◼ that work fine from Python.  However, besides little things like that, go-apologies is a fairly faithful translation of apologies from Python to Go.
