+++
author = "Ondrej Linek"
categories = ["začátky", "pokročilý"]
date = 2016-09-14T17:33:13Z
description = ""
draft = false
slug = "konfigurace-programu-v-golang"
tags = ["začátky", "pokročilý"]
title = "Konfigurace programu v #Golang"

+++

Dnes si povíme něco o tom, jak si nakonfigurovat takový program v *Go*.

Konfigurací zde myslím způsob, jak předat programu parametry, aby mohl běžet nezávisle. V zásadě máte několik možností a pojďme si o nich povědět detailněji.

1. parametry příkazové řádky
1. pomocí proměnných prostředí  
1. databáze
1. specializované programy

Tento výčet není konečný, ale prozatím si s ním vystačíme. 

### Příkazová řádka

Podívejme se na první možnost, parametry na příkazové řádce.

Pro tento případ má *Go* package `flags`. Možná už jste si sami všimli, že na rozdíl od jazyků odvozených od *C* (*C*, *C++*, *Java* apod.) tu funkce `main` nepřijímá pole parametrů. Pro připomenutí, takto vypadá `main` v *Go*.

```go
func main() {
	...
}
```

Pokud se chcete dostat k parametrům, tak si pomůžete mnohem elegantněji - posuďte sami.

```go
package main

import (
	"flag"
	"fmt"
)

var intValue = flag.Int("int-param", 0, "číselný parametr")
var stringValue = flag.String("string-param", "", "textový parametr")
var boolValue = flag.Bool("bool-param", false, "logický parametr")

func main() {
	flag.Parse()
	fmt.Printf("int-param:[%d], string-param:[%s], bool-param:[%t]", *intValue, *stringValue, *boolValue)
}
```

Tok programu je jednoduchý. Nejdříve inicializujeme  všechny parametry, které chceme použít a následně zavoláme `flag.Parse()`. To způsobí, že se do všech pointerů zapíšou skutečné hodnoty. Pokud během parsování nastane chyba, tak se spustí `flag.Usage`, což je proměnné funkce, co vytiskne nápovědu. Tu můžete vytisknout i vy ve vašem kódu, eventuálně tuto proměnnou nastavit na svoji funkci.

### Proměnné prostředí

[12-factor apps](https://12factor.net/config) nám říká, že konfigurace by měla být v proměnných prostředí. V *Go* není naprosto žádný problém se k nim dostat. Třeba takto:

```go
package main

import (
	"fmt"
	"os"
)

var stringValue = os.Getenv("STRING_VALUE")

func main() {
	fmt.Printf("String value:[%s]", stringValue)
}
```

### Zookeeper

Ukládání konfigurace do databáze je vcelku známý koncept, ale nebudu ho zde demonstrovat, neboť na práci s databází se chystám později. Místo toho přidám jen ukázku, jak si přečíst konfiguraci ze [Zookeepera](https://zookeeper.apache.org/).

```go
package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

const zkConnectionString = "127.0.0.1:2181"

func main() {
	if conn, _, err := zk.Connect(zkConnectionString, time.Second); err != nil {
		log.Panicf("Nemohu se připojit k Zookeeperovi: %s", err)
	} else {
		if data, _, err := conn.Get("/app/dev/db/url"); err != nil {
			log.Panicf("Nemohu číst z /app/dev/db/url: %s", err)
		} else {
			fmt.Printf("Načtena hodnota [%s]", string(data))
		}
	}
}
```

Samozřejmě tento příklad předpokládá, že vám Zookeeper běží na lokálním počítači, a že v něm jsou nahraná data, ale pro ilustraci to myslím stačí. Výhodou Zookeepera oproti ostatním je jistě to, že Zookeeper je schopný vás notifikovat o změnách v konfiguraci, což otevírá dveře možnosti v runtimu reflektovat změny konfigurace. Což je myslím snem snad každého správce high-availibility systému.