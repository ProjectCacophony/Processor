package serverlist

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

var (
	regexpKoreanChineseJapanese = regexp.MustCompile("[\u3040-\u30ff\u3400-\u4dbf\u4e00-\u9fff\uf900-\ufaff\uff66-\uff9f\u3131-\uD79D]+")
	regexpAlphanumeric          = regexp.MustCompile("[^a-zA-Z0-9]+")
)

func getSortName(input string) string {
	return strings.ToLower(
		regexpAlphanumeric.ReplaceAllString(
			unidecode.Unidecode(
				regexpKoreanChineseJapanese.ReplaceAllString(input, ""),
			), "",
		))
}

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
				return s.Servers[i].Server.TotalMembers > s.Servers[j].Server.TotalMembers
			}

		case SortByAddingDate:
			if s.Servers[i].Server.CreatedAt != s.Servers[j].Server.CreatedAt {
				return s.Servers[i].Server.CreatedAt.After(s.Servers[j].Server.CreatedAt)
			}

		case SortByAlphabetical:
			if s.Servers[i].SortName != s.Servers[j].SortName {
				return s.Servers[i].SortName < s.Servers[j].SortName
			}
		}
	}

	return false
}
