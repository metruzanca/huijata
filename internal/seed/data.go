package seed

const (
	matAcid                           = 0
	matAlcohol                        = 1
	matBlood                          = 2
	matBloodFungi                     = 3
	matBloodWorm                      = 4
	matCement                         = 5
	matLava                           = 6
	matMagicLiquidBerserk             = 7
	matMagicLiquidCharm               = 8
	matMagicLiquidFasterLevitation    = 9
	matMagicLiquidFasterLevitationAndMovement = 10
	matMagicLiquidInvisibility        = 11
	matMagicLiquidManaRegeneration    = 12
	matMagicLiquidMovementFaster      = 13
	matMagicLiquidProtectionAll       = 14
	matMagicLiquidTeleportation       = 15
	matMagicLiquidUnstablePolymorph   = 16
	matMagicLiquidUnstableTeleportation = 17
	matMagicLiquidWormAttractor       = 18
	matMaterialConfusion              = 19
	matMud                            = 20
	matOil                            = 21
	matPoison                         = 22
	matRadioactiveLiquid              = 23
	matSwamp                          = 24
	matUrine                          = 25
	matWater                          = 26
	matWaterIce                       = 27
	matWaterSwamp                     = 28
	matMagicLiquidRandomPolymorph     = 29
	matBone                           = 30
	matBrass                          = 31
	matCoal                           = 32
	matCopper                         = 33
	matDiamond                        = 34
	matFungi                          = 35
	matGold                           = 36
	matGrass                          = 37
	matGunpowder                      = 38
	matGunpowderExplosive             = 39
	matRottenMeat                     = 40
	matSand                           = 41
	matSilver                         = 42
	matSlime                          = 43
	matSnow                           = 44
	matSoil                           = 45
	matWax                            = 46
	matHoney                          = 47
	matWaterStatic                    = 48
	matWaterSalt                      = 49
	matMaterialDarkness               = 50
	matPeat                           = 51
	matFungisoil                      = 52
	matBloodCold                      = 53
	matAcidGas                        = 54
	matAcidGasStatic                  = 55
	matPoisonGas                      = 56
	matFungalGas                      = 57
	matRadioactiveGas                 = 58
	matRadioactiveGasStatic           = 59
	matMagicLiquidPolymorph           = 60
	matSteam                          = 61
	matSmoke                          = 62
	matSnowSticky                     = 63
	matRockStatic                     = 64
	matGoldBox2d                      = 65
	matSima                           = 66
	matVomit                          = 67
	matPeaSoup                        = 68
	matRockStaticRadioactive          = 69
	matMimicLiquid                    = 70
	matPoo                            = 71
	matVoidLiquid                     = 72
	matCheeseStatic                   = 73
	matGrassHoly                      = 74
	matMammi                          = 75
	matRottenMeatRadioactive          = 76
)

