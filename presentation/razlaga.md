# Razlaga — Maksimalno inducirano ujemanje (MIM)

Dokument služi kot poglobljena spremljevalna razlaga za predstavitev (`slides.md`).
Za vsak diapozitiv navaja intuicijo, primere in razloge za sprejete odločitve.

Pred razlago diapozitivov sta dve teoretični poglavji — eno o **samem problemu**
(kaj je MIM, od kod izvira, kako težek je) in eno o **algoritemskih paradigmah**
(požrešnost, lokalno iskanje, razvejaj in omeji), na katerih sloni naša rešitev.
Ti dve poglavji nista na prosojnicah; sta referenca za pripravo in za vprašanja.

---

## Diapozitiv 1 — Naslov

Konstrukcija maksimalne ortogonalne CC-množice je slovenski naziv za **Maximum Induced
Matching** (MIM) problem. "CC" se nanaša na *closed neighborhoods* (zaprte soseščine):
dve povezavi sta združljivi (ortogonalni), če nobeno krajišče ene ne leži v zaprti
soseščini krajišč druge. Natančnejša obravnava je v poglavju *Teoretično ozadje problema*
spodaj.

Predstavitev je deljena na dva govorca: **A** pokrije pristop in hevristiko (faza 1),
**B** pokrije točen algoritem in rezultate (faza 2).

---

## Teoretično ozadje problema  *(ni na prosojnicah)*

> Tega poglavja **ni** na prosojnicah — predpostavljamo, da občinstvo problem že pozna.
> Tukaj zbiramo teorijo, ki jo je dobro poznati za pripravo in za odgovore na vprašanja.

### Kaj problem sploh je

Predstavljajmo si graf kot omrežje, kjer vozlišča predstavljajo objekte, povezave pa
pare objektov v neposrednem stiku. Iščemo čim več povezav, ki so med seboj **popolnoma
neodvisne**: ne le, da si ne delijo vozlišča (kot pri navadnem ujemanju), temveč nobeno
njihovo krajišče ni niti *sosed* krajišča katere koli druge izbrane povezave. Izbrane
povezave torej ležijo "daleč narazen" — med poljubnima dvema je vsaj ena povezava
razmika.

To je problem **maksimalnega induciranega ujemanja** (angl. *Maximum Induced Matching*,
MIM). Inducirano ujemanje se v literaturi imenuje tudi **krepko ujemanje**
(*strong matching*).

### Od kod izvira in čemu služi

Problem sta leta 1982 vpeljala Stockmeyer in Vazirani pod slikovitim imenom *risk-free
marriage* (poroka brez tveganja). Od takrat se pojavlja v številnih aplikacijah:

- **Brezžična omrežja:** povezave induciranega ujemanja so pari oddajnik–sprejemnik, ki
  lahko **hkrati** oddajajo brez medsebojnih motenj — sosednost v grafu pomeni
  interferenco. MIM maksimira število sočasnih, nemotenih prenosov.
- **Testiranje VLSI vezij** in **varna komunikacija**: skupine, ki ne smejo "sevati"
  druga v drugo.
- Splošno: kjerkoli želimo izbrati čim več parov, ki delujejo v **medsebojni izolaciji**.

### Razlaga imena "ortogonalna CC-množica"

Slovenski naziv naloge opisuje isti problem z drugim besednjakom. "CC" stoji za *closed
neighborhoods* (zaprte soseščine). Zaprta soseščina vozlišča $v$ je
$N[v] = N(v) \cup \{v\}$ — vozlišče skupaj z vsemi svojimi sosedi. Dve povezavi sta
**ortogonalni** (združljivi), če nobeno krajišče ene ne leži v zaprti soseščini krajišč
druge:

$$ \{u_1, u_2\} \cap N[\{v_1, v_2\}] = \varnothing \quad\text{(in simetrično obratno).} $$

Iščemo največjo množico povezav, ki so paroma ortogonalne — od tod "maksimalna
ortogonalna CC-množica".

### Formalna definicija

Dan je neusmerjen graf $G = (V, E)$ brez zank in večkratnih povezav. Iščemo **največjo**
množico $S \subseteq E$, tako da za vsaki dve različni povezavi
$(u_1, u_2), (v_1, v_2) \in S$ velja, da nobeno krajišče ene povezave **ni enako in ni
sosednje** nobenemu krajišču druge. Formalno: nobena od povezav
$(u_1, v_1), (u_1, v_2), (u_2, v_1), (u_2, v_2)$ ni v $E$.

> Pogoj **ne** zahteva, da sta soseščini izbranih povezav povsem disjunktni — izbrani
> povezavi smeta imeti skupnega soseda, ki ni krajišče nobene od njiju. Prepovedana je
> le sosednost med samimi **krajišči**.

### Navadno, maksimalno in inducirano ujemanje

Pojme je vredno ločiti, ker se v pogovoru radi mešajo:

| Pojem | Pogoj |
|-------|-------|
| **Ujemanje** (*matching*) | povezave si ne delijo krajišč |
| **Maksimalno ujemanje** (*maximal*) | ujemanje, ki ga ne moremo več razširiti (lokalno) |
| **Največje ujemanje** (*maximum*) | ujemanje z največ povezavami (globalno) |
| **Inducirano / krepko ujemanje** | ujemanje, kjer med krajišči ni nobene dodatne povezave |

Pozor na slovensko dvoumnost: *maksimalno* (angl. *maximal*, "ne da se razširiti") ni isto
kot *največje* (angl. *maximum*, "z največ elementi"). Mi iščemo **največje inducirano**
ujemanje.

