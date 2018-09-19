+++
author = "Ondrej Linek"
categories = ["začátky"]
date = 2016-11-02T08:37:54Z
description = ""
draft = false
slug = "if-else-if-else-v-go"
tags = ["začátky"]
title = "if-else-if-else v #Go"

+++

Ve většině jazyků z rodiny *C* se větvění provádí pomocí `if-else`. Klasicky to funguje (např. v *Go*) takto.

```go
if cosi > dalsiCosi {
   log.Println("cosi je větší než dalsiCosi.")
} else {
   log.Println("cos je menší než dalsiCosi.")
}
```

Takto to funguje dobře. Co když ale chceme naspat něco většího? 

```go
if cosi > dalsiCosi && time.Weekday() == time.Monday {
   log.Println("Je pondělí a cosi je větší než dalsiCosi.")
} else if cosi > dalsiCosi && time.Weekday() == time.Tuesday  {
   log.Println("Je úterý a cosi je větší než dalsiCosi.")   
} else {
   log.Println("Není pondělí a úterý, nebo cosi je menší než dalsiCosi, nebo tak něco v tom smyslu.")
}
```

*Go* ovšem nabízí tuto variantu namísto sestavy `if-else-if-else`.

```go
switch {
   case cosi > dalsiCosi && time.Weekday() == time.Monday :
      log.Println("Je pondělí a cosi je větší než dalsiCosi.")
   case cosi > dalsiCosi && time.Weekday() == time.Tuesday :
      log.Println("Je úterý a cosi je větší než dalsiCosi.")
   default:
      log.Println("Něco úplně jiného.")
}
```

Osobně preferuji tento zápis, ale jako u všeho, volba je na vás.
