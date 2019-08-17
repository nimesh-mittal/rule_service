# Overview
Generic rule service to maintain and evaluate rules as REST service.

[![Go Report Card](https://goreportcard.com/badge/github.com/nimesh-mittal/rule_service)](https://goreportcard.com/report/github.com/nimesh-mittal/rule_service)

#### Supported Operators
| Operator | Supported |
|----------|-----------|
| '>'      | Yes       |
| <        | Yes       |
| >=       | Yes       |
| <=       | Yes       |
| ==       | Yes       |
| !=       | Yes       |
| in       | Yes       |
| not_in   | Yes       |

# Setup
- step 1: change MongoURL in config/config.yaml
- step 2: dep ensure
- step 3: go run main.go

# APIs

Postman Collection Link: https://github.com/nimesh-mittal/rule_service/blob/master/ruleset.postman_collection.json

- Create Ruleset
- List Ruleset
- Update Ruleset
- Delete Ruleset
- Get Ruleset By ID
- Evaluate Ruleset: Takes ruleset id and record and return matching rule with modified record

# Extendability
### Adding support for new data type
### Adding support for new operator
### Adding new rule strategy
Adding new strategy to pick one out of multiple matching rules

