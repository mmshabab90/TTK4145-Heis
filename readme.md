Morten Late Night Tales
=======================

Datastrukturer
--------------

Alle heiser har en kopi av en struct (eller noe) som inneholder:
- Alle ytre bestillinger
- For hver bestilling: Hvilken heis som er satt til å utføre den
- Liste over aktive heiser på nettverket
- Hvor lenge det er siden heisene har gitt lyd fra seg

Denne skal 'alltid' være lagret lokalt i alle levende heiser, og være identisk i alle heiser.

I tillegg skal alle heiser ha en separat, utelukkende lokal liste over indre bestillinger.

Når en ny ytre bestilling kommer skal følgende skje:
- Send bestillingen på nettverket
- Masterheisen legger den i sin kopi av køen
- Masterheisen ber alle heiser om kostnad for ny bestilling
- Masterheisen 'assigner' en heis til den nye bestillingen
- Masterheisen sender ut den oppdaterte køen
- Heisen som mottok bestillingen fortsetter å sende den frem til den får en oppdatert kø med bestillingen lagt til
- Heisene fortsetter å kjøre og håndtere bestillinger på samme måte som soloheis-koden

Kostfunksjon
------------

Kostnad én for hver etasje heisen må gå og hver etasje heisen skal stoppe i på veien mot målet.
