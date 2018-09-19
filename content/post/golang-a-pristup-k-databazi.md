+++
author = "Ondrej Linek"
categories = ["začátky", "pokročilý"]
date = 2016-09-27T05:31:55Z
description = ""
draft = false
slug = "golang-a-pristup-k-databazi"
tags = ["začátky", "pokročilý"]
title = "#Golang a přístup k databázi"

+++

V dnešní době nad neexistuje aplikace, která by nepřistupovala do databáze. Samozřejmě si můžeme ukládat data do souboru a mnohdy je to i efektivnější řešení, ale databáze nám luxusně umožní data ukládat a načítat a nestarat se přitom o detaily. V dnešní době rozdělujeme databáze do 2 proudů.

* SQL databáze (nebo také relační databáze) - tedy databáze, kde jsou data ukládána jako relace mezi entitami. Většina z nich také garantuje ACID (atomičnost operací, konzistenci stav, integritu dat a durabilitu uložení). Mezi nejznámější patři [MySql](https://www.mysql.com/), [PostgreSQL](https://www.postgresql.org/), [Oracle](https://www.oracle.com/database/index.html), MS SQL Server atd.
* NoSQL databáze - tedy databáze, které nevyžadují  striktní schéma (ad hoc ukládání dat) místo ACID používají [BASE](https://en.wikipedia.org/wiki/Eventual_consistency) (pouze základní dostupnost, měnící se (měkký) stav a eventuální konzistenci) - což samo osobě je slovní hříčka - prostě kyselost a zásaditost :-). Představily jsou [MongoDB](https://www.mongodb.com/), [Redis](http://redis.io/), [Apache Cassandra](http://cassandra.apache.org/), [ElesticSearch](https://www.elastic.co/products/elasticsearch) apod.

Dnes si povíme, jak přistupovat k SQL databázím. Už proto, že *Go* tyto databáze standardizuje, na rozdíl od NoSQL, kde standard není.

V první řadě, pro každou databázi potřebujeme ovladač. Pokud se budeme chtít připojit k PostgreSQL, tak si musíme stáhnout ovladač. A jelikož ovladač zasahuje do inicializace, ale sám o sobě nic nepřináší, použijeme onen magický *import pro vedlejší účinky*.

```go
_ "github.com/lib/pq"
```

Tím pádem se ovladač může zaregistrovat a `database/sql` o něm ví. A tím pádem je pak možno udělat něco jako toto a připojit se k databázi:

```go
url := "postgres://username:password@hostname/db_name"
db, _ := sql.Open("postgres", dataSourceName)
```

Jakmile máme připojení do databáze (v reálu počítám, si ověříte, že spojení dopadlo dobře, např. `db.Ping()`), můžeme začít úřadovat a pokládat dotazy. 

Nejdříve si ukážeme, jak načíst celý dotaz:

```go
rows, _ := db.Query("SELECT sloupec FROM tabulka WHERE a = $1", "b")
defer rows.Close()

sloupce := make(string, 0)
for rows.Next() {
   sloupec := ""
   err := rows.Scan(&sloupec)
   sloupce = append(sloupce, sloupec)
}
```


Nejdříve získáme řádky pomocí dotazu (`db.Query`). Následně si preventivně zaregistrujeme funkci `rows.Close()` jako defer, abysme za všech okolností řádky zabřely. Následně si projdeme všechny řádky a uložíme si je do slicu stringů. S tím už pak můžeme dělat cokoliv.

`database/sql` umožňuje načíst také jen jednu řádku pomocí `db.QueryRow`, ale osobně to nedoporučuji, jelikož to padá, když to nenajde ani jednu řádku. Pokud tedy chcete jednu řádku, pak doporučuji postup podobný prvnímu příkladu:

```go
for rows.Next() {
   sloupec := ""
   err := rows.Scan(&sloupec)
   sloupce = append(sloupce, sloupec)
   break
}
```

Prostě po první řádce smyčku ukončíte. 

V [minulém díle](https://go.ondralinek.cz/2016/09/24/testujeme-v-golang-cast-2/) jsme se podívali, jak testovat a demonstroval jsem tam, jak přistupovat k datům pomocí interfaců. Za tím si stojím a tím pádem nebudu opakovat kód, jak si takový interface udělat.

Mnozí z vás budou asi hned mít otázku: "Jak *Go* pracuju s `NULL`". A je to správná otázka. *Go* umožňuje mít nil, tedy pointer(ukazatel) na nic, ale co když nepoužíváme ukazatele? Popravdě funkce `Scan` ani neumí pracovat s ukazateli a ukazatele. Pokud tedy např. máte v databázi `VARCHAR`, tak musíte poslat jako hodnotu k uložení ukazatel na `string`. Pokud je ale hodnota sloupce `NULL`, pak *Go* nahlásí chybu. Pro tyto účely je tu ovšem `sql.NullString` apod. (`NullBool` ...). Vnitřně má si totiž uchovává dvě hodnoty. `Value` tedy string hodnoty a `Valid`, což je příznak, jestli byla hodnota nastavena. Pokud je `false`, pak byla hodnota `NULL`. Dá se použít jako přímá náhrada `string`, tedy např. jako v našem případě.

```go
sloupce := make(sql.NullString, 0)
for rows.Next() {
   sloupec := sql.NullString{}
   err := rows.Scan(&sloupec)
   if template.Valid {
         sloupce = append(sloupce, sloupec.Value)
   }
}
```

Možná vás v SQL dotazu zaujalo, že je tam podmínka `WHERE a = $1` a ptali jste se, co to je ono `$1`. Nuže je to jen zástupný znak. PostgreSQL tento znak nahradí řetězcem `"b"`, ale v praxi to bude nejspíš nějaká proměnná. Důvod je prostý, pokud bychom začali SQL řetězce sčítat (např. pomocí `ftm.Sprintf`, tak bychom si akorát zadělali na [SQL code injection](https://cs.wikipedia.org/wiki/SQL_injection).

Teď když víme o tom, jak si data z databáze načíst, si můžeme ukázat, jak data do databáze vložit.

```go
if result, err := db.Exec("INSERT INTO tabulka (a) VALUES($1)", "b"); err == nil {
   fmt.Printf("Vložena %d řádka s ID  %d", result.RowsAffected(), result.LastInsertId())
}
```

Opět jednoduché. V příkladu předpokládám, že `tabulka` má ID automaticky generované.

A na závěr ještě ukázka transakcí. [ACID databáze](https://cs.wikipedia.org/wiki/Datab%C3%A1zov%C3%A1_transakce) umožňují provádět několik dotazů v rámci transakcí, kdy se buď provede všechno, nebo nic. Typický příklad je ono okřídlené převod na účtu. Potřebuji v něm připsat peníze na účet příjemce a odepsat z účtu zasílatele. A buď se uloží obě změny, nebo žádná. Můžete se zeptat, proč by se neprovedli obě, inu, někdy dojde k chybě, nebo závadě (třeba nám vypnou proud, rozbije se disk apod.) a to, co se mělo provést společně se nám rozpojí. *Go* SQL má podporu transakcí a vypadá to následovně.

```go
částka, zákazník1, zákazník2 := 100, 1, 2
tx, _ := db.Begin()
if _, err := tx.Exec ("UPDATE ucet SET zustatek = zustatek - $1 WHERE zakaznik = $2", částka, zákazník1); err != nil {
   tx.Rollback()
   return err
}
if _, err := tx.Exec ("UPDATE ucet SET zustatek = zustatek + $1 WHERE zakaznik = $2", částka, zákazník2); err != nil {
   tx.Rollback()
   return err
}
return tx.Commit()
```

Jedná se o lehký úvod, ale asi máte představu, jak taková transakce vypadá.

Tímto prozatím uzavírám kapitolu seznámení se s SQL databázemi a těším se příště.
