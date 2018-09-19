+++
author = "Ondrej Linek"
categories = ["pokročilý"]
date = 2016-11-03T11:27:03Z
description = ""
draft = false
slug = "chyby-se-stacktrace-v-go"
tags = ["pokročilý"]
title = "Chyby se stacktrace v #Go"

+++

Snad každý kdo přichází do světa *Go* ze světa Javy se diví, proč chyby v *Go* nemají stacktrace. Vždyť přeci chci vědět kde se chyba stala. A mají pravdu, standardní chyby v *Go* stacktrace nemají. Ale to není proto, že by autoři *Go* byli snad sadisti, co se vyžívají ve zmatení uživatelů jazyka. Je to proto, že v *Go* jsou chyby jen interfacy. Tzn. cokoliv co implementuje interface `error` je automaticky chybou. A vzhledem k tomu, že interface `error` má jen jednu  funkci `Error() string`, tak je jeho implementace více než jednoduchá.

Většina lidí ze začátku používá běžné `errors.New()`, nebo `fmt.Errorf()`. Ty vám vrátí instanci `errorString`, který obsahuje jen a jen řetězec. Ten samozřejmě žádný stacktrace nemá. Tedy pokud byste ho tam nedali sami, nějak.

Mnohem zajímavější volbou se jeví `github.com/pkg/errors`. Tato package umí dvě základní věci. Když si vytvoříte novou chybu, tak vám k ní "přibalí" stackstrace a pokud už pracujete s jinou chybou, tak jí umí "obalit" a přidat stacktrace. 

```go
import "errors"

err:= errors.New("Nějaká chyba")
```

Toto je běžná chyba. Obalíme jí tedy, aby měla i stacktrace.

```go
import "errors"
import errors2 "github.com/pkg/errors"

err:= errors.New("Nějaká chyba")
err= errors2.Wrap(err, "Chyba se stacktracem") 
```

A můžeme si chybu hezky vytisknout.

```go
package main

import (
	"errors"
	"log"
	errors2 "github.com/pkg/errors"
)

func main() {
	err := errors.New("Nějaká chyba")
	err = errors2.Wrap(err, "Chyba se stacktracem")
	log.Printf("Máme tu chybu: %+v.", err) // `%+v` je důležité!
}
```

A takhle vypadá výstup:

```
defectus@sputnik go $ go run errors.go
2016/11/03 12:25:03 Máme tu chybu: Nějaká chyba
Chyba se stacktracem
main.main
	/Users/defectus/personal/sandpit/go/errors.go:11
runtime.main
	/usr/local/Cellar/go/1.7.1/libexec/src/runtime/proc.go:183
runtime.goexit
	/usr/local/Cellar/go/1.7.1/libexec/src/runtime/asm_amd64.s:2086.
```