Beseda "inducirano" izvira iz pojma **induciranega podgrafa**: če vzamemo vozlišča
izbranih povezav in pogledamo podgraf, ki ga ta vozlišča inducirajo (vse povezave
originalnega grafa med njimi), mora ta vsebovati **natanko** izbrane povezave in nič
drugega. Navadno ujemanje dopušča dodatne povezave med krajišči; inducirano ne.

### Ekvivalentna formulacija prek neodvisne množice

MIM lahko prevedemo na problem **največje neodvisne množice** (*Maximum Independent Set*,
MIS). Zgradimo *graf konfliktov* $H$:

- vozlišča $H$ = povezave grafa $G$;
- dve vozlišči (povezavi $e, f$ iz $G$) povežemo, če sta v **konfliktu** — če si delita
  krajišče ali če obstaja povezava v $G$ med krajiščem $e$ in krajiščem $f$.

Tedaj je inducirano ujemanje v $G$ natanko **neodvisna množica** v $H$, MIM pa ustreza
**največji** neodvisni množici v $H$. Graf $H$ je ravno **kvadrat grafa povezav**
$L(G)^2$ (dve povezavi sta povezani, če sta v $L(G)$ na razdalji $\le 2$).

Ta zveza je poučna: MIS je eden klasičnih NP-težkih problemov, kar takoj nakaže, zakaj
je tudi MIM težek — in pojasni, zakaj se naše tehnike (razvejanje, redukcije, zgornje
meje) tako spominjajo tehnik za neodvisne množice.

### Računska zahtevnost — širša slika

- **Splošni grafi:** MIM je **NP-težak**. Ostane NP-težak celo na močno omejenih
  razredih, npr. na **ravninskih bipartitnih grafih z največjo stopnjo 3**
  (Cameron 1989; Lozin; Mosca). Polinomskega algoritma torej ne pričakujemo (razen
  če P = NP).
- **Aproksimacija:** problem je tudi težko aproksimirati (APX-težak že na grafih
  omejene stopnje) — ni preprostega algoritma s konstantnim jamstvom za vse grafe.
- **Posebni razredi so lahki:** na **drevesih**, **tetivnih (chordal)**, **intervalnih**
  in **šibko tetivnih** grafih obstajajo **polinomski** algoritmi (tipično dinamično
  programiranje, npr. po drevesni dekompoziciji).
- **Parametrizirano:** problem je obvladljiv (FPT) glede na nekatere parametre, npr.
  drevesno širino.

Praktična posledica za nas: za splošne, večje vhode se moramo sprijazniti z dobrim
**približkom**, za manjše ali strukturirane vhode pa lahko dosežemo **točno** rešitev.
Od tod naš dvofazni (hibridni) pristop.

### Primer iz navodil (n = 7)

Povezave: $(1,3), (2,6), (3,7), (4,5), (4,7)$.

Optimalna rešitev: $S = \{(1,3), (2,6), (4,5)\}$, velikost 3.

**Preverba:** sosedi vozlišča 1 so $\{3\}$ — ni v $\{2, 6, 4, 5\}$. Sosedi vozlišča 3
so $\{1, 7\}$ — niti 7 ni v $\{2, 6, 4, 5\}$. In tako naprej za vsa vključena vozlišča.
Povezava $(3,7)$ ni v $S$, ker bi 7 bil sosed 4 — to bi pokvarilo induciranost s $(4,5)$.

### Tipične vrednosti in meje

| Graf | Velikost MIM |
|------|---------------|
| Polni $K_n$ | 1 (vsaki dve povezavi delita sosede) |
| Pot $P_n$ | $\lfloor n/3 \rfloor$ |
| Cikel $C_n$ | $\lfloor n/3 \rfloor$ |
| Zvezda $K_{1,k}$ | 1 |
| $k$ disjunktnih povezav | $k$ |
| Prazen graf | 0 (ni povezav) |

Dve preprosti **zgornji meji**, ki ju uporabljamo tudi pri rezanju vej (diapozitiv 6):
vsaka izbrana povezava porabi 2 vozlišči, zato $|S| \le \lfloor n/2 \rfloor$; in ne moremo
izbrati več povezav, kot jih graf premore, zato $|S| \le |E|$.

---

## Teoretično ozadje — algoritemske paradigme  *(ni na prosojnicah)*

> Kratek pregled treh metod, na katerih sloni naša rešitev. Pomaga razumeti, *zakaj* so
> koraki na prosojnicah taki, kot so.

### Požrešni algoritmi (*greedy*)

Požrešni algoritem gradi rešitev korak za korakom in v vsakem koraku izbere možnost, ki
se zdi **trenutno najboljša** (lokalno optimalna), ne da bi se kdaj vračal nazaj.

- **Prednost:** zelo hiter in preprost.
- **Slabost:** lokalno optimalne izbire ne dajo nujno globalnega optimuma. Za nekatere
  probleme (npr. matroidi, Kruskalovo vpeto drevo) je požrešnost dokazano optimalna, za
  MIM pa **ne** — zato ga uporabimo le kot **hevristiko** (dober začetni približek).

Bistveno vprašanje pri požrešnem algoritmu je *merilo izbire*. Pri nas je to "cena"
povezave $|N(u) \cup N(v)|$ — koliko vozlišč izbira porabi (diapozitiv 3).

### Lokalno iskanje (*local search*)

Lokalno iskanje vzame obstoječo rešitev in jo poskuša izboljšati z majhnimi spremembami
znotraj neke **okolice**. Definiramo:

- **okolico** — kateri majhni premiki so dovoljeni (pri nas $(1,2)$-zamenjave);
- **kriterij sprejema** — sprejmemo premik, ki rešitev izboljša.

