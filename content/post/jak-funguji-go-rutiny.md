+++
author = "Ondrej Linek"
categories = ["pokročilý"]
date = 2017-01-05T11:17:19Z
description = ""
draft = false
slug = "jak-funguji-go-rutiny"
tags = ["pokročilý"]
title = "Jak fungují Go rutiny"

+++

Dnes si povíme něco o *Go* rutinách. Všichni víme, že se používají, pokud chceme pracovat paralelně. Všichni určitě i víme, že *Go* rutiny nejsou thready (vlákna). Co tedy jsou *Go* rutiny a proč se používají?

Prvně, proč *Go* nepoužívá vlákna? *Go* vlákna používá, ale nepoužívá je 1:1 ke *Go* rutinám. Vlákno, jakkoliv se používají často, je má svoji režii. Režii na úrovni operačního systému (je to OS, kdo vlákna vytváří), režii na úrovni správy (je potřeba vlákna časovat a přepínat) a režii na úrovni paměti a registrů procesoru (je potřeba při přepínání vlákna ukládat prakticky všechny registry procesoru). Jak je zřejmé, pokud bychom chtěli mít tisíce Go rutin, tak by režie OS byla natolik veliká, že by se systém zhroutil.

Jak tedy *Go* řeší správu *Go* rutin? *Go* používá klasickou event loop (smyčku událostí), jak ji známe např. z *Node.js*. Jedná se tedy o kooperativní zpracování paralelismu. Výhodou je, že díky minimální režii je možné mít tisíce *Go* rutin a systém stále funguje. Nevýhodou je, že pokud nějaká *Go* rutina provádí dlouhé zpracování, tak doslova "vyhladoví" ostatní *Go* rutiny.

Asi si pokládáte otázku, jak tedy onen kooperativní multithreading funguje. Vždyť nikde ve vašem kódu neprovádíte nic, čím byste dali najevo, že je možné předat řízení další *Go* rutině. *Go* jde na to fikaně, a celé přepínání před námi schovává. V zásadě jsou tyto body předání kontroly následující:

* použití kanálu `<-` (čtení a zápis) 
* použití příkazu `go` (jakmile ho zavoláte, tak se předá řízení následující *Go* rutině, i když toto není 100% pravidlo)
* jakékoliv blokující volání OS (čtení/zápis z disku, síťové operace apod.)
* volání garbage collectoru
* explicitním volání [runtime.Gosched()](https://golang.org/pkg/runtime/#Gosched). Toto je mimochodem elegantní řešení v situaci, kdy *Go* rutina provádí náročné zpracování a nechce blokovat ostatní *Go* rutiny.

Pokud vás toto téma zajímá více, tak se neváhejte podívat [na tento dokument](http://www1.cs.columbia.edu/~aho/cs6998/reports/12-12-11_DeshpandeSponslerWeiss_GO.pdf) (v angličtině), který tuto problematiku detailně rozebírá.

Vyzbrojeni těmito znalostmi vám nic nebrání v tom, abyste si udělali program s tisícovkami *Go* rutin, který stále funguje jako švýcarské hodinky.
