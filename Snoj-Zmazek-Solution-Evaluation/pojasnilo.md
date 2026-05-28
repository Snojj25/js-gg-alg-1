# Vrednotenje začetne rešitve Snoj–Zmazek

Ocenjevalca: **Miha & Žiga**
Ocenjujemo: `solver.go` (avtorja Jure Snoj, Gal Zmazek).

Problem: **Konstrukcija največje ortogonalne CC-množice**

---

## Povzetek
Solver je "slab" tudi
takrat, ko vrne pravilen optimum, če pa za to porabi celotnih 300 s. To je
ključni vzvod. 


Teste razdelimo v tri jasne
kategorije:

**A. Fatalna napaka (1 test)**
- **Test 1**: solver panicira ob branju vhoda. Ne vrne nobenega odgovora,
  zato ne dobi nobenih točk.

**B. Časovna sabotaža — heuristik je
optimum, B&B 300 s ne prinese ničesar (8 testov)**

V vseh teh primerih solver vrne **isti odgovor** kot Faza 1 sama, le da
nepotrebno sežge celotnih 300 s na B&B drevesu, ki ga njegove zgornje
meje (`|V|/2`, `edgeCount`) ne znajo odrezati. Glede na merilo PDF-a
("ocenjujejo se glede na čas izvajanja") je to maksimalna možna kazen
za pravilen odgovor.

**C. Hkrati napačen *in* počasen odgovor (1 test)**
- **Test 3** (random p=0.45, n=200): heuristika obtiči pri 6, eksaktno
  pa najdemo 7. B&B v 300 s ne dohaja zaradi ~80-kratnega razvejitvenega
  faktorja → solver vrne **strogo podoptimalen** odgovor *in* porabi
  vseh 300 s. Najhujši posamezen rezultat.



**Zaključni vtis:** glavna šibkost solverja ni kvaliteta heuristične
rešitve (ta je dejansko zelo dobra) — problem je, da B&B faza s
trenutnimi mejami **ne zna prepoznati, da je heuristika že optimalna**,
in zato sistematično doseže `timed_out: true` na vseh netrivialnih
vhodih. Merilo PDF-a je čas izvajanja, zato je to neposreden poraz pri
ocenjevanju, čeprav je rezultat pogosto pravilen.

---
## Opis problema in rešitve
Iščemo največjo množico povezav `S ⊆ E` v neusmerjenem grafu, tako da za
poljubna para `(v₁, u₁), (v₂, u₂) ∈ S` velja
`(u₁, u₂), (u₁, v₂), (v₁, u₂), (v₁, v₂) ∉ E`. Vhod je matrika sosednosti
`n × n`, izhod pa število parov in seznam parov `1`-indeksiranih vozlišč.

Rešitev `solver.go` deluje v dveh fazah:
1. **Požrešno + lokalno (1,2)-iskanje**: izbira povezave z najmanjšim
   `|N(u) ∪ N(v)|`.
2. **Branch-and-bound**: razveji po vozlišču z najmanjšo stopnjo, reže z
   mejama `|aktivna vozlišča| / 2` in število preostalih povezav.

Matrika sosednosti je v bitnih nizih `Bitset = [16]uint64`, kar pomeni
**fiksno zgornjo mejo `maxN = 1024`**.


## Strukturne pomanjkljivosti, ki jih napadamo

- **Fiksna meja `maxN = 1024` brez preverjanja vhoda** → test 1 (panic).
- **Šibki zgornji meji `|V|/2` in `edgeCount`**: ne uporabi npr.
  matching-number sklepanja niti LP-relaksacije → testi 2, 3, 6.
- **Brez simetrijskih redukcij (twins, module)** → testi 4, 5, 7, 8.
- **Brez izkoriščanja strukture (ravninsko / dvodelno / mreža / DP)** →
  testa 5, 9.
- **Brez "dominating / high-degree universal vertex" redukcije** → test 10.
- **Brez detekcije optima**: ko Faza 1 že najde optimalno rešitev, B&B nima
  certifikata, da je `bestSize` optimalen, in ne sme zaključiti zgodaj —
  zato se vse 300 s zapravi na neuporabnem iskanju.

## Rezultat

| Test | Pričakovano vedenje |
|---|---|
| 1 | `panic: index out of range`, ni izhodne datoteke. **Fatalno.** |
| 2–10 | `timed_out: true` po 300 s, izhod = (skoraj) optimum, a najslabši možni čas. |

---

## Opisi posameznih testnih primerov

## 1. `01_nad_mejo_n1025.in` : prekoračitev fiksne meje

Cikel na **n = 1025** vozliščih. V `solver.go` je `maxN = 1024`,
`maxWords = 16`. Pri branju matrike solver kliče `adj[i].Set(j)` za
`i, j ∈ [0, 1024]`; pri `j = 1024` izvaja `b[j>>6] = b[16]`, kar je izven
polja `[16]uint64` in povzroči `panic: index out of range`. Edini test, kjer
solver sploh ne vrne odgovora.

