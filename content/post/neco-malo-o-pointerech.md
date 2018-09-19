+++
author = "Ondrej Linek"
categories = ["pokročilý"]
date = 2016-11-07T16:36:40Z
description = ""
draft = false
slug = "neco-malo-o-pointerech"
tags = ["pokročilý"]
title = "Něco málo o pointerech"

+++

Dnes jsem si uvědomil jednu věc. V *Go* se to množí pointery, ale vlastně mnozí čtenáři tohoto blogu, kteří přicházejí ke *Go* ze světa JavaScriptu, Javy a dalších jazyků bez pointerů (v Česku známé také jako ukazatele), mohou být mírně zmatení, o čem to tu stále hovořím.

Než se pustím do dalších příkladů, tak jen ve zkratce pro znalce *C* a podobných jazyků, které pointery mají. V *Go* se vždy přenášejí parametry hodnotou. A to včetně hodnoty pointerů. Rozumíme si :-)

Tak a teď přímo k věci. Co je to pointer. Pointer je číslo. Číslo a typ. Jaké číslo? `Int64`. Je to ukazatel na místo v paměti, kde leží data proměnné daného typu. Ukažme si kus kódu.

```go
package main

import (
	"fmt"
)

func main() {
	var číslo int = 1
	var pointerNaČíslo *int = &číslo
	var nilPointerNaČíslo *int
	fmt.Printf("%d\n", číslo)
	fmt.Printf("%d\n", *pointerNaČíslo)
	fmt.Printf("%d\n", *nilPointerNaČíslo)
}
```

První věc, program při běhu spadne. Ale to je v pořádku. Druhá věc. Pointery v akci. Program spadne, jelikož se snažíme číst z adresy paměti `0`, což je na 99.9999% mimo rozsah nám vyhrazený. `číslo` je místo v paměti, kde je hodnota `1`. 

Paměť počítače je přístupná pomocí adres. Začínáme na adrese `0x0` a končíme na maximální adrese, buď podle operačního systému a jeho správce paměti (memory management), nebo na hranici opravdové, fyzické paměti. Abychom si mohli do paměti něco uložit a následně vybrat, potřebujeme říct, kam toto cosi ukládáme a odkud chceme toto cosi číst. Kdykoliv deklarujeme proměnnou, tak kompilátor vytvoří pro danou proměnnou místo v paměti. A jelikož kompilátor ví, jak je proměnná velká (tedy kolik bytů paměti zabírá), tak ví, na jakou adresu může zapsat další proměnnou. A této znalosti kde proměnná sídlí se říká ukazatel.

Vezměme si příklad. Chci sečíst dvě proměnné. Jak bude vypadat skutečný kód vykonaný procesorem? Bude to něco na tento způsob.

```go

C := A + B
```

v pseudo-assembleru (jazyku procesoru)

```
načti data z adresy promměné A do registru 1
načti data z adresy promměné B do registru 2
k registru 1 přičti registr 2
ulož výsledek na adresu proměnné C
```

v assembleru pak 

```
MOV rax, [pointerA]
MOV r01, [pointerB]
ADD rax, r01
MOV [pointerC], rax
```

A proč tedy vlastně pracujeme s pointery? Vždyť podle tohoto nic jiného než pointery nejsou. To proto, že proměnné se v *Go* přenášejí hodnotou. Co to znamená? Znamená to, že když zavoláte jakoukoliv funkci, tak parametry se do ní překopírují. To znamená, že parametr ve volané funkci je kopií parametru z volající funkce. To s sebou nese dvě implikace. 

1. spotřebu paměti. U čísla se to snese, ale u větších objektů je to problém (texty, obrázky, řádky databáze apod.). Každý datový objekt je potřeba kompletně klonovat a po výstupu z funkce zase odstranit, což vede i k fragmentaci paměti.
1. Mnohdy potřebujeme, aby byl výsledek volání viditelný i ve volající funkci, chceme, aby byl postranní účinek viditelný i pro ostatní, aniž bychom museli používat výsledek.

V tom případě použijeme pointery na typ. Zdůrazňuji **na typ**. Bez toho by kompilátor netušil, jak se má chovat k cílové datové struktuře. Ukažme si na postranních účincích, jaký je rozdíl mezi ukazateli a hodnotami.

```go

package main

import (
	"fmt"
)

func main() {
	var a int = 0

	fmt.Printf("před voláním: %d.\n", a)
        pomocíHodnoty(a)
	fmt.Printf("pomocí hodnoty: %d.\n", a)
	pomocíPointerů(&a)
	fmt.Printf("pomocí pointeru: %d.\n", a)
}

func pomocíPointerů(pointerNaInt *int) {
	*pointerNaInt = 1
}

func pomocíHodnoty(hodnota int) {
	hodnota = 1
}
```

Jak vidíte, i když jsme ve funkci `pomocíHodnoty` měnili hodnotu `a`, tato změna nebyla navenek patrná. Pokud jsme to samé udělali pomocí pointeru, tak byla změna patrná i ve volající funkci.

V tomto osobně vidím výhodu *Go* proti jiným jazykům. Zde si můžeme říct, jak chceme hodnotu přenést a nést s tím spojená rizika. 

Na rozdíl od *C* zde ovšem chybí pointerová aritmetika. Pointer je možné získat, ale není možné jej měnit (zatím jsem nemluvili od `unsafe`, že ano :-)).

Na tomto místě bych rád požádal, pokud máte s ukazateli problémy, tak mi napište. Je pro mě těžké odhadnout, jak složité je pochopit jejich vnitřnosti.


