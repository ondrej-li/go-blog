+++
author = "Ondřej Linek"
categories = ["pokročilý", "názory"]
date = 2017-06-21T12:59:02Z
description = ""
draft = false
slug = "generics-v-golang"
tags = ["pokročilý", "názory"]
title = "Generics v Golang"

+++

Poslední dobou je možné sledovat na poli jazyka *Go* zajímavou bitvu. Jak se blíží Go verze `2.0` (a pozor, po `1.9` může klidně přijít `1.10`), tak se začínají objevovat nápady typu "přidejme do jazyka generické typy". Nejdřív jsem i já stál na generické straně barikády, ale musím uznat, že postupem času se spíš kloním ke skupině genericsů méně, nebo třeba vůbec ne. Možná bych chtěl na toto téma rozvířit debatu, takže zde je můj pohled na věc.

V první řadě je potřeba si objasnit, jaké druhy generics jsou. Nebudu vás trápit teorií, tu počítám, že buď víte, např. z Javy, nebo tušíte. Osobně vidím dva druhy generics:

* jako první vidím možnost definovat a používat funkce, které pracují nad množinou typů. Tedy např. mapovací funkce - na vstupu je pole hodnot neznámého typu + funkce, kterou chceme nad seznamem provést a na výstupu zpracovaný výsledek jistého typu (toho samého jako vstup, nebo úplně jiný). V *Haskellu* to vypadá následovně:

```haskell
map :: (a -> b) -> [a] -> [b]
```

* druhou možnost vidím v generických datových strukturách. Tedy například mít možnost si udělat řekneme množinu (Set) jakéhokoliv typu (tedy samozřejmě s tím, že daný typ musí být porovnatelný sám se sebou). V Javě je to např. 

```java
class ValueHolder<T> {
   private T value;
   public T getValue() {return value;};
}
```

Osobně mě si myslím, že druhý příklad, generické datové struktury, je pro *Go* naprosto mimo - i ostřílený Javař má problémy napsat podobné generické struktury. Filozofie Go je jasná, pokud je možnost udržet věci jednoduché, dělejme to jednoduše i za cenu více kódu, popř. ne tolik elegantního kódu.

První varianta, tedy mít možnost si psát a používat generické funkce, je sice zajímavá, a pro lidi možná i z počátku snadno použitelná, ale přeci jen i tam je značná komplexita. A to hlavně pro člověka, který programovat začíná.

Osobně se tedy kloním k jistému kompromisu na bázi generických funkcí. Stejně, jako například `range` může iterovat přes jakýkoliv typ slicu, nebo mapy (nebo kanálu), tak implementovat několik základních funkcí generic. 

V zásadě by se jednalo o funkce:

* `map([T], func(T) U) [U]` - tedy mapuj slice typu `T`, na každý element zavolej funkci, která na vstupu přijímá `T` a vrací `U` a jako výsledek operace vrať slice typu `U` (opět, `T` a `U` mohou být ty samé typy)
* `reduce([T], func(U, T) U, U) U` - tedy na vstupu pole typu `T`, funkce, která má parametry typu `U` a `T` a nakonec iniciální hodnotu opět typu `U`. Funkce projde elementy slicu `T`, pro každý element zavolá funkci se výsledkem předchozího volání a aktuální hodnotou. Ideální pro věci typu `suma` všech hodnot apod.
* `find([T], func(T) boolean) T` - funkce má na vstupu pole typu `T` a funkci, co přijímá `T` a vrací `boolean`. Funkce vrací první `T` pro který vstupní funkce vrátila `true`.
* `filter([T], func(T) boolean)[T]` - podobná funkce jako `find` s tím, že vrací všechny instance, která vstupní funkce označila jako `true`

Podobné funkce už jsou dostupné například pro řetězce, viz. [funkce map](http://godoc.org/strings#Map). Tyto funkce by byly stejné, ale obecné pro slicy jakéhokoliv typu.

Nadále by ale nebylo možné si podobné funkce vlastnoručně programovat. Jednalo by se tedy o kompromis, a já myslím rozumný kompromis, mezi nemít žádné generics a mít plnohodnotné generics. A hlavně, i takto by *Go* zůstalo jednoduchým jazykem.
