# Navodila za prevajanje in uporabo

## Zahteve

- Go 1.21 ali novejši
- Python 3 (za generiranje testov in preverjanje)

## Prevajanje

```bash
go build -o solver solver.go
```

To ustvari izvršljivo datoteko `solver`.

## Uporaba

```bash
./solver <vhodna_datoteka> <izhodna_datoteka> [časovna_omejitev_v_sekundah]
```

Tretji argument je opcijski in nastavi časovno omejitev za fazo Branch-and-Bound (privzeto 300 sekund).

Primeri:
```bash
./solver tests/example.in tests/example.out
./solver tests/random_200_sparse.in output.txt 60
```

Program izpiše na stderr vrstico z informacijami o rešitvi:
```
Solution size: 3, nodes: 11, time: 252.541us, timed_out: false
```

## Generiranje testnih primerov

```bash
python3 gen.py          # generira teste v tests/
python3 gen.py moji     # generira teste v mapo moji/
```

## Preverjanje rešitev

Za preverjanje posamezne rešitve:
```bash
python3 verify.py <vhodna_datoteka> <izhodna_datoteka>
```

Za preverjanje vseh testov naenkrat:
```bash
bash run_all.sh
```

Ta skripta požene solver na vseh testih v `tests/` in za vsakega preveri pravilnost z `verify.py`.

## Format vhoda

Prva vrstica vsebuje celo število n (število vozlišč). Sledi n vrstic, vsaka z n vrednostmi
0 ali 1, ločenimi s presledki — sosednostna matrika grafa.

## Format izhoda

Prva vrstica vsebuje število c (število najdenih parov). Sledi c vrstic, vsaka z dvema
številkama (1-indeksirani) — krajišči izbrane povezave.
