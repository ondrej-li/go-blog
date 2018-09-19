+++
author = "Ondrej Linek"
categories = ["začátky"]
date = 2016-08-30T20:50:18Z
description = ""
draft = false
slug = "proc-jazyk-go"
tags = ["začátky"]
title = "Proč jazyk Go"

+++

Než si povíme, proč programovat v *Go* a ne třeba v *Java*, tak odkaz na informace o jazyku *Go* na [wikipedia](https://cs.wikipedia.org/wiki/Go_(programovac%C3%AD_jazyk)). 

Po této trošce suché teorie něco praxe. Programovací jazyky se dělí do několika skupin. V rámci nich *Go* zapadá do těchto:

* kompilovaný - programy v *Go* se musí před během kompilovat, na rozdíl od jazyků typu *JavaScript*, nebo *Python*
* nativní - programy v *Go* se kompilují do nativního kódu. Programy tedy nepotřebují běhové prostředí, jako např. *Java* (ta má JVM), nebo *.net* (ten má .net framework)
* silně typovaný - program psaný v *Go* může použít pouze typy známé v době kompilace. To je rozdíl proti např. *JavaScript*, kde jsou typy známi až v době spuštění.
* strukturální - typy nemusí uvádět, s jakým typem jsou kompatibilní. Třeba *Java* jakožto nominální jazyk toto musí uvádět (`class b extends a`)
* funkcionální - v *Go* jsou funkce občany první kategorie. Funkce mohou přijímat funkce jako parametr a také je vracet. Toto např. *Java* do verze 8 neuměla.

Když tedy víme, co je *Go* zač, tak důvody proč si zvolit *Go*:

* *Go* je rychlý jazyk. Kompiluje se do nativního kódu, takže startuje bleskově. A není potřeba jej "zahřívat", běží rychle na první pokus.
* *Go* se rychle kompiluje, takže je tu značná úspora času.
* *Go* má širokou základní knihovnu. Jen málo věcí v ní nenajdete. Kromě toho tyto implementace jsou opravdu použitelné a mnoho programů si tak plně vystačí se standardní knihovnou.
* *Go* je poučené. Autoři jazyka se poučili ze zkušenosti jiných jazyků a implementovali jazyk tak, aby se vyhnuli nedostatkům jiných jazyků.
* *Go* má názor. Na rozdíl od jiných jazyků, které nechají uživatele dělat si co chce, *Go* podporuje, nebo přímo nařizuje, jak psát kód.
* *Go* je jazyk inženýrů. Při designu jazyka se přemýšlelo, což se nedá říci o mnoha jiných jazycích (radši žádný nejmenuji)
* *Go* má standardní ekosystém na který se dá spolehnout. Od formátování kódu, po kompilaci, vše dostanete v  základu.
* *Go* má jednu z nejlepších dokumentací. A plnou jasných a použitelných příkladů.
* *Go* se snadno testuje, neboť má testování zakomponováno přímo do jádra jazyka.


Má *Go* nějaké nevýhody? Má:

* *Go* je mnohých směrech omezené. Autoři se prostě rozhodli, že jsou věci, které programátoři nemohou použít. Např. generické programování, dědičnost, přetěžování apod.
* *Go* stále nevyřešilo management závislostí.

