package duckdb

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http" // Add this line to import the os package
	"strconv"
	"strings"
	"templrjs/pkg/config"
	"time"

	log "github.com/sirupsen/logrus"

	_ "github.com/cznic/ql/driver"

	"github.com/gin-gonic/gin"
	_ "github.com/marcboeker/go-duckdb"
)

var baseEntityCols = []string{"id", "created_at", "modified_at", "update_count"}

var (
	db *sql.DB
)

// Setup opens a database and saves the reference to `Database` struct.
func Setup() {

	log.Info("Initializing DuckDB")
	if err := initializeDB(); err != nil {
		log.Fatal(err)
	} else {
		log.Info("Initialized DuckDB")
	}

	log.Info("Bootstrapping DuckDB")
	log.Info("Creating base tables")
	ddlQueries := config.Config.Duckdb.Ddl_queries
	if err := CreateDefaultTables(db, ddlQueries); err != nil {
		log.Fatal(err)
	} else {
		log.Info("Created base tables")
	}
	log.Info("Importing external data and creating models")
	if err := bootup(); err != nil {
		log.Fatal(err)
	} else {
		log.Info("Imported external data and creating models")
	}
	log.Info("Bootstrapped DuckDB")

	//defer db.Close()
}

type Entity struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	UpdateCount int       `json:"update_count"`
}

func initializeDB() error {
	/*connStr := "stage.db"

	_, err := os.Stat("data.db")
	if os.IsNotExist(err) {
		// Create the database file if it doesn't exist
		_, err := os.Create("data.db")
		if err != nil {
			return err
		}
	}*/

	connStr := config.Config.Duckdb.Conn
	log.Info("Connecting to DuckDB", connStr)
	log.Debug(connStr)

	if connStr == "" {
		log.Info("Initializing DuckDB in memory")
		d, err := sql.Open("duckdb", "") //in memory
		if err != nil {
			return err
		}
		//return nil
		if err != nil {
			return err
		}
		db = d
		log.Info("Initialized DuckDB in memory")
	} else {
		log.Info("Initializing DuckDB in file")
		// This will create the database if it doesn't exist
		d, err := sql.Open("duckdb", connStr)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		//return nil
		if err != nil {
			return err
		}
		db = d
		log.Info("Initialized DuckDB in file")
	}
	return nil
}

// CALL postgres_attach('dbname=iotadb user=sqlreader password=2Petabytes host=svdorliotaheonode1.orl.eng.hitachivantara.com', source_schema='chinook', sink_schema='main')
var duckDBExtensions = []string{
	"s3",
	"json",
	"spatial",
	"httpfs",
	//"sqlite",
	//"postgres",
	// Add more extensions as needed
}

func bootup() error {

	// Install DuckDB extensions (if not already installed)
	for _, extension := range duckDBExtensions {
		installStatement := fmt.Sprintf("INSTALL %s;", extension)
		_, err := db.Exec(installStatement)
		if err != nil {
			return err
		}
	}

	// Load installed extensions
	for _, extension := range duckDBExtensions {
		loadStatement := fmt.Sprintf("LOAD %s;", extension)
		_, err := db.Exec(loadStatement)
		if err != nil {
			return err
		}
	}

	if err := seedDuckDB(db, config.Config.Duckdb.Seed_data_base_url); err != nil {
		log.Fatalf("Error seeding DuckDB: %v\n", err)
	}

	return nil
}

