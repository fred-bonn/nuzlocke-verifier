# Nuzlocke Verifier

This project simulates Pokémon battles, learns a per-state action policy for the player side, and can save that policy to disk for later replay.

## What has been added

### Battle simulation

- a single-battle engine ([battle_state_single.go](battle_state_single.go)) with a shared `battleState` interface ([battle_state_interface.go](battle_state_interface.go)) and a deep-cloning/reset path ([battle_state_cloning.go](battle_state_cloning.go)) for repeated iterations
- a priority-based action queue ([action_queue.go](action_queue.go)) with benchmarked comparisons against an older heap-based implementation ([action_queue_old.go](action_queue_old.go))
- move resolution, including damage calculation, accuracy/evasion, critical hits, and STAB/type effectiveness ([action_move.go](action_move.go), [action_move_helper.go](action_move_helper.go), [action_move_score_calculation.go](action_move_score_calculation.go))
- switch and move-replacement actions ([action_switch.go](action_switch.go), [action_replace.go](action_replace.go))
- full type chart effectiveness lookups ([lookup_types.go](lookup_types.go))
- weather (rain, sun, sandstorm, hail) with type immunities and per-turn effects ([lookup_weather.go](lookup_weather.go))
- field effects such as entry hazards ([lookup_field_effect.go](lookup_field_effect.go))
- non-volatile and volatile status ailments (burn, paralysis, poison/toxic, freeze, sleep, confusion, infatuation, trap, bound, leech seed, yawn) with ability/type-based immunities ([struct_ailment.go](struct_ailment.go))
- roughly 90 abilities affecting weather, stats, damage, and status interactions ([lookup_abilities.go](lookup_abilities.go))
- held items, including status-curing and stat-boosting berries, gems, choice items, and Leftovers ([struct_item.go](struct_item.go))
- Pokémon stat/stage calculations, HP tracking, and stat-changing effects ([struct_pokemon.go](struct_pokemon.go), [struct_pokemon_base.go](struct_pokemon_base.go), [struct_slot.go](struct_slot.go))
- base stat/EV/IV/nature balancing helpers used when constructing parties ([lookup_balancing.go](lookup_balancing.go), [lookup_stats.go](lookup_stats.go))
- per-battle and cross-run outcome statistics, including a safety-first scoring model ([battle_statistics.go](battle_statistics.go))

### AI opponents

- `rnbAi` ("risk and back off"), a heuristic AI that switches out of unsafe matchups ([struct_ai_rnb.go](struct_ai_rnb.go))
- `randomAi`, a baseline random action chooser ([struct_ai_random.go](struct_ai_random.go))
- `guidedAi`, an interactive AI that prompts a human for each action over stdin/stdout ([struct_ai_guided.go](struct_ai_guided.go))
- `learningAi`, reinforcement learning for the player trainer using discretized battle states ([struct_ai_learning.go](struct_ai_learning.go)), including:
  - safety-first reward scoring that penalizes losing and prioritizes survival
  - full-party survival rewards so the policy learns to keep the whole party alive
  - a negative penalty when the active opponent move has a non-zero crit chance to kill the player Pokémon
  - per-state action counts and decay so the policy keeps improving across many battle iterations
  - saturation detection to stop training once the policy stabilizes
- a static policy executor used when loading a saved policy with `--policy-file`, plus compatibility checking to reject policies that don't match the loaded party composition

### Data and persistence

- a Showdown-style party file parser ([internal/parser](internal/parser/parser.go)) with its own lexer and validation
- a PokeAPI client for fetching Pokémon and move data, with local JSON caching under [data/pokemon](data/pokemon) and [data/moves](data/moves) ([internal/pokeapi](internal/pokeapi/client.go))
- persisted policy storage under a `policies/` directory
- command-line support for training, saving, and reusing learned policies

## What might be missing

