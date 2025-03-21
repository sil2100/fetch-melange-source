# fetch-melange-source

A small Go application that attempts to fetch the source of a package. Sadly, for now this script is VERY naive and only performs the first git-checkout/fetch.

NOTE! Do not rely on this returning the full source for now! It's only useful for things like license linting.

## Usage

```
go build -o ./ fetch-melange-source.go
./fetch-melange-source ../os/vim.yaml source/
```

## Future ideas

Ideally, this app will be made smarter and smarter to fetch and patch all the sources, meaning that the resulting source would be more-or-less complete. My idea then was to move this into melange itself, adding it as a melange subcommand like `melange fetch-source` or similar.