## 2. `02_maxN_gost_1024.in` : največji dovoljeni n z gostoto ~15 %

Naključen graf na **1024 vozliščih**, `p = 0.15`. Faza 1 najde rešitev
velikosti ≈21 v nekaj sekundah. B&B nato pri min. stopnji > 100 odpre veje
brez kakršne koli redukcije; meja `|V|/2 = 512` pa je 25× večja od prave
vrednosti, torej **ne reže nič**. Solver porabi celotnih 300 s in vrne
heuristični odgovor.

## 3. `03_gost_random_200.in` : naključni graf pri fazi prehoda

`n = 200`, `p = 0.45`. Faza 1 najde 6, eksaktni B&B pa najde rešitev
velikosti **7** → **greedy je tu strogo podoptimalen**, ne le počasen.
Min. stopnja ≈ 80 → razvejitveni faktor 80 v vsakem koraku B&B. Meja
`|V|/2 = 100` ne reže, dokler ni rešitev že skoraj zgrajena. Klasičen
"timeout brez izboljšave".

## 4. `04_crown_40.in` : kronski graf (`K_{40,40}` brez ujemanja na diagonali)

39-regularen dvodelni graf. **Optimum je dokazano 2** in ga heuristika takoj najde.
Toda meji `|V|/2 = 40` in `edgeCount ≈ 1560` sta tako visoko nad 2, da B&B
ne more rezati. Zaradi popolne simetrije (vse povezave imajo isto vrednost
`|N(u) ∪ N(v)| = 80`) preišče tudi ekvivalentne permutacije znova in znova.
**Optimum + maksimalna časovna izguba.**

## 5. `05_hiperkocka_q8.in` : hiperkocka Q₈ (256 vozlišč, 8-regularna)

Vozliščno-tranzitiven dvodelni graf. Faza 1 najde 64 (validirano kot
veljavna CC-množica). Optimuma za Q₈ ne poznamo natančno; v vsakem
primeru je `≥ 64`. Vse veje so
strukturno enakovredne; B&B brez simetrijskih redukcij obišče ogromno
ekvivalentnih poddreves.

## 6. `06_dense_complement_160.in` : komplement redkega grafa

`n = 160`, gostota > 94 %, min. stopnja ~150. **Optimum je dokazano 2**
(eksaktni B&B v <0.1 s) in ju Faza 1 najde takoj. B&B pa ima v vsakem klicu ~150 include-kandidatov, in obe meji
(`|V|/2 = 80`, `edgeCount ≈ 12 000`) sta ogromno nad bestSize=2, zato se
nobena veja ne odreže. **Najhujše razmerje "trivialna heuristika vs.
neuporaben B&B"**.

## 7. `07_disjunktne_k4_x200.in` : 200 disjunktnih K₄

Optimum je natanko 200, Faza 1 ga najde. Toda v vsakem K₄ so vse stopnje
3, zato se ne sproži ne "isolated vertex" ne "isolated edge" redukcija.
B&B mora rekurzivno obdelati vsako od 200 strukturno enakih komponent
posebej. Brez "component decomposition" + "memoization" je to čisto
zapravljanje časa.

## 8. `08_twins_6x10.in` : graf z dvojčki (twin vertices)

6 skupin po 10 vozlišč (skupaj 60), znotraj skupine enaka soseščina, plus
naključni medskupinski vzorec in cikel znotraj vsake skupine. Faza 1 najde
12, kar je z eksaktnim B&B (2 s) **dokazano optimum**.
Brez twin-redukcije se B&B razveji ločeno po vsakem dvojčku, čeprav so vse
veje strukturno enake.

## 9. `09_mreza_24x24.in` : mreža 24×24 (576 vozlišč)

Ravninski dvodelni graf. Faza 1 najde 144 = `⌈576/4⌉` (znana formula za
mrežo, je hkrati zg. meja). Problem je
na mreži polinomsko rešljiv z DP, a solver tega ne ve. Meja `|V|/2 = 288`
je 2× večja od prave vrednosti → ne reže ničesar do zelo globoko v
drevesu.

## 10. `10_sotorov_graf_250.in` : "šotor" z enim skupnim centrom

Eno osrednje vozlišče stopnje 500, 250 trikotnikov okoli njega.
**Optimum je dokazano 250** (250 listnih povezav je dosegljivih, vključitev
centra blokira vse liste in da kvečjemu 1) in Faza 1 ga najde. B&B pa ko zaide v vejo, ki vključuje center, preskusi 500
include-vej zaporedoma. Brez "high-degree universal vertex" redukcije je
to katastrofa za čas.

