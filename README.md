
# Greenpeace Planet4 docker development environment

![Planet4](https://cdn-images-1.medium.com/letterbox/300/36/50/50/1*XcutrEHk0HYv-spjnOej2w.png?source=logoAvatar-ec5f4e3b2e43---fded7925f62)

## What is Planet4?

Planet4 is the NEW Greenpeace web platform

## What is this repository?

This repository contains needed files to set up a docker development environment that consists of:

*   [MySQL](https://hub.docker.com/_/mysql/) container as database engine
*   [Traeik](https://traefik.io) load balancing ingress controller
*   [OpenResty](https://openresty.org) dynamic web platform based on NGINX and Lua
*   [php-fpm](https://php-fpm.org/) high performance PHP FastCGI implementation
*   [Redis](https://redis.io/) key-value store caching FastCGI, object and session data
*   [PHPmyadmin](https://hub.docker.com/r/phpmyadmin/phpmyadmin/) for database administration

By default, the quickstart command `make` is all you'll need to pull all required images and spin up a load balanced nginx/php/redis/mysql web application with automatic SSL generation in the comfort of your own office.

*   Traefik listens on Port 80, load balancing requests to:
*   Two OpenResty reverse proxy servers, which cache FastCGI requests from
*   Three PHP-FPM application nodes, all backed by
*   A single Redis instance and
*   MySQL database server.
*   Self-signed SSL certificates, with HTTP > HTTPS redirection

## Quickstart

*Note this repository has been tested extensively on OSX, should be just fine in Ubuntu/Debian and is extremely unlikely to work in Windows even in their Ubuntu shell, any feedback will be welcome!*

```
# Clone the repository
git clone https://github.com/greenpeace/planet4-docker-compose

# Navigate to new directory
cd planet4-docker-compose

# Mac and Linux only
sudo echo "127.0.0.1  test.planet4.dev pma.planet4.dev planet.dev" >> /etc/hosts

# Start the application
make

# View log output
docker-compose logs -f
```

On first launch, the container bootstraps the installation with composer then after a short time (30 seconds to 1 minute) all services will be ready and responding to requests.

When you see the line `Starting service: openresty` you can navigate to: [https://test.planet4.dev](https://test.planet4.dev)

### Requirements

Firstly, requirements for running this development environment:

*   [Docker](https://docs.docker.com/engine/installation/)

For MacOS and Windows users docker installation already includes docker-compose
GNU/Linux users have to install docker-compose separately:

*   [Docker Compose](https://github.com/docker/compose/releases)


(optional)
*   [Make](https://www.gnu.org/software/make/) - Instructions for installing make vary, for OSX users `xcode-select --install` might work.

---

## Editing source code

By default, the Wordpress application is bind-mounted at
-   `./persistence/app/`
-   `./persistence/app/`

---

## <a name="login">Logging in</a>

### Administrator login

Backend administrator login is available at [https://test.planet4.dev/wp-admin/](https://test.planet4.dev/wp-admin/). An administrator user is created during first install with a randomly assigned password.

Login username is `admin`. To find the password enter the following in the project root (where docker-compose.yml lives):

```
make wppass
```

### Database access via phpMyAdmin

[phpmyadmin](https://hub.docker.com/r/phpmyadmin/phpmyadmin/) login: [https://pma.planet4.dev](https://pma.planet4.dev)

Enter the user values from `db.env` to login, or from bash prompt:

```
make pmapass
```
---

## Configuration

### Configuring WP-Stateless GCS bucket storage

[WP-Stateless](https://github.com/wpCloud/wp-stateless/) is installed and activated, however images will be stored locally until remote GCS storage is enabled in the administrator backend. [Log in](https://test.planet4.dev/wp-login.php) with details gathered [from here](#login) and navigate to [Media > Stateless Setup](https://test.planet4.dev/wp-admin/upload.php?page=stateless-setup).

You will need a Google account with access to GCS buckets to continue.

Once logged in:

*   Click 'Get Started Now'
*   Authenticate
*   Choose 'Planet-4' project (or a custom project from your private account)
*   Choose or create a Google Cloud Bucket - it's recommended to use a bucket name unique to your own circumstances, eg 'mynamehere-test-planet4-wordpress'
*   Choose a region close to your work environment
*   Skip creating a billing account (if using Greenpeace Planet 4 project)
*   Click continue, and wait a little while for all necessary permissions and object to be created.

Congratulations, you're now serving media files directly from GCS buckets!

### Configuring FastCGI cache purges

The Wordpress plugin [nginx-helper](https://wordpress.org/plugins/nginx-helper/) is installed to enable FastCGI cache purges. Log in to the backend as above, navigate to [Settings > Nginx Helper](https://test.planet4.dev/wp-admin/options-general.php?page=nginx) and click:
*   Enable Purge
*   Redis Cache
*   Enter `redis` in the Hostname field
*   Tick all checkboxes under 'Purging Conditions'

## Environment variables

This docker environment relies on the mysql official image as well as on the
[planet4-base](https://github.com/greenpeace/planet4-base) application image.

Both images provide environment variables which adjust aspects of the runtime configuration. For this environment to run only the database parameters such as hostname, database name, database users and passwords are required.

Initial values for this environment variables are dummy but are good to go for development porpoises. They can be changed in the provided `app.env` and `db.env` files, or directly in the [docker-compose.yml](https://docs.docker.com/compose/compose-file/#environment) file itself.

### Some useful variables
See [openresty-php-exim](https://github.com/greenpeace/planet4-docker/tree/develop/source/planet-4-151612/openresty-php-exim)

-   `NEWRELIC_LICENSE` set to the license key in your NewRelic dashboard to automatically receive server and application metrics
-   `PHP_MEMORY_LIMIT` maximum memory each PHP process can consume before being termminated and restarted by the scheduler
-   `PHP_XDEBUG_REMOTE_HOST` in development mode enables remote [XDebug](https://xdebug.org/) debugging, tracing and profiling

### Development mode

*@todo:
Document some of the useful builtin configuration options available in upstream docker images for debugging, including:*
-   XDebug remote debugging
-   Smarthost email delivery and interception
-   exec function limits
-   Memory and performance tweaks

---

## Notes

### Updating

To ensure you're running the latest version of both the infrastructure and the application:

```
make
```

This one simple command will ensure you're always running the latest version of the application and will perform the following:

*   stop the services
*   update the repository
*   delete the `persistence` directory
*   pull a new application image
*   restart docker compose

*Make sure you've pushed any changes you wish to keep first!*

Also, be aware that if you've only recently pushed new code to the repository there may be a delay of up to 30 minutes before the composer registry is updated.  You can always enter the relevant code directory and perform a `git pull` within the appropriate branch to speed things up.

```
make clean
make pull
make run
```

### Port conflicts

If you are running any other services on your local device which respond on port 80, you may experience errors attempting to start the environment. Traefik is configured to respond on port 80 in this application, but you can change it by editing the docker-compose.yml file as below:

```yaml
  traefik:
    ports:
      - "8000:80"
```

The first number is the port number on your host, the second number is mapped to port 80 on the openresty service container.  Now you can access the site at  [https://test.planet4.dev:8000](https://test.planet4.dev:8000) instead.

A more robust solution for hosting multiple services on port 80 is to use a reverse proxy  such as Traefik or [jwilder/openresty-proxy](https://github.com/jwilder/openresty-proxy) in a separate project, and use [Docker named networking](https://docs.docker.com/compose/networking/) features to isolate virtual networks.

## Traefik administration interface

Traefik comes with a simple admin interface accessible at [http://localhost:8080](http://localhost:8080).

## Performance

On a 2015 Macbook Pro, with a primed cache, this stack delivers 50 concurrent connections under siege with an average response time of 0.28 seconds and near-zero load on the php-fpm backend.

```bash
$ siege -c 50 -t 20s -b test.planet4.dev

Lifting the server siege...
Transactions:		          3374 hits
Availability:		          100.00 %
Elapsed time:		          19.39 secs
Data transferred:	        51.58 MB
Response time:		        0.28 secs
Transaction rate:	        174.01 trans/sec
Throughput:		            2.66 MB/sec
Concurrency:		          49.13
Successful transactions:  3374
Failed transactions:	    0
Longest transaction:	    2.31
Shortest transaction:	    0.00
```

To be continued...
