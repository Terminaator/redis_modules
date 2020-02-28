**Redis**

***Selgitus***

Selles projektis asub sentineli, redise conf ja moodulid.

Proxy asub https://git.mkm.ee/ehr/ehr-k8s-pipeline/redis-sentinel-proxy.

**Proxy**

Proxy on mõeldud vahelülina redise ja redist kasutada sooviva kliendi vahel. Proxy eesmärk on pidevalt hoida Redise masteriks kuulutatud instantsi, mille ta saab sentineli käest küsides.
Proxy suunab edasi liikluse Redis masteri vastu. Proxy võimaldab kontrollida nii redise sisendit ja väljundit.
Antud konstektis kontrollitakse väljundit ehr koodidega seoses ning vajadusel väärtustatakse uuesti Redise poolt pakutavad EHR koodid.
Samuti, kui peaks tekkima uus master, siis uus master algväärtustatakse EHR koodidega. Proxy küljes on http api koodide pärimiseks, mida saab kasutada väljaspool riigipilve või kui pole soovi kasutada
TCP socketit riigipilves.
Mõeldud on kasutamiseks redise frameworkidega. Samuti töötab tavalise tcp socketina, kuid siis tuleks tutvuda https://redis.io/topics/protocol.

**Sentinel**

Sentineli eesmärk on tagada Rediste instantside kättesaadavus, hoides ühe instantsidest masterina ning hallates synci nende vahel. See tagab teenuse pideva kättesaadavuse.

***Ühenduse loomine***

Ühenduse loomine asub siin https://git.mkm.ee/ehr/ehr-k8s-pipeline/redis-proxy

***Moodulid***
* yearmodule:
    * Käsklus - *YEAR*
    * Võti - *YEAR_KEY*
    * Võtme väärtus - *Aasta kaks viimast numbrit*
    * Lubatud võtme väärtus - *Hetkel oleva aasta kaks viimast numbrit*
    * Tagastab - *Tagastab aasta kaks viimast numbrit*
* proceduremodule:
    * Käsklus - *PROCEDURE_CODE*
    * Võti - *PROCEDURE_KEY*
    * Võtme väärtus - *Menetluse järjekorranumber*
    * Lubatud võtme väärtus - *Hetkel oleva menetluse järjekorranumber*
    * Tagastab - *Tagastab menetluse järjekorranumbri+1*
* documentmodule:
    * Käsklus - *DOCUMENT_CODE <dokumendi tüübi id>*
    * Võti - *DOCUMENT_KEY*
    * Võtme väljad - *Dokumendi tüübi id*
    * Võtme väljade väärtused - *Dokumendi tüübi id järjekorranumber*
    * Lubatud väljade väärtuse - *Hetkel oleva dokumendi tüübi id järjekorranumber*
    * Tagastab - *Tagastab dokumendi tüübi id järjekorranumbri+1*
* buildingmodule:
    * Käsklus - *BUILDING_CODE*
    * Võti - *EHR_CODE_SET_KEY*
    * Võtme väli - *EHR_CODE_SET_BUILDING_FIELD*
    * Võtme välja väärtus - *Ehitise järjekorranumber*
    * Lubatud välja väärtus - *Hetkel oleva ehitise järjekorranumber vahemikus 100000000-200000000*
    * Tagastab - *Tagastab ehitise järjekorranumbri+1*
* utilitybuildingmodule:
    * Käsklus - *UTILITY_BUILDING_CODE*
    * Võti - *EHR_CODE_SET_KEY*
    * Võtme väli - *EHR_CODE_SET_UTILITY_BUILDING_FIELD*
    * Võtme välja väärtus - *Rajatise järjekorranumber*
    * Lubatud välja väärtus - *Hetkel oleva rajatise järjekorranumber vahemikus 200000000-300000000*
    * Tagastab - *Tagastab rajatise järjekorranumbri+1*