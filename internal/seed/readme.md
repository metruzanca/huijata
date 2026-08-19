# Seed Prediction

Pure Go implementation of Noita seed prediction algorithms, providing:

- Alchemic Precursor recipes
- Lively Concoction recipes
- Fungal shifts table
- Mountain perk tables for all holy mountains

## References

This implementation was built by studying and referencing two open-source projects:

- **[noita-tools](https://github.com/Patile/noita-tools)** - The primary reference for understanding Noita's RNG implementation and perk generation algorithms.
- **[noita-telescope](https://github.com/AaronAsAChimp/noita-telescope)** - Additional reference for seed prediction approaches.

## Technical Notes

The RNG uses a pure Go port of Noita's `NollaPrng` algorithm. The perk deck generation was reverse-engineered by analyzing the noita-telescope JavaScript implementation.

Key implementation details:
- Perks data (`perks.json`) must be stored as a JSON array to ensure deterministic iteration order
- The `NollaPrng` implementation uses the same LCG constants as the game
- Perk spawn order is generated using a seeded shuffle followed by distance-based duplicate removal