func CreateDefaultTables(db *sql.DB, ddlQueries []string) error {
	// Check if the DB is nil
	if db == nil {
		return errors.New("nil DB provided to CreateDefaultTables")
	}

	if len(ddlQueries) == 0 {
		log.Fatal("DDL queries not set in config")
	}

	errors := []string{}
	for _, ddlQuery := range ddlQueries {
		log.Info("Executing query: %s", ddlQuery)
		_, err := db.Exec(ddlQuery)
		if err != nil {
			// If the error is just because the table exists, we can continue to the next query.
			if strings.Contains(err.Error(), "already exists") { // Modify this check according to the exact error message DuckDB gives for existing tables
				log.Info("Table already exists, skipping: %s", ddlQuery)
				continue
			}

			errors = append(errors, fmt.Sprintf("Error executing query: %s. Error: %s", ddlQuery, err.Error()))
			log.Error("An error occurred while executing query: %s. Error: %s", ddlQuery, err.Error())
		} else {
			log.Info("Successfully executed query: %s", ddlQuery)
		}
	}

	if len(errors) > 0 {
		log.Fatal("Failed to execute DDL queries: %s", strings.Join(errors, ", "))
		return fmt.Errorf("Failed to execute DDL queries: %s", strings.Join(errors, ", "))
	}

	return nil
}

func seedDuckDB(db *sql.DB, configBaseURL string) error {
	if db == nil {
		return errors.New("nil DB provided to seedDuckDB")
	}

	var tables []string

	tablesEnv := config.Config.Duckdb.Seed_tables_from_file
	if tablesEnv == "" {
		return errors.New("Seed data files not set")
	}

	tables = strings.Split(tablesEnv, ",")
	errors := []string{}

	for _, tablename := range tables {
		var fileNameWithoutExtension string
		var fileType string

		if strings.HasSuffix(tablename, ".csv") {
			fileType = "csv"
			fileNameWithoutExtension = strings.TrimSuffix(tablename, "."+fileType)
			log.Info("File Type: %s", fileType)
			log.Info("File Name without Extension: %s", fileNameWithoutExtension)
		} else if strings.HasSuffix(tablename, ".json") {
			fileType = "json"
			fileNameWithoutExtension = strings.TrimSuffix(tablename, "."+fileType)
			log.Info("File Type: %s", fileType)
			log.Info("File Name without Extension: %s", fileNameWithoutExtension)
		} else {
			errors = append(errors, fmt.Sprintf("Unsupported file type for %s", tablename))
			continue
		}

		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s AS SELECT * FROM read_%s_auto('%s/%s.%s');", fileNameWithoutExtension, fileType, configBaseURL, fileNameWithoutExtension, fileType)
		log.Info("Executing query: %s", query)
		_, err := db.Exec(query)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error executing query for %s: %s", tablename, err.Error()))
		} else {
			log.Info("Executed query: %s", query)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Failed to execute queries: %s", strings.Join(errors, ", "))
	}

	return nil
}

