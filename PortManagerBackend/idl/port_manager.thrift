namespace go api

struct BaseResp {
    1: required i32 code (api.json="code")
    2: required string msg (api.json="msg")
}

enum EquipmentType {
    SMALL_CALIBER_MAIN_GUN,
    MEDIUM_CALIBER_MAIN_GUN,
    LARGE_CALIBER_MAIN_GUN,
    BATTLESHIP_MAIN_GUN,
    TORPEDO,
    SUBMARINE_TORPEDO,
    ANTI_AIR_GUN,
    AUXILIARY,
    DIVE_BOMBER,
    TORPEDO_BOMBER,
    FIGHTER,
    AUGMENT
}

struct Equipment {
    1: required i64 id (api.json="id")
    2: required string name (api.json="name")
    3: required EquipmentType type (api.json="type")
}

struct EquipmentSlot {
    1: required list<EquipmentType> acceptableTypes (api.json="acceptableTypes")
    2: optional Equipment equippedItem (api.json="equippedItem")
}

struct Ship {
    1: optional i64 id (api.json="id")
    2: required string shipName (api.json="shipName")
    3: required string rarity (api.json="rarity")
    4: required string shipType (api.json="shipType")
    5: required string faction (api.json="faction")
    6: optional i32 level (api.json="level")
    7: optional i32 stars (api.json="stars") // 当前星级，满星后可装备兵装
    8: optional list<EquipmentSlot> equipmentSlots (api.json="equipmentSlots")
    9: optional EquipmentSlot augmentSlot (api.json="augmentSlot")
    10: optional string createdAt (api.json="createdAt")
    11: optional string updatedAt (api.json="updatedAt")
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

struct EquipShipReq {
    1: required i64 shipId (api.path="id")
    2: required i32 slotIndex (api.json="slotIndex")
    3: optional i64 equipmentId (api.json="equipmentId")
}

struct GetEquipmentListReq {
    1: optional EquipmentType type (api.query="type")
}

struct GetEquipmentListResp {
    1: required list<Equipment> equipments (api.json="equipments")
}

service PortManagerService {
    Ship CreateShip(1: Ship req) (api.post="/api/v1/ship")
    Ship GetShipByID(1: i64 id) (api.get="/api/v1/ship/:id")
    GetShipListResp GetShipList(1: GetShipListReq req) (api.get="/api/v1/ships")
    Ship UpdateShip(1: UpdateShipReq req) (api.put="/api/v1/ship/:id")
    BaseResp DeleteShip(1: i64 id) (api.delete="/api/v1/ship/:id")
    GetStatisticsResp GetStatistics() (api.get="/api/v1/stats")
    Ship EquipShip(1: EquipShipReq req) (api.put="/api/v1/ship/:id/equip")
    GetEquipmentListResp GetEquipmentList(1: GetEquipmentListReq req) (api.get="/api/v1/equipments")
}
