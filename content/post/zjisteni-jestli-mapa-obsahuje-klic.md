+++
author = "Ondrej Linek"
categories = ["pokročilý"]
date = 2016-10-20T08:05:29Z
description = ""
draft = false
slug = "zjisteni-jestli-mapa-obsahuje-klic"
tags = ["pokročilý"]
title = "Zjištění, jestli mapa obsahuje klíč"

+++

Každý potřebujeme čas od času zjistit, jestli je v mapě položka pod daným klíčem. V *Go* se to dělá velmi jednoduše.

```go
	mapa := map[string]string{
		"položka1": "hodnota1",
	}
	if hodnota1, ok := mapa["položka1"]; ok {
		fmt.Printf("Hurá, mapa má klíč 'položka1': [%s]\n", hodnota1)
	}
	if hodnota2, ok := mapa["položka2"]; !ok {
		fmt.Printf("Bohužel, klíč 'položka2' není v mapě. Posuďte sami: [%#v], [%s]", mapa, hodnota2)
	}
```

Mimochodem, podobný postup se používá i v situaci, kdy si chcete ověřit typ proměnné. 

```go
	cokoliv := interface{}(1)
	if _, ok := cokoliv.(int); ok {
		fmt.Println("Sláva, cokoliv je int!")
	}
	if _, ok := cokoliv.(string); ok {
		fmt.Println("Sláva, cokoliv je string!")
	} else {
		fmt.Println("Bohužel cokoliv není string")
	}
```