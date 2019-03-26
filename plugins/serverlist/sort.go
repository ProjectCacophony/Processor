package serverlist

type SortServers struct {
	SortBy  []SortBy
	Servers []*ServerToPost
}

func (s SortServers) Len() int {
	return len(s.Servers)
}

func (s SortServers) Swap(i, j int) {
	s.Servers[i], s.Servers[j] = s.Servers[j], s.Servers[i]
}

func (s SortServers) Less(i, j int) bool {
	for _, sortByItem := range s.SortBy {
		switch sortByItem {
		case SortByMemberCount:
			if s.Servers[i].Server.TotalMembers < s.Servers[j].Server.TotalMembers {
				return true
			}

		case SortByAddingDate:
			if s.Servers[i].Server.CreatedAt.After(s.Servers[j].Server.CreatedAt) {
				return true
			}

		case SortByAlphabetical:
			if s.Servers[i].Name > s.Servers[j].Name {
				return true
			}
		}
	}

	return false
}
