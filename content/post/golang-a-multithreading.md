+++
author = "Ondrej Linek"
categories = ["pokročilý"]
date = 2016-10-09T07:19:14Z
description = ""
draft = false
slug = "golang-a-multithreading"
tags = ["pokročilý"]
title = "#Golang a multithreading"

+++

Tento článek jsem si nechával na slavnostní chvíle a ta zrovna přišla. Dnes si povíme o to paralelismu v *Go*.

V první řadě, co je to multithreading? Je to schopnost programu, dělat několik věcí najednou. Tedy například stáhnout stránku a číst z terminálu. Řeknete si, to přeci není problém, jen co uživatel stiskne tlačítko, tak prostě na konzoli vypíšu co, napsal. Třeba následovně.

```go
   fmt.Print("Napiš něco: ")
   var input string
   fmt.Scanln(&input)
   fmt.Println(input)
   resp, _ := http.Get("http://www.google.cz/")
   _, err := io.Copy(os.Stdout, resp.Body)
   resp.Body.Close()
```

Samozřejmě to takhle jde, ale je to když budete chtít ještě během toho vypisovat např. aktuální čas? Nebo ukazovat, jak píše jiný uživatel na jiném počítači. Nakonec skončíte u nějaké podobné smyčky.

```go
for {
   for _, funkce := range seznamFunkcí {
      funkce()
   }
}
```

Prostě nekonečná smyčka, kdy neustále voláte jednu funkci za druhou. A tak dokolečka. Přitom dnešní operační systémy dovolují toto dělat naprosto transparentně. Říká se tomu multithreading - více vláknovost. Stejně jako vám běží na počítači najednou několik programů, tak vám může v programu běžet několik funkcí. Problém, z hlediska programátora, není ovšem v tom, spustit několik threadů (vláken) a zpracovávat tak několik činnosti najednou, problém je jak tyto činnosti efektivně sladit. Jak zajistit, aby se nám vlákna nepoprala (např tím, že přistupují ke stejnému zdroji) a uměla spolu komunikovat (tedy předat si bezpečně data). Mnoho jazyků multithreading podporuje přímo (třeba Java, C, Python), některé jen přímo nepodporují (třeba JavaScript).

Než si povíme, jak psát multithreadově v *Go* si neodpustím jen lehkou poznámku na téma kdy se multithreading vyplatí a jak.

* Pokud v programu blokujete. Čtení z disku, přístup k síti a podobné operace blokují. Tzn. program se zastaví a čeká, až přijde odpověď ze sítě, načtou se data z disku apod. Zde se vyplatí rozhodně vyplatí multithreading.
* Pokud máte na výpočet náročnou operaci a máte více procesorů. Vyplatí se si práci rozdělit a využít potenciál všech jader procesoru.

Kdy se multithreading nevyplatí.

* Pokud máte víc threadů, než jader procesoru a všechny thready opravdu pracují (tedy počítají něco). V takovém případě nejen že vám multithreading nepomůže ba naopak. Režie spojená s koordinací několika threadů váš program nepatrně zpomalí. 
* Pokud vlákna přistupují ke stejnému zdroji. Tedy scénář, kdy několik threadů přistupuje ke stejnému disku. 
* Pokud vlákna ve svém důsledku způsobí ten samý problém jako v předchozím bodě, ale na jiném stroji. Například přístup na vzdálený souborový server z několika vláken.

### Práce s více vlákny v *Go*

Jazyk *Go* byl od počátku vyvíjen tak, aby podporoval multithreading. Přímo v jádře jazyka je tak jeho podpora. To je samozřejmě rozdíl proti jiným jazykům, které multithreading umožňují, ale pouze na úrovni knihoven. 

V *Go* máme pro podporu multithreadingu tyto klíčová slova a operátory.

* `go` - spustí funkci v novém vlákně
```go
go DělejNěcoVeVlákně(1, 1) // funkce spuštěna ve vlákně
```
* `chan` - definuje channel - kanál - sloužící ke komunikaci mezí vlákny
```go
var Kanál chan int = make (chan int) // deklarace a inicializace kanálu typu `int`
```
* `select` - čte zprávy z několika kanálů naráz
```go
select {
   a := <- Kanál: 
      fmt.Println ("přijata zpráva ", a)
}
```
* `range` - smyčka pro čtení z kanálu. Tak prosté jak to zní.
```go
kanál := make(chan int)
...
for číslo := range kanál {
   ...
}
```
* `<-` čtení a zápis do kanálu. Podle polohy kanálu a operátoru se rozhodne, jestli chcete číst, nebo psát.  
```go
kanál := make(chan int)
kanál <- 1 // zapisujeme to kanálu
číslo := <- kanál // čteme z kanálu
```

Výpis podpory *Go* pro práci s více vlákny by nebyl kompletní, pokud bychom se nezmínili i o těchto funkcích.

* `make` vytvoří kromě sliců i kanály. Např. pokud chceme kanál typu `int` pak použijeme
```go
intKanál := make (chan int)
stringKanál := make (chan string, 128) // 128 určuje velikost bufferu 
```
Zde sluší připomenout, že zasílání do kanálu blokuje (zastaví činnost než si někdo data přijme). Pokud je ovšem kanálu bufferovaný, tak k onomu blokování dojde až potom, co se naplní buffer.

* `close` ten uzavře kanál. Pokud se uzavře kanál, tak všichni, kdo z něho čtou budou odpojeni. Tzn. uzavření kanálu je možno použít jako signál k pokračování práce. Pokud např. jedno vlákno pracuje a chce dát druhému vláknu najevo, že skončilo, tak to:
```go
hotovo := make (chan struct{})
go func() {
   defer close(hotovo)
   ...
}()
<- hotovo
fmt.Println("Přechozí vlákno skončilo")
```

Pokud pracujete s mnoha vlákny, tak se vám určitě bude hodně vědět, kolik jich třeba běží, nebo dokonce být notifikován, až všechny doběhnou. Pro tyto účely má *Go* přímo ve standardní knihovně čekací skupiny, neboli WaitGroup. To umožňuje inkrementovat a dekrementovat počet pracujících threadů (teoreticky i čehokoliv jiného). Typický příklad použití vypadá následovně:

```go
wg := &sync.WaitGroup{} // pozor, používáme odkaz, nebudeme tak přenášet data hodnotou!
go func(wg *sync.WaitGroup{}) {
   defer wg.Done()
   wg.Add(1)
   ... // cokoliv
}(wg)
go func(wg *sync.WaitGroup{}) {
   defer wg.Done()
   wg.Add(1)
   ... // a ještě něco
}(wg)
wg.Wait() // počkáme, až obě vlákna skončí
fmt.Println("A máme padla!")
```

Takže nezapomeňte - prakticky všechno lze v *Go* provádět ve vlákně. Ba co víc, prakticky všechno se v *Go* řeší ve vláknech! A jelikož má *Go* opravdu promyšlenou podporu vláken, tak pracovat paralelně je prostě radost!



