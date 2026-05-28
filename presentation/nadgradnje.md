# Nadgradnje na podlagi recenzij

Dokument opisuje **štiri ciljne izboljšave**, ki sva jih izvedla v `solver.go` po prejemu
dveh recenzij (Snoj–Zmazek interno ocenjevanje in poročilo "Slabi vhodni podatki").
Najprej globlja razlaga z intuicijo in primeri (za pripravo in vprašanja), nato kratek
**dvominutni govor** za predstavitev.

> Predpostavljamo, da osnovni algoritem (požrešno + lokalno iskanje + B&B z bitnimi nizi)
> je že razložen. To poglavje opisuje **le** dodatke nad tem temeljem.

---

## Kaj so recenzenti odkrili

Obe recenziji se ujemata v glavni ugotovitvi:

> *Hevristika je dovolj dobra — pogosto najde optimum. Težava je v B&B fazi: njene
> zgornje meje (`|V|/2`, `|E|`) so preohlapne, da bi **dokazale** optimalnost.
> Solver zato pogosto vrne pravilno rešitev, a porabi celotnih 300 s in zastavica
> `timed_out` je `true` — pri ocenjevanju, kjer šteje čas, je to maksimalna kazen.*

Konkretne šibke točke:

1. **Fiksna meja `maxN = 1024`** — graf s 1025 vozlišči je solver z `panic`om zrušil.
2. **Brez razgradnje po komponentah** — 200 disjunktnih $K_4$ obravnavamo kot eno
   eksponentno drevo, namesto 200 trivialnih.
3. **Preohlapna globalna zgornja meja `|V|/2`** — pri 341 disjunktnih trikotnikih
   $|V|/2 = 511$, dejanski optimum pa je 341.
4. **Brez certifikata optimalnosti** — tudi ko hevristika že najde optimum, B&B faza
   tega ne zna prepoznati in pretolče celotni časovni proračun.

Vsaka od štirih nadgradenj naslavlja **eno** od teh točk.

---

## Nadgradnja 1 — Dinamični bitni niz

### Težava

V prvotni različici:

```go
const maxN = 1024
const maxWords = maxN / 64           // = 16
type Bitset [maxWords]uint64         // fiksno polje 16 × uint64
var adj [maxN]Bitset
```

Branje vhoda kliče `adj[i].Set(j)` za $i, j \in [0, n)$. Pri $n = 1025$ Go izvede
`b[1024 >> 6] = b[16]` na polju dolžine 16 → `panic: index out of range`.

### Rešitev

`Bitset` postane **slice** dinamične dolžine, alociran enkrat pri branju vhoda:

```go
type Bitset []uint64

var nWords int   // = (n + 63) / 64, izračunan ob branju n

func newBitset() Bitset { return make(Bitset, nWords) }
```

### Subtilnost: globoke kopije

Pri fiksnem polju je `newActive := active` v Go **globoka** kopija (kopira vseh 16
besed). Pri slice-u pa je to **plitva** kopija (kopira le glavo slica, ne podatkov).
Brez popravka bi se kloni in originali povozili.

Zato sva dodala eksplicitno `Clone()`:

```go
func (b Bitset) Clone() Bitset {
    r := make(Bitset, len(b))
    copy(r, b)
    return r
}
```

In na vseh mestih, kjer je prej veljala semantika globokega kopiranja
(`newActive := active`, `combined = adj[u]`, `remove = adj[u]`), klicava `Clone()`.

### Tradeoff

- **Cena:** majhna alokacija na vsak klon. Velikost slica je še vedno
  $(n + 63) / 64$ besed, isto kot prej, le da je dinamično. Ko se v B&B ustvari
  podproblem, je `Clone()` `O(w)` z `w = \lceil n/64 \rceil`.
- **Dobiček:** **ni več trde meje** glede na velikost vhoda. Solver sprejme poljuben $n$.

### Učinek

| Test | Pred | Po |
|------|------|-----|
| `01_nad_mejo_n1025.in` | `panic: index out of range` | **VALID**, rešitev 341 |

Ta test je prej dal **nič točk** (solver ni izpisal datoteke). Sedaj vrne veljaven
rezultat. Ker je vhod cikel dolžine 1025, je teoretični optimum
$\lfloor 1025/3 \rfloor = 341$ in se ujema z najino rešitvijo.

---

## Nadgradnja 2 — Razgradnja po povezanih komponentah

### Ključno opažanje

**Inducirano ujemanje je aditivno preko komponent**: če sta $G$ in $H$ disjunktna grafa,
potem velja

