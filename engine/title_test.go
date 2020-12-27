package engine

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleTitle1 Title = Title{
	tconst: "tc1",
}

var sampleTitle2 Title = Title{
	tconst: "tc2",
}

func TestRandTitleCanOutputAllElements(t *testing.T) {
	assert := assert.New(t)
	rand.Seed(27)
	counts := map[string]int{
		sampleTitle1.tconst: 0,
		sampleTitle2.tconst: 0,
	}
	titles := []Title{sampleTitle1, sampleTitle2}
	for i := 0; i < 100; i++ {
		randomTitle := RandTitle(titles)
		counts[randomTitle.tconst]++
	}

	assert.Greater(counts["tc1"], 40, "Expected to get first title back reasonably often")
	assert.Greater(counts["tc2"], 40, "Expected to get second title back reasonably often")
	assert.Equal(len(counts), 2, "Expected to only ever get results for the input titles")
}
