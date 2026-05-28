---
marp: true
theme: default
paginate: true
size: 16:9
header: 'Konstrukcija maksimalne ortogonalne CC-množice'
footer: 'Algoritmi 2026'
style: |
  section { font-size: 26px; }
  h1 { color: #1a365d; }
  h2 { color: #2c5282; }
  code { background: #edf2f7; padding: 2px 6px; border-radius: 3px; }
  table { font-size: 22px; }
  .small { font-size: 20px; }
---

# Konstrukcija maksimalne ortogonalne CC-množice

### Maksimalno inducirano ujemanje

Govorca: **Speaker A** in **Speaker B**

Algoritmi, 2026

---

## Naš pristop

- **NP-težko** — že na ravninskih bipartitnih grafih z največjo stopnjo 3
- Točna rešitev za velike grafe ni dosegljiva v razumnem času

### Dvofazni pristop

```
┌─────────────────────┐     ┌──────────────────────┐
│  Faza 1: Hevristika │ ──> │  Faza 2: B&B (točno) │
│  (požrešni + LS)    │     │  z mejo iz hevristike│
│  ms — vedno konča   │     │  do časovne omejitve │
└─────────────────────┘     └──────────────────────┘
```

Hevristika nam **vedno** da spodnjo mejo (in s tem rešitev).
B&B jo **izboljša** do optimuma, kadar je čas dovolj.

---

## Faza 1a — Požrešni algoritem

V vsaki iteraciji izberemo povezavo z **najmanjšim stroškom**:

$$ \text{score}(u,v) = |N(u) \cup N(v) \cap \text{active}| $$

```
ponavljaj:
  za vsako povezavo (u,v) v aktivnem grafu: izračunaj score
  izberi (u,v) z najmanjšim score
  dodaj (u,v) v S
  odstrani N(u) ∪ N(v) iz active
dokler v aktivnem grafu ne ostane nobene povezave
```

**Intuicija:** $|N(u) \cup N(v)|$ pove, koliko vozlišč "porabimo".
Z minimiziranjem porabe ohranimo največ vozlišč za nadaljnje povezave.

Časovna zahtevnost: $O(m \cdot k \cdot w)$, kjer je $w = \lceil n/64 \rceil$.

---

## Faza 1b — Lokalno iskanje: (1,2)-zamenjave

```
ponavljaj:
  za vsako povezavo e v S:
    začasno odstrani e iz S
    poišči povezave, ki jih lahko dodamo
    če najdemo ≥ 2 povezave: sprejmi zamenjavo (+1)
dokler obstaja izboljšava
```

**Zakaj 1 ven, 2 noter?** Neto +1 v velikosti rešitve.

Sprostitev enega para krajišč pogosto odpre dve mesti, kjer je bilo prej eno.

Empirično: izboljša 10–30 % primerov pred B&B.

---

<!-- _header: '' -->

## ⇨ Preklop na govorca B

# Faza 2: Razvejaj in omeji

Sistematično preiskovanje vseh možnih induciranih ujemanj
**z rezanjem vej**, ki ne morejo izboljšati najboljše znane rešitve.

### Strategija razvejanja

V vsakem koraku izberemo vozlišče $v$ z **najmanjšo stopnjo** v aktivnem grafu.

- **Vključitvene veje** (po ena za vsakega soseda $w$): dodaj povezavo $(v,w)$ v $S$,
  odstrani $N(v) \cup N(w)$ iz aktivnega grafa.
- **Izključitvena veja**: vozlišča $v$ ne uporabimo, odstrani samo $v$.

Skupaj $\deg(v) + 1$ podproblemov — minimalna stopnja minimizira faktor razvejanja.

---

## Redukcije + rezanje vej

### Redukcije (varne za MIM)

1. **Izolirano vozlišče** (stopnja 0) → odstrani
2. **Izolirana povezava** (obe krajišči stopnje 1) → **prisilno v $S$**

> ⚠️ **Pozor:** Klasična redukcija "obeska" (vozlišče stopnje 1) **ni varna** za MIM!
> Protiprimer: pot $v - u - a - b$ z dodatno povezavo $u - c - d$ — prisila $(v,u)$
> da $|S|=1$, optimum pa je $\{(a,b),(c,d)\}$ z $|S|=2$.

### Zgornji meji za rezanje

$$ \text{ub} = \min\left(\left\lfloor \tfrac{|\text{active}|}{2}\right\rfloor,\ |E_\text{active}|\right) $$

Veja je odrezana, če $\text{curSize} + \text{ub} \le \text{best}$.

---

## Bitne strukture — ključ do hitrosti

```go
type Bitset [16]uint64    // do 1024 vozlišč
adj[v]: bitni niz sosedov vozlišča v
active: bitni niz trenutno aktivnih vozlišč
```

Operacije nad celotno soseščino tečejo v $O(n/64)$:

| Operacija | Pomen | Implementacija |
|-----------|-------|----------------|
| `active.AndNot(&adj[v])` | odstrani $N(v)$ | bitni AND-NOT |
| `a.And(&b).PopCount()` | $\lvert A \cap B \rvert$ | AND + POPCNT |
| `bs.PopCount()` | $\lvert S \rvert$ | strojni `POPCNT` |

**POPCNT** je en strojni ukaz na modernih procesorjih — kritičen za scoring v B&B.

---

## Rezultati

| Test | $n$ | Tip | $\lvert S \rvert$ | Čas | Točno? |
|------|-----|------|-----|------|--------|
| primer (PDF) | 7 | iz navodil | 3 | <1ms | ✅ |
| path_12 | 12 | pot | 4 | <1ms | ✅ |
| cycle_9 | 9 | cikel | 3 | <1ms | ✅ |
| petersen | 10 | Petersen | 3 | <1ms | ✅ |
| random_50 | 50 | redko ($p=0.05$) | 13 | 16 ms | ✅ |
| random_100 | 100 | srednje | 17 | 60 s | hevristika |
| random_200 | 200 | redko | 41 | 60 s | hevristika |
| random_500 | 500 | redko | 80 | 60 s | hevristika |

**Zid:** točno rešimo do $n \approx 50$ za srednjo gostoto, do $n \approx 200$ za zelo redke grafe.
Za večje vrne hevristika — vedno v milisekundah.

---

## Povzetek

- **Hibridni pristop**: hevristika garantira hitro rešitev, B&B jo izboljša do optimuma
- **Min-degree razvejanje + redukcije + dve zgornji meji** poskrbijo za učinkovito rezanje
- **Bitni nizi z `POPCNT`** dajo velik pospešek pri vseh operacijah nad soseščinami
- Iskreno o NP-težkosti: točno do $\sim 50$ vozlišč, naprej kvalitetna hevristika

### Možne izboljšave

- Boljša zgornja meja (klikno pokritje, frakcijsko LP)
- Paralelizacija B&B vej
- Inkrementalni scoring po posodobitvi `active`

---

# Vprašanja?

<br>

Hvala za pozornost!

<br>

Repozitorij: `red-group/` &nbsp;·&nbsp; Solver: `solver.go` &nbsp;·&nbsp; Dokumentacija: `docs/`
