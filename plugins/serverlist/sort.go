package serverlist

type ServersSorter struct {
	SortBy  []SortBy
	Servers []*ServerToPost
}

func (s ServersSorter) Len() int {
	return len(s.Servers)
}

func (s ServersSorter) Swap(i, j int) {
	s.Servers[i], s.Servers[j] = s.Servers[j], s.Servers[i]
}

func (s ServersSorter) Less(i, j int) bool {
	for _, sortByItem := range s.SortBy {
		switch sortByItem {
		case SortByMemberCount:
			if s.Servers[i].Server.TotalMembers != s.Servers[j].Server.TotalMembers {
				return s.Servers[i].Server.TotalMembers < s.Servers[j].Server.TotalMembers
			}

		case SortByAddingDate:
			if s.Servers[i].Server.CreatedAt != s.Servers[j].Server.CreatedAt {
				return s.Servers[i].Server.CreatedAt.After(s.Servers[j].Server.CreatedAt)
			}

		case SortByAlphabetical:
			if s.Servers[i].Name != s.Servers[j].Name {
				return s.Servers[i].Name > s.Servers[j].Name
			}
		}
	}

	return false
}