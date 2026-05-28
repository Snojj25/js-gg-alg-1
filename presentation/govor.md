# Govorni zapis — predstavitev MIM

Skupna ciljna dolžina: **5–6 minut** (≈ 30–35 sekund na diapozitiv).

> Predpostavljamo, da občinstvo problem **že pozna** — predstavitev je zato strogo
> o najinem algoritmu in pristopu, ne o definiciji problema.

| Govorec | Diapozitivi | Vsebina |
|---------|-------------|---------|
| **A** | 1 – 4 | uvod, pristop, hevristika |
| **B** | 5 – 10 | B&B, bitne strukture, rezultati, povzetek |

Format zapisa: za vsak diapozitiv (1) ciljni čas, (2) ali je celoten govor zapisan
ali samo iztočnice, (3) tekst. Vsako misel je dovoljeno preformulirati s svojimi
besedami — pomembno je, da ostanemo znotraj časa.

---

## Diapozitiv 1 — Naslov  &nbsp;·&nbsp;  Govorec **A**  &nbsp;·&nbsp;  ~20 s

**Celoten govor:**

> "Dober dan. Predstavljava vam najino rešitev naloge *Konstrukcija maksimalne
> ortogonalne CC-množice*, kar je v literaturi znano kot problem **maksimalnega
> induciranega ujemanja**. Ker problem vsi poznamo, se osredotočava na najin
> algoritem: jaz bom predstavil pristop in hevristični del rešitve,
> kolega [ime] pa nato točen algoritem in rezultate."

Premor 1 s. Prehod naprej.

---

## Diapozitiv 2 — Naš pristop  &nbsp;·&nbsp;  Govorec **A**  &nbsp;·&nbsp;  ~35 s

**Iztočnice:**

- Problem je **NP-težak** — celo na ravninskih bipartitnih grafih z največjo stopnjo 3.
- Točna rešitev za večje grafe v razumnem času torej *načelno* ni dosegljiva.
- Najini rešitvi gre za **dvofazno strategijo**:
  - **Faza 1**: hitra hevristika — vedno konča v milisekundah, vrne *spodnjo mejo*.
  - **Faza 2**: točen algoritem (razvejaj in omeji) — *izboljša* hevristiko do
    optimuma, kadar je čas dovolj.
- Hevristika je hkrati **varnostna mreža** in **začetna meja** za rezanje vej.

---

## Diapozitiv 3 — Požrešni  &nbsp;·&nbsp;  Govorec **A**  &nbsp;·&nbsp;  ~40 s

**Celoten govor:**

> "Požrešni del je preprost. Za vsako povezavo (u, v) v aktivnem grafu izračunamo
> 'ceno': **koliko vozlišč bi z izbiro porabili**. To je ravno velikost unije
> soseščin u in v.
>
> Izberemo povezavo z najmanjšo ceno, jo dodamo v rešitev S in iz aktivnega grafa
> odstranimo unijo N(u) in N(v) — torej u, v in vse njune sosede.
>
> **Intuicija**: ko porabimo malo vozlišč, jih ostane več za prihodnje povezave.
> Tako bo končna množica večja."

**Opcijska doplnitev (če čas dopusti):** "Časovna zahtevnost ene iteracije je
$O(m \cdot w)$, kjer je $w = \lceil n/64 \rceil$ — z bitnimi nizi je to v praksi
konstantno za $n$ do nekaj sto."

---

## Diapozitiv 4 — Lokalno iskanje  &nbsp;·&nbsp;  Govorec **A**  &nbsp;·&nbsp;  ~35 s

**Iztočnice:**

- Požrešni je lokalen — globalno lahko zgreši.
- **(1,2)-zamenjava**: vzamemo eno povezavo iz S, poskusimo dodati dve nazaj.
- Zakaj ravno 1 in 2? Najmanjši premik, ki **dejansko poveča** velikost rešitve.
- Postopek ponavljamo, dokler obstaja izboljšava.
- Eksperimentalno: izboljša **10–30 %** primerov pred B&B.

**Zaključek A (premostitev):** "S tem zaključujem prvo fazo. Hevristika nam vedno
da neko rešitev v nekaj milisekundah. Sedaj pa [ime] o točnem algoritmu."

---

## Diapozitiv 5 — Razvejaj in omeji  &nbsp;·&nbsp;  Govorec **B**  &nbsp;·&nbsp;  ~45 s

**Celoten govor:**

> "Hvala. Druga faza je točen algoritem **razvejaj in omeji**. Sistematično preiskujemo
> vse možne inducirane ujemanja, vendar pametno — *odrežemo* poddrevesa, ki ne morejo
> dati boljše rešitve od tiste, ki jo že imamo.
>
> V vsakem koraku izberemo vozlišče **z najmanjšo stopnjo** v aktivnem grafu in
> ustvarimo $d + 1$ podproblemov:
>
> - Za vsakega soseda w: **vključimo** povezavo (v, w) v S in odstranimo unijo
>   N(v) in N(w) iz aktivnega grafa.
> - Ali pa vozlišče v **izključimo** — ne uporabimo ga, odstranimo le njega.
>
> **Zakaj minimalna stopnja?** Število vej eksponentno raste s stopnjo. Z izbiro
> najmanjše stopnje *minimiziramo faktor razvejanja*. Hkrati najprej poskusimo
> vključitvene veje, ker povečajo trenutno velikost in s tem **takoj okrepijo
> mejo** za rezanje preostalega."

---

## Diapozitiv 6 — Redukcije + rezanje  &nbsp;·&nbsp;  Govorec **B**  &nbsp;·&nbsp;  ~50 s

