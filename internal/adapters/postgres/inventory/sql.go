package inventory

const qrySaveInventory = `
	INSERT INTO inventory (quantity, consumable_id, updated_at)
	VALUES (:quantity, :consumable_id, current_timestamp)`

const qryGetInventory = `
	SELECT i.*
	FROM inventory i
	WHERE consumable_id = $1`

const qryUpdateInventory = `
	UPDATE inventory
	SET quantity = :quantity,
		  updated_at = current_timestamp
	WHERE uid = :uid`

const qryDeleteInventory = `
	DELETE FROM inventory
	WHERE consumable_id = $1`

const qryDecreseQuantityInventory = `
UPDATE inventory
SET quantity = GREATEST(quantity - 1, 0)
WHERE consumable_id = :consumable_id;`

const qryIncreaseQuantityInventory = `
UPDATE inventory
SET quantity = quantity + 1
WHERE consumable_id = :consumable_id;`
