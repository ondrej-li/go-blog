+++
author = "Ondrej Linek"
categories = ["začátky"]
date = 2016-09-07T17:38:37Z
description = ""
draft = false
slug = "vyvoj-v-golang"
tags = ["začátky"]
title = "Vývoj v #Golang"

+++

Jelikož už známe syntaxi jazyka a umíme napsat např. `Hello world`, tak nastává čas si říct něco o tom, jak si nakonfigurovat prostředí, abychom mohli pracovat s *Go*.

Instalace je jednoduchá. Pokud jste na Macovi, pak doporučuji použít `homebrew`, i když si můžete stáhnout instalačky. Linux je podobný, `yum` nebo `apt` vám zajisté pomohou. Co se Windows týče, tam bohužel sloužit nemohu, ale na vlastní oči jsem viděl *Go* běžet na Windows, tak to určitě půjde.

Po instalaci si zadejte v shellu `go version`, abyste si ověřili, že 

a) Go vám funguje
b) máte nainstalovanou verzi, kterou chcete

Následně zkuste `go env`, což vypíše informace o vašem prostředí. Mezi všemi možnými věcmi jsou nejdůležitější dvě proměnné. 

* `GOROOT` - ta ukazuje na adresář, kde je *Go* nainstalováno
* `GOPATH` - ta ukazuje na místo, kde leží vaše projekty. 

A tím se plynule dostáváme k projektům a struktuře kódu, jako takového.

Když jsem říkal, že `GOPATH` ukazuje na místo, kde leží vaše projekty, tak nebyl překlep. Je to opravdu místo, kde běžně leží ==všechny== vaše projekty v *Go*. `GOPATH` si nastavíte sami, nejlépe v `.bash_profile`. Když budete u toho, tak si rovnou nastavte `PATH` na `$PATH:$GOPATH/bin`. Třeba v mém `.bash_profile` je:

```bash
...
export GOPATH=/personal/go
export PATH=$PATH:$GOPATH/bin
...
```

S takto upraveným profilem se můžete pustit do prvního projektu. A hned se naskýtá otázka - v čem editovat. Takže, pro opravdové programátory tu máme `vim`, `emacs`, `notepad` apod. Pro produktivní programátory, co mají 2000Kč, tu máme [IntelliJ Idea](https://www.jetbrains.com/idea/) s pluginem [Go](https://plugins.jetbrains.com/plugin/5047). Pro produktivní programátory, co neradi utrácí, tu máme [Eclipse](http://www.eclipse.org/) opět s pluginem pro [Go](https://marketplace.eclipse.org/content/goclipse). Pro milovníky `vim` (včetně mně), je tu [vim-go](https://github.com/fatih/vim-go).

Nyní si konečně můžeme začít nový projekt. První krok, je potřeba si udělat novou package, tedy domeček, pro váš projekt. Pokud plánujete mít vaši aplikaci např. na githubu, bitbucketu a pod, tak vřele doporučuji nazývat package jménem vašeho projektu. Tedy např. `github.com/lojza/prvni_projekt`. *Go* podporuje UTF-8 a je dobré toho využít, ale názvy souborů a packagí jsou výjimka. Tam se držte malých písmen a pište cestinou. Abyste si takovou package vytvořili, tak stačí zadat 

```bash
mkdir -p $GOPATH/src/github.com/lojza/prvni_projekt
cd $GOPATH/src/github.com/lojza/prvni_projekt
```

Teď když máte package si můžete vytvořit svůj první zdrojový soubor. Nazvěme ho `main.go` a může vypadat třeba takhle:

```go
package main

import "fmt"

func main () {
   fmt.Prinln ("Ahoj Lojzo")
}
```

A teď když jste stále v adresáři `$GOPATH/src/github.com/lojza/prvni_projekt`, zadejte `go install`. Pokud šlo vše hladce, tak by se nemělo nic vypsat, ale zkuste zadat `ls -al $GOPATH/bin` a měli byste vidět, kromě jiného i `prvni_projekt`, jako spustitelný soubor. A jelikož je `$GOPATH/bin` na vaší `PATH` (pamatujete - nastavili jsme tak `.bash_profile`), tak stačí zadat

```bash
prvni_projekt
```

a vypíše se

```bash
sputnik:prvni_projekt defectus$ prvni_projekt
Ahoj Lojzo
```

Adresářová struktura by měla vypadat následovně
![](/content/images/2016/09/prvni_projekt.png)

PS: jako bonus uvedu, ještě než si o tom povíme v samostatné kapitole, že pokud použijete ve vašem projektu odkaz na jinou knihovnu (třeba naprosto nepostradatelnou `github.com/stretchr/testify/assert`), tak stačí jen zadat `go get ./...` a Go za vás postahuje všechny závislosti. Prostě nádhera :-)