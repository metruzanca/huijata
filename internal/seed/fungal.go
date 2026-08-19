package seed

const maxShifts = 20
const maxFromMaterials = 6

type FungalShift struct {
	FlaskTo   bool
	FlaskFrom bool
	From      []string
	To        string
	GoldToX   string
	GrassToX  string
}

func PickFungal(worldSeed uint32, requestedMaxShifts int) []FungalShift {
	limit := maxShifts
	if requestedMaxShifts > 0 && requestedMaxShifts < maxShifts {
		limit = requestedMaxShifts
	}

	var shifts []FungalShift
	for iter := 0; iter < limit; iter++ {
		convertedAny := false
		for convertTries := 0; !convertedAny && convertTries < 20; convertTries++ {
			seed2 := 42345 + iter + 1000*convertTries
			pos := RandomPos{X: 9123, Y: int32(seed2)}

			randoms := NollaPrng{}
			randoms.SetWorldSeed(worldSeed)
			randoms.SetRandomSeed(89346, float64(seed2))

			from := pickFrom(worldSeed, &pos)
			to := pickTo(worldSeed, &pos)

			shift := FungalShift{
				From:      []string{},
				To:        MaterialName(to.Material),
				GoldToX:   MaterialName(matGold),
				GrassToX:  MaterialName(matGrassHoly),
			}

			if randoms.Random(1, 100) <= 75 {
				if randoms.Random(1, 100) <= 50 {
					shift.FlaskFrom = true
				} else {
					shift.FlaskTo = true
					if randoms.Random(1, 1000) != 1 {
						index := randoms.Random(0, int32(len(greedyMaterials)-1))
						shift.GoldToX = MaterialName(greedyMaterials[index])
						shift.GrassToX = MaterialName(matGrass)
					}
				}
			}

			for _, material := range from.Materials {
				if len(shift.From) == 0 || material != to.Material {
					shift.From = append(shift.From, MaterialName(material))
					convertedAny = true
				}
			}

			if convertedAny {
				shifts = append(shifts, shift)
			}
		}
	}

	return shifts
}

func pickFrom(worldSeed uint32, pos *RandomPos) FungalFromEntry {
	weightSum := 0.0
	for _, item := range fungalFromGroups {
		weightSum += item.Probability
	}
	val := randomNextF(worldSeed, pos, 0.0, weightSum)

	min := 0.0
	for _, item := range fungalFromGroups {
		max := min + item.Probability
		if val >= min && val <= max {
			return item
		}
		min = max
	}
	return fungalFromGroups[0]
}

func pickTo(worldSeed uint32, pos *RandomPos) FungalToEntry {
	weightSum := 0.0
	for _, item := range fungalToMaterials {
		weightSum += item.Probability
	}
	val := randomNextF(worldSeed, pos, 0.0, weightSum)

	min := 0.0
	for _, item := range fungalToMaterials {
		max := min + item.Probability
		if val >= min && val <= max {
			return item
		}
		min = max
	}
	return fungalToMaterials[0]
}

func FungalShifts(worldSeed uint32, limit int) []FungalShift {
	return PickFungal(worldSeed, limit)
}
