# CheckScope

Useful tool for checking if the URLs you provided are in or out of the scope of a Bug Bounty Program.

## Install

```bash
go install github.com/vitorfhc/checkscope@latest
```

## Usage

**Get everything in scope**

```bash
cat all-urls.txt | checkscope
```

**Get everything out-of-scope**

```bash
cat all-urls.txt | checkscope -r
```

**Define a diferent scope file**

```bash
cat all-urls.txt | checkscope -i scopefile.txt
```

**See other flags**

```bash
checkscope -h
```

## Scope file

The scope file (`scope.txt` by default) must contain the matching scope.

**Example**

```
*.foo.com
bar.com
*.zed.*
```

|Hostname|Matches|
|-|-|
|`https://xablau.foo.com/a/b/c?a=1`|Yes|
|`xablau.foo.com/a/b/c?a=1`|Yes|
|`abc.foo.com`|Yes|
|`foo.com`|Yes|
|`a.zed.b`|Yes|
|`sadad.net`|No|
|`sub.bar.com`|No|
