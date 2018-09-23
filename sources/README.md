# Sources

## 1. CheatSh

CheatSh implements `radium.Source` using [cheat.sh](https://cheat.sh) as the source of information.

### Examples:

```
radium query "append file in go"

radium query "append file" -a language:go

radium query "open socket in java" -a nocolor
```

### Attributes:

* `nocolor` - Specifying this will force cheat.sh to send non-colored snippets (appends `?T` query param)
* `language` - Programming language (recommended to use `append file in go` format instead of this tag)


