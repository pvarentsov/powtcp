# powtcp

This project is a simple example of [Proof of work (PoW)](https://en.wikipedia.org/wiki/Proof_of_work) protected TCP server. It implements chellenge-response protocol and uses [hashcash](https://en.wikipedia.org/wiki/Hashcash) algorithm.

<p align="center"> 
  <img src="assets/demo.png">
</p>


### Makefile

```bash
$ make help

Usage: make [command]

Commands:

 build-server          Build server app
 build-client          Build client app

 run-server            Run server app
 run-client            Run client app

 test                  Run tests
 fmt                   Format code
```