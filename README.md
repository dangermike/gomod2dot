# Go mod 2 dot

Run this in a Go project to get a dependency graph as a [dot digraph](https://graphviz.org/), which you can render using the `dot` command line tool. This thing also takes in an optional filter, which is a RE2 regex. If you want to see why a particular package was imported in a graphical way, this is a way to do that.

## Example

```bash
gomod2dot urfave/cli | dot -Tpng > go.mod.png
```
