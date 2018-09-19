+++
author = "Ondřej Linek"
categories = ["pokročilý"]
date = 2017-07-18T17:47:50Z
description = ""
draft = false
slug = "golang-a-kryptografie"
tags = ["pokročilý"]
title = "Golang a kryptografie"

+++

Dnes si povíme něco o tom, jak funguje kryptografie v Go. Go jako (relativně) nový programovací je v této oblasti daleko a podporuje prakticky všechny možné šifry a řešení. Co je ale na osvěžující je fakt, že Go se stále snaží dělat věci jednoduše, takže na rozdíl od jiných nejmenovaných jazyků se zde nesetkáme s milionem interfaců a object factory.

Kryptografie je opravdu obsáhlý prostor, takže my si ukážeme použitelné a praktické implementace kryptografie.

Jelikož již umíme nakonfigurovat naši aplikaci, aby používala HTTPS, tak přeskočíme tuto fázi a ukážeme si, jak vytvoříme certifikát, a jak ověříme, že protistrana nám dala platný certifikát.

V druhé části si zase ukážeme, jak podepsat zprávu a jak ověřit, že zpráva je platná.

## Certifikáty a jejich ověření

Certifikáty jsou všude. Stejně jako v reálném životě i v počítačové kryptografie je certifikát doklad, kde někdo certifikuje něco. U nás je tím někdo certifikační autorita a něco tzv. certifikační subjekt. Příkladem subjektu může být identita serveru, člověka, čehokoliv, co jsme schopni zabalit do X509 subject řádky. Samozřejmě je možné použít i pole z `X509v3 extensions` a tím dále specifikovat detaily certifikátu a certifikačního subjektu.

Pokud chceme podepsat certifikát, tak potřebujeme dvě věci - certifikační autoritu a tzv. certificate signing request (žádost o podepsání certifikátu). Aby věc nebyla jednoduchá, tak i certifikační autorita potřebuje certifikát. Takže je to tak trochu Hlava 22. Naštěstí je tu možnost použít tzv. self-signed certificate, tedy certifikát podepsaný sám sebou. Takový certifikát obecně nemá valnou cenu (je to jako byste vy řekli, že jste důvěryhodný - sobě věřit můžete, ale jak přesvědčíte ostatní, aby vám to věřili?), ale důvěra někde začít musí.

