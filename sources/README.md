# Sources

## 1. CheatSh

`cheatsh.CheatSh` implements `radium.Source` using [cheat.sh](https://cheat.sh) as the source of information.

### Examples:

```bash
radium query "append file in go"

radium query "append file" -a language:go

radium query "open socket in java" -a nocolor
```

### Attributes:

* `nocolor` - Specifying this will force cheat.sh to send non-colored snippets (appends `?T` query param)
* `language` - Programming language (recommended to use `append file in go` format instead of this tag)


## 2. Tldr

`sources.TLDR` implements `radium.Source` using [tldr](https://github.com/tldr-pages/tldr).

### Examples:

```bash
radium query "ls" --sources "tldr"

radium query "dir" --sources "tldr" -a platform:windows
```

### Attributes:

* `platform` - Specifying this will limit the lookup scope to platform specific directories in `tldr`


## 3. Wikipedia

`wikipedia.Wikipedia` implements `radium.Source` using [Wikipedia](https://wikipedia.org).

### Examples:

```bash
radium query "ls" --sources "wiki"

radium query "hindi" --sources "wikipedia" -a language:hi
```

### Attributes:

* `language` - Return results in specified language. (e.g., `-a language:hi` will look up in `https://hi.wikipedia.org/wiki`) 

## 4. DuckDuckGo

`duckduckgo.DuckDuckGo` implements `radium.Source` using [DuckDuckGo Instant Answers API](https://api.duckduckgo.com).

### Examples:

```bash
radium query "apple inc" --sources "ddg"

radium query "spacex" --sources "duckduckdgo" 
```

## 5. LearnXInYMinutes

`sources.LearnXInY` implements `radium.Source` using [learnxinyminutes](http://github.com/adambard/learnxinyminutes-docs).

### Examples:

```bash
radium query "dart" --sources "lxy"
```
