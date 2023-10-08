# Routes

## User Service

POST /login -> Login User, returns JWT token
POST /register -> Register User
DELETE /session -> Invalidates JWT token

(optional) GET /token -> Create an API Token for generating PDF´s

## Templating Service

List JSON
```json
{
    "id": 7,
    "name": "My table template",
    "created_at": "2023-10-08 00:00Z",
    "updated_at": "2023-10-08 00:00Z",
}
```

Detail JSON
```json
{
    "id": 7,
    "name": "My table template",
    "created_at": "2023-10-08 00:00Z",
    "updated_at": "2023-10-08 00:00Z",
    "elements": [
        { 
            "type": "rect"
            "x": 0, 
            "y": 0, 
            "width": 0, 
            "height": 0
        },
        {

            "type": "text"
            "x": 0, 
            "y": 0, 
            "width": 0, 
            "height": 0
            "content": "Hello world",
            "font": "JetBrainsMono"
            "size": 18
    ]
}
```

GET /templates -> All template
GET /templates/:id -> Detail template
POST /templates -> Create template
PUT /templates/:id -> Update template
DELETE /templates/:id -> Delete templates

## Pdf Creation Service

POST /render/:id -> Render a pdf by template and data
