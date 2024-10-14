# daily-dashboard-backend

Description Here...

## Installing Dependencies

### MongoDB

[Install MongoDB Community with Docker](https://www.mongodb.com/docs/manual/tutorial/install-mongodb-community-with-docker/)

```
docker pull mongodb/mongodb-community-server:latest
docker run -d -v ./mongodb:/data/db --name mongodb -p 27017:27017 mongodb/mongodb-community-server
```

If you have previously installed MongoDB locally, to disable it:

1. Press: Win + R
2. Type services.msc
3. Find the "MongoDB Server (MongoDB)" service and select "Stop"

### Redis

[DockerHub - Redis](https://hub.docker.com/_/redis)

```
docker pull redis
docker run -d --name redis -p 6379:6379 redis
```

### MSSQL Express

[Microsoft - Quickstart: Run SQL Server Linux container images with Docker](https://learn.microsoft.com/en-us/sql/linux/quickstart-install-connect-docker?view=sql-server-ver16&tabs=cli&pivots=cs1-bash)

```
docker run -e "ACCEPT_EULA=Y" -e "MSSQL_SA_PASSWORD=<REPLACE_THIS_WITH_YOUR_STRONG_PASSWORD>" -p 1433:1433 --name mssql-express -v ./mssql-data:/var/opt/mssql -d mcr.microsoft.com/mssql/server:2022-latest
```

Upon running the following command, MSSQL Express by default creates an account with the following details:

```
Username: "sa"
Password: <REPLACE_THIS_WITH_YOUR_STRONG_PASSWORD>
```

As a _"best practice"_, you would typically disable the `sa` account - [As Per Microsoft](https://learn.microsoft.com/en-us/sql/linux/quickstart-install-connect-docker?view=sql-server-ver16&tabs=cli&pivots=cs1-bash)

However, that is completely up to **you**, since we are intending to run it in an offline environment. If you wish to disable it, refer to **"Additional Notes"** at the end of the `README` for a guide to perform it, otherwise, carry on with the setup.

To allow the Daily-Dashboard-Backend to connect to your hosted MSSQL Express DB, specify the following details inside the `.env` file, and the application will automatically read in those details.

```
MSSQL_USERNAME=<YOUR_USERNAME>
MSSQL_PASSWORD=<YOUR_PASSWORD>
```

_Note: Ensure that you follow the exact naming of the environment variables._

### Go Packages

When developing the project, you may install packages from here:
https://pkg.go.dev/

## Additional Notes

### Disabling the "sa" account for MSSQL Express

We will be using [Microsoft - SQL Server Management Studio (SSMS)](https://learn.microsoft.com/en-us/sql/ssms/download-sql-server-management-studio-ssms?view=sql-server-ver16) to perform the following task. Install it if you do not have it yet.

#### Connect to our newly created MSSQL Express Database (If you followed the steps I mentioned above)

```
Server Name:    localhost,1433
Authentication: SQL Server Authentication
Login:          sa
Password:       <REPLACE_THIS_WITH_YOUR_STRONG_PASSWORD>
```

#### Create a new Account

1. Expand the `Server Instance`, and then expand the `Security` Folder
2. Right-click on the `Logins` folder and select `New Login`
3. In the pop-up, enter your `Login Name` (e.g. NewAdmin)
4. Select `SQL Server Authentication` and enter a `strong password`
5. Click OK to create the login

#### Give our newly created Account Administrator rights

1. Expand the `Server Instance`, and then expand the `Security` Folder
2. Right-click on the newly created `Login Name` (e.g., NewAdmin) and select `Properties`
3. In the pop-up, go to the `Server Roles` page
4. Check the `sysadmin` role to add the login to this role
5. Click OK to save the changes

#### Disable our old SA Account

1. Expand the `Server Instance`, and then expand the `Security` Folder
2. Right-click on the `sa` account and select `Properties`
3. In the pop-up, go to the `Status` page
4. Under `Login`, select `Disabled`
5. Click OK to save the changes
6. Disconnect from the Database and verify that your new account works, and the SA account is disabled