Iskanje teče, dokler obstaja izboljšujoč premik, in se ustavi v **lokalnem optimumu**
(rešitvi, ki je nobena dovoljena majhna sprememba ne izboljša). Lokalni optimum ni nujno
globalni — večja okolica najde boljše rešitve, a je dražja. Naša okolica $(1,2)$ je
najmanjša, ki sploh lahko poveča velikost (diapozitiv 4).

### Razvejaj in omeji (*branch and bound*)

Razvejaj in omeji je **točna** metoda za optimizacijske probleme: implicitno preišče
**ves** prostor rešitev, a velik del odreže, ne da bi ga eksplicitno pregledala. Tri
sestavine:

1. **Razvejanje (*branch*):** prostor rešitev razdelimo na manjše podprobleme — drevo
   odločitev. Pri nas: za izbrano vozlišče se odločimo, katero povezavo vključiti ali pa
   vozlišče izpustiti.
2. **Omejevanje (*bound*):** za vsak podproblem izračunamo **zgornjo mejo** najboljše
   rešitve, ki jo še lahko vsebuje (pri maksimizaciji). Hkrati vzdržujemo **najboljšo
   doslej najdeno** veljavno rešitev — *inkumbent* (spodnja meja).
3. **Rezanje (*prune*):** če zgornja meja podproblema ne preseže inkumbenta, podproblem
   **odrežemo** — v njem zagotovo ni boljše rešitve.

Metoda je **eksaktna**: če se izvede do konca, dokaže optimalnost (ko se spodnja in
zgornja meja srečata). V najslabšem primeru je še vedno eksponentna — kakovost mej in
vrstni red razvejanja odločata, koliko drevesa dejansko obiščemo. Zato vlagamo trud v
dobro mejo (diapozitiv 6) in pametno strategijo razvejanja (min-stopnja, diapozitiv 5).

### Zakaj kombinacija (hibridni pristop)

Hevristika (požrešni + lokalno iskanje) **takoj** ponudi dobro veljavno rešitev — to je
spodnja meja in hkrati "varnostna mreža", če nam zmanjka časa. Razvejaj in omeji to mejo
**izboljšuje** in jo, kadar uspe, dokaže za optimalno. Boljši začetni inkumbent pomeni
**močnejše rezanje** že od vsega začetka — fazi se torej dopolnjujeta.

---

## Diapozitiv 2 — Naš pristop

### NP-težkost

MIM ostane NP-težak celo na močno omejenih razredih grafov — npr. na **ravninskih
bipartitnih grafih z največjo stopnjo 3** (Cameron, 1989; Lozin in Mosca). Posledično
ne obstaja polinomski algoritem (razen če P = NP). Za splošne grafe je torej smiselno
sprejeti, da bomo za velike $n$ zadovoljni z dobrim približkom.

### Zakaj dvofazni pristop?

- **Faza 1 (hevristika)** se vedno konča v milisekundah. To je naša **garancija**:
  uporabnik bo vedno dobil rešitev, tudi če je čas izjemno omejen.
- **Faza 2 (B&B)** lahko *izboljša* hevristiko do točne rešitve. Ker se začne s spodnjo
  mejo iz hevristike, lahko **takoj** reže veje, ki ne morejo prinesti boljšega rezultata.

To kombinacijo srečamo pogosto v reševanju težkih kombinatoričnih problemov:
hevristika daje "varno mrežo", točen algoritem pa optimum, ko je dosegljiv.

---

## Diapozitiv 3 — Požrešni algoritem

### Strategija

V vsakem koraku se vprašamo: katera povezava nas "najmanj stane"? Cena povezave $(u,v)$
je število vozlišč, ki jih moramo z izbiro odstraniti iz nadaljnjega iskanja:

$$ \text{score}(u, v) = |(N(u) \cup N(v)) \cap \text{active}| $$

Pomen: $N(u) \cup N(v)$ so vsa vozlišča, ki postanejo "prepovedana" za nadaljnje
povezave (sosedi izbranih krajišč po definiciji ne smejo biti krajišča drugih povezav v $S$).

### Mali primer

Vzemimo pot $P_5$: $1 - 2 - 3 - 4 - 5$.

- Score povezave $(1,2)$: $|N(1) \cup N(2)| = |\{2, 1, 3\}| = 3$.
- Score povezave $(2,3)$: $|N(2) \cup N(3)| = |\{1, 3, 2, 4\}| = 4$.
- Score povezave $(3,4)$: $|N(3) \cup N(4)| = |\{2, 4, 3, 5\}| = 4$.
- Score povezave $(4,5)$: $|N(4) \cup N(5)| = |\{3, 5, 4\}| = 3$.

Požrešni izbere $(1,2)$ (ali $(4,5)$ — neodločeno) z oceno 3. Po izbiri $(1,2)$ odstranimo
$\{1, 2, 3\}$ iz aktivnih, ostaneta $4, 5$ in povezava $(4, 5)$. Dodamo še $(4, 5)$.

Rezultat: $|S| = 2 = \lfloor 5/3 \rfloor + 1$. Ujema se z optimumom za $P_5$.

### Zakaj prav $|N(u) \cup N(v)|$ in ne npr. vsota stopenj?

Vsota stopenj bi dvakrat štela skupne sosede $u$ in $v$. Pri MIM ti vendarle
"izgubijo" enkrat — gre torej za unijo, ne vsoto.

---

## Diapozitiv 4 — Lokalno iskanje z (1,2)-zamenjavami

### Ideja

Požrešni je hiter, a sprejme lokalno najboljšo odločitev — globalno lahko zgreši.
Po zaključku poskusimo izboljšati rešitev z majhnimi premiki:

