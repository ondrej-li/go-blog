+++
author = "Ondřej Linek"
categories = ["pokročilý"]
date = 2017-04-02T06:35:51Z
description = ""
draft = false
slug = "rozdil-mezi-kanalem-prazdne-struktury-a-kanalem-prazdneho-interfacu"
tags = ["pokročilý"]
title = "Rozdíl mezi kanálem prázdné struktury a kanálem prázdného interfacu"

+++

Mnohdy si člověk myslí, že něco je naprosto jasné. Tedy až do okamžiku, kdy dostanete dotaz, který vám ukáže, že nejasnosti jsou na každém rohu.

Zrovna teď se mně kolega ptal, jaký je rozdíl mezi kanálem prázdné struktury a kanálem prázdného interfacu.

Odpověď je, že zásadní. 

Prázdná struktura jako typ má omezené, ale legální použití. Sama o sobě nemůže nést data (nemá žádné pole, kam se dalo cokoliv uložit), ale už její přítomnost nám umožňuje ji použít na místech, kde prostě potřebujeme "nějakou" proměnnou. Jako další místo použití vidím jako možnost zavěsit na tuto strukturu receivera (příjemce) a použít takto vytvořenou proměnnou jako parametr typu interface.

Ale zpět k otázce. Co je to kanál typu prázdný interface. Pokud definujeme kanál takto 

```go
    c := make(chan interface{})
```

tak nám kanálem může téct cokoliv. Samozřejmě před použitím je potřeba udělat "type assertion" `hodnota.(typ)`, ale jinak nám takto může přijít cokoliv. 

Oproti tomu kanál prázdné struktury, tedy něco podobného

```go
    c := make(chan struct{})
```

nám neposkytne žádnou hodnotu. Skrze takový kanál nám neprojdou žádná data, ale i na takový kanál se aplikují základní pravidla kanálu. A to hlavně pravidlo blokování. Čtení z kanálu blokuje pokud je kanál prázdný. A o to nám mnohdy jde. Nepotřebujeme si předat žádná data, ale potřebujeme druhou stranu notifikovat, že se něco stalo. Na tomto kódu to bude zřejmé.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan struct{})
	go func(c chan struct{}) {
		<-c
		fmt.Println("hotovo")
	}(c)
	fmt.Println("blokuji")
	time.Sleep(time.Millisecond)
	fmt.Println("stále blokuji")
	time.Sleep(time.Millisecond)
	close(c) // zavíráme kanál, go rutina již není blokována
	fmt.Println("už neblokuji")
	time.Sleep(time.Millisecond)   
```

Nejdřív nadefinujeme kanál a spustíme rutinu, co z něho čte. Ta se samozřejmě nemůže hnout, jelikož kanál je prázdný. Po čase kanál uzavřeme, čímž rutinu odblokujeme a rutina se dokončí - tedy vypíše ono `hotovo`.

Pokud se ptáte, proč se `hotovo` vypíše až po `už neblokuji`, tak se podívejte [sem](https://go.ondralinek.cz/2017/01/05/jak-funguji-go-rutiny/) - na naší zablokovanou go rutinu se dostane čas až po skončení I/O operace, tedy právě onoho `už neblokuji`.

Na závěr zopakuji co jsem již psal - prázdná struktura v Go nezabírá žádné místo, proto se používá tam, kde potřebujeme proměnnou, ale nemáme zájem o data.