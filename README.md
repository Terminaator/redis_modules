**Redis**

***Selgitus***

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

Kasutusel on kaks proxyt (default-proxy, ehrcode-proxy). Mida saab kasutada redise frameworkidega. Default-proxy on mõeldud kasutamiseks kõikideks tegevusteks redis-ga, mis ei ole seotud ehr koodidega. Ehr koodide jaoks tuleks kaustada ehrcode-proxy't.

**Rancher'is secret nimega "ehr-redis-proxy".**

ehrcode-proxy teenuse ip on secretis EHRCODE_HOST ja port on EHRCODE_PORT

default-proxy teenuse ip on secretis DEFAULT_HOST ja port on DEFAULT_PORT


**ehrcode-proxy**
* ehrcode-proxy (mõeldud EHR koodide jagamiseks. Kontrollib pidevalt, kas väärtused on olemas ja vajadusel väärtustab algväärtused.)
    * api (https://devkluster.ehr.ee/api/redis/v1)
        * /building
        * /utilitybuilding
        * /procedure
        * /document/{doty}
        * täpsemalt [https://swaggerui.mkm.ee/](https://swaggerui.mkm.ee/)
    * socket (TCP on mõeldud kasutamiseks riigipilves, väljaspoolt ligi ei saa)
        *  host - <teenuse_nimi>.<namespace> -> ehrcode-proxy.dev-redis
        *  port - 9999 on mõeldud ehrcode-proxy jaoks
    * secret (ehr-redis-proxy)
        * EHRCODE_HOST proxy teenuse ip
        * EHRCODE_PORT proxy teenuse port
    
**näide ehrcode-proxy teenuse kasutamisest (python)**

*import redis*

*r = redis.Redis(host="ehrcode-proxy.dev-redis", port=9999)*

*result = r.execute_command("PROCEDURE_CODE")*

*print(result)*

*r.close()*

**default-proxy**
* default-proxy (mõeldud üldiseks kasutamiseks, ei soovita seda kasutada EHR koodide saamiseks)
    * socket (TCP on mõeldud kasutamiseks riigipilves, väljaspoolt ligi ei saa)
        *  host - <teenuse_nimi>.<namespace> -> default-proxy.dev-redis
        *  port - 9998 on mõeldud ehrcode-proxy jaoks
    * secret (ehr-redis-proxy)
        * DEFAULT_HOST proxy teenuse ip
        * DEFAULT_PORT proxy teenuse port  
    
**näide default-proxy teenuse kasutamisest (python)**

*import redis*

*r = redis.Redis(host="default-proxy.dev-redis", port=9998)*

*result = r.execute_command("PROCEDURE_CODE")*

*print(result)*

*r.close()*

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