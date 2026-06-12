package retrieve

import (
	"fmt"
	"math"
	"sort"

	"rag-demo/internal/types"
)

func Cosine(a, b []float64) float64 {
	var dot, na, nb float64
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func TopK(index types.Index, queryVec []float64, k int) []types.Doc {
	var scored []types.Scored
	for _, d := range index.Docs {
		scored = append(scored, types.Scored{
			Doc:   d,
			Score: Cosine(d.Vec, queryVec),
		})
	}

	fmt.Println(len(scored))

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	var res []types.Doc
	for i := 0; i < k && i < len(scored); i++ {
		res = append(res, scored[i].Doc)
	}
	return res
}
