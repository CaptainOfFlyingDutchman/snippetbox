### Start the app using the command

```bash
cd snippetbox
go run ./cmd/web
```

---

## Database Setup (Docker)

This project uses MySQL 9.x running inside a Docker container. Follow these steps to spin up the database and connect your development tools.

### 1. Run the MySQL Container

Choose **one** of the following methods to start the database. Both methods map the container to local host port `3306` to prevent conflicts with local MySQL installations.

#### Method A: Docker CLI Command (Fastest)

Run this command in your terminal to instantly pull the image and start the container with a managed volume:

```bash
docker run --name snippetbox-db \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=your_root_password \
  -v snippetbox_data:/var/lib/mysql \
  -d mysql:9.7.0
```

#### Method B: Docker Compose (Recommended)

If you prefer configuration files, create a `docker-compose.yml` file in your project root:

```yaml
version: "3.8"
services:
  db:
    image: mysql:9.7.0
    container_name: snippetbox-db
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: your_root_password
    volumes:
      - snippetbox_data:/var/lib/mysql

volumes:
  snippetbox_data:
```

Start it by running:

```bash
docker compose up -d
```

---

### 2. Configure Database Users & Permissions

Because Docker isolates the container, the application user must be configured to accept external connections from your host machine (WebStorm/Go/Node).

1. Access the running container's terminal:
   ```bash
   docker exec -it snippetbox-db mysql -u root -p
   ```
2. Enter the `your_root_password` you set in step 1.
3. Run the following SQL commands to initialize the database and the external `web` user:
   ```sql
   CREATE DATABASE IF NOT EXISTS snippetbox;
   CREATE USER IF NOT EXISTS 'web'@'%' IDENTIFIED BY 'pass';
   GRANT ALL PRIVILEGES ON snippetbox.* TO 'web'@'%';
   FLUSH PRIVILEGES;
   ```

---

### 3. WebStorm Connection Settings

To browse the database inside DataGrip, open the **Database** tool window (`View -> Tool Windows -> Database`) and add a new **MySQL** data source:

- **Host:** `localhost`
- **Port:** `3306`
- **User:** `web`
- **Password:** `pass`
- **Database:** `snippetbox`

Click **Test Connection** to verify WebStorm can communicate with the container, then click **Apply**.

### 4. DataGrip Connection Settings

To manage and query the database using DataGrip, follow these steps:

1. Open DataGrip and look at the **Database Explorer** tool window (usually on the left side).
2. Click the **`+` (New)** icon -> **Data Source** -> **MySQL**.
3. In the settings window that appears, fill out the details exactly as follows:
   - **Host:** `localhost`
   - **Port:** `3306`
   - **User:** `web`
   - **Password:** `pass`
   - **Database:** `snippetbox`
4. If you see a warning at the bottom saying _"Missing driver files"_, simply click the blue **Download** link next to it.
5. Click **Test Connection** to verify DataGrip can communicate with the container, then click **OK**.
