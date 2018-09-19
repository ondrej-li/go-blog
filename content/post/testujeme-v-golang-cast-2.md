+++
author = "Ondrej Linek"
categories = ["pokročilý"]
date = 2016-09-24T11:56:57Z
description = ""
draft = false
slug = "testujeme-v-golang-cast-2"
tags = ["pokročilý"]
title = "Testujeme v #Golang (část 2)"

+++

V první části povídání o testování jsme si ukázali vcelku jednoduchý test a spustili jsme si jej. 

Jedním z požadavků na testování je, aby testy doběhly co nejdříve. Je to proto, že chceme mít smyčku zpětné vazby co nejkratší. Nikdo nestojí o to, aby čekal více jak minutu na informaci, že je všechno v pořádku. Ale pak se naskytuje otázka - co v případě, že test závisí na něčem pomalém. Třeba databázi. Nebo volání externí služby přes HTTP. V tom případě tu máme celou plejádu možností a podíváme se na ně.

První, co zkušený programátor udělá je mockování. Jelikož, doufejme opodstatněně, předpokládá, že věci typu databáze prostě "fungují" (a nebo je to problém jiného, pokud ne), tak není potřeba ji skutečně testovat. Stačí, když víme, co z databáze dostaneme a přesvědčit váš kód, aby to akceptoval. V našem případě. Řekněme, že nše databáze má seznam uživatelů a my máme funkci, která nahraje jednoho podle `ID`. Tedy naprosto běžný scénář. My samozřejmě víme, jak vypadá takový záznam o uživateli a také nechceme během testu pokládat dotaz do databáze, protože to bude minimálně 10ms. A to neberu v potaz, že by tam musel očekávaný záznam být (ne ne, podle toho, co testujem). Ukažme si třeba, jak by daný kód mohl vypadat.

```go
package main

import (
   "sql"
   "errors"
)

type UserRecord struct {
   Id uint64
   Name string
   Password string
}

type UserManager interface {
   FindUserByUserName (user string) (*UserRecord, error)
}

type DB struct {
   DB *sql.DB
}

type Authenticator interface {
   Authenticate (username string, password) (bool, error)
}

type DefaultAuthenticator struct {
   db UserManager
}

func (db *DB) FindUserByUserName (user string) (*UserRecord, error) {
   rows := db.Query("SELECT username, password FROM user WHERE id = $1", id)
   defer rows.Close()
   for rows.Next() {
      user := UserRecord{Id:id}
      err := rows.Scan (&user.Id, &user.Name, &user.Password)
      return &user, err
   }
   return nil, errors.New("No user found")
}

func(a *DefaultAuthenticator) Authenticate (username string, password string) (bool, error) {
   userRecord, err := a.db.FindUserByUserName(username)
   if err != nil {
      return false, err
   }
   if userRecord.Password == password {
      return true, nil
   }
   return false, nil
}
```

Tento kód nám nejen implementuje načítání uživatele z databáze (doma nezkoušet, v praxi hesla do databáze neukládáme, ani pod nátlakem!), ale i ověřování hesla. Možná si řeknete, co to tam ten chlap dělá s interfacama? Interfacy zde používáme proto, že nám dovolují abstrakci. Zatímco pro běh programu použijeme `DB` jakožto implementátora interfacu `UserManager`, v testech si jej nahradíme naší implementací, která dělá to, co po ní chceme my. Tedy například nepůjde do databáze. Zrovna v tomto si pomůžeme knihovnou `github.com/stretchr/testify/mock`. A když ji spřáhneme s knihovnou `github.com/stretchr/testify/assert` ze stejné stáje, jsme neporazitelní. No posuďte sami:

```go
package main

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

type MockUserManager struct {
	mock.Mock
}

func (m *MockUserManager) FindUserByUserName (user string) (*UserRecord, error) {
	args := m.Called(user)
	return args.Get(0).(*UserRecord), args.Error(1)
}

func TestDefaultAuthenticator_Authenticate(t *testing.T) {
   mum := MockUserManager{}
   userRecord := &UserRecord{Id:uint64(1), Name: "user", Password: "password"}
   da := DefaultAuthenticator{db: mus}
   mum.On("FindUserByUserName", "user").Return(userRecord, nil)
   result, err := da.Authenticate("user", "password")
   mum.AssertExpectations(t)
   assert.NoError(t, err)
   assert.True(t, result)
}
```

Zde jsme si vytvořili "mock" interface `UserManager` se jménem `MockUserManager` a "implementovali" funkci `FindUserByUserName`. Ovšem naše verze nedělá nic jiného, než je porovná, že je opravdu volána s parametrem "user" a pokud ano, tak vrátí náš `userRecord` a žádnou chybu. Pokud test spustíme, tak si ověříme, že `Authenticate` funguje, ale nemusíme kvůli tomu jít do datábáze, takže test nám projdu kdykoliv (i když třeba vůbec žádnou databázi nemáme) a hlavně, opravdu rychle.

Možná si vzpomenete, jak jsme říkal, že testování je součást vývoje, která se nedodá až nakonec vývoje? Toto je onen případ. Tento kód je záměrně napsaný, aby byl dobře testovatelný. Samozřejmě bych mohl načítání uživatele dát do globální funkce, ale pak bych měl problém tuto funkci měnit pro účely testu.

Další možností, abychom urychlili testy je jejich paralelní provádění. V našich počítačích máme několik procesorů, tak proč je pořádně nevyužít?! 

```go
func TestParalelně(t *testing.T) {
  t.Parallel()
  // a test, jak je běžné
  ...
}
```

Pokud byste chtěli ovlivnit, v kolika procesech testy poběží, tak použijte parametr `go test -parallel 2`.

Některé testy ovšem chcete spouštět jen když jste například na síti, nebo pokud máte čas. K tomu slouží `t.Skip()`. V testu můžeme mít například něco podobného:

```go
func Test_Skip(t *testing.T) {
   if os.Getenv("INTEGRATION") == "" {
      t.Skip("test přeskakuji")
   }
}

```

Tak v poslední části povídání o testech si ukážeme, jak spouštět podtesty a další "vychytávky". Ovšem před tím se podíváme na něco mnohem zajímavější - práci s databázi. Zde jsem toto téma nakousli, ale ještě je potřeba jej pořádně dovyprávět.