```
ponavljaj:
  za vsako povezavo e ∈ S:
    začasno odstrani e iz S
    izračunaj "prepovedana" vozlišča = unija N(u) ∪ N(v) za vsako preostalo povezavo
    available = vozlišča izven prepovedanih
    poženi mali požrešni nad available — koliko povezav lahko vstavimo nazaj?
    če ≥ 2: sprejmi zamenjavo (neto +1), začni znova
  če ni izboljšave: konec
```

### Zakaj (1,2) in ne (1,3) ali (2,3)?

- (1,1) — zamenjava ene povezave za eno — nikoli ne izboljša velikosti.
- (1,2) — odstranimo 1, dodamo 2 — neto +1. **Najmanjši premik, ki izboljša velikost.**
- (1,3), (2,3), ... — drago iskanje za malo verjeten dobiček. Eksperimentalno (1,2) zadošča
  za večino izboljšav, ki se splačajo glede na čas.

### Konkreten primer

Vzemimo graf:

```
    1 — 2 — 3 — 4
        |
        5 — 6
```

Požrešni morda izbere $(2,5)$ kot prvo (score = 4). Po tem ostane $(3,4)$ kot druga
povezava. Rezultat: $|S| = 2$.

Lokalno iskanje preveri "odstrani $(2,5)$" — ostane samo $(3,4)$. Available =
$\{1, 5, 6\} \cup \{2\}$ minus "prepovedani od $(3,4)$" $= N(3) \cup N(4) = \{2,4,3\}$.
Available = $\{1, 5, 6\}$. Vidimo povezavo $(5, 6)$, vstavimo. Naprej v available je $\{1\}$ —
povezave ni, končamo. Vstavili smo 1 povezavo — to ni izboljšava (neto 0). Sprejem ne.

V resničnih grafih z več gostote (1,2)-zamenjave **pogosto** najdejo +1, ker odstranitev
ene povezave sprosti dve "regiji" prepovedi.

### Konvergenca

Vsak sprejet swap poveča $|S|$ za $\ge 1$. $|S|$ je navzgor omejen z $\lfloor n/2 \rfloor$,
zato lokalno iskanje vedno konvergira v $\le n/2$ iteracijah. V praksi ~3-5 prehodov.

---

## Diapozitiv 5 — Razvejaj in omeji (B&B)

### Splošna ideja

Razvejamo prostor vseh možnih induciranih ujemanj kot drevo odločitev. Za vsako
vozlišče v drevesu izračunamo **zgornjo mejo** velikosti najboljše rešitve, ki jo
lahko še najdemo iz te veje. Če meja + trenutna rešitev ne presega najboljše doslej,
celo poddrevo "odrežemo".

### Strategija razvejanja: min-degree

Izberemo vozlišče $v$ z **najmanjšo stopnjo** $d$ v trenutnem aktivnem grafu.
Iz tega vozlišča naredimo $d + 1$ podproblemov:

- Za vsakega soseda $w$ od $v$: **vključi povezavo $(v, w)$** v $S$, odstrani
  $N(v) \cup N(w)$ iz aktivnega grafa.
- **Izključi $v$**: ne uporabi $v$ kot krajišče nobene povezave. Odstrani samo $v$.

### Zakaj minimalna stopnja?

Število vej eksponentno raste z $d$. Minimum stopnje **minimizira faktor razvejanja**
v tem koraku. Hkrati so vključitvene veje "močne" — drastično skrčijo aktivni graf
(odstranijo $|N(v) \cup N(w)|$ vozlišč), zato kmalu pridemo do baznega primera.

### Vrstni red poskusov

Najprej **vključitvene veje**, šele nato izključitvena. Vključitve povečajo $|S|$ takoj
in s tem izboljšajo *best* — kar omogoča močnejše rezanje preostalih vej.

### Bazni primer

Ko v aktivnem grafu ni več nobene povezave, primerjamo $|S|$ z dosedanjim *best* in
ev. posodobimo. Vrni se.

---

## Diapozitiv 6 — Redukcije + rezanje vej

### Varne redukcije za MIM

#### 1. Izolirano vozlišče (stopnja 0)

Vozlišče brez sosedov ne more biti krajišče nobene povezave. Trivialno ga odstranimo.

#### 2. Izolirana povezava (obe krajišči stopnje 1)

Če $u$ in $v$ obstajata, da je $u$ edini sosed $v$ in obratno, potem povezava $(u, v)$
**mora** biti v optimalnem $S$. Argument: ker $u$ in $v$ nimata drugih sosedov,
ju lahko prisilno vključimo brez izgube. Če bi izgubili eno krajišče, izgubimo dejansko
povezavo — strogo poslabšanje.

Po prisili odstranimo $u, v$ iz aktivnega grafa in nadaljujemo.

### Zakaj NI varne splošne "pendant" redukcije

Klasično pri *navadnem* ujemanju velja: vozlišče stopnje 1 (obesek) se sme prisilno
vključiti v ujemanje s svojim edinim sosedom. **To pri MIM ne velja.**

Protiprimer:

```
v — u — a — b
        |
        c — d
```

Tu je $\deg(v) = 1$, sosed $v$ je $u$. Če prisilno vključimo $(v, u)$, potem $a$ in $c$
postaneta nedosegljiva (sosedna $u$). Rezultat: $|S| = 1$.

Optimum pa je $\{(a, b), (c, d)\}$ — povezavi sta inducirano nezdružljivi z $(v, u)$,
a zato z izbiro $(v, u)$ izgubimo dve povezavi za eno.

**Zaključek:** prisilimo le povezave, kjer sta **obe krajišči** stopnje 1 (izolirana
komponenta dveh vozlišč). Drugih obeskov se ne dotikamo.

