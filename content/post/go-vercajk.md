+++
author = "Ondřej Linek"
categories = ["začátky"]
date = 2017-11-11T06:52:47Z
description = ""
draft = false
slug = "go-vercajk"
tags = ["začátky"]
title = "Go vercajk"

+++

Jelikož stále nejsem schopný dát dohromady návod na šifrování v Go (skončil jsem u nástřelu, ale nedaří se mi ho dokončit), tak se posunu na další kapitolu a tou je váš/náš vercajk při práci s Go.

Každý jazyk má svoji sadu nářadí (tooling) se kterou se pracuje. Například Java. Pokud chci pracovat v Javě, tak je skoro jisté, že jako IDE použiji Ideu, nebo Eclipse, programy budu sestavovat pomocí mavenu, nebo gradlu, testovat budu např. JUnit, nejspíš se vykašlu na všechny možný experimenty a použiji Spring na všechno. Co je ale hlavní, máte relativně velký výběr a nic vás nebrání zkoušet jiné nástroje a postupy.

Go je ale jazyk s názorem a tak je mnohem vyhraněnější. Takže v čem a jak bude vývoj nejspíš probíhat?

#### IDE
Zde máte asi největší výběr. Osobně programuji ve [Visual studio Code](https://code.visualstudio.com/download), který má skvělý plugin pro práci s Go. Velmi dobře se dá pracovat i v [Goland](https://www.jetbrains.com/go/download/). Pro ty, co preferuji VIm je tu [vim-go](https://github.com/fatih/vim-go). Celkově se dá říct, že ohromnou výhodou Go je, že podpora editorů je řešena přes CLI nástroje. Tedy např. automatické doplňování řeší [gocode](https://github.com/nsf/gocode), importy [go-imports](https://godoc.org/golang.org/x/tools/cmd/goimports), formátování `go fmt` a podobně. Díky tomu je snadné tyto nástroje integrovat do libovolného editoru.

#### Testování

Toto téma jsme [již probírali](https://go.ondralinek.cz/2016/09/21/testujeme-v-golang-cast-1/), celkově tedy stačí říct, že Go má podporu pro psaní a spouštění testů přímo v jazyce, takže si stačí udělat funkce `Testxxxxx(t *testing.T)` spustit `go test ./...` a víte přesně na čem jste.

#### Dokumentace
Na psaní tu máme `go doc`. Každá veřejná funkce, typ a package by měl(a) mít svůj komentář začínající jménem funkce a tyto komentáře pak `go doc` "prosviští" a vystaví jako webový server. Tím se liší od jiných podobných nástrojů, které generují statické `HTML` stránky. 

#### Závislosti
Zde bych řekl je nejslabší místo celého Go ekosystému. Na jedné straně je velmi jednoduché si pomocí `go get` stáhnout nějakou závislost, na druhé straně pokud chcete mít závislost v nějaké verzi (a to z důvodu reprodukovatelnosti sestavení chcete), tak musíte použít nějaký nástroj mimo standardní balík Go - možností je několik, ale zdá se, že `go dep` by mohlo mít nakonec největší úspěch. Mezi jiným je to pak `go vendor` a třeba i `glide`.

#### Sestavení
Pokud chcete svůj program sestavit, tak je asi nejlepší použít dodávané `go build`. Samozřejmě, že v rámci sestavení budete chtít udělat i jiné kroky a celkem zajímavě nejvíc se ta tímto účelem používá klasické `make` - ano to `make`, co už používáme od dob jazyka `C`. Samozřejmě vám nic nebrání použít shellové skripty, nebo jakýkoliv jiný "sestavovač".

#### Generování kódu

Určitě víte, že Go nemá a nepodporuje makra (pro znalce `C` tedy nemá fázi prekompilace). Pokud tedy chcete změnit, nebo dokonce generovat kód před tím, než ho zkompilujete, tak je nejlépe se obrátit na `go generate`. Samozřejmě lze použít jiné nástroje, ale toto je standard. 

Jeden hezký příklad, jak použít generování kódu je náhrada za generics. Podívejte se jak tento problém řeší [například `genny`](https://github.com/cheekybits/genny) 

#### Interoperabilita

Toto je asi vhodné pro samostatný článek, ale jak říká klasik, žádný program není ostrov a je tedy potřeba se bavit i s jiným systémy. Tady bych řekl záleží hlavně na tom, jaké ty ostatní systémy jsou. Pokud jsou naštěstí napsané v Go, tak máte skoro vyhráno a můžete použít `gob` - serializaci Go struktur do binárního formátu. Pokud ne, tak je několik možností.
   
1. protobuf - nebo ještě lépe gRPC. Go má plnou podporu pro použití protobufferů, takže není problém si nechat vygenerovat kód z vašich `*.proto` souborů a následně použít dodávané serializéry/deserializéry.
2. json - můžete si myslet co chcete, ale json tu prostě je a jeho podpora je v Go velmi silná. Snad jen drobná rada, pokdu nevíte, v jakém formátu vám data přijdou, tak `map[string]interface{}` je nejlepší volba. 
3. xml - stále se používá, takže můžete využít Go podporu pro XML. Mapování na elementy, atributy a namespace se provádí pomocí strukturálních tagů

```go
element string `xml:"http://example.com/ns element"`
```

4. swagger - swagger je dobrý v tom, že je schopný vám z definice vygenerovat kostru vašeho API. Bohužel míru individualizace je celkem omezená, takže počítejte s routerem od projektu `gorilla` apod. Zase na druhou stranu vám nic nebrání v implementaci pomocí vašeho oblíbeného routeru.