package servicedefinition

const qrySaveServiceDefinition = `
	INSERT INTO service_definition (name, price, description, created_at)
	VALUES (:name, :price, :description, current_timestamp) RETURNING id;`

const qryGetServiceDefinition = `
	SELECT sd.*
	FROM service_definition sd
	WHERE uid = $1`

const qryListServiceDefinition = `
	SELECT sd.* FROM service_definition sd`

const qryTotalListServiceDefinition = `
	SELECT COUNT(*) FROM service_definition sd`

const qryUpdateServiceDefinition = `
	UPDATE service_definition
	SET name = :name,
			price = :price,
		  description = :description,
		  updated_at = current_timestamp
	WHERE uid = :uid`

const qryDeleteServiceDefinition = `
	DELETE FROM service_definition
	WHERE uid = $1`

const qryExistServiceDefinition = `
	SELECT count(*)
	FROM service_definition sd
	WHERE name = $1`

const qryGetServiceDefinitionID = `
	SELECT sd.*
	FROM service_definition sd
	WHERE id = $1`
