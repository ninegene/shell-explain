package catalog

import "strings"

// DescriptionWidth is the shared maximum prose width in compact blocks.
// Authored commands and indivisible inline literals are not split.
const DescriptionWidth = 100

func OptionText(section Section, width int) string {
	if width <= 0 || width > DescriptionWidth {
		width = DescriptionWidth
	}
	return strings.Join(section.Options, ", ") + "\n" + Plain(section.Markdown, width)
}

func OptionsText(sections []Section, width int) string {
	var blocks []string
	for _, section := range sections {
		if len(section.Options) > 0 {
			blocks = append(blocks, strings.TrimSuffix(OptionText(section, width), "\n"))
		}
	}
	if len(blocks) == 0 {
		return ""
	}
	return strings.Join(blocks, "\n\n") + "\n"
}
