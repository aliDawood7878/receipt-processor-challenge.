This application accepts receipts and returns points based on certain rules. It provides two endpoints. One endpoint processes a receipt and returns an id. Another endpoint returns points for an id.

Commands to build and run:
docker build -t receipt-processor .

docker run -p 8080:8080 receipt-processor

The service is available at:
http://localhost:8080

Use POST /receipts/process with a JSON payload that includes retailer, purchaseDate, purchaseTime, items, and total. For example below-

{
“retailer”: “Target”,
“purchaseDate”: “2022-01-01”,
“purchaseTime”: “13:01”,
“items”: [
{“shortDescription”: “Mountain Dew 12PK”, “price”: “6.49”},
{“shortDescription”: “Emils Cheese Pizza”, “price”: “12.25”},
{“shortDescription”: “Knorr Creamy Chicken”, “price”: “1.26”},
{“shortDescription”: “Doritos Nacho Cheese”, “price”: “3.35”},
{“shortDescription”: “Klarbrunn 12-PK 12 FL OZ “, “price”: “12.00”}
],
“total”: “35.35”
}

Send this payload to:
http://localhost:8080/receipts/process

Use a POST request and set Content-Type to application/json.
response: {“id”:“some-id”}

Use the returned id for the next endpoint.

Use GET /receipts/{id}/points to return points for the given id. For example:
http://localhost:8080/receipts/some-id/points
response: {“points”: 28}


If any field is not included in the POST request, the server returns 400. For example, remove total:
{
“retailer”: “Target”,
“purchaseDate”: “2022-01-01”,
“purchaseTime”: “13:01”,
“items”: [
{“shortDescription”: “Mountain Dew 12PK”, “price”: “6.49”}
]
}

This returns 400.

If purchaseDate does not match YYYY-MM-DD, the server returns 400. For example:

{
“retailer”: “Target”,
“purchaseDate”: “2022/01/01”,
“purchaseTime”: “13:01”,
“items”: [
{“shortDescription”: “Mountain Dew 12PK”, “price”: “6.49”}
],
“total”: “35.35”
}