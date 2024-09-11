## How to run

```sh
go mod vendor

docker build -t workspace3-image .
docker run -p 8080:4000 workspace3-image
```

I tried following `curl` commands in a different terminal:

```sh
curl http://0.0.0.0:8080/receipts/process \
-H "Content-Type: application/json" -H "Accept: application/json" \
--data-raw '{
    "retailer": "Walgreens",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "08:13",
    "total": "2.65",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
}'
# {"id":"ae06b97b-e034-4f51-bf62-456a81b7911c"}

curl http://0.0.0.0:8080/receipts/e2bdcb82-5884-41f0-a6e8-491a04dfc329/points \
 -H "Content-Type: application/json" -H "Accept: application/json"
# {"points":15}
```

```sh
curl http://0.0.0.0:8080/receipts/process \
-H "Content-Type: application/json" -H "Accept: application/json" \
--data-raw '{
    "retailer": "Target",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "13:13",
    "total": "1.25",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"}
    ]
}'
# {"id":"56701cec-50ad-4977-9f22-3697e2b6f93f"}

curl http://0.0.0.0:8080/receipts/56701cec-50ad-4977-9f22-3697e2b6f93f/points \
-H "Content-Type: application/json" -H "Accept: application/json"
# {"points":31}
```

```sh
curl http://0.0.0.0:8080/receipts/process \
-H "Content-Type: application/json" -H "Accept: application/json" \
--data-raw '{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}'
# {"id":"a9370346-1828-43b1-87ed-8942c8c54df9"}

curl http://0.0.0.0:8080/receipts/a9370346-1828-43b1-87ed-8942c8c54df9/points \
-H "Content-Type: application/json" -H "Accept: application/json"
# {"points":28}
```

```sh
curl http://0.0.0.0:8080/receipts/process \
-H "Content-Type: application/json" -H "Accept: application/json" \
--data-raw '{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}'
# {"id":"203b8c6e-4eac-4626-9467-73aecf4b072a"}

curl http://localhost:8080/receipts/203b8c6e-4eac-4626-9467-73aecf4b072a/points \
-H "Content-Type: application/json" -H "Accept: application/json"
# {"points":109}
