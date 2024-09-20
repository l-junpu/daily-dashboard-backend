# daily-dashboard-backend

## Configurations That May Be Required

1. Access the `SQL Server 2022 Configuration Manager` under `Programs\Microsoft SQL Server 2022\Configuration Tools`
2. Under `SQL Server Network Configuration`, go to `Protocols for SQLEXPRESS` and enable `TCP/IP`
3. Navigate back to `SQL Server Configuration Manager (Local)` and open up the `properties` page for `SQL Server Browser`
4. Navigate to the `Service` tab and set the `Start Mode` to `Automatic`
5. Restart `SQL Server (SQLEXPRESS)` and you should be able to establish a connection to the Database
