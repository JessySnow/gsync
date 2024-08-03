# Gsync

This is a GitHub Release sync tool based on a B+ tree index implemented in GoLang, created just for fun.

## Features

- Synchronize releases from GitHub repositories.
- Efficient indexing with a B+ tree structure.

## Installation

- build
```shell
cd cmd && go build -o gsync
```
- create a config file
> Create a configuration file named gsync.json in the same directory as the executable file.```json
 
```json
{
  "root_dir": "a dir already exists",
  "sync_interval": 30,
  "sync_gap": 10,
  "repos": [
    {
      "owner": "repo's owner",
      "repoName": "repo's name"
    }
  ]
}
```
`sync_interval`: The time interval between two synchronizations, measured in minutes.

`sync_gap`: The time interval between the synchronization of two repositories, measured in minutes.

- run