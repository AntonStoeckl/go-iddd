package mocks

//go:generate mockery -name Customer  -dir ../../domain -output .
//go:generate mockery -name CustomersWithPersistance -dir ../ -output .
//go:generate mockery -name Command  -dir ../../../shared -output .
