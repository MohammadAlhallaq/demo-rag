package chunk

import (
	"strings"
)

// func Chunk(text string, tokenSize, overlap int) []string {
// 	separators := []string{"\n\n", "\n", ". ", " "}

// 	var split func(text string, seps []string) []string
// 	split = func(text string, seps []string) []string {
// 		if len(seps) == 0 {
// 			return []string{text}
// 		}
// 		chunks := []string{text}
// 		sep := seps[0]
// 		var result []string
// 		for _, c := range chunks {
// 			if len(strings.Fields(c)) <= tokenSize {
// 				result = append(result, c)
// 				continue
// 			}
// 			parts := strings.Split(c, sep)
// 			for _, p := range parts {
// 				p = strings.TrimSpace(p)
// 				if p == "" {
// 					continue
// 				}
// 				if len(strings.Fields(p)) <= tokenSize {
// 					result = append(result, p)
// 				} else {
// 					result = append(result, split(p, seps[1:])...)
// 				}
// 			}
// 		}

// 		var merged []string
// 		buf := ""
// 		for _, c := range result {
// 			if buf == "" {
// 				buf = c
// 			} else if len(strings.Fields(buf+" "+c)) <= tokenSize {
// 				buf += " " + c
// 			} else {
// 				merged = append(merged, buf)
// 				buf = c
// 			}
// 		}
// 		if buf != "" {
// 			merged = append(merged, buf)
// 		}

// 		return addOverlap(merged, overlap)
// 	}

// 	return split(text, separators)
// }

func chunk(text string, tokenSize, overlap int) []string {
	paragraphs := strings.Split(text, "\n\n")
	var chunks []string

	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if len(strings.Fields(p)) <= tokenSize {
			chunks = append(chunks, p)
		} else {
			for _, s := range splitByLimit(p, tokenSize) {
				chunks = append(chunks, s)
			}
		}
	}

	return addOverlap(chunks, overlap)
}

func Chunks(docs map[string]string, tokenSize, overlap int) map[string][]string {
	result := make(map[string][]string, len(docs))
	for name, text := range docs {
		result[name] = chunk(text, tokenSize, overlap)
	}
	return result
}

func splitByLimit(text string, tokenSize int) []string {
	var result []string
	for _, s := range strings.Split(text, ". ") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if len(strings.Fields(s)) <= tokenSize {
			result = append(result, s)
		} else {
			words := strings.Fields(s)
			for len(words) > 0 {
				if len(words) <= tokenSize {
					result = append(result, strings.Join(words, " "))
					break
				}
				result = append(result, strings.Join(words[:tokenSize], " "))
				words = words[tokenSize:]
			}
		}
	}
	return result
}

func addOverlap(chunks []string, overlap int) []string {
	if overlap <= 0 || len(chunks) <= 1 {
		return chunks
	}

	result := make([]string, len(chunks))
	result[0] = chunks[0]

	for i := 1; i < len(chunks); i++ {
		prevWords := strings.Fields(chunks[i-1])
		if len(prevWords) > overlap {
			prevWords = prevWords[len(prevWords)-overlap:]
		}
		result[i] = strings.Join(prevWords, " ") + " " + chunks[i]
	}

	return result
}
