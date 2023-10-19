# Line Number

`line-number` takes in some newline deliminated text and outputs those same lines but with a line number gutter on the left side.

## Usage:

Using stdin/out:

```
% cat go.mod
module line-number

go 1.21.2

% cat go.mod | line-number
 1 | module line-number
 2 |
 3 | go 1.21.2
```

With files:

```
% line-number -in in.file -out out.file
```