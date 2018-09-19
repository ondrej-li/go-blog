+++
author = "Ondřej Linek"
date = 2018-04-04T07:27:24Z
description = ""
draft = true
slug = "sablony-v-go"
title = "Šablony v Go"

+++

Po dlouhé době si povíme něco nového, tentokrát o šablonách. V mnoha situacích je potřeba šablon, tedy připravených kusů textu, do kterých se doplňují/nahrazují různé informace, většinou dodané za běhu aplikace. Typickým příkladem mohou být smlouvy, dopisy, emaily, formuláře atd. atd. Vždy v šabloně najdete statický text a místa, kde se podle situace doplní dynamický text.


Například:


Vážený ==pane/paní== __Nováková__,


Dne __1.1.2018__ jsme vás informovali o úspěšném vyřízení vaší žádosti __12345__ ze dne __1.12.2017__. Vaší stavbu v obci __Liptákov__ postavíme již zítra. Pokud vám termín nevyhovuje, kontaktujte svého zástupce ==pana/paní== __Zbořila__ na čísle __765432109__.


Děkujeme a přejeme hezký den,


Technická podpora firmy Postavil a Zbořil.

Tučně vyznačené texty jsou dynamické - mění se pro každou zprávu. Zvýrazněné texty jsou také dynamické, ale mohou se řešit genericky (lze dovodit, kdo je pán a kdo je paní).

Zatímco mnoho jazyků nenabízí přímo ve svém jádru podporu pro šablony, a člověk si tedy musí stahovat knihovny, či jiné externí závislosti, _Go_ (jako vždy) podporu nabízí. V základu podporuje dva druhy šablon - HTML a prostý text.

Hlavním rozdílem mezi nimi je skutečnost, že HTML šablony přímo ve svém jádru řeší bezpečnost HTML - tedy nahrazování znaků (escapování). Pro běžné použití, jak výše zmíněné dopisy, je tedy vhodné textové šablony z balíčku `"text/template"`.

Pokud si chcete vytvořit a spustit šablonu, potřebujete obecně 3 věci:

* šablonu
* data 
* spouštěč šablon

Asi nejkratší ukázka, jak si šablonu vytvořit, doplnit data a spustit by mohl vypadat následovně:

```go
package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func main() {
	if tmplt, err := template.New("test").Parse(`toto jsou {{ . }}`); err == nil {
		buff := bytes.NewBufferString("")
		if err = tmplt.Execute(buff, "data"); err == nil {
			fmt.Println(buff.String())
			return
		}
		return
	}
}

```

Vytvoříme si novou šablonu (`toto jsou {{ . }}`), připravíme data (v našim  případě prostý řetězec `data`) a spustíme (`tmplt.Execute`). 

Jakýkoliv příklad použití je jen variací na toto téma. Takže si povíme více o tom, jak šablony fungují a co s nimi lze dělat.

Prvně přístup k datům. Na to se používá konstrukce `{{ . }}`. Ona tečka říká "použij data z daného kontextu". Kontext se mění podle toho, jak se navigujete v dodaných datech. Pokud jsou data jako v našem případě jen jedna hodnota, tak s tím toho moc neuděláme, ale pokud použije strukturu, nebo mapu, tak už je to něco jiného. Pokud máme například mapu 

```go
data := map[string]string{"klíč1":"hodnota1", "klíč2":"hodnota2"}
```

pak můžeme přistoupit ke klíči `klíč1` pomocí `{{ .klíč1 }}`. Pokud by to snad byla mapa v mapě, tak by se zápis zanořoval např. `{{ .klíč1.subklíč1}}` Platí, že ke kořenu kontextu je vždy možno se destat pomocí znaku `$`.



