package mocks

//go:generate mockery -name Customer  -dir ../../domain -output . -note "+build test"
//go:generate mockery -name PersistableCustomers -dir ../ -output . -note "+build test"
//go:generate mockery -name StartsRepositorySessions -dir ../ -output . -note "+build test"
//go:generate mockery -name Command  -dir ../../../shared -output . -note "+build test"
