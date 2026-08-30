# Nuzlocke Verifier

This project simulates Pokémon battles, learns a per-state action policy for the player side, and can save that policy to disk for later replay.

## What has been added

The project now includes:

- reinforcement learning for the player trainer using discretized battle states
- safety-first reward scoring that penalizes losing and prioritizes survival
- full-party survival rewards so the policy learns to keep the whole party alive
- a negative penalty when the active opponent move has a non-zero crit chance to kill the player Pokémon
- per-state action counts and decay so the policy keeps improving across many battle iterations
- persisted policy storage under a `policies/` directory
- a static policy executor used when loading a saved policy with `--policy-file`
- compatibility checking to reject invalid policies for the wrong party composition
- command-line support for training, saving, and reusing learned policies

## CLI usage

The main entrypoint is in [main.go](main.go). Supported flags are:

| Flag | Meaning |
| --- | --- |
| `-p`, `--player-learning-ai` | Train the player using the learning AI while the opponent keeps the RNB AI logic |
| `-f`, `--policy-file <path>` | Load a saved policy JSON and use it as a static policy for the player |
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

Load a saved policy and use it statically:

```bash
go run . -f policies/player__vs__rnb_trainer_1.json data/player.txt data/rnb_trainer_1.txt -i 1
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
- `policy`: the state-to-action map
- `scores`: learned action scores per state
- `counts`: action usage counts per state
- `created_at`: save timestamp
- `version`: policy version
- `metadata`: source input file names

Example shape:

```json
{
  "player_party": ["horsea", "staravia"],
  "opponent_party": ["dwebble", "sandygast"],
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

Compatibility is checked before a policy is used. If the saved policy does not match the loaded party composition, the program exits with an error instead of running with invalid state.

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

## Notes

- saved policies are stored under the new `policies/` directory
- the directory is ignored by Git via [.gitignore](.gitignore)
- safety-first training and crit-risk penalties are enforced in the learning logic in [struct_ai_learning.go](struct_ai_learning.go)
- the combat reward model is summarized in [battle_statistics.go](battle_statistics.go)

