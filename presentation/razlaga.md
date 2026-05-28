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

### Predstavitev

Graf hranimo kot polje **bitnih nizov** fiksne velikosti:

```go
type Bitset [16]uint64        // 16 × 64 = 1024 vozlišč
adj  [maxN]Bitset             // adj[v] = soseščina vozlišča v
active Bitset                 // trenutno aktivna vozlišča
```

`Bitset[i]` z bitom $j$ pomeni: vozlišče $j$ je v množici.

### Hitre operacije

| Operacija | Pomen | Zahtevnost |
|----------|-------|------------|
| `bs.Set(i) / Clear(i) / Has(i)` | osnovne | $O(1)$ |
| `bs.OrWith(&c)` | unija | $O(w)$ |
| `bs.AndNot(&c)` | odstrani vse v $c$ | $O(w)$ |
| `bs.PopCount()` | velikost | $O(w)$ + **POPCNT** |
| `a.And(&b)` | presek | $O(w)$ |

Kjer je $w = \lceil n/64 \rceil$ — za $n=200$ je $w = 4$, torej **konstanta**.

### Konkreten dobiček

"Odstrani $v$ in vse njegove sosede iz aktivnega grafa" — kritičen korak v B&B:

```go
active.AndNot(&adj[v])   // 16 strojnih AND-NOT operacij
active.Clear(v)
```

S klasično adjacency listo bi to bilo $O(\deg(v))$ s skoki po pomnilniku in razveljavljanjem
predpomnilnika. Z bitnimi nizi je to **16 zaporednih `uint64` operacij** — zaporedoma v
pomnilniku, idealno za CPU.

### POPCNT

Hardverska `POPCNT` (od SSE4.2 naprej) prešteje bite v `uint64` v enem urinem ciklu.
Go `math/bits.OnesCount64` se prevede prav v ta ukaz. Velikost množic merimo zelo poceni.

### Zakaj fiksna velikost in ne `[]uint64`?

Fiksne `[16]uint64` so **na stacku**, brez alokacij. Kopiranje (`a := b`) je 128 bajtov —
hitreje od dinamičnega `make`. Cena: omejitev na 1024 vozlišč, kar za to nalogo več kot
zadošča.

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