func CreateEntity(c *gin.Context) {
	entityType := c.Param("entity")

	// Parse JSON data from request body
	var entityData map[string]interface{}
	if err := c.ShouldBindJSON(&entityData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction to manage the ID sequence and entity insertion
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the current ID from the sequence table and update it
	var nextID int
	err = tx.QueryRow("SELECT next_id FROM id_sequence WHERE entity_name = ?", entityType).Scan(&nextID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If the entity doesn't exist in the sequence table, create it
			_, err = tx.Exec("INSERT INTO id_sequence (entity_name, next_id) VALUES (?, 1)", entityType)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			nextID = 1
		} else {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Construct the column list and placeholders for the dynamic query
	columns := []string{"id"}
	placeholders := []string{"?"}
	values := []interface{}{nextID}
	for columnName, value := range entityData {
		columns = append(columns, columnName)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	columns = append(columns, "created_at", "modified_at", "update_count")
	placeholders = append(placeholders, "?, ?, ?")
	values = append(values, time.Now(), time.Now(), 0)

	// Construct the dynamic SQL query for entity insertion
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", entityType, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	// Insert the entity using the retrieved ID and JSON data
	_, err = tx.Exec(query, values...)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the sequence table with the next ID
	_, err = tx.Exec("UPDATE id_sequence SET next_id = next_id + 1 WHERE entity_name = ?", entityType)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": nextID})
}

func CreateEntityCore(entityType string, entityData map[string]interface{}) (int, error) {

	// Start a transaction to manage the ID sequence and entity insertion
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// Retrieve the current ID from the sequence table and update it
	var nextID int
	err = tx.QueryRow("SELECT next_id FROM id_sequence WHERE entity_name = ?", entityType).Scan(&nextID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If the entity doesn't exist in the sequence table, create it
			_, err = tx.Exec("INSERT INTO id_sequence (entity_name, next_id) VALUES (?, 1)", entityType)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			nextID = 1
		} else {
			tx.Rollback()
			return 0, err
		}
	}

	// Construct the column list and placeholders for the dynamic query
	columns := []string{"id"}
	placeholders := []string{"?"}
	values := []interface{}{nextID}
	for columnName, value := range entityData {
		columns = append(columns, columnName)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	columns = append(columns, "created_at", "modified_at", "update_count")
	placeholders = append(placeholders, "?, ?, ?")
	values = append(values, time.Now(), time.Now(), 0)

	// Construct the dynamic SQL query for entity insertion
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", entityType, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	log.Info("Executing query: %s", query)
	// Insert the entity using the retrieved ID and JSON data
	//log.Info("Values: %s", values)
	_, err = tx.Exec(query, values...)
	if err != nil {
		log.Info("Error executing query: %s", err.Error())
		tx.Rollback()
		return 0, err
	}

	// Update the sequence table with the next ID
	_, err = tx.Exec("UPDATE id_sequence SET next_id = next_id + 1 WHERE entity_name = ?", entityType)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func GetEntities(c *gin.Context) {
	entityType := c.Param("entity")
	var columns string

	// Get the list of columns to be retrieved from the query parameters
	colsParam, hasColsParam := c.GetQuery("select")
	if hasColsParam {
		columns = colsParam
	} else {
		columns = "*"
	}

	// Sorting
	var orderBy string
	orderParam, hasOrderParam := c.GetQuery("order")
	if hasOrderParam {
		orderBy = orderParam
	}

	// Construct the base query for total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", entityType)

	// Parse where conditions from query parameters
	var whereConditions []string
	for key, values := range c.Request.URL.Query() {
		if key != "select" && key != "limit" && key != "offset" && key != "order" {
			operator := "="
			if strings.HasSuffix(key, ".eq") {
				operator = "="
				key = strings.TrimSuffix(key, ".eq")
			} else if strings.HasSuffix(key, ".gt") {
				operator = ">"
				key = strings.TrimSuffix(key, ".gt")
			} else if strings.HasSuffix(key, ".gte") {
				operator = ">="
				key = strings.TrimSuffix(key, ".gte")
			} else if strings.HasSuffix(key, ".lt") {
				operator = "<"
				key = strings.TrimSuffix(key, ".lt")
			} else if strings.HasSuffix(key, ".lte") {
				operator = "<="
				key = strings.TrimSuffix(key, ".lte")
			} else if strings.HasSuffix(key, ".neq") {
				operator = "<>"
				key = strings.TrimSuffix(key, ".neq")
			} else if strings.HasSuffix(key, ".like") {
				operator = "ILIKE"
				key = strings.TrimSuffix(key, ".like")
			} // Add more cases for other operators

			// Construct filter conditions for each value
			for _, value := range values {
				// Handle spaces in the value by quoting it
				quotedValue := strings.ReplaceAll(value, "'", "''") // Escape single quotes
				whereConditions = append(whereConditions, fmt.Sprintf("LOWER(%s) %s LOWER('%s')", key, operator, quotedValue))
			}
		}
	}

	// Add WHERE clause if conditions are present
	if len(whereConditions) > 0 {
		countQuery = fmt.Sprintf("%s WHERE %s", countQuery, strings.Join(whereConditions, " OR "))
	}

	// Execute the count query
	var totalCount int
	err := db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Construct the base query for data retrieval
	baseQuery := fmt.Sprintf("SELECT %s FROM %s", columns, entityType)

	// Add WHERE clause if conditions are present
	if len(whereConditions) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(whereConditions, " OR "))
	}

	// Add ORDER BY clause if sorting is requested
	if orderBy != "" {
		baseQuery = fmt.Sprintf("%s ORDER BY %s", baseQuery, orderBy)
	}

	// Add LIMIT and OFFSET clauses
	query := fmt.Sprintf("%s LIMIT %d OFFSET %d", baseQuery, limit, offset)

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var entities []map[string]interface{}
	for rows.Next() {
		entity := make(map[string]interface{})
		values := make([]interface{}, len(columnNames))
		valuePointers := make([]interface{}, len(columnNames))
		for i := range columnNames {
			valuePointers[i] = &values[i]
		}

		err := rows.Scan(valuePointers...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for i, colName := range columnNames {
			entity[colName] = values[i]
		}
		entities = append(entities, entity)
	}

	// Calculate the current page based on offset and limit
	currentPage := (offset / limit) + 1

	// Calculate the next offset
	nextOffset := (currentPage * limit) + 1

	start := offset + 1
	end := offset + limit
	if end > totalCount {
		end = totalCount
	}

	// Construct the response object with additional information
	response := gin.H{
		"total_rows":   totalCount, // Updated to use the total count from the count query
		"limit":        limit,
		"offset":       offset,
		"current_page": currentPage,
		"next_offset":  nextOffset,
		"start":        start,
		"end":          end,
		"data":         entities,
	}

	c.JSON(http.StatusOK, response)
}

func UpdateEntity(c *gin.Context) {
	entityType := c.Param("entity")
	id := c.Param("id")

	// Get JSON payload from the request body
	var updatedEntity map[string]interface{}
	if err := c.ShouldBindJSON(&updatedEntity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Build the update query and arguments
	var updateQuery string
	var args []interface{}
	for key, value := range updatedEntity {
		if key != "id" {
			updateQuery += key + " = ?, "
			args = append(args, value)
		}
	}
	updateQuery += "modified_at = ?, update_count = update_count + 1"
	args = append(args, time.Now())

	// Execute the update query
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", entityType, updateQuery)
	args = append(args, id)
	_, err := db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entity updated successfully"})
}

func DeleteEntity(c *gin.Context) {
	entityType := c.Param("entity")
	id := c.Param("id")

	// Delete the entity from the database
	_, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", entityType), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entity deleted successfully"})
}

func ExecuteSelectQuery(query string) ([]map[string]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for rows.Next() {
		entity := make(map[string]interface{})
		values := make([]interface{}, len(columnNames))
		valuePointers := make([]interface{}, len(columnNames))
		for i := range columnNames {
			valuePointers[i] = &values[i]
		}

		err := rows.Scan(valuePointers...)
		if err != nil {
			return nil, err
		}

		for i, colName := range columnNames {
			entity[colName] = values[i]
		}
		result = append(result, entity)
	}

	return result, nil
}

func ExecuteDDLQuery(query string) error {
	println(query)
	_, err := db.Exec(query)
	return err
}

func ExecuteCustomQuery(c *gin.Context) {
	var queryData struct {
		Query string `json:"query" binding:"required"`
	}

	if err := c.ShouldBindJSON(&queryData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.HasPrefix(strings.ToUpper(queryData.Query), "SELECT") {
		// Handle SELECT query
		results, err := ExecuteSelectQuery(queryData.Query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := gin.H{
			"total_rows": len(results),
			"data":       results,
		}

		c.JSON(http.StatusOK, response)
	} else {
		// Handle DDL query
		print("Executing DDL query " + queryData.Query)
		err := ExecuteDDLQuery(queryData.Query)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusOK, gin.H{"error": err.Error(), "message": "Query execution failed"})
			return
		}
		print("Executed DDL query ")
		c.JSON(http.StatusOK, gin.H{"message": "Query executed successfully"})
	}
}
