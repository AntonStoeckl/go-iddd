rsconf = {
	_id: "rs0",
	members: [
		{_id: 0, host: "mongodb0:27017"},
		{_id: 1, host: "mongodb1:27017"},
		{_id:2, host: "mongodb2:27017"}
	]
}
rs.initiate(rsconf);
rs.conf();