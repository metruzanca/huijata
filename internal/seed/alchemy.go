package seed

func PickAlchemy(worldSeed uint32) [6]int {
	result := [6]int{0, 0, 0, 0, 0, 0}
	rng := NewNollaPrng(float64(worldSeed)*0.17127000 + 1323.59030000)
	for i := 0; i < 5; i++ {
		rng.Next()
	}

	pickForOutput(&rng, worldSeed, result[0:3])
	pickForOutput(&rng, worldSeed, result[3:6])
	return result
}

func pickForOutput(rng *NollaPrng, worldSeed uint32, output []int) {
	materials := [4]int{0, 0, 0, 0}
	materialCount := 0
	pickMaterials(rng, &materials, &materialCount, liquids, 3)
	pickMaterials(rng, &materials, &materialCount, alchemyMaterials, 1)

	shuffleRng := NewNollaPrng(float64(worldSeed>>1) + 12534.0)
	i := materialCount - 1
	for i >= 0 {
		limit := float64(i + 1)
		index := int(shuffleRng.Next() * limit)
		tmp := materials[i]
		materials[i] = materials[index]
		materials[index] = tmp
		i--
	}

	rng.Next()
	rng.Next()

	output[0] = materials[0]
	output[1] = materials[1]
	output[2] = materials[2]
}

func contains(items []int, len int, item int) bool {
	for i := 0; i < len; i++ {
		if items[i] == item {
			return true
		}
	}
	return false
}

func pickMaterials(rng *NollaPrng, materials *[4]int, materialCount *int, source []int, count int) {
	counter := 0
	failed := 0
	for counter < count && failed < 99999 {
		index := int(rng.Next() * float64(len(source)))
		picked := source[index]
		if !contains(materials[:], *materialCount, picked) {
			materials[*materialCount] = picked
			*materialCount++
			counter++
		} else {
			failed++
		}
	}
}

func AlchemyRecipe(worldSeed uint32) (lc []string, ap []string) {
	result := PickAlchemy(worldSeed)
	lc = []string{
		MaterialName(result[0]),
		MaterialName(result[1]),
		MaterialName(result[2]),
	}
	ap = []string{
		MaterialName(result[3]),
		MaterialName(result[4]),
		MaterialName(result[5]),
	}
	return
}