$$ \mathrm{MIM}(G \cup H) = \mathrm{MIM}(G) + \mathrm{MIM}(H). $$

Razlog: nobena povezava iz $G$ ne more biti v konfliktu (delitev krajišča ali sosednost
krajišča) s povezavo iz $H$, ker med grafoma ni nobene povezave.

### Posledica za B&B

Brez razgradnje B&B išče v prostoru, ki je **kartezijski produkt** podproblemov vsake
komponente — eksponenten v številu komponent. Z razgradnjo je celotni čas
**vsota**, ne produkt.

Primer: 200 disjunktnih $K_4$. Vsak $K_4$ ima MIM = 1, kar je trivialno najti. Brez
razgradnje B&B išče v prostoru s ~$5^{200}$ kandidatnimi konfiguracijami in nikoli ne
zaključi. Z razgradnjo: 200 majhnih problemov, vsak rešen v mikrosekundah, skupna
rešitev 200.

### Implementacija

```go
func findComponents(active Bitset) []Bitset {
    // BFS/DFS po bitnem sosednem grafu.
    // Vrne seznam Bitsetov, vsak predstavlja eno komponento.
}
```

V `main()` zamenjava: namesto enega klica `solve(active)` čez celoten graf, za vsako
komponento:

```go
for _, comp := range comps {
    edges := solveComponent(comp)
    allEdges = append(allEdges, edges...)
}
```

`solveComponent` neodvisno požene hevristiko + B&B na komponenti in vrne njene
povezave.

### Učinek

| Test | Število komponent | Pred | Po |
|------|------------------|------|-----|
| `07_disjunktne_k4_x200.in` | 200 | 200, TLE 5 s | 200, **112 ms** |
| `08_twins_6x10.in` | 3 | 12, 149 ms | 12, **1.8 ms** |
| `bad_disjoint_triangles.in` (341 $K_3$) | 341 | 341, TLE 5 s | 341, **185 ms** |

Test 7 je dramatičen: iz "5 s in nepravilno označena časovna odpoved" v "**0 vej B&B**" —
vse komponente reši že prva faza in zgornja meja (nadgradnja 4) certificira optimalnost
takoj.

---

## Nadgradnja 3 — Tesnejša zgornja meja po komponentah

### Težava z globalno mejo

Prvotni meji `|V|/2` in `|E|` se računata čez **celotno** aktivno množico. Za graf z
veliko komponentami je to ohlapno.

Primer: 341 disjunktnih trikotnikov, $|V| = 1023$.

- Globalna meja: $\lfloor 1023/2 \rfloor = 511$.
- Dejanski optimum: 341.
- Razmerje: 1.5× preveč. Pogoj $\text{best} + \text{ub} \le \text{best}$ se sproži
  šele, ko iz aktivne množice odstranimo več kot 340 vozlišč — globoko v drevesu.

### Rešitev: meja po komponentah

Ker zdaj komponente obravnavamo neodvisno, je smiselno računati mejo **na ravni
komponente**:

$$ \mathrm{ub}(C) = \min\Bigl(\lfloor |V_C| / 2 \rfloor,\ |E_C|\Bigr) $$

Za trikotnik $K_3$: $|V_C|/2 = \lfloor 3/2 \rfloor = 1 = \mathrm{MIM}(K_3)$. **Točno!**

Za 341 trikotnikov skupna meja: $341 \times 1 = 341 = \mathrm{MIM}$. **Tesno.**

### Dodatek: detekcija polnega grafa

Za poljuben polni graf $K_t$ velja $\mathrm{MIM}(K_t) = 1$ (vsaki dve povezavi delita
soseda — celotno $V$).

Detekcija je trivialna: komponenta velikosti $t$ je $K_t$, če in samo če ima natanko
$\binom{t}{2}$ povezav:

```go
func isClique(active Bitset) bool {
    k := active.PopCount()
    if k < 2 { return false }
    return edgeCount(active) == k*(k-1)/2
}
```

Pri pozitivnem odgovoru: $\mathrm{ub}(C) = 1$.

### Primer učinka: 200 disjunktnih $K_4$

Brez tega:

- $|V_C| / 2 = 2$, $|E_C| = 6$ → ub = 2.
- Hevristika najde 1, ub = 2 → B&B mora preveriti.

S klikovno detekcijo:

- isClique($K_4$) = true → ub = 1.
- Hevristika najde 1 = ub → B&B preskočimo (glej nadgradnjo 4).

### Pomembna lekcija

V prvi različici tega popravka sva poskusila še **zahtevnejšo "edge-packing" mejo**:
greedy izbira povezave z najmanjšim $|N[u] \cup N[v]|$, štetje izbir kot zgornja meja.

