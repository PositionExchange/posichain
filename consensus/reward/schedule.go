package reward

import (
	"sort"
	"time"

	shardingconfig "github.com/PositionExchange/posichain/internal/configs/sharding"
	"github.com/PositionExchange/posichain/internal/utils"
	"github.com/PositionExchange/posichain/numeric"
	"github.com/PositionExchange/posichain/shard"
)

type pair struct {
	ts    int64
	share numeric.Dec
}

var (

	// schedule is the Token Release Schedule of Posichain
	releasePlan = map[int64]numeric.Dec{
		// 2022
		mustParse("2022-Jul-31"): numeric.MustNewDecFromStr("1.000000000000000"),
	}

	sorted = func() []pair {
		s := []pair{}
		for k, v := range releasePlan {
			s = append(s, pair{k, v})
		}
		sort.SliceStable(
			s,
			func(i, j int) bool { return s[i].ts < s[j].ts },
		)
		return s
	}()
)

func mustParse(ts string) int64 {
	const shortForm = "2006-Jan-02"
	t, err := time.Parse(shortForm, ts)
	if err != nil {
		panic("could not parse timestamp")
	}
	return t.Unix()
}

// PercentageForTimeStamp ..
func PercentageForTimeStamp(ts int64) numeric.Dec {
	if shard.Schedule.GetNetworkID() != shardingconfig.MainNet {
		return numeric.MustNewDecFromStr("1")
	}

	bucket := pair{}
	i, j := 0, 1

	for range sorted {
		if i == (len(sorted) - 1) {
			if ts < sorted[0].ts {
				bucket = sorted[0]
			} else {
				bucket = sorted[i]
			}
			break
		}
		if ts >= sorted[i].ts && ts < sorted[j].ts {
			bucket = sorted[i]
			break
		}
		i++
		j++
	}

	utils.Logger().Info().
		Str("percent of total-supply used", bucket.share.Mul(numeric.NewDec(100)).String()).
		Str("for-time", time.Unix(ts, 0).String()).
		Msg("Picked Percentage for timestamp")

	return bucket.share
}
