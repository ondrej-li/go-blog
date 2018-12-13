+++
author = "Ondřej Linek"
categories = ["pokročilý"]
date = 2018-12-13T11:17:07Z
description = ""
draft = false
slug = "mage"
tags = ["pokročilý"]
title = "Mage - nástroj pro sestavování (nejen) Go programů"

+++

Sastavování (kompilace + linkování) v Go je jedna z věcí, kde vládne organizovaný chaos. Samozřejmě `go build` prostě funguje. Ale v dnešní době je potřeba dělat mnohem víc věcí, než prosté `go build`. Bohužel nějaký standardní způsob neexistuje. Jelikož jsem odkojený Un*xem, tak preferuji nástroje, které se samy nabízejí - `make`. Člověk si udělá `makefile` a je spokojený - jako [například tento](https://github.com/defectus/glutton/blob/master/makefile). Nevýhodou makefilů je, že ne všichni jim rozumí, a i ti, kdo jim rozumí, obecně neznají všechny finty. Ruku na srdce, kdo ví, co dělá na příklad toto `go fmt $$(go list ./... | grep -v /vendor/)`.

Několik chytrých hlav si tedy sedlo a řeklo, že by bylo potřeba si udělat něco, co by více vyhovovalo porgramátorům v #Go. A tak vznikl [mage](https://magefile.org/). Je to nástroj na sestavování programů, který je ovšem vlastně #Go + doplňující knihovny. Spustí se jednoduše pomocí `mage` a když jde vše dobře, tak na konci je sestavený program. Nebo něco jiného, co chcete dělat.

Nejdřív si tedy `mage` nainstalujeme.
```bash
go get -u -d github.com/magefile/mage\
cd $GOPATH/src/github.com/magefile/mage\
go run bootstrap.go
```

Zkuste si teď na konzoli zadat `mage -version` a měli byste vidět smysluplný výstup.

Nyní se ponořte do projektu, pro který byste rádi vytvořili `magefile` a zadejte `mage -init`. Po chvilce byste měli vidět v tomto adresáři `magefile.go`. Zajdet `mage -l` a měli byte vidět něco jako 

```
Targets:
  build          A build step that requires additional params, or platform specific steps for example
  clean          up after yourself
  install        A custom install step if you need your bin someplace other than go/bin
  installDeps    Manage your deps, or running package managers.
```

Mage operuje s tím, co známe v `make` jako cíle - targets. Stejně jako v `make` i zde mohou cíle záviset na jiných cílech. Samotný cíl je vlastně `POGoF` - plain old #Go function. Tedy až na to, že musí mít signaturu `func Target() error`, `func Target()`, `func Target(context.Context) error`, `func Target(context.Context)`. **Samotný magefile musí mít build tag `// +build mage`**. 

Pokud chceme, aby cíl závisel na jiném cíli, stačí jen přidat `mg.Deps(Target)` kde target je funkce jiného cíle. Samozřejmě je potřeba mít naimportovaný `import "github.com/magefile/mage/mg"`.

Pokud už tedy znáte základy `mage` není problém si ukázat, jak se v `mage` dělají běžné úkony.

### Kompilace Go
Toto je podle celkem realistický target `Build`. Tedy dobrý startovací blok, který nejen sestaví appku, ale dá jí i trochu fazónu. Jen pozor

```go
import "github.com/magefile/mage/sh"
import "time"
import "unicode"
import "golang.org/x/text/transform"
import "golang.org/x/text/unicode/norm"


var ldflags = "-s -w -X main.AUTHOR=${AUTHOR} -X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH} -X main.TAG=${TAG} -X main.BUILDTIME=${BUILDTIME}"
var goexe="go"
func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	branch, _ := sh.Output("git", "rev-parse", "--abbrev-ref", "HEAD")
	author, _ := sh.Output("git", "log", "-1", "--pretty=format:'%an'")
	version, _ := sh.Output("git", "describe", "--tags", "--abbrev=0")
   Goarch := "amd64"
   Goos := "linux"
	return map[string]string{
		"COMMIT":      normalizeString(hash),
		"BRANCH":      normalizeString(branch),
		"AUTHOR":      normalizeString(author),
		"VERSION":     normalizeString(version),
		"BUILDTIME":   time.Now().Format("2006-01-02T15:04:05Z0700"),
		"GOARCH":      Goarch,
		"GOOS":        Goos,
		"CGO_ENABLED": "0",
	}
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

func normalizeString(input string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, input)
	result = strings.Replace(result, " ", "_", -1)
	return result
}

func Build() error {
   return sh.RunWith(flagEnv(), goexe, "build", "-ldflags", ldflags, "-o", "nazev_appky", "github.com/organizace/projekt/cmd/profil")
}
```

### Fmt projektu

Každý správný projekt by měl mít všechny soubory správně formátované. Skript na to vypadá komplikovaně, ale věřím, že po chvilce to celé začne dávat smysl :-)

