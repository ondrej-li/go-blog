+++
author = "Ondřej Linek"
categories = ["pokročilý"]
date = 2017-05-12T07:35:51Z
description = ""
draft = false
slug = "vyctove-typy"
tags = ["pokročilý"]
title = "Výčtové typy"

+++

Určitě to víte - Go nepodporuje výčtové typy - tedy `enum`. Go umí numerické sekvence pomocí `iota`, ale to se opravdu nerovná enumu. Když jsem se nedávno na projektu vrátil ke staršímu kódu (který jsem napsal), tak jsem si uvědomil, že v jedné chvíli jsem vlastně použil něco, co může fungovat jako náhrada za výčty. No posuďte sami.

```go
package main

import (
	"fmt"
)

type AnswerEnum struct {
	string
}

var (
	Yes AnswerEnum = AnswerEnum{"Yes"}
	No  AnswerEnum = AnswerEnum{"No"}
)

func main() {
	cond := []AnswerEnum{AnswerEnum{"Yes"}, AnswerEnum{"No"}}
	for i, answer := range cond {
		switch answer {
		case Yes:
			fmt.Println("na indexu", i, "je hodnota", Yes.string)
		case No:
			fmt.Println("na indexu", i, "je hodnota", No.string)
		}
	}
}
```

To co se zde děje je zjevné. Nejdříve si vytvoříme nový typ typu `AnswerEnum`, který je kompozicí typu `string`. Jelikož lze stringy mezi sebou porovnávat, tak je pak možné použít tyto hodnoty v bloku `switch`.

Takový enum je tedy porovnatelný s jinou hodnotou, ale přitom se jedná o samostatný typ - tedy nelze poslat jako parametr jiný typ (enum) byť se stejnou hodnotou.