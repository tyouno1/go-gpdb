package library

import (
	"os/exec"
	log "../../core/logger"
	"../objects"
	"strconv"
	"../../core/arguments"
	"os"
)

// Execute DB Query
func execute_db_query(query_string string, to_write bool, file_name string) ([]byte, error) {

	// Set GPDB Environment
	err := SourceGPDBPath()
	if err != nil { return []byte(""), err}

	// Execute command
	master_port := strconv.Itoa(objects.GpInitSystemConfig.MasterPort)
	cmd := exec.Command("psql", "-p", master_port , "-d", "template1", "-Atc", query_string)

	// If request to file, then write to o/p file
	if to_write {
		outfile, err := os.Create(file_name)
		if err != nil { return []byte(""), err }
		defer outfile.Close()
		cmd.Stdout = outfile
	}

	// Start the execution of command.
	err = cmd.Start()
	if err != nil { return []byte(""), err }
	err = cmd.Wait()
	if err != nil { return []byte(""), err }

	return []byte(""), nil
}

// Check if the database is healthy
func IsDBHealthy() error {

	// Query string
	log.Println("Ensuring the database is healthy...")
	query_string := "select 1"

	// Execute string
	_, err := execute_db_query(query_string, false, "")
	if err != nil { return err }

	return nil
}

// Build uninstall script
func CreateUnistallScript(t string) error {

	// Uninstall script location
	UninstallFileDir := arguments.UninstallDir + "uninstall_" + arguments.RequestedInstallVersion + "_" + t
	log.Println("Creating Uninstall file for this installation at: " + UninstallFileDir)

	// Query
	query_string := "select $$ps -ef|grep postgres|grep -v grep|grep $$ ||  port " +
			"|| $$ | awk '{print $2}'| xargs -n1 /bin/kill -11 &>/dev/null $$ " +
			"from gp_segment_configuration union select $$rm -rf /tmp/.s.PGSQL.$$ || port || $$*$$ " +
		        "from gp_segment_configuration union select $$rm -rf $$ || fselocation from pg_filespace_entry"

	// Execute the query
	_, err := execute_db_query(query_string, true, UninstallFileDir)
	if err != nil { return err }

	return nil
}
