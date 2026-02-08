package consumable

const qrySaveConsumable = `
	INSERT INTO consumable (name, price, code, description, created_at)
	VALUES (:name, :price, :code, :description, current_timestamp) RETURNING id;`

const qryGetConsumable = `
	SELECT c.*
	FROM consumable c
	WHERE uid = $1`

const qryListConsumable = `
	SELECT c.* FROM consumable c`

const qryTotalListConsumable = `
	SELECT COUNT(*) FROM consumable c`

const qryUpdateConsumable = `
	UPDATE consumable
	SET name = :name,
			price = :price,
		  description = :description,
		  updated_at = current_timestamp
	WHERE uid = :uid`

const qryDeleteConsumable = `
	DELETE FROM consumable
	WHERE uid = $1`

const qryExistConsumable = `
	SELECT COUNT(*)
	FROM consumable c
	WHERE name = $1`

const qryGetConsumableByID = `
	SELECT c.*
	FROM consumable c
	WHERE id = $1`
