// Package rules contains all linting rules applied by tagguard.
package rules

// editDistance computes the Damerau-Levenshtein distance between two strings.
// Unlike plain Levenshtein, this also counts transpositions (e.g. "jsno"→"json")
// as a single edit, making it much better for catching typos.
func editDistance(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	// d[i][j] = distance between a[:i] and b[:j]
	d := make([][]int, la+1)
	for i := range d {
		d[i] = make([]int, lb+1)
	}
	for i := 0; i <= la; i++ {
		d[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		d[0][j] = j
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			d[i][j] = min3(
				d[i-1][j]+1,      // deletion
				d[i][j-1]+1,      // insertion
				d[i-1][j-1]+cost, // substitution
			)
			// Transposition: a[i-2]==b[j-1] && a[i-1]==b[j-2]
			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				if d[i-2][j-2]+cost < d[i][j] {
					d[i][j] = d[i-2][j-2] + cost
				}
			}
		}
	}

	return d[la][lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// closestMatch finds the closest string in candidates to target.
// Returns the best match and true if the distance is within threshold,
// or empty string and false if no close match found.
func closestMatch(target string, candidates []string, threshold int) (string, bool) {
	best := ""
	bestDist := threshold + 1

	for _, c := range candidates {
		d := editDistance(target, c)
		if d < bestDist {
			bestDist = d
			best = c
		}
	}

	if bestDist <= threshold {
		return best, true
	}
	return "", false
}