var materialNames = map[int]string{
	matAcid:                           "acid",
	matAlcohol:                        "alcohol",
	matBlood:                          "blood",
	matBloodFungi:                     "blood_fungi",
	matBloodWorm:                      "blood_worm",
	matCement:                         "cement",
	matLava:                           "lava",
	matMagicLiquidBerserk:             "magic_liquid_berserk",
	matMagicLiquidCharm:               "magic_liquid_charm",
	matMagicLiquidFasterLevitation:    "magic_liquid_faster_levitation",
	matMagicLiquidFasterLevitationAndMovement: "magic_liquid_faster_levitation_and_movement",
	matMagicLiquidInvisibility:        "magic_liquid_invisibility",
	matMagicLiquidManaRegeneration:    "magic_liquid_mana_regeneration",
	matMagicLiquidMovementFaster:      "magic_liquid_movement_faster",
	matMagicLiquidProtectionAll:       "magic_liquid_protection_all",
	matMagicLiquidTeleportation:       "magic_liquid_teleportation",
	matMagicLiquidUnstablePolymorph:   "magic_liquid_unstable_polymorph",
	matMagicLiquidUnstableTeleportation: "magic_liquid_unstable_teleportation",
	matMagicLiquidWormAttractor:       "magic_liquid_worm_attractor",
	matMaterialConfusion:              "material_confusion",
	matMud:                            "mud",
	matOil:                            "oil",
	matPoison:                         "poison",
	matRadioactiveLiquid:              "radioactive_liquid",
	matSwamp:                          "swamp",
	matUrine:                          "urine",
	matWater:                          "water",
	matWaterIce:                       "water_ice",
	matWaterSwamp:                     "water_swamp",
	matMagicLiquidRandomPolymorph:     "magic_liquid_random_polymorph",
	matBone:                           "bone",
	matBrass:                          "brass",
	matCoal:                           "coal",
	matCopper:                         "copper",
	matDiamond:                        "diamond",
	matFungi:                          "fungi",
	matGold:                           "gold",
	matGrass:                          "grass",
	matGunpowder:                      "gunpowder",
	matGunpowderExplosive:             "gunpowder_explosive",
	matRottenMeat:                     "rotten_meat",
	matSand:                           "sand",
	matSilver:                         "silver",
	matSlime:                          "slime",
	matSnow:                           "snow",
	matSoil:                           "soil",
	matWax:                            "wax",
	matHoney:                          "honey",
	matWaterStatic:                    "water_static",
	matWaterSalt:                      "water_salt",
	matMaterialDarkness:               "material_darkness",
	matPeat:                           "peat",
	matFungisoil:                      "fungisoil",
	matBloodCold:                      "blood_cold",
	matAcidGas:                        "acid_gas",
	matAcidGasStatic:                  "acid_gas_static",
	matPoisonGas:                      "poison_gas",
	matFungalGas:                      "fungal_gas",
	matRadioactiveGas:                 "radioactive_gas",
	matRadioactiveGasStatic:           "radioactive_gas_static",
	matMagicLiquidPolymorph:           "magic_liquid_polymorph",
	matSteam:                          "steam",
	matSmoke:                          "smoke",
	matSnowSticky:                     "snow_sticky",
	matRockStatic:                     "rock_static",
	matGoldBox2d:                      "gold_box2d",
	matSima:                           "sima",
	matVomit:                          "vomit",
	matPeaSoup:                        "pea_soup",
	matRockStaticRadioactive:          "rock_static_radioactive",
	matMimicLiquid:                    "mimic_liquid",
	matPoo:                            "poo",
	matVoidLiquid:                     "void_liquid",
	matCheeseStatic:                   "cheese_static",
	matGrassHoly:                      "grass_holy",
	matMammi:                          "mammi",
	matRottenMeatRadioactive:          "rotten_meat_radioactive",
}

func MaterialName(id int) string {
	if name, ok := materialNames[id]; ok {
		return name
	}
	return ""
}

var liquids = []int{
	matAcid,
	matAlcohol,
	matBlood,
	matBloodFungi,
	matBloodWorm,
	matCement,
	matLava,
	matMagicLiquidBerserk,
	matMagicLiquidCharm,
	matMagicLiquidFasterLevitation,
	matMagicLiquidFasterLevitationAndMovement,
	matMagicLiquidInvisibility,
	matMagicLiquidManaRegeneration,
	matMagicLiquidMovementFaster,
	matMagicLiquidProtectionAll,
	matMagicLiquidTeleportation,
	matMagicLiquidUnstablePolymorph,
	matMagicLiquidUnstableTeleportation,
	matMagicLiquidWormAttractor,
	matMaterialConfusion,
	matMud,
	matOil,
	matPoison,
	matRadioactiveLiquid,
	matSwamp,
	matUrine,
	matWater,
	matWaterIce,
	matWaterSwamp,
	matMagicLiquidRandomPolymorph,
}

var alchemyMaterials = []int{
	matBone,
	matBrass,
	matCoal,
	matCopper,
	matDiamond,
	matFungi,
	matGold,
	matGrass,
	matGunpowder,
	matGunpowderExplosive,
	matRottenMeat,
	matSand,
	matSilver,
	matSlime,
	matSnow,
	matSoil,
	matWax,
	matHoney,
}

