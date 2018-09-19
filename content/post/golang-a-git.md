+++
author = "Ondrej Linek"
categories = ["začátky"]
date = 2016-09-12T19:02:00Z
description = ""
draft = false
slug = "golang-a-git"
tags = ["začátky"]
title = "#Golang a Git"

+++

Čas od času budete potřebovat svoje dílo uložit někam jinam, než na svůj počítač. Mnozí budou dokonce pomýšlet na sdílení svého kódu s ostatními např. na [githubu](www.github.com). K tomu je samozřejmě nejpříhodnější použít `git`. A jak nejlépe na to?

V první řadě je potřeba mít nainstalovaný git. Majitelé maců to mají jednoduché, tam je většinou git hned v základu. Ostatní si nainstalují git buď pomocí `apt-get`, nebo `yum`.

Nyní máte dvě možnosti. Už jste si napsali nějaký kód a chcete ho dát např. na github, nebo jste ještě nic nenapsali, ale chcete kód (který teprve napíšete) např. na github.

V obou případěch si založte [nový repozitář](https://github.com/new) např. na githubu.
Před prním použitím gitu doporučuji několik nastavení, co ulehčí život. 
`git config --global user.name "Jan Novák"`
`git config --global user.email jan.novak@gmail.com`

### Mám již kód, chci jej publikovat

* Nejdříve přejděte do adresáře package vaší aplikace.
* Proveďte `git init`. To vytvoří git repozitář v package vaší aplikace.
* Přidejte všechny soubory k správě pomocí gitu příkazem `git add .`
* Potvrďte změny pomocí `git commit -m "Můj první commit"`
* Přidejte si github jako vzdálený repozitář: `git remote add origin https://github.com/user/repo.git` (nahrďte `https://github.com/user/repo.git` za cestu k vašemu repozitáři)
* A publikujte změny příkazem `git push origin master`. Pokud pracujete jen s jedním vzdáleným repozitářem, pak `git push -u origin master` vám projekt vypublikuje a nastaví `origin master` jako hlavní repozitář.

### Nemám kód a chci jej napsat a publikovat
* Vytvořte si adresář `$GOPATH/src/github/user` např `mkdir -p $GOPATH/src/github/user`.
* Nyní proveďte `git clone https://github.com/user/projekt1.git`, kde `https://github.com/user/projekt1.git` je váš repozitář na githubu.
* Napište kostru vašeho programu.
* Potvrďte změny pomocí `git commit -m "Můj první commit"`
* A pošlete kódu na github příkazem `git push`.

Osobně doporučuji tyto dvě nastavení na úrovni projektu:   
`git config pull.rebase true`   
`git config rebase.autoStash true`   
Tím se nastaví default při pull na rebase (více v [manuálu](https://git-scm.com/book/cs/v1/V%C4%9Btve-v-syst%C3%A9mu-Git-P%C5%99eskl%C3%A1d%C3%A1n%C3%AD))
A také se před pull automaticky provede stash, tedy uloží lokální změny a znovu se aplikují po dokončení pull.

Ale aby jsem naprosto nezazdil to hlavní, co by si měl odnést nováček v #Golang, co chce použít git:

**Do gitu dávejte jen `$GOPATH/src/github.com/vase_jmeno/projekt`!**
