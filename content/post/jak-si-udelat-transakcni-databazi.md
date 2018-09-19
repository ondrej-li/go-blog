+++
author = "Ondřej Linek"
categories = ["pokročilý"]
date = 2017-04-13T12:29:57Z
description = ""
draft = false
slug = "jak-si-udelat-transakcni-databazi"
tags = ["pokročilý"]
title = "Jak si udělat transakční databázi"

+++

Určitě jste si toho všimli - Go má (podle mě) velmi dobrou podporu relačních databází. A to včetně věcí typu transakce. Pokud chcete ovšem transparentně přistupovat k databázi tak, aby vám bylo jedno, jestli dotaz provádíte v rámci transakce, nebo mimo ni, tak máte problém. Standardní API totiž rozlišuje mezi přímým dotazem do databáze a mezi dotazem v rámci transakce.

Abychom tuto "podružnost" eliminovali, a také si ukázali jak si udělat abstrakční vrstvu nad konkrétní implementací, tak si napíšeme zaobalení (wrapper) databáze. Bude mít stejné funkce na provádění základních dotazů jako má standardní Go API. 

Radši hned upozorním, že tento wrapper nebude 1:1 s Go API, tzn. na místech kde se striktně vyžaduje DB struktura z Go API se neobejdete od rozbalení (unwrapping) původního Go API.

První věc je definovat si kontrakt. To uděláme pomocí definice interface. Třeba takto 

```go
type DatabaseAccessor interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}
```

Jelikož ještě potřebuje API pro práci s transakcemi, tak si ho také vytvoříme jako interface. Např.

```go
type Transactionable interface {
	DatabaseAccessor
	NewTx() error
	Begin() (*sql.Tx, error)
	Rollback()
	Commit()
}
```

Tento interface velmi připomíná API, které Go nabízí pro přístup k SQL databázím s benefitem možnosti začínat a ukončovat transakce. Nad tímto interfacem si vytvoříme svoji implementaci přístup databázi.

Teď si vytvoříme strukturu, která bude zapouzdřovat jak databázi, tak transakci.

```go
type DB struct {
	*sql.DB
	*sql.Tx
}
```

Struktura si drží odkaz na databázi a na transakci. Myšlenka je taková, že podle toho, jestli odkaz na transakci je, nebo není nil, se rozhodneme, jestli operaci (definovanou jako `DatabaseAccessor`) provedeme přímo na databázi, nebo na transakci.

V praxi to vypadá následně:

```go
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if db.Tx != nil {
		rows, err := db.Tx.Query(query, args...)
		return rows, err
	}
	rows, err := db.DB.Query(query, args...)
	return rows, err
}

func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	if db.Tx != nil {
		row := db.Tx.QueryRow(query, args...)
		return row
	}
	row := db.DB.QueryRow(query, args...)
	return row
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if db.Tx != nil {
		rows, err := db.Tx.Exec(query, args...)
		return rows, err
	}
	result, err := db.DB.Exec(query, args...)
	return result, err
}
```

To hlavní je zřejmé, před každou operací se podíváme, jestli pracujeme nad transakcí, nebo ne.

Pro samotné správu transakcí máme tento kód. Ten vytvoří novou transakci, a provede commit, nebo rollback. Přesně co potřebujeme.

```go
func (db *DB) Rollback() {
	if db.Tx != nil {
		db.Tx.Rollback()
		db.Tx = nil
		return
	}
}

func (db *DB) Commit() {
	if db.Tx != nil {
		db.Tx.Commit()
		db.Tx = nil
		return
	}
}

func (db *DB) NewTx() error {
	var err error
	if db.Tx != nil {
		db.Tx, err = db.DB.Begin()
	}
	return err
}
```

Jediný problém nastává, pokud bychom chtěli toto řešení použít v konkurenčně, tedy v situaci, kdy naše API může být použito více go rutinami. Proto si vytvoříme ještě jednu, možná ne zrovna pěknou, funkci, která vrací novou strukturu DB s nastartovanou transakcí. 

```go
func NewTxOnDB(db Transactionable) (*DB, error) {
	if _, ok := db.(*DB); !ok {
		return nil, errors.Errorf("can't convert Transactionable to *DB")
	}
	tx, err := db.(*DB).DB.Begin()
	newDB := &DB{
		DB:    db.(*DB).DB,
		Tx:    tx,
	}
	return newDB, err
}
```

I tak si musíme dávat pozor, abychom k této nově vytvořené struktuře nepřistupovali z více go rutin. 