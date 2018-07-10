# pretty

A library for formatting and printing objects to the command line in Golang.
Currently, it only has support for tables.

## Licensing

This library is MIT-licensed.

## Example

```go
table, err := pretty.NewPrettyTable(
  pretty.NewColumnDef("Name"),
  pretty.NewColumnDef("Type"))
if err != nil {
	return err
}

table.SetHeader("People")

table.AddRow("Noel", "Human")
table.AddRow("David", "Cyborg")
table.AddRow("Pranava", "Crusher")

table.Print()
```

```
--------
 People |
+---------+---------+
| Name    | Type    |
+---------+---------+
|    Noel |   Human |
|   David |  Cyborg |
| Pranava | Crusher |
+---------+---------+
```

## Testing

Run `go test -vet="" -short -v ./...`.

## Get involved

We are happy to receive bug reports, fixes, documentation enhancements, and
other improvements.

Please report bugs via the
[github issue tracker](https://github.com/rubrikinc/pretty/issues).

Contributions will be accepted only after the execution of a Contributor License Agreement.
