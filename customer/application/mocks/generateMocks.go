package mocks

//go:generate mockery -name Customers -dir ../ -output . -note "+build test"
//go:generate mockery -name StartsCustomersSession -dir ../ -output . -note "+build test"
//go:generate mockery -name Command  -dir ../../../shared -output . -note "+build test"
