+++
author = "Ondrej Linek"
categories = ["pokročilý"]
date = 2016-09-21T06:08:12Z
description = ""
draft = false
slug = "testujeme-v-golang-cast-1"
tags = ["pokročilý"]
title = "Testujeme v #Golang (část 1)"

+++

Testování pro mnoho lidí (včetně mě) není nic zábavného, ale po letech programování připouštím, že je lehčí mít pořádné testy, než neustále kompilovat a spouštět aplikaci a pomocí debuggeru se snažit pochopit, co je špatně. Už ta ztráta času. I když je *Go* hodně rychlé, kompilace pořád představuje pár vteřin, následně spustit aplikaci, dostat na místo, kde si chcete něco ověřit a následně konečně to ověřit. Mnohdy to vyžaduje např. zásahy do databáze, otevírání browseru a navigace na stránky apod. Prostě a jednoduše, je to  neefektivní. A hlavně, jsme programátoři, automatizujeme co se dá! 

Proto si lidé píšou testy - rychle a opakovaně si jimi potvrzují, jestli jimi napsaný kód vyhovuje požadavkům. Testy se tedy skládají z těchto fází:

* inicializace - nastavení stavu do podoby, v jaké chceme mít kód pod testem. Většinou se jedná o nastavení proměnných apod., ale může se jednat třeba i o nahrání dat do databáze, spuštění jiného programu atd.
* exekuce - spustíme kód, který testujeme. Zde rozdělujeme dva základní druhy testů - unit testy, kdy testujeme (ideálně) pouze jednu funkci (jednotku kódu) a integrační testy, kdy testujeme souhru několika funkcí, nebo systémů.
* validace - zde si ověříme, že finální stav odpovídá našim předpokladům. Finálním stavem je myšlena především návratová hodnota, ale může se jednat o cokoliv.
* deinicializace - zvláště u integračních testů je potřeba vrátit stav do výchozí podoby (např. smazat databázi apod.)

*Go* má jednu podstatnou výhodu, proti mnohým jiným jazykům. *Go* bylo od začátku navrhované tak, aby se v něm dalo pohodlně testovat. Jinak řečeno, tam kde potřebujete v jiných jazycích knihovny a utility (např. *JUnit*, *NUnit*, *Karma*+*Jasmine*), tam si v *Go* vystačíte jen se základními package. *Go* má již v základu knihovnu pro testování (lakonicky pojmenovaná package `testing`) a spouštěč testů `go test`.

Než si ukážeme, jak psát testy v *Go*, tak si neodpustím pár poznámek. Ačkoliv je možné psát testy na jakýkoliv kód, je rozumné psát kód tak, aby byl dobře testovatelný. Nebo jinak řečeno, testovatelnost je vlastnost, se kterou počítáte od první řádky vašeho programu, ne něco, co přidá na závěr. Stejně si uvědomte, že testujete program **především** pro sebe, ne pro vašeho šéfa (pokud nějakého máte) a už vůbec pro zákazníka (ten mnohdy ani neví, co to testování je :-)).

Nuže, náš první test!

```go
// soubor main_test.go
package main

import "testing"

func TestSuperFunkce(t *testing.T) {
   if výsledek := SuperFunkce(1, 1); výsledek != 2 {
      t.Error("Toto jsem tedy nečekal!")
   }   
}
```

A k němu příslušný zdroják
```go
// soubor main.go
package main

func main () {
}

func SuperFunkce(a, b int) int {
   return a - b
}
```

A celý test spustíme pomocí `go test -v`. Výsledek by měl být:

```sh
defectus@sputnik go $ go test -v
=== RUN   TestSuperFunkce
--- FAIL: TestSuperFunkce (0.00s)
	main_test.go:7: Toto jsem tedy nečekal!
FAIL
exit status 1
FAIL	_/Users/defectus/personal/sandpit/go	0.005s
```

A když změníme naši `SuperFunkce`, aby čísla sčítala, tak nám test hezky projde:

```sh
defectus@sputnik go $ go test -v
=== RUN   TestSuperFunkce
--- PASS: TestSuperFunkce (0.00s)
PASS
ok  	_/Users/defectus/personal/sandpit/go	0.011s
```

Nyní si probereme, co jsme to vlastně stran testování udělali. Prvně, podle konvence, pokud se zdrojový soubor jmenuje `a.go`, tak testy k němu jsou v `a_test.go`. Technicky mohou být testy kdekoliv, co končí `*_test.go`, ale pak se v tom člověk těžko vyzná.
Dále, test funkce `func A`, se jmenuje `func TestA` a vždy přijímá parametr `*testing.T`. Ten se opět podle konvence jmenuje `t`. Pokud testujete funkci se příjemcem (receiver), tedy např. `func (a *A) Přičti(b *A)`, tak její název bude `func TestA_Přičti()`. Opět, je to konvence, ale váš kód bude hůř pochopitelný, pokud se jí nebudete držet.

Myslím, že toto bude stačit mnohým, ale příště se podíváme na další aspekty testování - mockování (tedy předstírání činnosti závislého kódu), několik vychytávek stran běhu testů apod.