V Go si takový certifikát vytvoříte velmi jednoduše. Ale nejdřív, jak to tak bývá, je potřeba [oholit jaka](https://en.wiktionary.org/wiki/yak_shaving) a to v podobě tvorby asymetrických šifrovacích klíč. V asymetrické kryptografii existuje vždy pár šifrovacích klíčů - veřejný a soukromý. Veřejný klíč můžete volně distribuovat po světě, zatímco soukromý klíč by měl za každých okolností zůstat tajný. Věc se má tak, že pomocí veřejného klíče jste schopni data zašifrovat, ale odšifrovat je může jen vlastník soukromého klíče. 

**Obecně tedy platí, že když dám entitě svůj veřejný klíč, ona jím zašifruje zprávu a já jsme následně schopný tuto zprávu svým soukromým klíčem rozšifrovat a vrátit oné entitě v čitelné formě, tak prokazuji, že jsem opravdu vlastníkem soukromého klíče k tomuto veřejnému klíči.**

Funguje to i obráceně. Zašifrujete zprávu pomocí soukromého klíče a necháte druhou stranu tuto zprávu rozšifrovat pomocí veřejného klíče. Tímto způsobem se "podepisujete" a dokazujete, že jste to vy. Radši upozorňuji - takováto zpráva není bezpečně zašifrovaná - veřejný klíč může mít kdokoliv, pokud chcete bezpečně přenést zprávu, tak šifrujte veřejným klíčem.

Takže prvně tvorba klíčů. Použijeme moderní metodu [eliptických křivek](https://cs.wikipedia.org/wiki/Kryptografie_nad_eliptick%C3%BDmi_k%C5%99ivkami), která vyžaduje kratší klíče než zastaralejší (ale stále bezpečná) [metoda RSA](https://cs.wikipedia.org/wiki/RSA). Radši znovu opakuji, při vhodné volbě délky klíče (RSA 2048bit a ECC 256bit) se jedná o **metodu neprolomitelnou se současnými výpočetními metodami**. Pro zajímavost, nejdelší prolomená ECC byla 109bit a 300 počítačů ji počítalo 3 měsíce.

```go
privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
```

S tímto klíčem pak můžeme sestavit certifikát a nechat si jej podepsat - tedy podepsat sám sebou. Podepsání certifikátu je proces, kdy certifikační autorita vypočítá [hash](https://cs.wikipedia.org/wiki/Ha%C5%A1ovac%C3%AD_funkce) certifikátu zakóduje jej svým privátním klíčem. Ověření pak spočívá v tom, že vezmete onen hash a rozkódujete jej veřejným klíčem certifikační autority. Tím ověříte, že podpis je autentický. Kromě toho si uděláte hash celého certifikátu a musí být stejný jako hash deklarovaný (zakódovaný) certifikační autoritou - tím ověříte, že certifikát nikdo neupravil.

Nejdřív se tedy uděláme "kostru" certifikátu, tedy všechny data pro výsledný certifikát.

```go
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "Our CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * 60 * time.Minute),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA: true,
	}
```

Certifikát je tedy platný jeden rok (což je dost málo na kořenový certifikát, o hierarchii certifikátů si možná povíme dále) a má náhodné sériové číslo. To by mělo být unikátní, i když se to v praxi dá velmi těžko ověřit. A teď přijde to nejlepší - podepíšeme si tento certifikát sami sebou.

```go
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
        caCert, _ := x509.ParseCertificate(derBytes)
```

Výsledkem je [DER](https://cs.wikipedia.org/wiki/Basic_Encoding_Rules#K.C3.B3dov.C3.A1n.C3.AD_DER) kódovaný certifikát, který si hned převedeme zpátky na strukturu certifikátu. Vzhledem k tomu, že certifikáty si předáváme většinou jako PEM kódované soubory, což je vlastně DER jako BASE64 zabalený mezi "BEGIN CERTIFICATE" a "END CERTIFICATE" tak si náš certifikát převedeme na PEM a můžeme ho klidně poslat do světa.

```go
	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCert.Raw,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)
```

Rekapitulace - máme ECC privátní klíč naší certifikační autority a její certifikát. Tedy vše, co potřebujeme, pokud chceme mít možnost vydávat svoje vlastní certifikáty. 

Bystří z vás se již ptají, co pak může následovat. Stejným způsobem můžeme vydávat další certifikáty, stačí jen vyměnit čtvrtý parametr volání `CreateCertificate` za podepsaný certifikát (`caCert`) a vše je vyřešeno. Popravdě, v praxi si neposíláme nepodepsané certifikáty. V praxi se posílají žádosti o podpis certifikátu. Ten obsahuje všechny pole certifikátu kromě podpisu. A takový CSR si vytvoříme snadno. Opět potřebujeme privátní klíč, ale to je asi tak vše.

```go
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: "unique identifier",
		},
		SignatureAlgorithm: x509.ECDSAWithSHA512,
		PublicKey:          privateKey.PublicKey,
	}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)

    derBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	pemBlock := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: derBytes,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)
```

Jakmile máme požadavek na certifikát tak není problém si jej dát podepsat. Překopírujeme si hodnoty z CSR do kostry certifikátu a podepíšeme privátním klíčem certifikační autority.

```go
	template := &x509.Certificate{
		Subject:               csr.Subject,
		SerialNumber:          big.NewInt(random.New(random.NewSource(time.Now().UnixNano())).Int63()),
		NotBefore:             now.Add(time.Hour * -24),
		NotAfter:              time.Hour * time.Duration(365) * 24,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:               false,
		PublicKey:          csr.PublicKey,
		SignatureAlgorithm: x509.ECDSAWithSHA512,
	}
     derBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, csr.PublicKey, privateKey)
     clientCert, _ := x509.ParseCertificate(derBytes)
```

Pokud chcete, můžete si certifikát převést do PEM podoby pomocí `pem.EncodeToMemory`, jak bylo již ukázáno výše.

Posledním krokem je ověření certifikátu. Jak jste si již přečetli, certifikát se ověří tak, že se vezme jeho podpis spolu s veřejným klíčem certifikační autority a podpis se zkontroluje. Je dobré si také ověřit časovou platnost certifkátu (`notBefore`/`notAfter`). *Go* má v sobě přímo zabudovanou podporu ověřování certifikátů, takže se nám ověření smrskne na jednu řádku:

```go
if err := caCert.CheckSignatureFrom(clientCert); err != nil {
     log.Printf("certifikát není vydaný naší certifikační autoritou")
}
```

Samozřejmě můžeme ověřovat mnohem detailněji. Na to slouží metoda `Verify`. Sama o sobě funguje bez problému, jen v našem případě bychom museli vytvořit nový pool kořenových autorit, jelikož naše autorita není uvedena v systémovém poolu. Ovšem nic nám nebrání ji přidat do systémového pool důvěryhodných certifikátů. Výhodou metody `Verify` je, že testuje i časovou platnost, platnost domény atd.

```go
	caCert, clientCert := &x509.Certificate{}, &x509.Certificate{}
	pool := x509.NewCertPool()
	pool.AddCert(caCert)
	opts := x509.VerifyOptions{
		Roots: pool,
	}

	if _, err := clientCert.Verify(opts); err != nil {
		log.Printf("certifikát není platný %+v", err)
	}
}
```

Tolik zatím k tvorbě certifikátů. V dalším článku si ukážeme, jak si podepíšeme zprávu pomocí soukromého klíče a jak na druhé straně ověříme platnost této zprávy. Tedy to co děláme při ověřování platnosti certifikátu. Stejně tak si ukážeme obrácený postup. Zakódujeme zprávu veřejným klíčem a rozšifrujeme soukromým klíčem.
