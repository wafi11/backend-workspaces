package templates

var (
	QueryCreate = `
		INSERT INTO templates (
			name,
			display_name,
			description,
			category,
			git_repo_url,
			git_branch,
			helm_chart_path,
			dockerfile_path,
			default_cpu_request,
			default_cpu_limit,
			default_memory_request,
			default_memory_limit,
			default_replicas,
			requires_database,
			default_database_type,
			requires_redis,
			requires_rabbitmq,
			default_port,
			env_vars_schema,
			tags,
			features,
			icon_url,
			is_active,
			is_featured
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24
		) RETURNING id, created_at, updated_at
	`
	queryListWithCursor = `SELECT
		id, name, display_name, description, category, version,
		git_repo_url, git_branch,default_port,icon_url,
		is_active, is_featured, created_at, updated_at, deleted_at
	FROM templates
	WHERE deleted_at IS NULL
		AND (
			$1::timestamp IS NULL 
			OR created_at < $1::timestamp
			OR (created_at = $1::timestamp AND id < $2)
		)
	ORDER BY created_at DESC, id DESC
	LIMIT $3`

	queryGetByID = `
		SELECT
			id, name, display_name, description, category, version,
			git_repo_url, git_branch, helm_chart_path, dockerfile_path,
			default_cpu_request, default_cpu_limit, default_memory_request, default_memory_limit, default_replicas,
			requires_database, default_database_type, requires_redis, requires_rabbitmq, default_port,
			env_vars_schema, tags, features, icon_url, screenshot_urls,
			is_active, is_featured, created_at, updated_at
		FROM templates
		WHERE id = $1 AND deleted_at IS NULL
	`

	// queryGetByName = `
	// 	SELECT
	// 		id, name, display_name, description, category, version,
	// 		git_repo_url, git_branch, helm_chart_path, dockerfile_path,
	// 		default_cpu_request, default_cpu_limit, default_memory_request, default_memory_limit, default_replicas,
	// 		requires_database, default_database_type, requires_redis, requires_rabbitmq, default_port,
	// 		env_vars_schema, tags, features, icon_url, screenshot_urls,
	// 		is_active, is_featured, created_at, updated_at, deleted_at
	// 	FROM templates
	// 	WHERE name = $1 AND deleted_at IS NULL
	// `

	// queryList = `
	// 	SELECT
	// 		id, name, display_name, description, category, version,
	// 		git_repo_url, git_branch, helm_chart_path, dockerfile_path,
	// 		default_cpu_request, default_cpu_limit, default_memory_request, default_memory_limit, default_replicas,
	// 		requires_database, default_database_type, requires_redis, requires_rabbitmq, default_port,
	// 		env_vars_schema, tags, features, icon_url, screenshot_urls,
	// 		is_active, is_featured, created_at, updated_at, deleted_at
	// 	FROM templates
	// 	WHERE deleted_at IS NULL
	// 	ORDER BY is_featured DESC, created_at DESC
	// 	LIMIT $1 OFFSET $2
	// `

	// queryCount = `
	// 	SELECT COUNT(*) FROM templates WHERE deleted_at IS NULL
	// `

	// queryUpdate = `
	// 	UPDATE templates SET
	// 		display_name = COALESCE($1, display_name),
	// 		description = COALESCE($2, description),
	// 		category = COALESCE($3, category),
	// 		git_repo_url = COALESCE($4, git_repo_url),
	// 		git_branch = COALESCE($5, git_branch),
	// 		helm_chart_path = COALESCE($6, helm_chart_path),
	// 		default_cpu_request = COALESCE($7, default_cpu_request),
	// 		default_cpu_limit = COALESCE($8, default_cpu_limit),
	// 		default_memory_request = COALESCE($9, default_memory_request),
	// 		default_memory_limit = COALESCE($10, default_memory_limit),
	// 		default_replicas = COALESCE($11, default_replicas),
	// 		requires_database = COALESCE($12, requires_database),
	// 		default_database_type = COALESCE($13, default_database_type),
	// 		requires_redis = COALESCE($14, requires_redis),
	// 		requires_rabbitmq = COALESCE($15, requires_rabbitmq),
	// 		env_vars_schema = COALESCE($16, env_vars_schema),
	// 		tags = COALESCE($17, tags),
	// 		features = COALESCE($18, features),
	// 		icon_url = COALESCE($19, icon_url),
	// 		is_active = COALESCE($20, is_active),
	// 		is_featured = COALESCE($21, is_featured),
	// 		updated_at = CURRENT_TIMESTAMP
	// 	WHERE id = $22 AND deleted_at IS NULL
	// 	RETURNING updated_at
	// `

	// queryDelete = `
	// 	UPDATE templates
	// 	SET deleted_at = CURRENT_TIMESTAMP
	// 	WHERE id = $1 AND deleted_at IS NULL
	// `

	// queryHardDelete = `
	// 	DELETE FROM templates WHERE id = $1
	// `

	// queryListByCategory = `
	// 	SELECT
	// 		id, name, display_name, description, category, icon_url,
	// 		tags, features, requires_database, is_featured, created_at
	// 	FROM templates
	// 	WHERE category = $1 AND is_active = true AND deleted_at IS NULL
	// 	ORDER BY is_featured DESC, created_at DESC
	// `

	// queryListFeatured = `
	// 	SELECT
	// 		id, name, display_name, description, category, icon_url,
	// 		tags, features, requires_database, is_featured, created_at
	// 	FROM templates
	// 	WHERE is_featured = true AND is_active = true AND deleted_at IS NULL
	// 	ORDER BY created_at DESC
	// 	LIMIT $1
	// `

	// querySearchByTag = `
	// 	SELECT
	// 		id, name, display_name, description, category, icon_url,
	// 		tags, features, requires_database, is_featured, created_at
	// 	FROM templates
	// 	WHERE $1 = ANY(tags) AND is_active = true AND deleted_at IS NULL
	// 	ORDER BY is_featured DESC, created_at DESC
	// `

	// querySearchByFeature = `
	// 	SELECT
	// 		id, name, display_name, description, category, icon_url,
	// 		tags, features, requires_database, is_featured, created_at
	// 	FROM templates
	// 	WHERE $1 = ANY(features) AND is_active = true AND deleted_at IS NULL
	// 	ORDER BY is_featured DESC, created_at DESC
	// `
)
