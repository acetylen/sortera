# Sortera

Sorterar alla filer som hittas i nuvarande katalog och underkataloger 
baserat på när de skapades, och städar upp efter sig.
En fil som skapats 2020-02-23 hamnar i "2020/Februari".

Detta program är till för att hjälpa en släkting att sortera sitt
bildbibliotek utan att behöva krångla med var saker är.

Steg:

1. Leta upp dubletter och ta bort alla utom den första 
   (bara om `--delete-duplicates` är satt)
2. Hitta alla filer vars katalog inte är på formen 
   `skapad_år/skapad_månad`, skapa katalogerna om de inte finns, och flytta
   filerna dit
3. Hitta och ta bort alla nu tomma kataloger

Om allt inte blir sorterat är det bara att köra igen, eftersom programmet
inte bryr sig om saker som ligger där de ska.
