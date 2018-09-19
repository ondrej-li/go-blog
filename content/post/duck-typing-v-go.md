+++
author = "Ondrej Linek"
categories = ["začátky"]
date = 2017-01-25T15:36:16Z
description = ""
draft = false
slug = "duck-typing-v-go"
tags = ["začátky"]
title = "Duck typing v #GO"

+++

Možná už jste někdy slyšeli termín "`duck typing`". Pro ty co náhodou ne, tak jen krátce uvedu, že se jedná o aplikaci známého pořekadla "kváká to jako kachna, vypadá to jako kachna, tak je to kachna". 

Představte si tyto dva typy.

```go
type Člověk struct {
   Jméno string
}

type Mazlíček struct {
   Jméno string
}
```

Je jasné, že se jedné o dva typy. `Člověk` a `Mazlíček`. Oba jsou ovšem stejné - mají (sice jen jedno) pole - `Jméno`.

Mnoho jazyků vám neumožní napsat 
```go
   člověk, mazlíček := Člověk{}, Mazlíček {}
   člověk = mazlíček
```

A to ani *Go*. Co ale můžete udělat je prosté přetypování.

```go
    člověk, mazlíček := Člověk{Jméno: "Ondra"}, Mazlíček{Jméno: "Ká"}
    člověk = Člověk(mazlíček)
```

Podle mě se jedná o rozumný kompromis mezi implicitním duck typingem (kde nemusí člověk dělat nic) a mezi nemožností (kde musíme ručně přiřadit pole mezi proměnnými).

Zajímavé na tom je, že toto se vztahuje i na interfacy. V *Go* není potřeba explicitně označovat typy jako implementátory interfacu. Dokavaď typ implementuje všechny funkce interfacu, tak je jeho implentátorem. Co se tím bohužel nese je, že pokud změníte interface, tak může celá řada typů přestat být implementátory tohoto interfacu a naopak. Pokud změníme funkci přestaneme také být implementátory interface - a popravdě, mnohdy lidí ani neví, že neúmyslně implementují nějaký interface.

Vemte tento příklad.

```go

type Implementátor struct {}

func (i *Implementátor) String() string {
    return "Implementátor"
} 
```

Náš typ `Implementátor` se, aniž bychom to někde uváděli, stal implementátorem interfacu `fmt.Stringer`. A teď co čert nechtěl, jste si uvědomili, že potřebujete vracet i `error`. Provedete patřičnou změnu.

```go

type Implementátor struct {}

func (i *Implementátor) String() (string, error) {
    return "Implementátor", nil
} 
```

A už interface `fmt.Stringer` neimplementujete. 