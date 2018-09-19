+++
author = "Ondrej Linek"
categories = ["pokročilý"]
date = 2016-10-24T19:15:30Z
description = ""
draft = false
slug = "objektove-programovani-vs-golang"
tags = ["pokročilý"]
title = "Objektové programování vs. #Golang"

+++

Dnes se podíváme na *Go* tak trochu filozoficky. Je, nebo není *Go* objektový jazyk, jak v tomto smyslu vypadá jeho srovnání se v současnosti nejpopulárnějším jazykem Java?

*Go* je opravdu tak trochu jako chytrá horákyně. Je funkcionální, ale ne zcela, je objektový, ale ne zcela. Než začne polemizovat o objektovosti *Go*, pojďme si připomenout, co dělá jazyk objektovým.

* Třídy a objekty. Vše jsou objekty a objekty jsou instance tříd.
* Zapouzdření. Data a operace nad daty jedno jsou.
* Dědičnost. Třídy dědí svoje vlastnosti a potomek třídy je specializací svého rodiče.
* Polymorfismus. Potomek dovede nahradit svého předka.
* Přetěžování. Metoda může mít stejné jméno, ale jiné parametry.

Nuže. *Go* nemá třídy a objekty. Má ovšem struktury a funkce, které nad nimi operují, tzv. přijímače (receivers). Jelikož volání přijímače nad strukturou se označuje operátorem tečky, tak to v kódu vypadá stejně, jako např. v Javě. 

```go
package main

import "fmt"

type Objekt struct {
	Jméno string
}

func (o *Objekt) Metoda() {
	fmt.Printf("Zavolána metoda na %s.\n", o.Jméno)
}

func main() {
	o := Objekt{Jméno: "objekt"}
	o.Metoda()
}
```

Zápis volání vypadá úplně stejně jako v Javě. Co je potřeba mít na paměti je, že *Go* nemá konstruktory. Inicializace se provádí při přiřazení. Většinou se to obchází přes arbitrární funkci `New` (např. `func New() *Objekt {...}` ), ale není to podmínkou, člověku se meze nekladou.

*Go* nemá klasické zapouzdření. K položkám struktury se dá buď přistoupit odkudkoliv, nebo jen v rámci balíčku (viz. malá/velká počáteční písmena názvu položky struktury). Není nic jako `public`, `protected` a `private`. 

Co se týče dědičnosti, tam je to v *Go* komplikované. *Go* totiž podporuje kompozici. Struktura může "absorbovat" jiné struktury. Třeba takto

```go
package main

import "fmt"

type A struct {
	Proměnná1 string
}

type B struct {
	Proměnná2 string
}

type C struct {
	A
	B
}

func main() {
	proměnná := C{A{Proměnná1: "A"}, B{Proměnná2: "B"}}
	fmt.Printf("Výsledek: [%#v]\n", proměnná)
	fmt.Printf("proměnná.Proměnná1: %s\n", proměnná.Proměnná1)
	fmt.Printf("proměnná.Proměnná2: %s\n", proměnná.Proměnná2)
}
```

Struktura `C` je kompozicí `A` a `B` a má tak položky z obou struktur. V tomto směru jde tedy *Go* dál, než Java. Ta nepodporuje dědění z více tříd, i když s pomocí default metod interfaců se dostává hodně blízko.

*Go* nemá dědičnost. Nepodporuje žádnou hierarchii tříd, nebo v případě *Go* struktur. Nejsou tu potomci a předci. Tím pádem odpadá polymorfismus. Tedy do určité míry. Nemůžeme udělat toto.

```go
package main

import "fmt"

type A struct {
	Proměnná1 string
}

type B struct {
	Proměnná2 string
}

type C struct {
	A
	B
}

func main() {
	a := A{Proměnná1: "A"}
	Funkce(a)
	c := C{A{Proměnná1: "A"}, B{Proměnná2: "B"}}
	Funkce(c)
}

func Funkce(a A) {
	fmt.Printf("a.Proměnná1 %s\n", a.Proměnná1)
}
```

Jakkoliv `C` obsahuje `Proměnná1` kompozicí z `A`, tak *Go* nepovolí jeho použití namísto `A`. *Go* ovšem podporuje interfacy. A to implicitně. Tzn. nemusíte psát, že struktura implementuje interface, jako v Javě, stačí jen interface implementovat. Jako tady.

```go
package main

import "fmt"

type I interface {
	ŘekniAhoj()
}

type A struct{}

func (a *A) ŘekniAhoj() { fmt.Println("Ahoj z A") }

type B struct{}

func (b *B) ŘekniAhoj() { fmt.Println("Ahoj z B") }

func main() {
	a := &A{}
	Funkce(a)
	b := &B{}
	Funkce(b)
}

func Funkce(i I) {
	i.ŘekniAhoj()
}
```

Všimněte si, že ani `A`, ani `B` nikde neříkají, že implementují interface `I`, ale tím, že obě mají přijímač na `ŘekniAhoj()`, jsou automaticky implementátory interfacu `I`. Proto je pak můžeme použít při volání `Funkce`. Je to sice mnohem praktičtější, nemusíte nikde psát, co všechno vaše struktura implementuje, ale na druhou stranu ani nevíte, kdo všechno implementuje váš interface.

Co se přetěžování týče, v tomto směru vás *Go* zklame. Nic takového nepodporuje. Funkce s jedním názvem používá jednu sadu parametrů. Toto například nejde.

```go
package main

import "fmt"

func main() {
	Funkce(10)
	Funkce("Ahoj")
}

func Funkce(i string) {
	fmt.Printf("i: %s\n", i)
}

func Funkce(i int) {
	fmt.Printf("i: %d\n", i)
}

func Funkce(i int, j string) {
	fmt.Printf("i: %d; j: %s\n", i, j)
}
```

*Go* nepodporuje ani genericsy, tedy obecné typy. A proto si příště ukážeme, jaké máme  v tomto směru možnosti. Jen napovím, sílu hledej v `interface{}` :-)