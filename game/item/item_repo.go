package item

import "database/sql"

type ItemRepo struct {
	database *sql.DB
}

func NewItemRepo(db *sql.DB) ItemRepo {
	return ItemRepo{db}
}

func (ir *ItemRepo) LoadItemDefinitions() ([]Definition, error) {
	rows, err := ir.database.Query("SELECT * FROM items_definitions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var definitions []Definition
	for rows.Next() {
		definition, err := fillItemDefinitionData(rows)
		if err != nil {
			return nil, err
		}

		definitions = append(definitions, definition)
	}

	return definitions, nil
}

func (ir *ItemRepo) LoadPublicItemDataByModel(modelName string) ([]publicItem, error) {
	stmt, err := ir.database.Prepare("SELECT * FROM items_public WHERE room_model = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(modelName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var publicItems []publicItem
	for rows.Next() {
		item, err := fillPublicRoomItemData(rows)
		if err != nil {
			return nil, err
		}

		publicItems = append(publicItems, item)
	}

	return publicItems, nil
}

func (ir *ItemRepo) LoadPublicRoomItemData() ([]publicItem, error) {
	rows, err := ir.database.Query("SELECT * FROM items_public")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var publicItems []publicItem
	for rows.Next() {
		item, err := fillPublicRoomItemData(rows)
		if err != nil {
			return nil, err
		}

		publicItems = append(publicItems, item)
	}

	return publicItems, nil
}

func fillItemDefinitionData(rows *sql.Rows) (Definition, error) {
	var (
		ID, spriteID, length, width, tradable, recyclable                            int
		sprite, name, desc, behaviorData, interaction, color, drinkIDData, maxStatus string
		topHeight                                                                    float64
	)
	if err := rows.Scan(
		&ID,
		&sprite,
		&spriteID,
		&name,
		&desc,
		&color,
		&length,
		&width,
		&topHeight,
		&maxStatus,
		&behaviorData,
		&interaction,
		&tradable,
		&recyclable,
		&drinkIDData,
	); err != nil {
		return Definition{}, err
	}

	return NewItemDefinition(
		ID,
		length,
		width,
		sprite,
		name,
		desc,
		behaviorData,
		interaction,
		color,
		drinkIDData,
		topHeight,
		Itob(tradable),
		Itob(recyclable),
	), nil
}

func fillPublicRoomItemData(rows *sql.Rows) (publicItem, error) {
	var behaviorData string

	item := publicItem{}
	err := rows.Scan(
		&item.id,
		&item.roomModel,
		&item.sprite,
		&item.x,
		&item.y,
		&item.z,
		&item.rotation,
		&item.topHeight,
		&item.length,
		&item.width,
		&behaviorData,
		&item.currentProgram,
		&item.teleportTo,
		&item.swimTo,
	)
	if err != nil {
		return publicItem{}, err
	}
	item.behavior = behaviorFromString(behaviorData)
	return item, nil
}

// Itob converts an int to a bool.
func Itob(i int) bool {
	return i != 0
}
