![main](https://github.com/Beheadedstraw/YARC/actions/workflows/Server.yml/badge.svg?branch=main) ![main](https://github.com/Beheadedstraw/YARC/actions/workflows/Client.yaml/badge.svg?branch=main)
![dev](https://github.com/Beheadedstraw/YARC/actions/workflows/Dev-Server.yaml/badge.svg?branch=dev) ![dev](https://github.com/Beheadedstraw/YARC/actions/workflows/Dev-Client.yaml/badge.svg?branch=dev) 
# YARC (Yet Another Redis Clone)

## About This Project
This project is mostly just me bumbling around with golang. Right now it's extremely bare bones that it's essentially a socket server that caches data in a key:value format using maps inside of slices. This can be used with a raw socket client like telnet/nc or programmatic sockets also. If I ever get around to it ill make a few libraries that makes it easier to use.

## Current Commands
Right now the only commands that are available are SET, GET, DEL. These commands do work but there's no validation on the data itself. It's literally bare bones.

## Querying The Server
There's no connection string or anything to get started. The current method of setting data is `0 set {key} {value}` where `0` is the DB store, `{set}` is the command, `{key}` for what key you want to store it in, and of course the `{value}`. You can't have duplicate keys in the same DB as `set` will just overwrite the key. `0 get key` will retrieve the value along with `0 del key` for deleting it.

Right now there's a barebones client that will query the localhost, you can run this as `./yarc-cli 0 set key value`.

## Under The Hood
This uses goroutines for each client and channels to maintain sync with the main DB store across the routines. I have no idea if there's a better way to do this, but from what I've read this is the easiest and most effective way to maintain sync across threads. If someone else has a better way I'm all ears. 

## Future Plans
Not much, just get it stable, add a few more commands, maybe throw in actual JSON support, connection strings, TLS, etc.