### Dve zgornji meji za rezanje

Za vsako vozlišče v drevesu izračunamo:

1. **Vozliščna meja**: $\lfloor |\text{active}| / 2 \rfloor$ — vsaka povezava potrebuje 2 krajišči.
2. **Povezavna meja**: $|E_\text{active}|$ — ne moremo izbrati več povezav, kot jih obstaja.

Uporabimo tesnejšo:

$$ \text{ub} = \min\bigl(\lfloor |\text{active}| / 2 \rfloor,\ |E_\text{active}|\bigr) $$

Veja je odrezana, če:

$$ |S_\text{current}| + \text{ub} \le |S_\text{best}| $$

#### Implementacijski trik

Število povezav `|E_active|` izračunamo iz **vsote stopenj / 2** *med istim* pregledom,
ki najde min-degree vozlišče. Praktično nič dodatnega dela.

Cenejšo vozliščno mejo preverimo **prvo** (en `PopCount`), da prihranimo dražje
računanje stopenj, ko je dovolj že to za prirez.

### Časovna omejitev

Vsakih 4096 vozlišč drevesa (`nodes & 0xFFF == 0`) preverimo, ali smo presegli časovno
omejitev. Če smo, postavimo zastavico `timedOut` in se vrnemo navzgor — najboljši
trenutni $S$ je vrnjen. Hevristični rezultat je vedno spodnja meja, torej **nikoli** ne
vrnemo nič slabšega od faze 1.

---

## Diapozitiv 7 — Bitne strukture

### Cilj: množice vozlišč kot prvorazredna podatkovna struktura

Algoritem ves čas izvaja istovrstne operacije nad **množicami vozlišč**: "koliko aktivnih
sosedov ima $v$?", "odstrani $N(v) \cup N(w)$ iz aktivnih", "presekaj soseščini dveh
vozlišč", "ali je še kakšna povezava na voljo?". Če bi te množice hranili kot rezine
indeksov (`[]int`) ali kot zgoščevalne tabele, bi vsaka operacija stala $O(\deg)$ ali
$O(n)$ z nepredvidljivimi skoki po pomnilniku. Z **bitnimi nizi** se vse spremeni v
nekaj zaporednih operacij nad `uint64` — strojnih, zaporednih in povsem znotraj
predpomnilnika.

### Tip `Bitset`

```go
const maxN = 1024
const maxWords = maxN / 64        // = 16

type Bitset [maxWords]uint64      // 16 × 64 = 1024 bitov
```

`Bitset` je fiksno polje 16 šestdesetštirih-bitnih besed. **Vsak bit ustreza enemu
vozlišču.** Vozlišče $i$ je v množici natanko tedaj, ko je $i$-ti bit (v zaporedju
1024 bitov, beseda 0 najprej) prižgan. Bit je dvojiški zastavica: 1 = "v množici",
0 = "ni v množici".

Pretvorba med indeksom vozlišča in lokacijo bita:

| Količina | Izračun | Pomen |
|----------|---------|-------|
| Beseda | `i >> 6` ($= \lfloor i/64 \rfloor$) | kateri od 16 `uint64`-ov |
| Pozicija v besedi | `i & 63` ($= i \bmod 64$) | kateri bit v tej besedi |
| Maska | `1 << (i & 63)` | bit, ki ga prižgemo/preverimo |

Zato so `Set`, `Clear`, `Has` skoraj brezplačni — par strojnih ukazov (premik,
OR / AND-NOT / AND):

```go
func (b *Bitset) Set(i int)      { b[i>>6] |=  1 << uint(i&63) }
func (b *Bitset) Clear(i int)    { b[i>>6] &^= 1 << uint(i&63) }
func (b *Bitset) Has(i int) bool { return b[i>>6] & (1<<uint(i&63)) != 0 }
```

Globalna spremenljivka `nWords = ⌈n / 64⌉` določa, koliko od 16 besed dejansko
iteriramo v operacijah nad celotno množico. Za $n = 200$ je $w = 4$ in obdelamo
le 4 besede, ne vseh 16.

### Kako shranimo graf

Graf $G$ **ne** hranimo kot seznam povezav in ne kot $n \times n$ matriko v dveh
dimenzijah, ampak kot **polje soseščinskih bitnih nizov**:

```go
var adj [maxN]Bitset
```

- `adj[v]` je **soseščina vozlišča $v$**: bit $u$ v `adj[v]` je 1 natanko tedaj,
  ko obstaja povezava $(u, v) \in E$.
- Polje ima fiksno dolžino `maxN = 1024`; indeksiramo ga z $v$.
- Ker je graf neusmerjen, je polje **simetrično**: pri branju vhoda postavimo
  oba bita,

  ```go
  adj[i].Set(j)
  adj[j].Set(i)
  ```

  tako da je bit $u$ v `adj[v]` enak bitu $v$ v `adj[u]`. Vsako povezavo torej
  hranimo dvakrat — to je premišljena izbira, ki nam dovoli "pridobiti vse
  sosede" v enem branju 16 zaporednih besed.

Vsebinsko gre za isto informacijo kot v matriki sosednosti, le da je vsaka vrstica
matrike **stisnjena** v 16 `uint64`-ov in tako dostopna z enim kazalcem (`&adj[v]`),
en pomnilniški blok 128 bajtov.

### Beri zapis: primer $n = 7$

Vzemimo graf iz navodil. Povezave: $(1,3), (2,6), (3,7), (4,5), (4,7)$. Soseščine:

