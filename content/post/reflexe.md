+++
author = "Ondřej Linek"
categories = ["pokročilý"]
date = 2018-10-04T07:21:00Z
description = "Reflexe - co, na co a jak"
draft = false
slug = "reflexe"
tags = ["pokročilý"]
title = "Reflexe v #Go"

+++

V tomto příspěvku se podíváme na to, co to je reflexe, na co se používá a hlavně jak se používá. Hned z kraje jedna moudrá hláška

Zřejmě je vždy lepší než chytré ... a reflexe není nikdy zřejmá.

Takže varování hned na úvod - reflexe ano, ale musíte k tomu mít dobrý důvod.

O co se tedy jedná. Reflexe je schopnost zjišťovat si informace a typech a hodnotách dat ze běhu programu. Ti kdo se s reflexí ještě nesetkali so možná ptají proč. Vždyť přeci kompilátor ví, jakého typu hodnoty jsou a za běhu není problém se k nim dostat. To je sice pravda, ale v okamžíku, kdy do svého kódu přidáte abstrakci, tedy interfacy, a nedej bože prázdné interfacy, už ani vy, ani kompilítor neví, jakého typu je proměnná a tím pádem nemáte možnost se ani zeptat na její hodnotu.

A v tomto okažiku přichází na scénu reflexe. Ta vám přesně umožní si za běho prohramu zjistit cokoliv o jakékoliv proměnné. Tím cokoliv jsou myšleny hlavně tyto věci:

* typ (`type`) - v tomto pojetí se jedná především o jméno typu. Tedy `int`, `string`, `main.TestStruct`. Poslední jmenovaná by byla vaše structura např. `type TestStruct struct{}` v package `main`.
* druh (`kind`) - druh je obecné typ proměnné. Tedy místo jména typu, dostanete jeho druh. Je to `int`, `string`, `pointer` nebo snad `struct`?
* hodnota (`value`) - práce s hodnotou dané proměnné. Jak její čtení, tak její zápis.
* tagy (`tag`) - pomocí reflexe si můžete vytáhnout informaci o značkách.

Příklady použití jsou přesně v situacích, kdy u jakéhokoliv důvodu musíte, nebo chcete podporovat neznámé typy dat a přitom s nimi skutečně pracovat - číst je a měnit je. V mém osobním případě o šlo o scénář, kdy jsou chtěl načítat nastavení z proměnných prostředí (environmentu), ale nechtěl jsem otrocky volat `os.Getenv("XZY")`, ale nechat samotnou strukturu nastavení definovat, jaká proměnná prostředí má být použíta a také jaká je hodnota pole, pokud proměnná prostředí neexistuje. Mezi jiné použítí vidím parsování dat, ukládání a načítání strukturovaných dat (třeba z databáze, disku) apod.

Pokud chcem začít pracovat s reflexí, tak na to máme v #Go dva vstupní body.

* `reflect.TypeOf(interface{})` - vrací `interface reflect.Type{...}` který použijete k získání informací o typu proměnné.
* `reflext.ValueOf(intercace{}` - vrací `struct reflect.Value{...}` a ta slouží ke čtení a zápisu hodnoty proměnné.

takový drobný příklad by mohl vypadat takto

```go
package main

import "reflect"
import "fmt"

func ref(v interface{}) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)
	fmt.Println("Typ je:", typ.Name())
	fmt.Println("Je hodnota pointer:", val.Kind() == reflect.Ptr)
	fmt.Println("Je hodnota struct:", val.Elem().Kind() == reflect.Struct) // Elem() dereferencuje pointer
	fmt.Println("Hodnota prvního pole struktury:", val.Elem().Field(0).String())
	for i := 0; i < val.Elem().NumField(); i++ {
		fmt.Println("Nalezeno pole:", typ.Elem().Field(i).Name, "typu:", val.Elem().Field(i).Type().Name(), "s tagem:", typ.Elem().Field(i).Tag.Get("tag"))
	}
	// nastavíme novou hodnotu prvního pole 
	val.Elem().Field(0).SetString("newval")
}

func main() {
	s := &struct {
		V string `tag:"test"`
	}{V: "val"}
	ref(s)
	fmt.Printf("struktura po spuštění: %+v\n", s)
}
```

Tento příklad ukazuje základy reflexe s tím, že se nepouští do věcí typu spoustění metod, tvorba funkcí apod. Rozhodně si myslím, že je to dobrý základ pro studování reflexe v #Go.
