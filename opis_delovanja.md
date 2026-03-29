# Opis delovanja

## Problem

Imamo neusmerjen graf G = (V, E). Iščemo čim večjo množico parov sosednjih vozlišč S ⊆ E,
pri čemer nobeno vozlišče iz nekega para ni sosednje z nobenim vozliščem iz drugega para v S.
To je t.i. problem maksimalnega induciranega ujemanja (Maximum Induced Matching).

Problem je NP-težek, zato za večje grafe ne moremo zagotoviti optimalnega rezultata v razumnem času.

## Pristop

Rešitev deluje v dveh fazah.

### Faza 1: Požrešni algoritem + lokalno iskanje

Najprej zgradimo začetno rešitev s požrešnim algoritmom. V vsakem koraku izberemo
povezavo (u, v), ki "porabi" najmanj sosednjih vozlišč — torej tisto, kjer je |N(u) ∪ N(v)| minimalen.
Ko izberemo povezavo, iz grafa odstranimo vsa sosednja vozlišča obeh krajišč.

Nato poskusimo rešitev izboljšati z lokalnim iskanjem: za vsako povezavo v trenutni rešitvi
preverimo, ali jo lahko zamenjamo z dvema novima povezavama (t.i. (1,2)-zamenjava). Če uspemo,
smo rešitev povečali za 1 in ponovimo postopek.

Ta faza se zaključi v milisekundah tudi za grafe z več sto vozlišči.

### Faza 2: Razvejaj in omeji (Branch-and-Bound)

Nato poženemo točen algoritem, ki sistematično preiskuje vse možne inducirane ujemanja.
Začetna požrešna rešitev služi kot spodnja meja za rezanje vej.

Algoritem izbere vozlišče z najmanjšo stopnjo in se razveji:
- za vsako sosednje vozlišče poskusi vključiti to povezavo,
- ali pa vozlišče izključi iz nadaljnjega iskanja.

Pred vsakim razvejanjem se uporabita dve redukciji:
- odstranitev izoliranih vozlišč (stopnja 0),
- vsiljevanje izoliranih povezav (obe krajišči imata stopnjo 1).

Za rezanje vej se izračunata dve zgornji meji: ⌊|aktivna vozlišča| / 2⌋ in število preostalih povezav.
Če trenutna rešitev + zgornja meja ne preseže najboljše znane rešitve, vejo odrežemo.

Če algoritem ne konča v časovni omejitvi (privzeto 300s), vrne najboljšo rešitev, ki jo je našel.

## Podatkovne strukture

Sosednostna matrika je shranjena kot polje bitnih nizov (`uint64` polja). To omogoča, da
operacije nad soseščinami (unija, presek, odstranitev) tečejo v O(n/64) namesto O(n),
kar opazno pospeši tako požrešni algoritem kot preiskovanje.