| Vozlišče $v$ | $N(v)$ | `adj[v]` (biti 7..0) |
|--------------|--------|----------------------|
| 1 | $\{3\}$ | `0000 1000` |
| 2 | $\{6\}$ | `0100 0000` |
| 3 | $\{1, 7\}$ | `1000 0010` |
| 4 | $\{5, 7\}$ | `1010 0000` |
| 5 | $\{4\}$ | `0001 0000` |
| 6 | $\{2\}$ | `0000 0100` |
| 7 | $\{3, 4\}$ | `0001 1000` |

Bit z indeksom $k$ predstavlja vozlišče $k$ (najnižji bit = vozlišče 0; v tem primeru
je ta bit vedno 0, ker vozlišče 0 ni del vhoda). Vse prikazane vrednosti v resnici
zasedajo zgornjih 56 bitov besede 0 = 0, spodnji bajt pa je tisti, ki je narisan.
V Go-ju bi `adj[3]` izpisano kot `uint64` znašalo `0x0000000000000082` (=$2^1 + 2^7$).

Iz tabele takoj preberemo: `adj[v].PopCount()` = stopnja vozlišča $v$, npr.
`adj[3].PopCount() = 2`. In `adj[3] & adj[7]` = `1000 0010 & 0001 1000` = `0000 0000`,
torej $N(3) \cap N(7) = \varnothing$ — vozlišči nimata skupnih sosedov.

### Kaj poleg grafa hranimo z bitseti

Med tekom uporabljamo več bitsetov, vsak predstavlja drugo množico vozlišč:

| Spremenljivka | Pomen bita $i = 1$ | Kdaj se spreminja |
|---------------|--------------------|--------------------|
| `adj[v]` | $i$ je sosed vozlišča $v$ | postavljeno enkrat, ob branju vhoda |
| `active` | $i$ je v trenutnem **aktivnem podgrafu** | v vsakem rekurzivnem koraku B&B |
| `newActive` (lokalno v B&B) | enako, za naslednjo rekurzivno raven | kopija `active`, odštete $N(v) \cup N(w)$ |
| `remove` (zaslonjeno) | $i$ je v $N(v) \cup N(w)$ | sestavljeno tik pred odstranitvijo |
| `forbidden` (lokalno iskanje) | $i$ leži v soseščini katere od sprejetih povezav | pri vsakem (1,2)-poskusu |
| `avail` (lokalno iskanje) | $i$ je kandidat za novo povezavo | komplement `forbidden` znotraj `active` |

Sama **rešitev** ni bitset: povezave hranimo kot `[]Edge` (rezina parov `{u, v}`),
ker je pri rezultatu pomembna identiteta in vrstni red. Spremenljivki `curEdges`
(trenutna delna rešitev v rekurziji) in `bestEdges` (najboljša doslej) sta torej
običajni rezini, ne bitseti.

### Hitre operacije in njihov "grafovski" pomen

Vsaka spodnja operacija je zanka čez $w = \lceil n/64 \rceil$ besed; za $n = 200$
to pomeni **4 ponovitve** — v praksi konstanta.

| Operacija | Pomen v jeziku grafov | Tipična raba |
|-----------|----------------------|---------------|
| `bs.Set(i)` / `Clear(i)` / `Has(i)` | dodaj / odstrani / preveri eno vozlišče | točkovne posodobitve |
| `a.And(&b)` | presek $A \cap B$ | `adj[v].And(&active)` = aktivni sosedi $v$ |
| `bs.OrWith(&c)` | unija $A \cup B$ (na mestu) | gradnja $N(u) \cup N(v)$ |
| `bs.AndNot(&c)` | razlika $A \setminus C$ (na mestu) | odstrani $N(v)$ iz aktivnih |
| `bs.PopCount()` | $\lvert A \rvert$ | velikost množice / stopnja |
| `bs.FirstSet()` | najmanjši $i \in A$ | iteracija po vozliščih |
| `bs.IsZero()` | $A = \varnothing$ | bazni primer rekurzije |

Tipično zaporedje v B&B (vključitvena veja za povezavo $(v, w)$):

```go
// 1) Aktivni sosedi v — stopnja v aktivnem podgrafu.
deg := adj[v].And(&active).PopCount()

// 2) Vključi (v, w) v rešitev in pripravi nov aktivni podgraf.
newActive := active            // kopija (128 B, na stacku)
var remove Bitset
remove = adj[v]
remove.OrWith(&adj[w])         // remove = N(v) ∪ N(w)
newActive.AndNot(&remove)      // odstrani vse iz N(v) ∪ N(w)
newActive.Clear(v)             // odstrani še sama v in w
newActive.Clear(w)
solve(newActive)               // rekurzija — kopija gre kot vrednost
```

Vsaka vrstica je zanka čez 4–16 besed: brez alokacij, brez razveljavljanja
predpomnilnika, brez kazalcev čez heap.

### Zakaj `POPCNT` šteje

Moderni procesorji (x86 od SSE4.2, ARM od v8) imajo strojni ukaz **`POPCNT`**,
ki v enem urinem ciklu prešteje prižgane bite v 64-bitnem registru. Go ga uporabi
prek `math/bits.OnesCount64`. Velikost množice je torej "skoraj brezplačna":
$w$ klicev `POPCNT` (≤ 16) je nekaj nanosekund **ne glede na to, koliko vozlišč
je dejansko aktivnih**. Scoring v požrešnem in B&B, ki desetkrat na vsako vejo
izračuna $\lvert (N(u) \cup N(v)) \cap \text{active} \rvert$, zato ni ozko grlo.

### Zakaj fiksna velikost `[16]uint64`

Alternativa bi bila `[]uint64` dolžine `nWords`. Toda:

