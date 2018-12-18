+++
author = "Ondřej Linek"
categories = ["pokročilý"]
date = 2018-12-17T08:13:37Z
description = ""
draft = false
slug = "mage-2"
tags = ["pokročilý"]
title = "Mage - kompilace"

+++

V minulém díle jsme se podívali na nástroj [mage](https://magefile.org/) a jeho použití pro sestavování #Go programů. Tentokrát se podíváme na jednu specialitku, která asi nebyla předem zřejmá, ale rozhodně stojí za prozkoumání.

Jak z minula víte, `mage` potřebuje tzv. `magefile`. Magefile je vlastně go-skript na stereoidech. Pokud tedy chceme sestavit projekt, který používá `magefile`, tak zadáme na příkazovou řádku něco jako `mage` a máme hotovo (když je tedy definovaný defaultní cíl). Pokud bychom ovšem chtěli, aby se náš projekt takto pěkně spustil na našem CD/CI serveru, tak budeme muset přidat ještě něco - instalaci samotného `mage`. Takže náš skript by vypadal jako 

```bash
go get -u -d github.com/magefile/mage\
cd $GOPATH/src/github.com/magefile/mage\
go run bootstrap.go
mage
```

Samozřejmě jsou tu i jiné varianty - například použít docker, kde v image by byl už `mage` přítomný, nebo pokud máte dedikovaný server, tak se `mage` nainstalovat přímo na něj.

Mnohem pohodlnější je ale možnost si z `magefile` udělat spustitelný soubor. Prostě exáč. Ruku na srdce, jak často měmíte svůj ma(k|g)efile? A mezi tím, než ho zase změníte si můžete svůj magefile zkompilovat. Velmi jednoduše, stačí jen být v adresáři, kde máte svůj `magefile` a zadat 

```bash
mage -compile ./build
```

Samotný `mage -compile` ignoruje `GOOS` a `GOARCH` - pokud chcete například na MacOS vyrobit `build` pro Linux (což je mimochodem realná varianta - většina CD/CI stojí na Linuxu), tak musíte zadat něco jako 

```bash
mage -compile -goos=linux -goarch=amd64 ./build
```

kde `./build` je spustitelný soubor, který nahrazuje `mage`. Takto zkompilovaný soubor přijímá tyto parametry:

* `-l` - seznam cílů (targets)
* `-t` - timeout sestavování 
* `-v` - "ukacaný" výstup, tedy nejden chyby, ale i výpis ze `stdout` a `stderr` spuštěných programů.

poslední parametr může být cíl, který chcete explicitně spustit - třeba `clean` - záleží na vás.

Mage je v tomto směru minimalistický, takže počítám, že více se o něm rozepisovat nebudu, v každém případě si myslím, že stojí za zvážení.