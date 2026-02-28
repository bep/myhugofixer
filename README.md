# myhugofixer

A CLI tool that prints version-specific Hugo deprecation and migration guides as Markdown.

## Installation

```sh
go install github.com/bep/myhugofixer@latest
```

## Usage

```sh
# Print all available fix guides
myhugofixer

# Print fixes starting from a specific version
myhugofixer -low v0.110.0

# Print fixes up to a specific version
myhugofixer -high v0.146.0

# Print fixes for a specific version range (inclusive)
myhugofixer -low v0.110.0 -high v0.146.0
```

## Use with Claude

Example usage:

```
claude
Run myhugofixer --low v0.146.0 and apply the Hugo upgrades described. Commit when done.
```

Or without interaction:

```
claude -p "Run myhugofixer --low v0.146.0 and apply the Hugo upgrades described. Commit when done."
```

 Note: The workspace trust dialog is skipped when Claude is run with the -p mode. Only use this flag in directories you trust.
