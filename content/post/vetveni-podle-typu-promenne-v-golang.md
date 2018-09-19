+++
author = "Ondrej Linek"
categories = ["pokročilý"]
date = 2016-10-13T05:34:53Z
description = ""
draft = false
slug = "vetveni-podle-typu-promenne-v-golang"
tags = ["pokročilý"]
title = "Větvení podle typu proměnné v #Golang"

+++

Někdy je potřeba si větvit program podle toho, jakého je proměnná typu. Typický příklad je `interface{}`. Jelikož se jedná o prázdný interface, tak jím může být cokoliv. Co když ale chcete jinak pracovat s `int` a jinak se `string`?

Následující příklad vám to ukáže.

```go
proměnná := inteface{}(1) // přetypujeme int 1 na interface{}
switch proměnná.(type) { // podle typu proměnné
case int:
	fmt.Println("Je to int!") 
default:
	printString("Je to záhada")     
}
```

Jak vidno, prosté.