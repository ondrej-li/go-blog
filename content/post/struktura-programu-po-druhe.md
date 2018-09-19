+++
author = "Ondrej Linek"
categories = ["začátky"]
date = 2016-09-07T20:34:34Z
description = ""
draft = false
slug = "struktura-programu-po-druhe"
tags = ["začátky"]
title = "Struktura programu po druhé"

+++

Jak se říká, žádný člověk není ostrov, tak ani žádný program není ostrov. Programy závisí na knihovnách jazyka a na dalších knihovnách. *Go* řeší tedy dva problémy. Jak závislost získat a jak závislost použít.

Aby mohlo *Go* závislost získat, musí vědět, kde se závislost nachází. A k tomu slouží právě posledně zmiňované jméno package. Pokud máme např. package `github.com/stretchr/testify`, tak je naprosto zřejmé, že ji můžeme stáhnout z [githubu](https://github.com/stretchr/testify). *Go* nám v tom ale rádo pomůže. Stačí jen říct `go get github.com/stretchr/testify` a závislost je v našem `GOPATH`. 

Jakmile tedy máme všechny závislosti na `GOPATH`, tak můžeme v klidu onu závislost použít. Např. takto

```go linenumbers
package main

import "./dalsi_package"
import nadrazena "../"
import "testing"
import "github.com/stretchr/testify/assert"

func TestNěco(t *testing.T) {
   ...
}
```

V předchozím příkladu jsme to trochu přehnali, ale zase je zde vidět jak

* referencujeme podřízenou package `dalsi_package`. Pokud by soubor byl v package `$GOPATH/src/github.com/lojza/prvni_projekt` tak bychom tato package byla v `$GOPATH/src/github.com/lojza/prvni_projekt/dalsi_package`
* referencujeme nadřízenou package `../`. Ta by byla v `$GOPATH/src/github.com/lojza/`. V našem souboru ji ovšem budeme používat jako `nadrazena`. Tedy např. `nadrazena.NejakaFunkce()`
* Dále pak odkazujeme na package z SDK.
* Na závěr pak použijeme externí referenci. Tu si ovšem musíme nejdřív obstarat pomocí `go get`.

Pokud přicházíte na *Go* z Javy tak budeme mít z počátku tendenci vytvářet mnoho packagí. V *Go* se to tak běžně nedělá. I relativně velké programy mají jen několik packagí.

Každý zdroják v *Go* musí začínat určením package. Podle konvence to je buď poslední část za lomítkem (v našem případě tedy `prvni_projekt`, nebo pokud se jedná o spustitelný program, tak to bude `main`. V jedné package (adresáři) musí všechny zdrojáky uvádět jedno jméno package!

A jako bonus viditelnost. *Go* jako takové nemá zapouzdření jak jej známe z objektových jazyků, ale používá jednoduchou konvenci pro "viditelnosti". Vše co začíná malými písmeny je viditelné pouze v dané packagi. A naopak vše začínající velkými písmeny je vidět i vně package. Tedy

```go
package "subpackage"

func neviditelnáFunkce() {}
func ViditelnáFunkce(){}
var neviditelnáProměnná string
var ViditelnáProměnná string
type ViditelnáStruktura struct {
   neviditelnéPole string
   ViditelnéPole string
}
```