Zdela se je intuitivna, a je bila **nepravilna**. Štetje izbir je v resnici *spodnja*
meja (vsaka izbira predstavlja eno povezavo veljavnega induciranega ujemanja). Z napačno
mejo `ub` je solver na grafu `random_100_medium` zgodaj prekinil B&B in vrnil 15 namesto
17. Po odkritju regresije sva to mejo odstranila in pustila le dokazljivo varne
($\lfloor |V_C|/2 \rfloor$, $|E_C|$, klikovna detekcija).

Ta epizoda potrjuje pomemben princip: **zgornja meja v B&B mora biti vedno pravilna**
($\ge$ dejanski optimum). Če je premajhna, B&B reže veljavne rešitve in vrne nepravilen
odgovor.

---

## Nadgradnja 4 — Zgodnja prekinitev pri LB $\ge$ UB

### Težava

Tudi z odlično mejo: če $\mathrm{LB} = \mathrm{UB}$ na korenu, prvotna implementacija B&B
**ne preneha** takoj. Vstopi v rekurzijo, izvede redukcije, izračuna mejo, in šele
*znotraj* preverja `len(curEdges) + ub <= bestSize`. Pri korenu je `curEdges` prazen in
`bestSize = LB`, torej se `0 + UB <= LB` sproži *šele ob naslednjem klicu*. V vmesnem
času pa lahko zapravi tisoče vej.

### Rešitev

V `solveComponent` najprej izračunamo `LB` (iz hevristike) in `UB`. Če sta enaka,
B&B **sploh ne poženemo**:

```go
heur := greedySolve(comp)
heur = localSearch(heur, comp)
lb := len(heur)
ub := upperBound(comp)

if lb >= ub {
    return heur          // certificiran optimum, brez B&B
}

// drugače: poženi B&B na komponenti
solveBB(comp.Clone())
```

### Intuicija

Hevristika **dokaže** $\mathrm{MIM}(C) \ge \mathrm{LB}$ (vrne veljavno rešitev te
velikosti). Meja **dokaže** $\mathrm{MIM}(C) \le \mathrm{UB}$. Če sta enaki, je
$\mathrm{MIM}(C) = \mathrm{LB} = \mathrm{UB}$ in optimum je certificiran.

Pred tem popravkom je hevristika dala optimalno *kvantiteto*, a brez certifikata, in
B&B je iskal "nekaj boljšega". Sedaj sta meja in spodnja meja **dialog**:
ko se srečata, se algoritem ustavi.

### Učinek

| Test | Pred | Po | Komentar |
|------|------|-----|----------|
| `10_sotorov_graf_250.in` | 250, 47 ms, 1 vozlišče B&B | 250, 47 ms, **0 vozlišč** | Že prej hitro; sedaj certificirano brez vsakega B&B koraka. |
| `07_disjunktne_k4_x200.in` | 200, TLE | 200, 112 ms, **0 vozlišč** | Vsaka komponenta certificirana takoj. |
| `bad_disjoint_triangles.in` | 341, TLE | 341, 185 ms, **0 vozlišč** | 341 trikotnikov × 0 vej B&B. |

Pomemben vzorec: ko se vse tri nadgradnje (komponente + tesna meja + zgodnja prekinitev)
sestavijo, **najtežji disjunktni primeri postanejo trivialni**, ker certifikat optimuma
pride že iz prve faze.

---

## Skupni rezultati

Vse stare teste solver še vedno reši pravilno (13/13 v `tests/`).

**Recenzentova problematična gradiva (5 s časovne omejitve):**

| Vir | Test | Pred | Po | Pospešek |
|------|------|------|------|----------|
| Snoj–Zmazek | 01 cikel n=1025 | **panic** | 341, 5 s | crash → pravilno |
| Snoj–Zmazek | 07 disjunktni K₄ × 200 | 200, **TLE** | 200, **112 ms** | ≈ 45× |
| Snoj–Zmazek | 08 dvojčki 6×10 | 12, 149 ms | 12, **1.8 ms** | ≈ 80× |
| Snoj–Zmazek | 10 šotor | 250, 47 ms (1 vozlišče) | 250, 47 ms (0 vozlišč) | enako, certificirano |
| Eval | disjunktni K₃ × 341 | 341, **TLE** | 341, **185 ms** | ≈ 27× |

**Težki primeri, kjer pristop ne pomaga** (vse so ene same velike povezane komponente
brez strukture, ki bi jo razgradnja ali enostavna meja izkoristila):

