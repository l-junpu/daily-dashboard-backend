# daily-dashboard-backend

Description Here...

## Installing Dependencies

### MongoDB

```
docker pull mongodb/mongodb-community-server:latest
docker run -d -v ./mongodb:/data/db --name mongodb -p 27017:27017 mongodb/mongodb-community-server
```

If you have previously installed MongoDB locally, to disable it:

1. Press: Win + R
2. Type services.msc
3. Find the "MongoDB Server (MongoDB)" service and select "Stop"

### Redis

```
docker pull redis
docker run -d --name redis -p 6379:6379 redis
```

### MSSQL Express

Install MSSQL Express here:
https://www.microsoft.com/en-sg/sql-server/sql-server-downloads

Additional configurations that may be required:

1. Access the `SQL Server 2022 Configuration Manager` under `Programs\Microsoft SQL Server 2022\Configuration Tools`
2. Under `SQL Server Network Configuration`, go to `Protocols for SQLEXPRESS` and enable `TCP/IP`
3. Navigate back to `SQL Server Configuration Manager (Local)` and open up the `properties` page for `SQL Server Browser`
4. Navigate to the `Service` tab and set the `Start Mode` to `Automatic`
5. Restart `SQL Server (SQLEXPRESS)` and you should be able to establish a connection to the Database

### Go Packages

When developing the project, you may install packages from here:
https://pkg.go.dev/
