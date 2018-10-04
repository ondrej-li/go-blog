+++
author = "Ondřej Linek"
categories = ["pokročilý"]
date = 2018-10-04T07:21:00Z
description = ""
draft = false
slug = "reflexe"
tags = ["pokročilý"]
title = "Reflexe v #Go"

+++

V tomto příspěvku se podíváme na to, co to je reflexe, na co se používá a hlavně jak se používá. Hned z kraje jedna moudrá hláška

### Zřejmě je vždy lepší než chytré ... a reflexe není nikdy zřejmá.

Takže varování hned na úvod - reflexe ano, ale musíte k tomu mít dobrý důvod.

O co se tedy jedná. Reflexe je schopnost zjišťovat si informace a typech a hodnotách dat za běhu programu. Ti kdo se s reflexí ještě nesetkali se možná ptají proč. Vždyť přeci kompilátor ví, jakého typu hodnoty jsou a za běhu není problém se k nim dostat. To je sice pravda, ale v okamžiku, kdy do svého kódu přidáte abstrakci, tedy interfacy, a nedej bože prázdné interfacy, už ani vy, ani kompilátor neví, jakého typu je proměnná a tím pádem nemáte možnost se ani zeptat na její hodnotu.

A v tomto okažiku přichází na scénu reflexe. Ta vám přesně umožní si za běhu programu zjistit cokoliv o jakékoliv proměnné. Tím cokoliv jsou myšleny hlavně tyto věci:

* typ (`type`) - v tomto pojetí se jedná především o jméno typu. Tedy `int`, `string`, `main.TestStruct`. Poslední jmenovaná by byla vaše structura např. `type TestStruct struct{}` v package `main`.
* druh (`kind`) - druh je obecný typ proměnné. Tedy místo jména typu, dostanete jeho druh. Je to `int`, `string`, `pointer` nebo snad `struct`?
* hodnota (`value`) - práce s hodnotou dané proměnné. Jak její čtení, tak její zápis.
* tagy (`tag`) - pomocí reflexe si můžete vytáhnout informaci o značkách.
* funkce (`method/function`) - za pomoci reflexe jste schopni se dotazovat na metody, spouštět je, nebo dokonce vytvářet funkce.

Příklady použití jsou přesně v situacích, kdy z jakéhokoliv důvodu musíte, nebo chcete podporovat neznámé typy dat a přitom s nimi skutečně pracovat - číst je a měnit je. V mém osobním případě šlo o scénář, kdy jsem chtěl načítat nastavení z proměnných prostředí (environmentu), ale nechtěl jsem otrocky volat `value := os.Getenv("XZY")`, místo toho jsem chtěl nechat samotnou strukturu nastavení definovat, jaká proměnná prostředí má být použita a také jaká je hodnota pole, pokud proměnná prostředí neexistuje. Mezi jiné použití vidím parsování dat, ukládání a načítání strukturovaných dat (třeba z databáze, disku) apod.

Pokud chceme začít pracovat s reflexí, tak na to máme v #Go dva vstupní body.

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
	}{V: "val"} // s je pointer na anonymní strukturu
	ref(s)
	fmt.Printf("struktura po spuštění: %+v\n", s)
}
```

Tento příklad ukazuje základy reflexe s tím, že se nepouští do věcí typu spoustění metod, tvorba funkcí apod. Rozhodně si myslím, že je to dobrý základ pro studování reflexe v #Go.
