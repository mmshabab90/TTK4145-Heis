Morten Late Night Tales
=======================

Datastrukturer
--------------

Alle heiser har kopier av:
- Struct med alle eksterne bestillinger i systemet og hvilke heiser hver bestilling er gitt til
- Liste over aktive heiser på nettverket og hvor lenge det er siden de har gitt lyd fra seg

Denne skal 'alltid' være lagret lokalt i alle levende heiser, og være identisk i alle heiser på et nettverk.

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

Antakelser
----------

- Det er kool om to heiser bestilles til samme etasje for å ekspedere bestillinger i ulike retninger. Den som henter folk som vil opp slukker kun lyset for bestilling opp, og vice versa.