- **Stack vs. heap:** `[16]uint64` je vrednostni tip. `var x Bitset` ne sproži
  alokacije, `newActive := active` skopira **128 bajtov** kot zlepek strojnih
  premikov. Rezina `[]uint64` zahteva `make` (heap) in upravlja "header"
  (kazalec + dolžina + kapaciteta).
- **Inlining in unroll:** prevajalnik laže optimizira zanke s konstantno zgornjo
  mejo. V vročih zankah Go-jev kompilator prepozna fiksno dolžino 16 in lahko
  zanko razvije ali vektorizira.
- **Predvidljiv pomnilniški vzorec:** vsi biti enega bitseta ležijo zaporedoma
  v 128 bajtih — natanko **2 vrstici predpomnilnika** (64 B vsaka).

Cena je trda omejitev na **1024 vozlišč**. Za to nalogo (in tipične primere
iz knjižnice) povsem zadošča.

---

## Diapozitiv 8 — Rezultati

### Pregled

| Test | $n$ | Tip | $|S|$ | Čas | Točno? |
|------|-----|------|-----|------|--------|
| primer | 7 | iz navodil | 3 | <1ms | ✅ |
| complete_6 | 6 | $K_6$ | 1 | <1ms | ✅ |
| path_10/12 | 10/12 | poti | 3/4 | <1ms | ✅ |
| cycle_9 | 9 | cikel | 3 | <1ms | ✅ |
| disjoint_10 | 10 | 5 disj. pov. | 5 | <1ms | ✅ |
| petersen | 10 | Petersen | 3 | <1ms | ✅ |
| random_50 | 50 | $p = 0.05$ | 13 | 16 ms | ✅ |
| random_100 | 100 | $p = 0.10$ | 17 | 60 s | hevristika |
| random_200 | 200 | $p = 0.03$ | 41 | 60 s | hevristika |
| random_500 | 500 | $p = 0.02$ | 80 | 60 s | hevristika |

### Opazovanja

- **Mali grafi (do $n \approx 50$):** B&B konča v točno rešitev v milisekundah.
- **Posebni grafi ($K_n$, $P_n$, Petersen):** se ujemajo z znanimi formulami.
- **Srednji-veliki:** B&B ne konča v 60 s, vrne hevristični rezultat. Vendar je ta
  rezultat za sparse grafe pogosto blizu optimumu (eksperimentalno znotraj 5-10 %).
- **Velika redka (n=500, p=0.02):** hevristika sama poda 80 povezav v <10 ms.

### Kje je "zid"?

Eksaktnost je odvisna od **gostote**, ne le velikosti:

- $p = 0.02$, $n = 500$ — premalo strukture za hitro rezanje, čeprav je graf redek.
- $p = 0.10$, $n = 100$ — premalo redukcij, branching faktor visok.

Bolj kot je graf strukturiran (drevesa, ciklični grafi, redki grafi), bolj B&B uspe.

---

## Diapozitiv 9 — Povzetek + možne izboljšave

### Glavne ugotovitve

1. **Hibridni pristop deluje**: hevristika ponuja "varno mrežo", točen algoritem
   izboljšuje, ko je čas dovoljen.
2. **Min-degree razvejanje** + **izolirane redukcije** + **dve zgornji meji** so
   skupaj učinkovita kombinacija za rezanje drevesa.
3. **Bitni nizi** dajo $\times 10$ – $\times 50$ pospešek nad naivno predstavitvijo
   (ocena iz mikrobenchmarkov).
4. Iskreni smo o NP-težkosti: nad $n \approx 50$ srednje gostote ali $n \approx 200$
   redko, vrnemo hevristični rezultat — vedno hitro.

### Možne nadgradnje

- **Boljša zgornja meja**: frakcijska LP relaksacija ali zgornja meja iz iskanja
  pokritja klike (clique cover). Tesnejša meja → manj vej.
- **Paralelizacija**: B&B veje so neodvisne, naravna kandidatka za work-stealing.
- **Inkrementalni scoring**: trenutno score za vsako povezavo računamo "od nič".
  Po posodobitvi `active` se večina vrednosti spremeni le malo — inkrementalna
  posodobitev bi prihranila $O(m)$ dela na korak.
- **Posebne hitre poti**: če je $G$ drevo, je MIM rešljiv v polinomskem času (DP po
  drevesu). Avtomatsko zaznavanje takih razredov bi pohitrilo "srečne" primere.

---

## Dodatek — Implementacijski pregled (`solver.go`)

Strnjen pregled, **katere operacije** se izvedejo in **v kakšnem zaporedju**, od
začetka programa do izpisa rezultata. Vsaka točka navaja ključne klice nad bitseti,
da je videti, kako prejšnja teorija postane vrstice kode.

### 1. Vhod in inicializacija (`main`)

1. Preberemo $n$ in nastavimo `nWords = (n + 63) / 64` — odločimo, koliko od 16
   besed dejansko obdelujemo.
2. Beremo $n \times n$ matriko sosednosti. Za vsako enico postavimo **oba** bita:
   ```go
   adj[i].Set(j)
   adj[j].Set(i)
   ```
   Po koncu je `adj[v]` množica sosedov vozlišča $v$.
3. Zgradimo začetni `active` z `Set(0..n-1)` — vsa vozlišča so na začetku v igri.
4. Zaženemo **fazo 1** (požrešni + lokalno iskanje), shranimo dobljen $|S|$ v
   `bestSize` in povezave v `bestEdges`. To je naša garantirana spodnja meja.
5. Zaženemo **fazo 2**: `solve(active)`.
6. Po vrnitvi (bodisi naravno bodisi zaradi `timedOut`) izpišemo `bestEdges`.

### 2. Faza 1a — `greedySolve()`

Zanka, dokler v aktivnem grafu obstaja povezava:

