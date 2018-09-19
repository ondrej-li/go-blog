+++
author = "Ondřej Linek"
categories = ["pokročilý"]
date = 2017-09-19T06:07:45Z
description = ""
draft = false
slug = "throttling-v-go"
tags = ["pokročilý"]
title = "Throttling v Go"

+++

Throttling je způsob, jak omezit počet běžících vláken. Pokud chceme, abychom nezahltili druhý systém (subsystém) nekontrolovaným množství paralelních volání, potřebujeme způsob, jak omezit počet běžících go rutin.

Určitě je mnoho způsobů, jak toho docílit, nabízím zde jeden, který se mi osvědčil.

V zásadě vycházím z toho, že zápis do kanálu blokuje, pokud je překrečena kapacita bufferu kanálu.

Tedy

```go
ch := make(chan int, 1)
ch <- 1 // OK
ch <- 2 // blokuje, dokuď někdo kanál nepřečte
```

Tento konkrétní příklad vám neprojde, samotné Go detekuje deadlock (zápis bez čtení v jediné go rutině) a radši program zabije.  

Jak se to dá využít pro throttling programu? Tak, že si uděláte kanál s bufferem o velikosti maximálního počtu paralelních zpracování a při vstupu do kritické části do kanálu zapíšete. Naopak při výstupu z kanálu přečtete (nejlépe přes `defer`). Ukázka:

```go
func throttler() {
	ch := make(chan int, 2) // nastavíme kapacitu na 2
	defer close(ch)         // good practice - zavřeme kanál před opuštění funkce

	for i := 1; i <= 10; i++ {
		ch <- 0 // zápis blokuje, pokud překračuje kapacitu kanálu
		go func(i int) {
			defer func() { <-ch }() // na závěr přečteme data z kanálu, tím ho odblokujeme
			fmt.Printf("provádím %d\n", i)
		}(i)
	}
}
```

