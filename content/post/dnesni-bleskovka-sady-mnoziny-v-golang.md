+++
author = "Ondrej Linek"
categories = ["začátky"]
date = 2017-01-10T08:48:24Z
description = ""
draft = false
slug = "dnesni-bleskovka-sady-mnoziny-v-golang"
tags = ["začátky"]
title = "Dnešní bleskovka - sady (množiny) v Golang"

+++

Dnes opět jen telegraficky. *Go* nepodporuje sady (množiny) a tak tu máme drobnou obezličku, jak si takovou sadu vytvořit. Využijeme toho, že mapy mají unikátní klíče a jsou netříděné. Že to zní povědomě? Ano! Vždyť to je přece množina.

Takže než abychom si mazali ruce nějakým bastlením kódu si uděláme např. toto:

```go
package main

import (
	"fmt"
)

func main() {
	množinaIntů := map[int]struct{}{}
	množinaIntů[2] = struct{}{}
	množinaIntů[1] = struct{}{}
	množinaIntů[3] = struct{}{}
	množinaIntů[4] = struct{}{}
	množinaIntů[1] = struct{}{}
	// jen si mapu prosvištíme, ať víme co v ní je
        for key := range množinaIntů {
		fmt.Println("Položka množiny:", key)
	}
	// zjištění jestli množina obsahuje položku
	if _, ok := množinaIntů[1]; ok {
		fmt.Println("Množina obsahuje položku 1.")
	}
}
```

Padl tu dotaz, proč je použit typ `map[int]struct{}`. Tedy mapu s klíčem `int` a položkami typu prázdná struktura. Důvod je prostý, prázdná struktura zabírá (logicky) v paměti 0 bajtů. 

To je pro dnešek vše :-)


