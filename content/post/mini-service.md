+++
author = "Ondřej Linek"
categories = ["názory"]
date = 2018-10-17T13:11:00Z
description = ""
draft = false
slug = "mini-service"
tags = ["názory"]
title = "Mini služby"

+++

Píšety monolity? Pak možná máte jeden problém, ale za to velký.

Píšete micro-servicy? Pak možná máte mnoho problémů, ale menších.

Co takhle psát miniservicy? To byste měli trochu běžných problémů. Prostě zlatá střední cesta.

Co je to miniservice. Je to něco většího než mikroservice, má řekněme 10-20 endpointů, popř. stejný počet gRPC, či jiných endpointů. Zachovává si tedy určité rysy microservice, ale přitom se z ní nestává monolit. Správa farmy miniservice je tedy přeci jen jednodušší, přitom když jedna spadne, je možné ji rychle nahodit a nespadne kvůli tomu zbytek služeb.

V čem takové miniservice psát? Samozřejmě v #Go! Miniservice v #Go bude mít něco kolem 5000 řádek kódu (proto mini), bude startovat do 100ms, paměť si nealokuje více jak 20MB RAM.

A jak takový stack může vypadat? Tak všichni víme, že největší kamarád aplikací v #Go je docker, potažmo kubernetes (a spříznění kamarádi dle uvážení). Na komunikaci a běžné služby (TLS, load balaning atd.) se stále dá řešit přes Istio apod.

Jinak řečeno, běh aplikací se výrazně neliší od microservice s přidaným bonusem, že komponent je méně a obsluhují tedy kompletnější spektrum funkcionality.

A na závěr, ano, to je přesně, co dělám :-)