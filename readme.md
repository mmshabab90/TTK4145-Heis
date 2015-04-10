Morten Late Night Tales
=======================

Det egentlige spørsmålet er: Er Morten på jordet?

Datastrukturer
--------------

Alle heiser har en kopi av en struct (eller noe) som inneholder:
- Alle ytre bestillinger (totalt 6 mulige bestillinger, kanskje heller en 1x6-liste heller enn en 2x3-matrise?)
- For hver bestilling: Hvilken heis som er satt til å utføre den
- Liste over aktive heiser på nettverket
- Hvor lenge det er siden heisene har gitt lyd fra seg

Denne skal 'alltid' være lagret lokalt i alle levende heiser, og være identisk i alle heiser.

I tillegg skal alle heiser ha en separat, lokal liste over indre bestillinger.

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

Kostnad to for hver etasje heisen må gå og hver etasje heisen skal stoppe i på veien mot målet. Kostnad én for bevegelse fra mellom etasjer til en etasje. Summer kostnader, sammenlikn heiser, velg den beste (mega-duh). Bruk IP eller noe for å velge heis om flere har samme kostnad. Kanskje sjekk hvilken som har færrest ordre totalt?
