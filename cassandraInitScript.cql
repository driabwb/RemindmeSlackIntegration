CREATE KEYSPACE RemindMeSlackIntegration WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
CREATE TABLE RemindMeSlackIntegration.Messages ( alertTime timestamp, messages text, date int, user_name text, Primary Key(date, alertTime)) with CLUSTERING ORDER BY (alerttime ASC ) ;
CREATE TABLE RemindMeSlackIntegration.Users ( ID text, name text, Primary Key(ID));