type FungalFromEntry struct {
	Probability float64
	Materials   []int
}

type FungalToEntry struct {
	Probability float64
	Material    int
}

var fungalFromGroups = []FungalFromEntry{
	{Probability: 1.0, Materials: []int{matWater, matWaterStatic, matWaterSalt, matWaterIce}},
	{Probability: 1.0, Materials: []int{matLava}},
	{Probability: 1.0, Materials: []int{matRadioactiveLiquid, matPoison, matMaterialDarkness}},
	{Probability: 1.0, Materials: []int{matOil, matSwamp, matPeat}},
	{Probability: 1.0, Materials: []int{matBlood}},
	{Probability: 1.0, Materials: []int{matBloodFungi, matFungi, matFungisoil}},
	{Probability: 1.0, Materials: []int{matBloodCold, matBloodWorm}},
	{Probability: 1.0, Materials: []int{matAcid}},
	{Probability: 0.4, Materials: []int{matAcidGas, matAcidGasStatic, matPoisonGas, matFungalGas, matRadioactiveGas, matRadioactiveGasStatic}},
	{Probability: 0.4, Materials: []int{matMagicLiquidPolymorph, matMagicLiquidUnstablePolymorph}},
	{Probability: 0.4, Materials: []int{matMagicLiquidBerserk, matMagicLiquidCharm, matMagicLiquidInvisibility}},
	{Probability: 0.6, Materials: []int{matDiamond}},
	{Probability: 0.6, Materials: []int{matSilver, matBrass, matCopper}},
	{Probability: 0.2, Materials: []int{matSteam, matSmoke}},
	{Probability: 0.4, Materials: []int{matSand}},
	{Probability: 0.4, Materials: []int{matSnowSticky}},
	{Probability: 0.05, Materials: []int{matRockStatic}},
	{Probability: 0.0003, Materials: []int{matGold, matGoldBox2d}},
}

var fungalToMaterials = []FungalToEntry{
	{Probability: 1.0, Material: matWater},
	{Probability: 1.0, Material: matLava},
	{Probability: 1.0, Material: matRadioactiveLiquid},
	{Probability: 1.0, Material: matOil},
	{Probability: 1.0, Material: matBlood},
	{Probability: 1.0, Material: matBloodFungi},
	{Probability: 1.0, Material: matAcid},
	{Probability: 1.0, Material: matWaterSwamp},
	{Probability: 1.0, Material: matAlcohol},
	{Probability: 1.0, Material: matSima},
	{Probability: 1.0, Material: matBloodWorm},
	{Probability: 1.0, Material: matPoison},
	{Probability: 1.0, Material: matVomit},
	{Probability: 1.0, Material: matPeaSoup},
	{Probability: 1.0, Material: matFungi},
	{Probability: 0.8, Material: matSand},
	{Probability: 0.8, Material: matDiamond},
	{Probability: 0.8, Material: matSilver},
	{Probability: 0.8, Material: matSteam},
	{Probability: 0.5, Material: matRockStatic},
	{Probability: 0.5, Material: matGunpowder},
	{Probability: 0.5, Material: matMaterialDarkness},
	{Probability: 0.5, Material: matMaterialConfusion},
	{Probability: 0.2, Material: matRockStaticRadioactive},
	{Probability: 0.02, Material: matMagicLiquidPolymorph},
	{Probability: 0.02, Material: matMagicLiquidRandomPolymorph},
	{Probability: 0.15, Material: matMagicLiquidTeleportation},
	{Probability: 0.1, Material: matMimicLiquid},
	{Probability: 0.01, Material: matUrine},
	{Probability: 0.01, Material: matPoo},
	{Probability: 0.01, Material: matVoidLiquid},
	{Probability: 0.01, Material: matCheeseStatic},
}

var greedyMaterials = []int{
	matBrass,
	matSilver,
	matRadioactiveLiquid,
	matPeaSoup,
	matAcidGas,
	matPoo,
	matMammi,
	matRottenMeatRadioactive,
	matVomit,
}
