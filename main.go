package main

import (
	"os"
	"rule_service/commons"
	"rule_service/infra"
	"rule_service/ruleset"
)

const (
	defaultPort              = "8080"
	defaultRoutingServiceURL = "http://localhost:7878"
)

func main() {
	// initialise application
	os.Setenv("ENVIRONMENT", commons.PROD)
	serverContext := infra.New()

	// Register Routes
	rulesetContext := ruleset.NewRulesetContext()
	defer rulesetContext.SafeClose()

	serverContext.Mount("/api/v1/rulesets", rulesetContext.NewRulesetRouter())

	// Start server
	serverContext.StartServer(":3000")
}
