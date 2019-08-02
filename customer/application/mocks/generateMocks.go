package mocks

//go:generate mockery -name Customer  -dir ../../domain -output .
//go:generate mockery -name PersistableCustomers -dir ../ -output .
//go:generate mockery -name StartsRepositorySessions -dir ../ -output .
//go:generate mockery -name Command  -dir ../../../shared -output .