Hlavní část je "šumu" je vlastně jen listování souborů, které se mají poslat `gofmt`.

```go
// importy si doplňte sami :-)

// vrací cesty package pro daný projekt jako slice řetězců
func packages() ([]string, error) {
	var err error
	pkgsInit.Do(func() {
		var s string
		s, err = sh.Output(goexe, "list", "./...")
		if err != nil {
			return
		}
		pkgs = strings.Split(s, "\n")
		for i := range pkgs {
			pkgs[i] = "." + pkgs[i][pkgPrefixLen:]
		}
	})
	return pkgs, err
}

// true pokud je go ve verzi 1.11 (kvůli `go mod`)
func isGoLatest() bool {
	return strings.Contains(runtime.Version(), "1.11")
}

// Run gofmt linter
func Fmt() error {
	if !isGoLatest() {
		return nil
	}
	pkgs, err := packages()
	if err != nil {
		return err
	}
	failed := false
	first := true
	for _, pkg := range pkgs {
		files, err := filepath.Glob(filepath.Join(pkg, "*.go"))
		if err != nil {
			return nil
		}
		for _, f := range files {
			s, err := sh.Output("gofmt", "-l", f)
			if err != nil {
				fmt.Printf("ERROR: running gofmt on %q: %v\n", f, err)
				failed = true
			}
			if s != "" {
				if first {
					fmt.Println("The following files are not gofmt'ed:")
					first = false
				}
				failed = true
				fmt.Println(s)
			}
		}
	}
	if failed {
		return errors.New("improperly formatted go files")
	}
	return nil
}
```

### Testování

Bez toho to prostě nejde. Následující cíl (target) vám spustí testy a vypočítá nad nimi i pokrytí (coverage). Samozřejme `race` přepínač je nastavený (bez toho fakt nejde :-))

```go
import "github.com/magefile/mage/sh"

func buildTags() string {
	if envtags := os.Getenv("VASE_APPKA_TAGY"); envtags != "" {
		return envtags
	}
	return "none"
}

// Run tests with race detector
func TestRace() error {
	return sh.Run(goexe, "test", "-race", "-coverprofile=coverage.txt", "-covermode=atomic", "-v", "-tags", buildTags(), "./...")
}
```

### Vet 
Mnohdy jde kód zkompilovat, ale jsou v něm skryté chyby, nesmysly apod. se kterými by prohram neměl opustit dveře vačí programátorské dílny (brlohu). `go vet` je najde a nahlásí. A vrátí chybu, pokud tomu tak je (z vlastní zkušenosti - 99% toho co najde prostě potřebuje opravit).

```go
import "github.com/magefile/mage/sh"

//  Run go vet linter
func Vet() error {
	if err := sh.Run(goexe, "vet", "./..."); err != nil {
		return fmt.Errorf("error running go vet: %v", err)
	}
	return nil
}
```

### Lint

Pokud nechcete použít `go vet` nebo chcete najít i další chyby a nesrovnalosti, můžete použít i jiné lintery. Zde je příklad `golint`.

```go
import "github.com/magefile/mage/sh"

// Run golint linter
func Lint() error {
	// tato funkce je ukázána ve `Fmt`
   pkgs, err := packages()
	if err != nil {
		return err
	}
	failed := false
	for _, pkg := range pkgs {
		if _, err := sh.Exec(nil, os.Stderr, nil, "golint", pkg); err != nil {
			fmt.Printf("ERROR: running go lint on %q: %v\n", pkg, err)
			failed = true
		}
	}
	if failed {
		return errors.New("errors running golint")
	}
	return nil
}
```

### Bonus
Jako bonus je ukázka jak porovnáním času modifikace můžete zjistit, jestli se vůbec něco změnilo a tudíž jestli je potřeba něco delat.

```go
import "github.com/magefile/mage/target"

func isBuildNeeded() (bool, error) {
	return target.Dir("nazev_exace", "pkg", "cmd")
}

// Formats, Lints, Tests and Builds the application.
func Build() error {
	if needed, err := isBuildNeeded(); !needed && err == nil {
		log.Printf("Build not required.")
		return nil
   }
   ...
}
```

To je pro dnešek všechno, hodně zdaru při sestavování vašich aplikací. Na úplný závěr jeden můj [magefile.go](https://github.com/defectus/glutton/blob/master/magefile.go).
