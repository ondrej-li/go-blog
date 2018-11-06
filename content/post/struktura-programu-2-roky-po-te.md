+++
author = "Ondřej Linek"
categories = ["začátky"]
date = 2018-11-06T06:23:00Z
description = ""
draft = false
slug = "struktura-programu-2-roky-po-te"
tags = ["začátky"]
title = "Struktura programu 2 roky po té"

+++

Jsou to již více než dva roky, co jsem se v tomto [článku](post/struktura-programu-po-druhe/) rozepsal nad strukturou programu v #Go. A jelikož s časem rostou zkušenosti myslím, že je čas popsat, co se od té doby změnilo a jak strukturuju programy teď.

V první řadě nám #Go podporuje moduly, tedy od verze 1.11. Takže když začnu nový projekt, okamžitě dělám `go mod init`. Moduly ovšem nelze použít na `$GOPATH`, takže je potřeba se jí vyhnout. Ovšem nic vám nebrání pracovat jako dříve, být na `$GOPATH` a používat na správu závislostí něco jiného - třeba `dep`.

V rootu projektu jako takového nemám žádné zdrojové kódy. Je tam ale `README.md`, `makefile` (popř. jiný buildovací skript), `CONTRIBUTING.md` a licence. Mimochodem licence je zásadní dokument už proto, že žijeme v právním státě. Pokud by snad chtěl někdo váš projekt použít, upravit apod. tak se bez něj neobejdete. Pokud nepotřebujete jinak doporučuji svobodnou licenci, jejich seznam najdete na [na stránkách FSF](https://www.gnu.org/licenses/license-list.html).

Při návrhu rozmístění souborů po adresářích se zkuste vžít do pozice člověka, který se jde na váš program podívat a nemá šajn, jak jsou zdrojové kódy organizované. Adresáře by měly být jasně pojmenované, třeba místo kde uchováváte dockerfily by se měl jmenovat `docker`, `dockerfiles` apod. To samé pro balíčkovací systémy, dokumentaci apod.

Samotný #Go kód je dobré dát do dvou hlavních adresářů. 

* `cmd` - tam umístnit všechen kód, který produkuje spustitelné soubory. I když bude jen jeden, tak doporučuji udělat adresář, se jménem spustitelného souboru. Např. `/cmd/super_appka/main.go`.

* `pkg` - tady umístnit soubory, které jsou v zásadě "knihovnou" dané aplikace. Tedy věci, které samy o sobě nejsou aplikací, ale které poskytují služby, bez nichž by vaše aplikace nefungovala. Jistou výhodou je, že se pak opradu tyto soubory tváři jako knihovna vaší aplikace. Už když je importujete, tak import vypadá jako `import "github.com/moje_organizace/super_appka/pkg/super_kalkulacka"`. Stále se osobně vyhýbám vytvářené velkého množství package - je to čistě pragmatická volba - nechci skončit jako #Java a mít stovky package adresářů ve více jak 3 úrovních.

Pokud už použijete buildovací nástroj, pak opět doporučuji nazývat jednotlivé akce rozumě - např. `make all`, `make test`, `make deploy`.

Pokud chcete vidět projekt, který je podle mě celkem slušným musterem tak zkuste datadog agenta: https://github.com/DataDog/datadog-agent.