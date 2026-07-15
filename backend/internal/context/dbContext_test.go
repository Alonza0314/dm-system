package context_test

import "testing"

var testDbContextCases = []struct {
	name   string
	key    string
	valuse string
}{
	{
		name:   "Data 1",
		key:    "key1",
		valuse: "value1",
	},
	{
		name:   "Data 2",
		key:    "key2",
		valuse: "value2",
	},
	{
		name:   "Data 3",
		key:    "key3",
		valuse: "value3",
	},
}

var (
	testDbCollections string = "test_db_collections"
)

func TestDbContext(t *testing.T) {
	t.Run("Single Data Operations", testSingleDataOperations)
	t.Run("Multiple Data Operations", testMultipleDataOperations)
}

func testSingleDataOperations(t *testing.T) {
	for _, tc := range testDbContextCases {
		if err := dmCtx.Save(testDbCollections, tc.key, tc.valuse); err != nil {
			t.Errorf("Failed to save data for %s: %v", tc.name, err)
		}

		exist, err := dmCtx.Exist(testDbCollections, tc.key)
		if err != nil {
			t.Errorf("Failed to check existence for %s: %v", tc.name, err)
		} else if !exist {
			t.Errorf("Data for %s does not exist after saving", tc.name)
		}
	}

	for _, tc := range testDbContextCases {
		value, err := dmCtx.Load(testDbCollections, tc.key)
		if err != nil {
			t.Errorf("Failed to load data for %s: %v", tc.name, err)
		} else if value != tc.valuse {
			t.Errorf("Loaded value for %s is incorrect: got %s, want %s", tc.name, value, tc.valuse)
		}
	}

	for _, tc := range testDbContextCases {
		if err := dmCtx.Remove(testDbCollections, tc.key); err != nil {
			t.Errorf("Failed to delete data for %s: %v", tc.name, err)
		}

		exist, err := dmCtx.Exist(testDbCollections, tc.key)
		if err != nil {
			t.Errorf("Failed to check existence after deletion for %s: %v", tc.name, err)
		} else if exist {
			t.Errorf("Data for %s still exists after deletion", tc.name)
		}
	}
}

func testMultipleDataOperations(t *testing.T) {
	dataMap := make(map[string]string)
	for _, tc := range testDbContextCases {
		dataMap[tc.key] = tc.valuse
	}

	if err := dmCtx.SaveAll(testDbCollections, dataMap); err != nil {
		t.Errorf("Failed to save multiple data: %v", err)
	}

	for _, tc := range testDbContextCases {
		exist, err := dmCtx.Exist(testDbCollections, tc.key)
		if err != nil {
			t.Errorf("Failed to check existence for %s: %v", tc.name, err)
		} else if !exist {
			t.Errorf("Data for %s does not exist after saving", tc.name)
		}
	}

	allData, err := dmCtx.LoadAll(testDbCollections)
	if err != nil {
		t.Errorf("Failed to load all data: %v", err)
	} else {
		for _, tc := range testDbContextCases {
			value, ok := allData[tc.key]
			if !ok {
				t.Errorf("Key %s not found in loaded data", tc.key)
			} else if value != tc.valuse {
				t.Errorf("Loaded value for %s is incorrect: got %s, want %s", tc.name, value, tc.valuse)
			}
		}
	}

	if err := dmCtx.RemoveAll(testDbCollections); err != nil {
		t.Errorf("Failed to delete all data: %v", err)
	}

	for _, tc := range testDbContextCases {
		exist, err := dmCtx.Exist(testDbCollections, tc.key)
		if err != nil {
			t.Errorf("Failed to check existence after deletion for %s: %v", tc.name, err)
		} else if exist {
			t.Errorf("Data for %s still exists after deletion", tc.name)
		}
	}
}
