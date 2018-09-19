+++
author = "Ondrej Linek"
categories = ["pokročilý"]
date = 2016-09-17T09:07:31Z
description = ""
draft = false
slug = "kompilace-golang-v-dockeru"
tags = ["pokročilý"]
title = "Kompilace #Golang v Dockeru"

+++

Tentokrát kratší vsuvka. 

Potřeboval jsem na svém  Macovi zkompilovat program pro 🐧 Linux. Samozřejmě mě jako první napadl cross compiler, ale to není ono. Nehledě na to, že bych pak stejně nemohl program vyzkoušet. Samozřejmě by s trochou cviku šlo zkusit nějaký VirtualBox apod., ale mě napadlo elegantnější řešení - [Docker](https://www.docker.com)!

Nejdřív si nainstalujete na počítač docker - podívejte se na [instalační stránky](https://www.docker.com/products/overview) a vyberte si instalaci podle svého systému. V ideálním případě by toto mělo fungovat:

```bash
defectus@sputnik ~ $ docker -v
Docker version 1.12.1, build 6f9534c

defectus@sputnik ~ $ docker run --rm hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world

c04b14da8d14: Pull complete
Digest: sha256:0256e8a36e2070f7bf2d0b0763dbabdd67798512411de4cdcf9431a1feb60fd9
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
```

Pokud ano, tak tady je příkaz, kterým si zkompilujete svoje *Go* zdrojáky do Linuxu:

`docker run --rm -v "$PWD":/go -w /go golang:1.7 go build -v -o nazev_souboru cesta/k/vasi/package`

Einfach :-)

