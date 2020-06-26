# gitget
Browse through files and directories of git repository and download them to your current working directory (`$PWD`).

# Installation
```sh
go get github.com/FalcoSuessgott/gitget
```

# Usage
```
$ gitget https://github.com/golang/example
```

# Features
* checkout branches
* recursively copiying directories
* read git-urls from Clipboard Buffer

# Demo
```sh
$ gitget https://github.com/golang/example.git
Fetching https://github.com/golang/example.git

Enumerating objects: 150, done.
Total 150 (delta 0), reused 0 (delta 0), pack-reused 150

Checking out the only branch: master
? Select files and directories to be imported  [Use arrows to move, space to select, type to filter]
  [ ]  [27] │   │   ├── lookup.go
  [ ]  [28] │   ├── nilfunc
  [ ]  [29] │   │   ├── main.go
  [ ]  [30] │   ├── pkginfo
  [ ]  [31] │   │   ├── main.go
  [ ]  [32] │   ├── skeleton
  [ ]  [33] │   │   ├── main.go
  [ ]  [34] │   ├── typeandvalue
  [ ]  [35] │   │   ├── main.go
  [ ]  [36] │   ├── weave.go
  [ ]  [37] └── hello
  [ ]  [38] │   ├── hello.go
  [ ]  [39] └── outyet
  [ ]  [40] │   ├── Dockerfile
  [ ]  [41] │   ├── containers.yaml
  [ ]  [42] │   ├── main.go
  [ ]  [43] │   ├── main_test.go
  [ ]  [44] └── stringutil
  [ ]  [45] │   ├── reverse.go
  [ ]  [46] │   ├── reverse_test.go
> [x]  [47] └── template
  [ ]  [48]     └── image.tmpl
  [ ]  [49]     └── index.tmpl
  [ ]  [50]     └── main.go
  [ ]  [51]

Fetched the following files and directories:
/home/morelly_t1/git/gitget
└── template
    └── image.tmpl
    └── index.tmpl
    └── main.go
```