- `02_maxN_gost_1024.in` (gosti naključni)
- `03_gost_random_200.in` (faza prehoda)
- `05_hiperkocka_q8.in`
- `09_mreza_24x24.in`
- `bad_random_3reg.in`, `bad_random_4reg.in`

Ti zahtevajo LP-relaksacijo (frakcijsko ujemanje, blossom) ali algoritme specifične za
graf (DP na ravninskih, bipartitnih). To je presegalo obseg te iteracije, a ostaja
možna naslednja nadgradnja.

---

## Dvominutni govor za predstavitev

> Predpostavljamo, da je algoritem (hevristika + B&B + bitni nizi) že razložen.
> Ta segment **doda** štiri ciljne izboljšave na podlagi recenzij. Cilj: ~2 minuti.

---

### Diapozitiv "Nadgradnje" — Govorec **B** — ~60 s

**Celoten govor:**

> "Po predaji sva prejela dve recenziji. Obe sta se strinjali v eni ugotovitvi:
> hevristika je dobra — pogosto najde optimum — ampak B&B s prvotnimi mejami tega
> **ne zna dokazati** in zato zapravi celoten časovni proračun. Glede na merilo
> naloge, ki šteje čas, je to izguba.
>
> Na to sva odgovorila s štirimi ciljnimi popravki.
>
> **Prvič: dinamični bitni niz.** Fiksna meja 1024 vozlišč je solver s 1025 vozlišči
> zrušila. Sedaj velikost alociramo glede na vhod.
>
> **Drugič, najpomembnejše: razgradnja po povezanih komponentah.** Inducirano ujemanje
> je *aditivno* — MIM celotnega grafa je vsota MIM-ov posameznih komponent. 200
> disjunktnih $K_4$ smo prej obravnavali kot eno eksponentno drevo; sedaj jih
> rešujemo neodvisno.
>
> **Tretjič: tesnejša zgornja meja na ravni komponente.** Računamo $|V_C|/2$ za vsako
> komponento posebej, dodali smo **detekcijo polnih grafov** $K_t$, kjer je MIM
> zagotovo 1.
>
> **Četrtič: zgodnja prekinitev.** Ko po hevristiki velja spodnja meja je enaka
> zgornji, je optimum **certificiran** in B&B sploh ne poženemo."

---

### Diapozitiv "Učinek" — Govorec **B** — ~50 s

**Celoten govor:**

> "Rezultati so dramatični tam, kjer struktura izpostavi razgradnjo:
>
> Test s 1025 vozlišči, ki je program **zrušil**, sedaj vrne pravilno rešitev 341.
>
> 200 disjunktnih $K_4$: iz **časovne odpovedi** v **110 milisekund**, brez ene same
> raziskane veje B&B.
>
> 341 disjunktnih trikotnikov: iz časovne odpovedi v **185 milisekund**.
>
> Test z dvojčki: iz 150 v **1.8 milisekunde**.
>
> Trije težki primeri — gost naključni n=200 pri fazi prehoda, hiperkocka in mreža —
> ostanejo časovno težki, ker so to ene same velike povezane komponente brez
> strukture, ki bi jo katera od najinih nadgradenj izkoristila. Te zahtevajo
> linearno-programsko relaksacijo ali metode specifične za graf — kar je
> **naslednji** korak, ne ta."

**Premor 2 s. Možen prehod na vprašanja.**

---

### Tehnične opombe za govorca

1. Glavna pripoved naj **ne bo** "tehnologija" ampak "zakaj": recenzent je odkril, da
   hevristika *je* dobra — težava je certifikat. Vse štiri izboljšave skupaj omogočijo,
   da meja in spodnja meja "se srečata" že na korenu.
2. Pri **drugi nadgradnji** poudarite aditivnost: to je *naravna* lastnost MIM-a, ki je
   prvotni solver ni izkoristil. Občinstvo bo vprašalo: "kdaj se to ne pomaga?" Odgovor:
   ko je graf ena velika povezana komponenta.
3. Pri **tretji nadgradnji** lahko omenite epizodo z napačno mejo: "intuicija je bila
   napačna — meja je bila premajhna in solver je vrnil suboptimalno rešitev". To je
   poučno: dokazljivost > intuicija.
4. Številke govorijo same: ne berite vseh vrstic tabele, izpostavite samo najbolj
   dramatične (panic → pravilno, TLE → 100 ms, ali "0 vej B&B").

---

### Časovni proračun

| Element | Čas |
|---------|-----|
| Diapozitiv "Nadgradnje" | ~60 s |
| Diapozitiv "Učinek" | ~50 s |
| Premor / vprašanja | ~10 s |
| **Skupaj** | **~2 min** |
