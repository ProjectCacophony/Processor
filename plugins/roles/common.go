package roles

func getPoolLimitsFromCategories(categories []*Category) map[string]int {

	// pool name =>
	poolMap := make(map[string]int)

	for _, cat := range categories {
		if cat.Pool == "" {
			continue
		}

		if _, ok := poolMap[cat.Pool]; ok {
			if cat.Limit > poolMap[cat.Pool] {
				poolMap[cat.Pool] = cat.Limit
			}
		} else {
			poolMap[cat.Pool] = cat.Limit
		}
	}

	return poolMap
}
