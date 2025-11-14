//go:build postgres

package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"FZU-DB-ALPortManager/model"
	"FZU-DB-ALPortManager/query"

	"FZU-DB-ALPortManager/pkg/constants"
)

func Init() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		constants.PostgresHost,
		constants.PostgresUser,
		constants.PostgresPassword,
		constants.PostgresDBName,
		constants.PostgresPort,
		constants.PostgresSSLMode,
		constants.PostgresTimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[DB] Failed to connect to database: %v", err)
	}

	query.SetDefault(db)

	err = db.AutoMigrate(
		&model.ShipInfo{},
		&model.Equipment{},
		&model.ShipEquipment{},
	)
	if err != nil {
		log.Fatalf("[DB] Failed to auto migrate: %v", err)
	}

	// 初始化装备数据
	initEquipments(db)

	log.Println("[DB] Database connected and initialized successfully")
}

// initEquipments 初始化装备数据
func initEquipments(db *gorm.DB) {
	equipments := []model.Equipment{
		// 小口径主炮（驱逐炮）
		{Name: "双联100mm98式高射炮改T0", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "双联装138.6mm主炮Mle1934T0", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "双联装127mm高平两用炮Mk12T3", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "试作型双联装137mm高平两用炮Mk1T0", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "双联装120mm高平两用炮Mark XIT0", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "双联装114mm高平两用炮Mark IVT0", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "单装113mm高平两用炮Mark IVT0", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "双联100mm98式高射炮T3", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "127mm连装炮（D型）T0", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "双联装128mmSKC41高平两用炮改T0", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "试作型双联装130mm主炮Model1936T0", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "双联装135mm主炮Model1938T3", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "B-13 双联装130mm主炮B-2LMT3", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "138.6mm单装炮Mle1929T3", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "12磅长管炮T3", Type: "SMALL_CALIBER_MAIN_GUN"},
		{Name: "旧式重火炮T3", Type: "SMALL_CALIBER_MAIN_GUN"},

		// 中口径主炮（轻巡炮）
		{Name: "试作型四联装152mm主炮T0", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "三联装152mm主炮Mk16T0", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "试作型三联装152mm高平两用炮Mk17T0", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "试作型三联装152mm主炮T0", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "155mm三连装炮T3", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "试作型155mm三连装炮改T0", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "试作型三联装150mm五式高平两用炮T0", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "试作型双联装SKC28式150mm主炮改T0", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "试作型双联装TbtsKC42T式150mm主炮T0", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "三联装152mm主炮Model1934T3", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "试作型三联装152mm主炮Model1936T0", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "B-38 三联装152mm主炮MK-5T3", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "试作型B-1-P 三联装180mm主炮Model1932改T0", Type: "MEDIUM_CALIBER_MAIN_GUN"},
		{Name: "三联装152mm主炮Mle1930(高爆弹)T0", Type: "MEDIUM_CALIBER_MAIN_GUN"},

		// 大口径主炮（重巡炮）
		{Name: "试作型三联装234mm主炮T0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "试作型203mm/55三连装主炮T0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "试作型三联装203mmSKC主炮改T0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "三联装203mm主炮Mk15T0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "试作型双联装234mm主炮T0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "试作型三联装203mm主炮Mark IXT0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "试作型三联装203mm主炮Mark XT0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "试作型203mm(3号)连装炮T0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "双联装203mmSKC主炮T3", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "试作型三联装203mmSKC主炮T0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "双联203mmSKC主炮改T0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "双联203mm主炮Model1927T3", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "试作型三联装254mm主炮Model1939T0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "试作型三联装240mm主炮T0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "试作型三联装203mm舰炮T0", Type: "LARGE_CALIBER_MAIN_GUN"},
		{Name: "双联装203mm主炮Mle1931T3", Type: "LARGE_CALIBER_MAIN_GUN"},

		// 战列炮
		{Name: "三联装406mm主炮Mk7T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型双联装457mm主炮MkAT0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型三联装406mm主炮Mk.IIT0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "460mm三连装炮T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型四联装305mmSKC39主炮T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型三联装406mm主炮Model1940改T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型三联装406mm主炮MkDT0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型三联装406mm/45主炮Mk7T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型双联装406mm主炮Mk4T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型三联装406mm主炮Mk6改T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "四联装356mm主炮T3", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型三联装381mm主炮T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "三联装406mm主炮T3", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "双联装381mm主炮改T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "410mm连装炮(三式弹)T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型410mm三连装炮T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型双联装410mm主炮Mod.AT0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "双联380mmSKC主炮T3", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型双联装406mmSKC主炮T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型三联装305mmSKC39主炮T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型三联装380mmSKC主炮T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "三联装381mm主炮Model1934T3", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型三联装406mm主炮Model1940T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "B-37 三联装406mm主炮MK-1T3", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "四联装380mm主炮Mle1935T3", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型三联装380mm主炮Mle1935T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "试作型三联装406mm/50主炮T0", Type: "BATTLESHIP_MAIN_GUN"},
		{Name: "旧式半加农炮T3", Type: "BATTLESHIP_MAIN_GUN"},

		// 水面鱼雷
		{Name: "五联装610mm鱼雷T0", Type: "TORPEDO"},
		{Name: "五联装533mm磁性鱼雷T3", Type: "TORPEDO"},
		{Name: "五联装533mm鱼雷T3", Type: "TORPEDO"},
		{Name: "五联装533mm鱼雷Mk17T0", Type: "TORPEDO"},
		{Name: "533mm鱼雷Mk35(4连发射)T0", Type: "TORPEDO"},
		{Name: "五联装533mm鱼雷Mark IXT0", Type: "TORPEDO"},
		{Name: "四联装610mm鱼雷T3", Type: "TORPEDO"},
		{Name: "四联装610mm鱼雷改T0", Type: "TORPEDO"},
		{Name: "四联装533mm磁性鱼雷T3", Type: "TORPEDO"},
		{Name: "四联装533mm磁性鱼雷改T0", Type: "TORPEDO"},
		{Name: "五联装533mm磁性鱼雷T2", Type: "TORPEDO"},
		{Name: "试作型四联装533毫米鱼雷Si 270T0", Type: "TORPEDO"},
		{Name: "四联装550mm鱼雷改（弹药调整）T0", Type: "TORPEDO"},

		// 潜艇鱼雷
		{Name: "潜艇用95式纯氧鱼雷改T0", Type: "SUBMARINE_TORPEDO"},
		{Name: "潜艇用Mark 16鱼雷T3", Type: "SUBMARINE_TORPEDO"},
		{Name: "潜艇用Mark 28鱼雷T0", Type: "SUBMARINE_TORPEDO"},
		{Name: "潜艇用Mark 12鱼雷-菲里T0", Type: "SUBMARINE_TORPEDO"},
		{Name: "潜艇用Mark 20 S鱼雷-彼得T0", Type: "SUBMARINE_TORPEDO"},
		{Name: "潜艇用95式纯氧鱼雷T3", Type: "SUBMARINE_TORPEDO"},
		{Name: "潜艇用96式纯氧鱼雷T3", Type: "SUBMARINE_TORPEDO"},

		// 防空炮
		{Name: "双联装76mmRF火炮Mk37T0", Type: "ANTI_AIR_GUN"},
		{Name: "双联装57mmL/60博福斯对空机炮Mle1951T0", Type: "ANTI_AIR_GUN"},
		{Name: "双联装76mmRF火炮Mk27T0", Type: "ANTI_AIR_GUN"},
		{Name: "四联40mm博福斯对空机炮T3", Type: "ANTI_AIR_GUN"},
		{Name: "双联装127mm高平两用炮Mk12(定时引信)T0", Type: "ANTI_AIR_GUN"},
		{Name: "双联装134mm高炮T0", Type: "ANTI_AIR_GUN"},
		{Name: "八联装40mm砰砰炮T3", Type: "ANTI_AIR_GUN"},
		{Name: "双联装113mm高射炮T3", Type: "ANTI_AIR_GUN"},
		{Name: "双联装40mm博福斯STAAGT0", Type: "ANTI_AIR_GUN"},
		{Name: "双联装40mm博福斯海兹梅耶T0", Type: "ANTI_AIR_GUN"},
		{Name: "六联装40mm博福斯对空机炮T0", Type: "ANTI_AIR_GUN"},
		{Name: "双联装134mm高炮(定时引信)T0", Type: "ANTI_AIR_GUN"},
		{Name: "九六式25mm三连装暴风避盾机炮T0", Type: "ANTI_AIR_GUN"},
		{Name: "100mm连装高炮T0", Type: "ANTI_AIR_GUN"},
		{Name: "127mm连装高角炮改T0", Type: "ANTI_AIR_GUN"},
		{Name: "试作型五式40mm高射机关炮T0", Type: "ANTI_AIR_GUN"},
		{Name: "127mm连装高角炮改(定时引信)T0", Type: "ANTI_AIR_GUN"},
		{Name: "80mm98式连装高炮T0", Type: "ANTI_AIR_GUN"},
		{Name: "双联105mmSKC高炮T3", Type: "ANTI_AIR_GUN"},
		{Name: "双联105mmSKC高炮改进型T0", Type: "ANTI_AIR_GUN"},
		{Name: "双联105mmSKC高炮改进型(定时引信)T0", Type: "ANTI_AIR_GUN"},
		{Name: "试作型四联装30mm机炮T0", Type: "ANTI_AIR_GUN"},
		{Name: "试作型55mm Gerät 58防空炮T0", Type: "ANTI_AIR_GUN"},
		{Name: "90mm单装高角炮Model1939T3", Type: "ANTI_AIR_GUN"},
		{Name: "试作型双联90mm高角炮Model1939T0", Type: "ANTI_AIR_GUN"},
		{Name: "试作型六联装Scotti20mm机炮Model1941T0", Type: "ANTI_AIR_GUN"},
		{Name: "双联37mm高射炮Mle1936T0", Type: "ANTI_AIR_GUN"},

		// 设备
		{Name: "维修工具T3", Type: "AUXILIARY"},
		{Name: "海军部火控台T0", Type: "AUXILIARY"},
		{Name: "Z旗T0", Type: "AUXILIARY"},
		{Name: "九三式纯氧鱼雷T3", Type: "AUXILIARY"},
		{Name: "开拓者奖章T0", Type: "AUXILIARY"},
		{Name: "华盛顿海军条约T0", Type: "AUXILIARY"},
		{Name: "舰艇维修设备T3", Type: "AUXILIARY"},
		{Name: "液压弹射装置T3", Type: "AUXILIARY"},
		{Name: "改良声呐T3", Type: "AUXILIARY"},
		{Name: "异世界冒险终端T0", Type: "AUXILIARY"},
		{Name: "炽烈之歌T0", Type: "AUXILIARY"},
		{Name: "活力之歌T0", Type: "AUXILIARY"},
		{Name: "闪耀之歌T0", Type: "AUXILIARY"},
		{Name: "引力舞鞋T0", Type: "AUXILIARY"},
		{Name: "星云舞裙T0", Type: "AUXILIARY"},
		{Name: "Alizarin应援毛巾T0", Type: "AUXILIARY"},
		{Name: "Cyanidin应援毛巾T0", Type: "AUXILIARY"},
		{Name: "小海狸中队队徽T0", Type: "AUXILIARY"},
		{Name: "珍珠之泪T0", Type: "AUXILIARY"},
		{Name: "超重弹T0", Type: "AUXILIARY"},
		{Name: "白鹰精英损管T0", Type: "AUXILIARY"},
		{Name: "高性能火控雷达T0", Type: "AUXILIARY"},
		{Name: "SG雷达T3", Type: "AUXILIARY"},
		{Name: "高性能舵机T0", Type: "AUXILIARY"},
		{Name: "速运高速无人机T0", Type: "AUXILIARY"},
		{Name: "纳尔逊的旗语T0", Type: "AUXILIARY"},
		{Name: "6CRH穿甲弹T0", Type: "AUXILIARY"},
		{Name: "高性能对空雷达T0", Type: "AUXILIARY"},
		{Name: "小王冠T0", Type: "AUXILIARY"},
		{Name: "治愈系猫爪T0", Type: "AUXILIARY"},
		{Name: "一式穿甲弹T0", Type: "AUXILIARY"},
		{Name: "九三式纯氧鱼雷T2", Type: "AUXILIARY"},
		{Name: "VH装甲钢板T0", Type: "AUXILIARY"},
		{Name: "约定的证明T0", Type: "AUXILIARY"},
		{Name: "FuMO 25T0", Type: "AUXILIARY"},
		{Name: "533mm磁性鱼雷T3", Type: "AUXILIARY"},
		{Name: "改良型水下进气管T0", Type: "AUXILIARY"},
		{Name: "改良蓄电池阵列T0", Type: "AUXILIARY"},
		{Name: "四神之印T0", Type: "AUXILIARY"},
		{Name: "梅之语T0", Type: "AUXILIARY"},
		{Name: "妖精魔法海报T0", Type: "AUXILIARY"},
		{Name: "天使之羽T0", Type: "AUXILIARY"},
		{Name: "珍贵货物箱T0", Type: "AUXILIARY"},
		{Name: "随机单词生成器T0", Type: "AUXILIARY"},
		{Name: "晃悠悠T0", Type: "AUXILIARY"},
		{Name: "智慧模块T0", Type: "AUXILIARY"},
		{Name: "组徽T0", Type: "AUXILIARY"},
		{Name: "Gamers的证明T0", Type: "AUXILIARY"},
		{Name: "玉米灯笼T0", Type: "AUXILIARY"},
		{Name: "鮟鱇肝T0", Type: "AUXILIARY"},
		{Name: "觉醒宝珠T0", Type: "AUXILIARY"},
		{Name: "心之钥匙T0", Type: "AUXILIARY"},
		{Name: "偶像手环T0", Type: "AUXILIARY"},
		{Name: "征战巨坦T0", Type: "AUXILIARY"},
		{Name: "古立特圣剑T0", Type: "AUXILIARY"},
		{Name: "爆裂钻孔机T0", Type: "AUXILIARY"},
		{Name: "苍穹喷射机T0", Type: "AUXILIARY"},
		{Name: "戴拿爆能加农T0", Type: "AUXILIARY"},
		{Name: "煌翼炎龙T0", Type: "AUXILIARY"},
		{Name: "炙烈炎烧T0", Type: "AUXILIARY"},
		{Name: "结晶冰精T0", Type: "AUXILIARY"},
		{Name: "震耳雷球T0", Type: "AUXILIARY"},
		{Name: "涡旋风精T0", Type: "AUXILIARY"},
		{Name: "万灵药剂T0", Type: "AUXILIARY"},
		{Name: "神秘的羽衣T0", Type: "AUXILIARY"},
		{Name: "默示录T0", Type: "AUXILIARY"},
		{Name: "N/AT0", Type: "AUXILIARY"},
		{Name: "创世之槌T0", Type: "AUXILIARY"},
		{Name: "泡云弹车T0", Type: "AUXILIARY"},
		{Name: "形意口琴T0", Type: "AUXILIARY"},
		{Name: "忍者大师徽章T0", Type: "AUXILIARY"},
		{Name: "忍者装束T0", Type: "AUXILIARY"},
		{Name: "忍术卷轴T0", Type: "AUXILIARY"},
		{Name: "咻咻料理君T0", Type: "AUXILIARY"},
		{Name: "嗡嗡倾听君T0", Type: "AUXILIARY"},
		{Name: "嘘嘘隐身君T0", Type: "AUXILIARY"},
		{Name: "嘻嘻加班君T0", Type: "AUXILIARY"},
		{Name: "嘭嘭速生君T0", Type: "AUXILIARY"},
		{Name: "噗噗氛围君T0", Type: "AUXILIARY"},
		{Name: "辉辉光环君T0", Type: "AUXILIARY"},
		{Name: "宏伟光辉T0", Type: "AUXILIARY"},
		{Name: "高级魔导书T0", Type: "AUXILIARY"},
		{Name: "最终陨石T0", Type: "AUXILIARY"},
		{Name: "神药球T0", Type: "AUXILIARY"},
		{Name: "地狱立方体T0", Type: "AUXILIARY"},
		{Name: "天恩浑仪T0", Type: "AUXILIARY"},
		{Name: "芙拉米T0", Type: "AUXILIARY"},
		{Name: "刺猬弹T0", Type: "AUXILIARY"},
		{Name: "九四式40厘米炮部件T0", Type: "AUXILIARY"},
		{Name: "帝王蟹（附带小票）T0", Type: "AUXILIARY"},
		{Name: "木龙雕塑T0", Type: "AUXILIARY"},

		// 战斗机
		{Name: "试作舰载型La-9T0", Type: "FIGHTER"},
		{Name: "F2A水牛(萨奇队)T0", Type: "FIGHTER"},
		{Name: "F4U(VF-17海盗中队)T0", Type: "FIGHTER"},
		{Name: "F6F地狱猫T3", Type: "FIGHTER"},
		{Name: "F7F虎猫T0", Type: "FIGHTER"},
		{Name: "F8F熊猫T0", Type: "FIGHTER"},
		{Name: "试作型XF5U飞碟T0", Type: "FIGHTER"},
		{Name: "F6F地狱猫(HVAR搭载型)T0", Type: "FIGHTER"},
		{Name: "海喷火FR.47T0", Type: "FIGHTER"},
		{Name: "海毒牙T3", Type: "FIGHTER"},
		{Name: "海怒T0", Type: "FIGHTER"},
		{Name: "海大黄蜂T0", Type: "FIGHTER"},
		{Name: "零战五二型T3", Type: "FIGHTER"},
		{Name: "烈风T3", Type: "FIGHTER"},
		{Name: "紫电改二T0", Type: "FIGHTER"},
		{Name: "试作型紫电改四T0", Type: "FIGHTER"},
		{Name: "Me-155A舰载战斗机T3", Type: "FIGHTER"},
		{Name: "试作舰载型BF-109GT0", Type: "FIGHTER"},
		{Name: "试作舰载型FW-190 A-6/R6T0", Type: "FIGHTER"},
		{Name: "试作舰载型FW-190 G-3/R1T0", Type: "FIGHTER"},

		// 轰炸机
		{Name: "AD-1天袭者T0", Type: "DIVE_BOMBER"},
		{Name: "试作舰载型天雷T0", Type: "DIVE_BOMBER"},
		{Name: "SBD无畏(麦克拉斯基队)T0", Type: "DIVE_BOMBER"},
		{Name: "实验型XSB3C-1T0", Type: "DIVE_BOMBER"},
		{Name: "BTD-1毁灭者T3", Type: "DIVE_BOMBER"},
		{Name: "萤火虫T0", Type: "DIVE_BOMBER"},
		{Name: "萤火虫(1771中队)T0", Type: "DIVE_BOMBER"},
		{Name: "梭鱼(831中队)T0", Type: "DIVE_BOMBER"},
		{Name: "彗星T3", Type: "DIVE_BOMBER"},
		{Name: "彗星一二型甲T0", Type: "DIVE_BOMBER"},
		{Name: "彗星二一型T0", Type: "DIVE_BOMBER"},
		{Name: "试作舰载型Su-2T0", Type: "DIVE_BOMBER"},
		{Name: "LN.401T0", Type: "DIVE_BOMBER"},

		// 鱼雷机
		{Name: "飞龙T0", Type: "TORPEDO_BOMBER"},
		{Name: "试作型旗鱼T0", Type: "TORPEDO_BOMBER"},
		{Name: "TBD蹂躏者(VT-8中队)T0", Type: "TORPEDO_BOMBER"},
		{Name: "TBM复仇者(VT-18中队)T0", Type: "TORPEDO_BOMBER"},
		{Name: "XTB2D-1天空海盗T0", Type: "TORPEDO_BOMBER"},
		{Name: "剑鱼(818中队)T0", Type: "TORPEDO_BOMBER"},
		{Name: "梭鱼T3", Type: "TORPEDO_BOMBER"},
		{Name: "火把T0", Type: "TORPEDO_BOMBER"},
		{Name: "火冠T0", Type: "TORPEDO_BOMBER"},
		{Name: "流星T3", Type: "TORPEDO_BOMBER"},
		{Name: "流星改T0", Type: "TORPEDO_BOMBER"},
		{Name: "试作型彩云(舰攻型)T0", Type: "TORPEDO_BOMBER"},
		{Name: "Ju-87 D-4T0", Type: "TORPEDO_BOMBER"},
		{Name: "试作型VIT-2(VK-107)T0", Type: "TORPEDO_BOMBER"},
		{Name: "试作型VIT-2（模式调整）T0", Type: "TORPEDO_BOMBER"},
		{Name: "BR.810T0", Type: "TORPEDO_BOMBER"},

		// 兵装（AUGMENT）- 满星才能使用
		// 驱逐兵装
		{Name: "双剑", Type: "AUGMENT"},
		{Name: "单手锤", Type: "AUGMENT"},
		// 轻巡兵装
		{Name: "铁剑", Type: "AUGMENT"},
		{Name: "手弩", Type: "AUGMENT"},
		// 重巡兵装
		{Name: "大剑", Type: "AUGMENT"},
		{Name: "骑枪", Type: "AUGMENT"},
		// 战列/战巡兵装
		{Name: "指挥刀", Type: "AUGMENT"},
		{Name: "轻弩", Type: "AUGMENT"},
		// 轻航/航母兵装
		{Name: "猎弓", Type: "AUGMENT"},
		{Name: "权杖", Type: "AUGMENT"},
		// 维修兵装
		{Name: "维修手弩", Type: "AUGMENT"},
		// 潜艇兵装
		{Name: "短剑", Type: "AUGMENT"},
		{Name: "若无", Type: "AUGMENT"},
	}

	for _, equip := range equipments {
		// 检查装备是否已存在
		var count int64
		db.Model(&model.Equipment{}).Where("name = ?", equip.Name).Count(&count)
		if count == 0 {
			// 装备不存在，创建它
			if err := db.Create(&equip).Error; err != nil {
				log.Printf("[DB] Failed to create equipment %s: %v", equip.Name, err)
			} else {
				log.Printf("[DB] Created equipment: %s", equip.Name)
			}
		}
	}

	log.Println("[DB] Equipment initialization completed")
}