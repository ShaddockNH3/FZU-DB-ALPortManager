namespace go api

struct BaseResp {
    1: required i32 code (api.json="code")
    2: required string msg (api.json="msg")
}

struct Ship {
    1: optional i64 id (api.json="id")
    2: required string shipName (api.json="shipName")
    3: required string rarity (api.json="rarity")
    4: required string shipType (api.json="shipType")
    5: required string faction (api.json="faction")
    6: optional i32 level (api.json="level")
    7: optional string createdAt (api.json="createdAt")
    8: optional string updatedAt (api.json="updatedAt")
}

struct GetShipListReq {
    1: optional i32 page = 1 (api.query="page")
    2: optional i32 pageSize = 10 (api.query="pageSize")
    3: optional string shipName (api.query="shipName")
    4: optional string rarity (api.query="rarity")
    5: optional string shipType (api.query="shipType")
    6: optional string faction (api.query="faction")
}

struct GetShipListResp {
    1: required i64 total (api.json="total")
    2: required list<Ship> ships (api.json="ships")
    3: required i32 page (api.json="page")
    4: required i32 pageSize (api.json="pageSize")
}

struct StatItem {
    1: required string name (api.json="name")
    2: required i64 count (api.json="count")
}

struct GetStatisticsResp {
    1: required i64 totalShips (api.json="totalShips")
    2: required list<StatItem> factionStats (api.json="factionStats")
    3: required list<StatItem> rarityStats (api.json="rarityStats")
    4: required list<StatItem> shipTypeStats (api.json="shipTypeStats")
}

struct UpdateShipReq {
    1: required i64 id (api.path="id") 
    2: required Ship ship (api.body="ship")
}

service PortManagerService {
    Ship CreateShip(1: Ship req) (api.post="/api/v1/ship")
    Ship GetShipByID(1: i64 id) (api.get="/api/v1/ship/:id")
    GetShipListResp GetShipList(1: GetShipListReq req) (api.get="/api/v1/ships")
    Ship UpdateShip(1: UpdateShipReq req) (api.put="/api/v1/ship/:id")
    BaseResp DeleteShip(1: i64 id) (api.delete="/api/v1/ship/:id")
    GetStatisticsResp GetStatistics() (api.get="/api/v1/stats")
}