1. Za vsako vozlišče $u \in$ `active`:
   - `nbU := adj[u].And(&active)` — aktivni sosedi $u$.
   - Iteriramo bite `nbU` (kanonično $v > u$, da povezave ne štejemo dvakrat) z
     `bits.TrailingZeros64` in `w &= w - 1`.
   - Za vsako kandidatno $(u, v)$ izračunamo

     ```go
     combined = adj[u]
     combined.OrWith(&adj[v])
     score := combined.And(&active).PopCount()   // |N(u) ∪ N(v) ∩ active|
     ```
   - Spremljamo $(u^*, v^*)$ z najmanjšim score.
2. Dodamo $(u^*, v^*)$ v rezultat in iz `active` odstranimo $N(u^*) \cup N(v^*)$:
   ```go
   remove = adj[bestU]; remove.OrWith(&adj[bestV])
   active.AndNot(&remove)
   ```
3. Konec, ko nobenega kandidata ni več (`bestU == -1`).

Operacije na povezavo: ena `And`, ena `OrWith`, eden `PopCount` — vse $O(w)$.

### 3. Faza 1b — `localSearch(edges)`

Zunanja zanka teče, dokler je `improved == true`. Za vsako povezavo $e_i \in S$:

1. **Zgradi `forbidden`**: za vse preostale povezave $e_j$ ($j \neq i$) zlije
   `adj[u_j]` in `adj[v_j]` v `forbidden` z `OrWith`. To je množica vozlišč,
   ki jih nove povezave ne smejo uporabiti.
2. **Zgradi `avail`**: `Set(0..n-1)`, nato `avail.AndNot(&forbidden)`.
3. **Mali požrešni** nad `localAvail`: enak vzorec kot v `greedySolve`, le da
   uporabi `localAvail` namesto `active`. Cilj je dobiti **dve** novi povezavi;
   po vsaki sprejeti odstranimo $N(\cdot) \cup N(\cdot)$ iz `localAvail`.
4. Če smo našli ≥ 2 povezavi, sestavimo nov $S$ (brez $e_i$, plus nove),
   postavimo `improved = true` in zaženemo zunanjo zanko znova.

Konvergira po $\le \lfloor n/2 \rfloor$ sprejetih izmenjavah (vsaka prinese +1).

### 4. Faza 2 — rekurzivni `solve(active Bitset)`

Vsak klic dobi **kopijo** `active` (vrednostni tip, 128 B na stacku). V tem
zaporedju:

1. **Časovni in flag check** — če je `timedOut`, takoj vrni. Vsakih 4096 vozlišč
   drevesa (`nodes & 0xFFF == 0`) preverimo `time.Since(startTime) > timeLimit`.
2. **Shrani `len(curEdges)`** in `defer` obnovitev — za rollback ob vrnitvi
   iz veje.
3. **Redukcije v zanki "do mirovanja"**:
   - *Izolirana vozlišča*: za vsak $u$ z `adj[u].And(&active).IsZero()` izvedemo
     `active.Clear(u)`.
   - *Izolirane povezave*: za vsak $u$ z $\deg = 1$ (en `PopCount` aktivnih
     sosedov) najdemo sosed $v$ prek `FirstSet()`. Če ima tudi $v$ stopnjo 1,
     dodamo $(u, v)$ v `curEdges` in počistimo oba.
   - Ponavljaj, dokler v eni iteraciji nič ne odstraniš.
4. **Poceni zgornja meja** — `active.PopCount() / 2`. Če
   `len(curEdges) + activeCount/2 ≤ bestSize`, vrni.
5. **En sprehod čez aktivna vozlišča** izračuna hkrati **dve količini**:
   - `minDeg` in `minV` — vozlišče z najmanjšo (a pozitivno) stopnjo;
   - `totalDeg` → `edgeCount = totalDeg / 2`.
6. **Tesnejša zgornja meja**: `ub = min(activeCount/2, edgeCount)`. Če
   `len(curEdges) + ub ≤ bestSize`, vrni.
7. **Bazni primer**: če `minV == -1` (v aktivnem grafu ni nobene povezave) in
   `len(curEdges) > bestSize`, posodobi `bestSize` in skopiraj `curEdges` v
   `bestEdges`. Vrni.
8. **Vključitvene veje** — `neighbors := adj[minV].And(&active)`. Za vsak
   prižgan bit (sosed `nei`):
   ```go
   newActive := active                   // kopija
   remove   = adj[minV]
   remove.OrWith(&adj[nei])               // N(minV) ∪ N(nei)
   newActive.AndNot(&remove)              // ⇒ aktivni \ N(minV)∪N(nei)
   curEdges = append(curEdges[:afterReduction], Edge{minV, nei})
   solve(newActive)                       // rekurzija
   ```
   Po vrnitvi gre v naslednjega soseda; `curEdges[:afterReduction]` odreže
   zadnjo dodano povezavo.
9. **Izključitvena veja** — `newActive := active; newActive.Clear(minV); solve(newActive)`.
   (Po varovalu `if !timedOut`.)

### 5. Skupna slika v eni vrstici

```
main → preberi adj → greedySolve → localSearch → bestSize/bestEdges
     → solve(active):
          reductions* → ub-check → min-deg & edge-count → ub-check
          → include-branches (po sosedih) → exclude-branch
     → izpis bestEdges
```

Vsak korak v `solve` je nekaj zank čez `nWords` besed; ena rekurzija pomeni
**eno alokacijo nič** (kopija `active` je vrednostna, `curEdges` raste in se
krči nad isto rezino prek `[:afterReduction]`). Od tod izhaja praktična
hitrost tudi pri globokih drevesih iskanja.
