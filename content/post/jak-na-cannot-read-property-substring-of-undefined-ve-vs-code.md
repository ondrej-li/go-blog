+++
author = "Ondrej Linek"
categories = ["začátky"]
date = 2017-01-10T10:38:06Z
description = ""
draft = false
slug = "jak-na-cannot-read-property-substring-of-undefined-ve-vs-code"
tags = ["začátky"]
title = "Jak na `Cannot read property 'substring' of undefined` ve VS Code"

+++

Pokud jako já používáte [Visual Studio Code](https://code.visualstudio.com/) pro práci s projekty v *Go* tak jsme možná narazili na problém při formátování zdrojových souborů (standardně `cmd + alt + L` na MacOs). Místo vytouženého zformátování se zobrazí 

`Cannot read property 'substring' of undefined`

Nuže problém je v chybějící prázdné řádce na konci souboru. Stačí tedy jen přidat na konec souboru jednu prázdnou řádku a problém je vyřešen.

Např. 

```go
package main

func main () {
....
} // následuje prázdná řádka

```

Více na [Githubu](https://github.com/Microsoft/vscode-go/issues/690).