+++
author = "Ondřej Linek"
categories = ["pokročilý"]
date = 2017-03-14T19:47:03Z
description = ""
draft = false
slug = "predavani-poli-slicu-hodnotou-2"
tags = ["pokročilý"]
title = "Předávání polí (sliců) hodnotou a na co si dát pozor"

+++

Když jsme se bavili o slicech tak nezazněla jedna podstatná informace. Jak vlastně vypadá taková struktura slicu.

Vytvořme si tedy slice. 

```go
package main

import "fmt"

func main() {
	slice := []int{}
	fmt.Println("délka:", len(slice))
	fmt.Println("maximální délka:", cap(slice))
}
```

Délka i maximální délka jsou shodné, tedy 0. 

Teď si k danému slicu přidejme jedno položku a zkusme to znovu.

```go
package main

import "fmt"

func main() {
	slice := []int{}
	slice = append(slice, 0)
	fmt.Println("délka:", len(slice))
	fmt.Println("maximální délka:", cap(slice))
}
```

Tentokrát je délka 1 a maximální délka je 2. Co to znamená? `append` nám nejen přidal do pole jednu položku, ale také nám nechal ve slicu "fuku" pro jedno další číslo. Tzn. pokud přidáme další položku, tak nám `append` nemusí zvětšovat pole, prostě přidá položku na konec již alokovaného pole.

Teď to možná zní zmateně - bavíme se o slicech a najednou pole? 

Je to tím, že slice je vlastně pole na stereoidech. Interně se jedná o pointer na pole (tedy prní element pole), a informace o délce a kapacitě slicu. Pokud tedy předáváme slice hodnotou, tak se nepřenáší celé pole, ale přenese se pouze ona "hlavička" pole. Podívejme se na tento příklad.

```go
package main

import "fmt"

func main() {
	slice1 := []int{1, 2, 3}
	slice2 := slice1
	fmt.Println("délka1:", len(slice1), ";délka2", len(slice2))
	fmt.Println("maximální délka1:", cap(slice1), ";maximální délka2:", cap(slice2))
	slice2[0] = 42
	fmt.Println("element1:", slice1[0], "element2:", slice2[0])
}
```

Vytvořili jsme jeden slice a následně jeho kopii. Pokud by se jednalo o přenášení hodnotou, pak by nastavení hodnoty jednoho slicu nemělo ovlivnit druhý slice. Jak je ale zřejmé, došlo k ovlivnění! Podívejme se na druhý příklad.

```go
package main

import "fmt"

func main() {
	slice1 := []int{1, 2, 3}
	slice2 := slice1
	fmt.Println("délka1:", len(slice1), ";délka2", len(slice2))
	fmt.Println("maximální délka1:", cap(slice1), ";maximální délka2:", cap(slice2))
	slice2 = append(slice2, 4, 5)
	slice2[0] = 42
	fmt.Println("element1:", slice1[0], "element2:", slice2[0])
}
```

Tentokrát změna jednoho slicu neovlivnila druhý. A otázka zní proč? Je v tom nějaký systém? Je! Tímto zlomem je překročení kapacity. Jakmile je vyžadováno přerozdělení pole, ve kterém jsou data uložena, tak `Go` interně volá `copy` a vytvoří nové pole. Původní pole je stále alokované pro původní proměnnou.

Co z toho plyna za ponaučení? Hned několik:

* memory management je složitá věc. Pokud víte, že váš slice poroste, pomozte `Go` tím, že místo nic neříkajícího `slice:= []int{}` uděláte `slice := make([]int, 0, 100)`. Nastavením kapacity zamezíte tomu, aby `Go` muselo při každém volání `append` alokovat nové pole dat.
* i když předáváte pole hodnotou, změna dat může být viditelná "za hranicí" funkce, ve které slice měníte.
* toto se vztahuje i na práce s pásmy. Pokud napíšeme `slice2 := slice1[2:]` tak `slice2` stále ukazuje do pole `slice1`. Změny v něm budou viditelné i v našem novém `slice2` (a naopak).