**Celoten govor (slide je dolg, govor mora poudariti protiprimer):**

> "Pred razvejanjem uporabimo dve **varni redukciji**:
>
> Prva: vozlišča stopnje 0 odstranimo — ne morejo biti del nobene povezave.
>
> Druga: če sta dve vozlišči **obe stopnje 1** in sta si edina soseda, je njuna
> povezava *prisilno* v optimalni rešitvi.
>
> **Pomembno:** *klasična* redukcija obesek — vozlišče stopnje 1 — pri MIM **ni
> varna**. Hiter protiprimer: imamo pot v–u–a–b in pri u še vejo u–c–d. Če prisilno
> vključimo (v, u), izgubimo dostop do a in c — končna rešitev je velikosti 1.
> Optimum pa sta (a,b) in (c,d) — velikost 2. Zato prisilimo le **izolirane povezave**.
>
> Za rezanje vej računamo **dve zgornji meji**: število aktivnih vozlišč deljeno z 2,
> in število preostalih povezav. Tesnejšo uporabimo. Če trenutna velikost plus meja
> ne preseže najboljše doslej, vejo odrežemo."

---

## Diapozitiv 7 — Bitne strukture  &nbsp;·&nbsp;  Govorec **B**  &nbsp;·&nbsp;  ~40 s

**Iztočnice:**

- Sosedstvo predstavimo kot **bitne nize fiksne velikosti** — `[16]uint64`, do 1024 vozlišč.
- Operacije nad celotno soseščino tečejo v **O(n/64)**.
- Kritičen korak v B&B: `active.AndNot(&adj[v])` — z eno vrstico odstranimo vozlišče
  in vse njegove sosede.
- **POPCNT** je strojni ukaz, ki šteje bite v `uint64` v enem urinem ciklu — uporabimo
  ga za hitre velikosti množic.
- Vse na stacku, brez alokacij. Pospešek nad naivno listo: **eksperimentalno faktor 10–50**.

---

## Diapozitiv 8 — Rezultati  &nbsp;·&nbsp;  Govorec **B**  &nbsp;·&nbsp;  ~45 s

**Celoten govor:**

> "Poglejmo še rezultate. Vsi specifični grafi — poti, cikli, Petersenov graf —
> rešimo točno v manj kot milisekundi. Petersenov graf ima MIM velikosti 3, kar
> se sklada z literaturo.
>
> Za naključne grafe je slika bolj zanimiva. Pri **n = 50, redko** rešimo točno
> v 16 milisekundah. Pri **n = 100 srednje gostote** ali **n = 200 redko** B&B
> ne konča v 60 sekundah — vrnemo **hevristični rezultat**, ki ga imamo iz prve
> faze. Hevristika sama vedno konča v nekaj milisekundah.
>
> Zid eksaktnosti je torej odvisen tako od velikosti kot od **gostote**:
> bolj kot je graf strukturiran, dlje seže B&B."

---

## Diapozitiv 9 — Povzetek  &nbsp;·&nbsp;  Govorec **B**  &nbsp;·&nbsp;  ~30 s

**Iztočnice:**

- **Hibridni pristop**: hevristika kot varna mreža + B&B kot izboljševalec.
- **Min-degree razvejanje + redukcije + dve zgornji meji** — učinkovito rezanje.
- **Bitni nizi s POPCNT** — velik praktični pospešek.
- Iskreno o NP-težkosti: točno do ~50 vozlišč pri srednji gostoti, ali ~200 redko.
- Možne izboljšave: tesnejša LP-meja, paralelizacija, inkrementalni scoring.

---

## Diapozitiv 10 — Vprašanja  &nbsp;·&nbsp;  Govorec **B** (ali oba)  &nbsp;·&nbsp;  do konca

**Celoten govor:**

> "Hvala za pozornost. Vesela bova vprašanj."

---

## Skupni časovni proračun

| Govorec | Diapozitivi | Predviden čas |
|---------|-------------|---------------|
| A | 1, 2, 3, 4 | 20 + 35 + 40 + 35 = **130 s** ≈ 2:10 |
| B | 5, 6, 7, 8, 9, 10 | 45 + 50 + 40 + 45 + 30 + ε = **210 s** ≈ 3:30 |
| **Skupaj** | | **5:40** (pred Q&A) |

Pri vajah skrajšajte:
- Diapozitiv 6 (redukcije) lahko skrajšate, če izpustite **najmanjši del** protiprimera
  in se zanesete na sliko.
- Diapozitiv 3 (požrešni) lahko skrajšate brez "opcijske doplnitve".

Realističen cilj med oddajo: **5:00 ± 0:30 + Q&A.**

---

## Tehnična navodila za govorca

1. **Govorita počasi.** Pri 5-minutni predstavitvi je skušnjava tlačiti — raje
   izpustite stavek, kot da pohitite.
2. **Premori med diapozitivi.** 1 sekunda premora dovoli občinstvu, da se preusmeri.
3. **Pri prehodu A → B** (po diapozitivu 4): A naj eksplicitno reče "in sedaj
   [kolega/ica] o točnem algoritmu", da je preklop jasen.
4. **Pri diapozitivu 6** (protiprimer obeskove redukcije): to je **najzanimivejši
   tehnični del**. Vredno je vsak stavek povedati jasno — to bo najverjetneje vir
   vprašanj.
5. **Pri diapozitivu 8** (rezultati): ne berite vrstic tabele eno za drugo.
   Povzemite vzorec.
