package ship

// 稀有度分三档决定星级：普通 1 起 4 星封顶，精锐以上更高。
// 一张表两个函数，避免"初始星级"和"满星星级"两处 switch 各写一遍。
type starRule struct {
	initial int
	max     int
}

var starRules = map[string]starRule{
	"普通":   {initial: 1, max: 4},
	"稀有":   {initial: 2, max: 5},
	"精锐":   {initial: 2, max: 5},
	"超稀有":  {initial: 3, max: 6},
	"海上传奇": {initial: 3, max: 6},
	"最高方案": {initial: 3, max: 6},
	"决战方案": {initial: 3, max: 6},
}

const (
	fallbackInitialStars = 1
	fallbackMaxStars     = 3
)

// InitialStars 是新船入港时的星级；不认识的稀有度按最低档处理。
func InitialStars(rarity string) int {
	if rule, ok := starRules[rarity]; ok {
		return rule.initial
	}
	return fallbackInitialStars
}

// MaxStars 是该稀有度的星级上限。
func MaxStars(rarity string) int {
	if rule, ok := starRules[rarity]; ok {
		return rule.max
	}
	return fallbackMaxStars
}

// CanEquipAugment 判断当前星级是否已解锁兵装槽（满星才开）。
func CanEquipAugment(rarity string, stars int) bool {
	return stars >= MaxStars(rarity)
}
