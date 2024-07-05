package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"

	initi "easyauthapi/configs"
)

func init() {
	configFile := flag.String("config", "app.env", "Name of the config file (without extension)")
	initi.LoadConfigViper("../", configFile)
	initi.ConnectDB()
}

func DBOutputToSQL() string {
	cmd := exec.Command("pg_dump",
		"-U", initi.UseConfig.PGUSER,
		"-d", initi.UseConfig.PGDATABASE,
		"-h", initi.UseConfig.PGHOST,
		"-p", initi.UseConfig.PGPORT,
	)

	// Set the password environment variable for pg_dump
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", initi.UseConfig.PGPASSWORD))

	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running pg_dump: %v", err)
	}

	return string(output)
}

func DBInputFromSQL(SQL string) {
	initi.DB.Exec(SQL)
}

func main() {
}
