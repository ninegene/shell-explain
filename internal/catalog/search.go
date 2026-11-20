package catalog

import (
	"sort"
	"strings"
)

// ExampleHit is one authored example that matched a search query.
type ExampleHit struct {
	Command        string `json:"command"`
	Title          string `json:"title"`
	Implementation string `json:"implementation"`
	Description    string `json:"description"`
	Short          string `json:"short,omitempty"`
	Code           string `json:"code"`

	rank, order int
	hay, lead   string
}

// exampleIndex returns every example with code, in catalog and authored
// order, parsing the Markdown once.
func (c *Catalog) exampleIndex() []ExampleHit {
	c.examplesOnce.Do(func() {
		for _, e := range c.Entries {
			seen := map[string]bool{}
			// Common sections come first so a shared example keeps the "common" label.
			sections := append([]Section(nil), e.Sections...)
			sort.SliceStable(sections, func(i, j int) bool {
				return sections[i].Implementation == "common" && sections[j].Implementation != "common"
			})
			for _, s := range sections {
				if !s.IsExamples() {
					continue
				}
				for _, ex := range ParseExamples(s.Markdown) {
					if ex.Code == "" || seen[ex.Code] {
						continue
					}
					seen[ex.Code] = true
					lead := strings.ToLower(ex.Description + " " + ex.Short)
					c.examples = append(c.examples, ExampleHit{
						Command: e.ID, Title: e.Title, Implementation: s.Implementation,
						Description: ex.Description, Short: ex.Short, Code: ex.Code,
						order: len(c.examples),
						lead:  lead,
						hay:   strings.ToLower(e.ID) + " " + lead + " " + strings.ToLower(ex.Code),
					})
				}
			}
		}
	})
	return c.examples
}

// SearchExamples returns the examples whose command name, description, short
// label and code together contain every word of query, case-insensitively.
// Examples whose description or label contains all the words rank before those
// matching only through the code; ties keep catalog and authored order. A limit
// of zero or less means no limit. Nothing is executed.
func (c *Catalog) SearchExamples(query string, limit int) []ExampleHit {
	terms := strings.Fields(strings.ToLower(query))
	hits := []ExampleHit{}
	if len(terms) == 0 {
		return hits
	}
	for _, h := range c.exampleIndex() {
		all, inLead := true, true
		for _, t := range terms {
			if !strings.Contains(h.hay, t) {
				all = false
				break
			}
			if !strings.Contains(h.lead, t) {
				inLead = false
			}
		}
		if !all {
			continue
		}
		if !inLead {
			h.rank = 1
		}
		hits = append(hits, h)
	}
	sort.SliceStable(hits, func(i, j int) bool {
		if hits[i].rank != hits[j].rank {
			return hits[i].rank < hits[j].rank
		}
		return hits[i].order < hits[j].order
	})
	if limit > 0 && len(hits) > limit {
		hits = hits[:limit]
	}
	return hits
}
