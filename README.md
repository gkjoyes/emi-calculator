# emi-calculator

API is to calculate the break-even (the point of balance) where it makes equal sense to either rent or buy a good.

## Requirements

- [Docker](https://docs.docker.com/install/linux/docker-ce/debian/)

- [Docker Compose](https://docs.docker.com/compose/install/)

## Setup

- Clone the repository [`go get github.com/gkjoyes/emi-calculator`]

- Switch to project directory and run [`docker-compose up`]

## API

- [POST] `http://127.0.0.1:5200/emi-calculator`

- The sample request body is shown below:

    ```json
    {
        "down_payment":2000,
        "total_amount":2000000,
        "interest_rate":10,
        "property_tax":2000,
        "property_transfer_tax":25000,
        "years_expected_to_live":20
    }
    ```

- This API will return monthly EMI.
