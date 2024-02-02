# Apologies Go Library

This is a Go library that implements a game similar to the [Sorry](https://en.wikipedia.org/wiki/Sorry!_(game)) board game.

_This is a work-in-progress._

I'm in the process of implementing this as a learning exercise. My goal is to re-implement something that is functionally similar to the original Python version ([apologies](https://github.com/pronovic/apologies)), except in GoLang.  This problem is interesting enough that it should help me get my head around the language.

_Note:_ while the code is functionally equivalent, it's organized somewhat differently, to reflect the differences in the languages.  A Python module (`mymodule.py`) roughly maps to a Go package (i.e. the subdirectory `mymodule`), but I've renamed some things and moved other code around to work better with Go's naming conventions.  For instance, `game.py` was mostly moved to the `model` package, and some of the functionality in `rules.py` was moved into `model` and `reward`.  
