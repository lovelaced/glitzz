package util

var colors = map[string]string{
	"NORMAL":     "\u000f",
	"BOLD":       "\u0002",
	"UNDERLINE":  "\u001f",
	"REVERSE":    "\u0016",
	"WHITE":      "\u00030",
	"BLACK":      "\u00031",
	"DARK_BLUE":  "\u00032",
	"DARK_GREEN": "\u00033",
	"RED":        "\u00034",
	"BROWN":      "\u00035",
	"GREEN":      "\u00039",
}

func Greentext(colorable string) string {
	result := colors["GREEN"] + colorable
	return result
}

func Normaltext(colorable string) string {
	result := colors["NORMAL"] + colorable
	return result
}

func Returntonormal(colorable string) string {
	result := colorable + colors["NORMAL"]
	return result
}

func Boldtext(colorable string) string {
	result := colors["BOLD"] + colorable
	return result
}
