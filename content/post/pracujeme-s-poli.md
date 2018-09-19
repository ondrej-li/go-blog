+++
author = "Ondrej Linek"
categories = ["začátky", "pokročilý"]
date = 2016-10-17T11:27:00Z
description = ""
draft = false
slug = "pracujeme-s-poli"
tags = ["začátky", "pokročilý"]
title = "Pracujeme s poli"

+++

Práce s poli patří k základům prakticky všeho. Pokud mám něčeho více jak jedna, tak se z toho stává pole.
Běžně rozlišujeme tři typy polí:

* obyčejné pole. V *Go* se deklaruje jako `[n]T` např.
```go
   pole := [10]int
```
Všimněte si, že `pole` je typu `[10]int`. Tzn. velikost pole se nedá změnit. Pokud tedy předem víte, kolik
prvků potřebujete a nechcete aby se to dalo změnit, pak je pole právě pro vás.

* seznam. V *Go* se mu říká `slice`. Deklaruje se jako `[]T` např. `[]int`. Pokud chcete, můžete si u něj nastavit jeho
počáteční kapacitu. Ovšem nikoliv při deklaraci, ale při inicializaci. Tedy např.
```go
   seznam := make([]int, 10) // počáteční velikost pole je 10

```

* množina. Mnoho jazyků přímo podporuje množiny, anlicky `set`. *Go* mezi takové jazyky nepatří, a tak musíte vzít za vděk 
nějakou knihovnou , která tuto funkcionalitu nabízí. Jen pro přesnost, množiny nejsou setříděné a neumožňují duplikáty.

Teď nejběžnějši operace nad poli, nebo slicy. Prvně si pole vytvoříme.

```go
   pole := [10]int // prázdné pole
   pole2 := [2]int{1, 2} // pole vyplněné prvky 1 a 2
   proměnná := pole2[0] // proměnná má hodnotu 1, indexujeme od 0
   for index, číslo := range pole2 { // vypíšeme obsah pole
      fmt.Printf("%d:%d", index, číslo)
   }
   
```

Pokud začneme používat slicy, tak si užijeme mnohem více legrace. Popravdě, na pole v *Go* moc nenarazíte.
Můžeme samozřejmě polemizovat, co je šetrnější k paměti a co se rychlosti přístupu k prvkům týče, ale mikrooptimalizace
tohoto typu už jsou přeci jen v dnešní době přežité.

```go
	slice := []int{} // prázdný slice, musí se inicializovat
	fmt.Printf("slice: %#v\n", slice)
	slice2 := []int{1, 2}              // slice inicializovaný na hodnoty 1 a 2
	proměnná := slice2[0]              // proměnná má hodnotu 1, indexuje se od 0
	for index, číslo := range slice2 { // vypíšeme obsah pole
		fmt.Printf("%d:%d\n", index, číslo)
	}
```

Do teď to bylo copy-paste, teď ale postupujeme na další level.

Přidání položky do slicu
```go
	slice3 := append(slice2, 3) // přidáme prvek ke slicu, původní slice2 zůstal jaký byl!
	fmt.Printf("slice3: %#v\n", slice3)
	proměnná = slice3[2] // proměnná má hodnotu 3
	fmt.Printf("proměnná: %#v\n", proměnná)
```

Vytvoření sub-slicu, tedy části původního slicu od začátku do indexu.
```go
	slice4 := slice3[:2] // slice4 obsahuje 1, 2 (od začátku slice3, do indexu 2 včetně).
	fmt.Printf("slice4: %#v\n", slice4)
```

Vytvoření dalšího sub-slicu od indexu do konce.
```go
	slice5 := slice3[1:] // slice 5 obsahuje 2, 3 (od indexu 1 do konce)
	fmt.Printf("slice5: %#v\n", slice5)
```

Pokud chceme odstranit element ze slicu, tak na to máme tuto fintu
```go
	slice6 := append(slice3[:1], slice3[2:]...) // slice6 obsahuje 1 a 3 (2 jsme vymazali, pokud se to tak dá říct)
	fmt.Printf("slice6: %#v\n", slice6)
```

Pokud chceme naopak chceme to slicu přidat, tak se to dělá následovně.
```go
	slice7 := append(slice6[:1], append([]int{2}, slice6[1:]...)...)
	fmt.Printf("slice7: %#v\n", slice7)   
```

Možná se mýlím, ale řekl, bych že více se toho s poli asi dělat nedá. Nebo snad ano? Pokud od něčem víte, dejte mi vědět!

