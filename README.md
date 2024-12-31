# bitaksi-casestudy

## Prerequisites

Ensure that you have Docker and Docker Compose installed.

## Setup

1.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd <project_directory>
    ```

2.  **Start Docker Compose:**

    Navigate to the root of your project where the `docker-compose.yml` file is located and run:

    ```bash
    docker-compose up -d
    ```

## API Usage

After starting the Docker environment, APIs should be accessible at `http://localhost:<port>`. Replace `<port>` with the port exposed in `docker-compose.yml` file.

*  **Example**

    Send a POST request to `/locations` to get nearest driver location

    ```bash
    curl --location 'http://localhost:9600/api/v1/match/driver?radius=10000000000' \
        --header 'Content-Type: application/json' \
        --header 'Accept: application/json' \
        --header 'Authorization: your-jwt-token' \
        --data '{
        "type": "Point",
        "coordinates": [
            40.4,
            29.2
        ]
        }'
    ```

## Teardown

To stop the Docker Compose environment:

```bash
docker-compose down