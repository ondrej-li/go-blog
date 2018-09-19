+++
author = "Ondrej Linek"
categories = ["začátky"]
date = 2016-09-30T11:00:50Z
description = ""
draft = false
slug = "bleskovka-o-ukazatelych-na-interface"
tags = ["začátky"]
title = "Bleskovka o ukazatelých na interface"

+++

Dnes si dáme první bleskovou zprávu.

### Nepoužívejte v kódu ukazatel (pointer) na interface. 

**Důvod.**

Jakákoliv struktura, která implementuje interface, jej implementuje jak jako samotná struktura, tak jako ukazatel na strukturu.

Jen si musíte dát pozor na `příjemce` (receivery). Pokud příjemce předpokládá ukazatel (pointer), tak musíte volat daný příjemce na ukazateli. V následujícím příkladu je vidět, jak příjemce implementuje funkci `Funkce` na ukazateli a tím pádem je potřeba použít k jeho volání `iface1 = &Secti{}`. Tedy dereferencovat ukazatel na strukturu `Secti`.

```go
package main

import "fmt"

type Iface interface {
	Funkce(a, b int) int
}

type Secti struct{}

func (s *Secti) Funkce(a, b int) int {
	return a + b
}

func main() {
	var iface1 Iface
	iface1 = &Secti{}
	Soucet(iface1, 1, 1)
}

func Soucet(iface Iface, a, b int) {
	fmt.Println(iface.Funkce(a, b))
}
```