**Redis**

***Proxy***
* api - mõeldud väljaspoole riigipilve
* clients - saab vastavalt kliendi käest väärtused, vajalik redise väärtustamiseks
* init - algväärdustab rakenduse
* main - proxy käivitamiseks vajalik
* proxy - redise socket, kontrollib käsklusi, mis tuleb tcp socketi pealt, ainult riigipilvest
* redis - hoiab redis pool'i
* sentinel - hoiab ühendust sentineliga, saab vastava redis masteri

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

**Proxy**
* Socket - Lubatud ainult riigipilve sees, kättesaadav namespace järgi
    * Mõeldud redise frameworkide jaoks, luuakse tcp ühendus
    * Lubatud käsklused:
        * GET
        * HGETALL
        * PING
        * DOCUMENT_CODE
        * PROCEDURE_CODE
        * BUILDING_CODE
        * UTILITY_BUILDING_CODE
        * EVAL (piiratud, oleneb frameworkist, kui framework ei võimalda kasutada custom commande)
            * return redis.call('PROCEDURE_CODE')
            * return redis.call('UTILITY_BUILDING_CODE')
            * return redis.call('DOCUMENT_CODE', 'doty_id')
* Api - Avatud väljaspoole (Oracle jaoks)
    * Token peab olema headeris X-Session-Token
        * /building - *Tagastab ehitise järjekorranumbri+1*
        * /utilitybuilding - *Tagastab rajatise järjekorranumbri+1*
        * /procedure - *Tagastab menetluse järjekorranumbri+1*
        * /document/{doty} - *Tagastab dokumendi tüübi id järjekorranumbri+1*