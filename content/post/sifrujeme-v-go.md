+++
author = "Ondřej Linek"
categories = ["pokročilý"]
date = 2017-07-21T14:03:09Z
description = ""
draft = true
slug = "sifrujeme-v-go"
tags = ["pokročilý"]
title = "Šifrujeme v Go"

+++

Jak jsem sliboval, dnes si povíme něco šifrování dat. První otázka zní, bude nám to vůbec k něčemu? Vždyť není problém si data bezpečně přenášet pro SSL (secure socket layer - bezpečná vrstva nad sokety). To je sice pravda a mnohdy potřebujete přenést, a hlavně uložit, data mimo síťovou vrstvu. V takovém případě už to nejde jinak, než za použití kryptografie.

Z minula už víme, že máme šifry symetrické a asymetrické. Symetrické šifry jsou takové, kde k zašifrování a odšifrování použijeme stejný klíč. Obecně tedy pracují na bázi sdíleného tajemství. Problémem symetrických šifer je, jak si navzájem vyměnit onem klíč. Buď je natvrdo zadrátovaný, nebo vychází z nějakého vzorce (třeba měnícího se v čase), ale platí, že jakmile se někdo ke klíči dostane, tak nám může hezky číst vše, co jsme si zašifrovali. Takže platí, že všechny strany participující na symetrickém šifrování musí mít a používat nejen stejnou šifru, ale i stejný klíč. Představme si ale situaci, kdy k našemu tajnému systému chceme přibrat i dalšího člena, nebo nám unikl šifrovací klíč a je všem účastíkům potřeba doručit klíč nový. To je problematické, jelikož kdokoliv nás může odposlouchávat a nový klíč si stáhnout. Narozdíl od toho asymetrické šifry mají klíče dva, soukromý a veřejný. Veřejný klíč můžeme dát komukoliv a pokud jím kdokoliv cokoliv zašifrujeme tak si můžeme být jisti, že odšifrovat to může jen někdo, kdo vlastní příslušný soukromý klíč. To se obzvláště hodí v situaci, kdy si chceme vyměnit klíče pro symetrickou komunikaci. Scénář je totiž následující.

Pošlu protistraně svůj veřejný klíč.
Protistrana mi pošle svůj veřejný klíč.
Pomocí veřejného klíče protostrany zašifruju nový symetrický klíč a pošlu jej protistraně.
Protistrana jej odšifruje svým soukrommý klíčem.
Protistrana zašifruje takto přijatý symetrický klíč mým veřejným klíčem a pošle mi ho zpět.
Rozšifruju svým soukromým klíčem zprávu protistrany. Pokud se shoduje s klíčem, který jsem protistraně poslal vím, že si rozumíme a shodli jsme se na klíči.

Následně můžeme ověřit, že zvolený symetrický klíč funguje (což se běžně dělá), ale tento a další kroky již vynechám.

Asi si pokládáte otázku, proč se vůbec budeme zdržovat nějakou symetrickou šifrou a přenosem klíče. Podle všeho jsou asymetrické šifry jednodušší a bezpečnější. Odpověď je jednoduchá - symetrické šifry jsou řádově rychlejší a méně náročné na výkon procesoru. 

A teď konečně ke kódu. Nejdřív symetrické šifrování. Jako šifru si zvolíme AES, ale samozřejmě nikomu nebráním použít jakoukoliv jinou, jeho srdci blízkou.

```go
func encrypt(key, text []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    b := base64.StdEncoding.EncodeToString(text)
    ciphertext := make([]byte, aes.BlockSize+len(b))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }
    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
    return ciphertext, nil
}
```