- ability and item coverage is broad but not exhaustive; several common competitive abilities/items (e.g. Multiscale, Protean, priority-negating abilities) aren't implemented yet
- no support for double battles, only single 1v1 slots per side
- no support for held-item recovery (e.g. Recycle) or move-based item removal (Knock Off's item-destroying side effect is untested beyond ailment application)
- the guided AI only supports move selection prompts, not interactive switching
- no persistence/versioning migration path if the saved policy JSON schema changes
- test coverage for the battle statistics reward shaping is present but largely integration-style; see [battle_statistics_test.go](battle_statistics_test.go)

## CLI usage

The main entrypoint is in [main.go](main.go). Supported flags are:

| Flag | Meaning |
| --- | --- |
| `-p`, `--player-learning-ai` | Train the player using the learning AI while the opponent keeps the RNB AI logic |
| `-g`, `--player-guided-ai` | Prompt for the player's action each turn (implies `--verbose`) |
| `-f`, `--policy-file <path>` | Load a saved policy JSON and use it as a static policy for the player; uses the player/opponent parties embedded in the policy, so `<player_showdown> <opponent_showdown>` must be omitted |
| `-s`, `--save-policy` | Save the learned policy to `policies/` using the two input file names |
| `-i`, `--iterations <n>` | Number of battle repetitions to run for training or statistics |
| `-w`, `--weather <0..4>` | Weather override: `0=none`, `1=rain`, `2=sun`, `3=sandstorm`, `4=hail` |
| `-v`, `--verbose` | Verbose logging |

Examples:

```bash
go run . -p -i 250 data/player.txt data/rnb_trainer_1.txt
```

Train and save a policy:

```bash
go run . -p -s -i 250 data/player.txt data/rnb_trainer_1.txt
```

Load a saved policy and use it statically. The player and opponent parties come from the policy file itself, so no input files are given:

```bash
go run . -f policies/player__vs__rnb_trainer_1.json -i 1
```

The long-form flags remain supported too, so these are equivalent:

```bash
go run . --player-learning-ai --iterations 250 data/player.txt data/rnb_trainer_1.txt
```

The policy file is named using the input filenames:

```text
policies/<player_name>__vs__<opponent_name>.json
```

For example, the default battle pair produces:

```text
policies/player__vs__rnb_trainer_1.json
```

## Saved policy format

Policies are stored as JSON and live under the `policies/` directory. The files include:

- `player_party`: player Pokémon names
- `opponent_party`: opponent Pokémon names
- `player_party_file`: the raw contents of the player's Showdown party file (e.g. `cat data/player.txt`)
- `opponent_party_file`: the raw contents of the opponent's Showdown party file
- `policy`: the state-to-action map
- `scores`: learned action scores per state
- `counts`: action usage counts per state
- `created_at`: save timestamp
- `version`: policy version
- `metadata`: source input file names

Embedding the party file contents means a saved policy is self-contained: loading it with `--policy-file` never needs the original `data/player.txt` / `data/rnb_trainer_1.txt` files on disk.

Example shape:

```json
{
  "player_party": ["horsea", "staravia"],
  "opponent_party": ["dwebble", "sandygast"],
  "player_party_file": "Horsea @ Oran Berry\nLevel: 17\nModest Nature\nAbility: Swift Swim\n- Bubble Beam\n- Twister\n",
  "opponent_party_file": "Dwebble @ Salac Berry\nLevel: 17\nAdamant Nature\nAbility: Sturdy\n- Knock Off\n- Sticky Web\n",
  "policy": {
    "{\"player_pokemon\":\"horsea\",...}": [
      "move:bubble beam",
      "switch:staravia"
    ]
  },
  "scores": {
    "{\"player_pokemon\":\"horsea\",...}": {
      "move:bubble beam": 142.5,
      "switch:staravia": -45.2
    }
  },
  "counts": {
    "{\"player_pokemon\":\"horsea\",...}": {
      "move:bubble beam": 17,
      "switch:staravia": 4
    }
  }
}
```

Compatibility is checked before a policy is used. If the saved policy does not match the loaded party composition, or does not have `player_party_file`/`opponent_party_file` embedded, the program exits with an error instead of running with invalid state.

## Input party format

The parser accepts Showdown-style party files, like the examples in [data/player.txt](data/player.txt) and [data/rnb_trainer_1.txt](data/rnb_trainer_1.txt).

### Basic format

```text
Pokemon Name @ Item
Level: N
Nature Nature
Ability: Ability
Status: Status   (optional)
HP: N            (optional)
IVs: 31 Atk / 31 Def / 31 SpA / 31 SpD / 31 Spe
- Move 1
- Move 2
- Move 3
- Move 4
```

### Example

```text
Horsea @ Oran Berry
Level: 17
Modest Nature
Ability: Swift Swim
- Bubble Beam
- Twister
```

### Notes

- each Pokémon starts with a name line
- `@ Item` is optional but common
- `Level`, `Nature`, and `Ability` are required
- status, HP, and IVs are optional sections
- moves are listed as `- Move Name`
- the parser normalizes names, strips punctuation, and reads the move list into the battle state

The actual parser logic lives in [internal/parser/parser.go](internal/parser/parser.go) and [internal/parser/lexer.go](internal/parser/lexer.go). The project validates the file structure before constructing the Pokémon party.

## Behavior of loaded policies

When a policy is loaded via `--policy-file`, it is used as a static policy, not as a live learning AI. That means:

- it does not accumulate new battle outcomes during evaluation
- it does not decay or reset the learned scores while the battle runs
- it uses the saved score table to select the best action for the current state
- it remains deterministic for a given input file and policy file

This is specifically important when comparing a single-battle loaded run to a long training run. A single loaded battle with a saved policy is expected to behave consistently, while the training AI continues to improve over many iterations.

## Building and testing

The [Makefile](Makefile) wraps the common developer commands:

```bash
make test        # go test ./...
make build        # runs tests, then builds ./bin/myprog
make run          # build, then run the default player vs. rnb_trainer_1 battle verbosely
make brief        # same, without verbose output
make rain         # build, then run under forced rain weather (also: sun, sandstorm, hail)
make 250          # build, then run 250 iterations (any positive integer works as a goal)
```

## Notes

- saved policies are stored under the `policies/` directory, which is ignored by Git via [.gitignore](.gitignore)
- safety-first training and crit-risk penalties are enforced in the learning logic in [struct_ai_learning.go](struct_ai_learning.go)
- the combat reward model is summarized in [battle_statistics.go](battle_statistics.go)

