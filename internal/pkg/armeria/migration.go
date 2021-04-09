package armeria

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// SchemaVersion defines the current version of the schema. If the file system is using an older version, a
// migration will be performed.
const SchemaVersion int = 6

// schemaVersionOnDisk reads the schema version from disk and returns it as an int.
func schemaVersionOnDisk() int {
	b := Armeria.storageManager.ReadFile("schema-version")
	sv, err := strconv.Atoi(string(b))
	if err != nil {
		Armeria.log.Fatal("error parsing schema version from disk", zap.Error(err))
	}

	return sv
}

// writeSchemaVersionToDisk writes the schema version to disk.
func writeSchemaVersionToDisk(version int) {
	v := strconv.Itoa(version)
	b := []byte(v)

	err := ioutil.WriteFile(Armeria.dataPath+"/schema-version", b, 0644)
	if err != nil {
		Armeria.log.Fatal("error writing schema version to disk", zap.Error(err))
	}
}

// VerifySchemaVersion verifies that the game server can run on the current schema version.
func verifySchemaVersion() {
	sv := schemaVersionOnDisk()

	if SchemaVersion < sv {
		Armeria.log.Fatal("cannot downgrade schema version",
			zap.Int("installed", sv),
			zap.Int("desired", SchemaVersion),
		)
	}

	Armeria.log.Info("schema version check",
		zap.Int("installed", sv),
		zap.Int("desired", SchemaVersion),
	)

	if SchemaVersion > sv {
		Armeria.log.Fatal("schema has been upgraded; perform a migration first")
	}
}

// migrateCharacters handles migrations for characters.
func migrateCharacters(to int) {
	s := struct {
		Characters []*Character `json:"characters"`
	}{}

	b, err := ioutil.ReadFile(Armeria.dataPath + "/characters.json")
	if err != nil {
		Armeria.log.Fatal("error reading characters.json", zap.Error(err))
	}

	err = json.Unmarshal(b, &s)
	if err != nil {
		Armeria.log.Fatal("error unmarshalling characters.json", zap.Error(err))
	}

	for _, c := range s.Characters {
		switch to {
		case 2:
			// set UnsafeLastSeen to now
			c.UnsafeLastSeen = time.Now()
		case 3:
			// set UnsafeSettings to an initialized map
			c.UnsafeSettings = map[string]string{}
		}

		Armeria.log.Info("character migration successful",
			zap.String("name", c.UnsafeName),
		)
	}

	b, err = json.Marshal(s)
	if err != nil {
		Armeria.log.Fatal("error marshalling characters.json", zap.Error(err))
	}

	err = ioutil.WriteFile(Armeria.dataPath+"/characters.json", b, 0644)
	if err != nil {
		Armeria.log.Fatal("error writing characters.json", zap.Error(err))
	}
}

// migrateMobs handles migrations for mobs.
func migrateMobs(to int) {
	s := struct {
		Mobs []*Mob `json:"mobs"`
	}{}

	b, err := ioutil.ReadFile(Armeria.dataPath + "/mobs.json")
	if err != nil {
		Armeria.log.Fatal("error reading mobs.json", zap.Error(err))
	}

	err = json.Unmarshal(b, &s)
	if err != nil {
		Armeria.log.Fatal("error unmarshalling mobs.json", zap.Error(err))
	}

	for _, m := range s.Mobs {
		for _, mi := range m.Instances() {
			switch to {
			case 4:
				// set UnsafeInventory to an initialized inventory
				mi.UnsafeInventory = NewObjectContainer(0)
			}
		}

		Armeria.log.Info("mob migration successful",
			zap.String("name", m.UnsafeName),
		)
	}

	b, err = json.Marshal(s)
	if err != nil {
		Armeria.log.Fatal("error marshalling mobs.json", zap.Error(err))
	}

	err = ioutil.WriteFile(Armeria.dataPath+"/mobs.json", b, 0644)
	if err != nil {
		Armeria.log.Fatal("error writing mobs.json", zap.Error(err))
	}
}

// migrateMobs handles migrations for items.
func migrateItems(to int) {
	s := struct {
		Items []*Item `json:"items"`
	}{}

	b, err := ioutil.ReadFile(Armeria.dataPath + "/items.json")
	if err != nil {
		Armeria.log.Fatal("error reading items.json", zap.Error(err))
	}

	err = json.Unmarshal(b, &s)
	if err != nil {
		Armeria.log.Fatal("error unmarshalling items.json", zap.Error(err))
	}

	for _, i := range s.Items {
		switch to {
		case 6:
			if _, exists := i.UnsafeAttributes["rarity"]; exists {
				i.UnsafeAttributes["rarity"] = "common"
			}
		}

		for _, ii := range i.Instances() {
			switch to {
			case 6:
				if _, exists := ii.UnsafeAttributes["rarity"]; exists {
					ii.UnsafeAttributes["rarity"] = "common"
				}
			}
		}

		Armeria.log.Info("item migration successful",
			zap.String("name", i.UnsafeName),
		)
	}

	b, err = json.Marshal(s)
	if err != nil {
		Armeria.log.Fatal("error marshalling items.json", zap.Error(err))
	}

	err = ioutil.WriteFile(Armeria.dataPath+"/items.json", b, 0644)
	if err != nil {
		Armeria.log.Fatal("error writing items.json", zap.Error(err))
	}
}

// migrateLedgers handles migrations for ledgers.
func migrateLedgers(to int) {
	if to == 5 {
		lm := &LedgerManager{
			dataFile:      fmt.Sprintf("%s/ledgers.json", Armeria.dataPath),
			UnsafeLedgers: []*Ledger{},
		}
		lm.SaveLedgers()
		Armeria.log.Info("initial ledger created successfully")
	}
}

// Migrate performs a sequential data migration.
func Migrate() {
	sv := schemaVersionOnDisk()

	Armeria.log.Info("migration starting",
		zap.Int("installed", sv),
		zap.Int("desired", SchemaVersion),
	)

	for i := sv + 1; i <= SchemaVersion; i++ {
		Armeria.log.Info("migrating to next schema version", zap.Int("version", i))
		migrateCharacters(i)
		migrateMobs(i)
		migrateLedgers(i)
		migrateItems(i)
	}

	writeSchemaVersionToDisk(SchemaVersion)

	Armeria.log.Info("migration complete")
}
