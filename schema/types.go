package schema

const redirect = `
{
  "$schema": "http://json-schema.org/draft-07/schema",
  "$id": "swerve-redirect22232",
  "type": "object",
  "examples": [
    {
      "redirect_from": "mydomain.com",
      "path_map": [
        {
          "from": "/sport",
          "to": "/de/sport"
        }
      ],
      "redirect_to": "bild.de",
      "promotable": false,
      "code": 301
    }
  ],
  "required": [
    "redirect_from",
    "redirect_to",
    "code"
  ],
  "properties": {
    "redirect_from": {
      "$id": "#/properties/redirect_from",
      "type": "string",
      "pattern": "^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$"
    },
    "path_map": {
      "$id": "#/properties/path_map",
      "type": [
        "array",
        "null"
      ],
      "additionalItems": false,
      "items": {
        "type": "object",
        "required":["from","to"],
        "properties": {
          "from2": {
            "type": "string",
            "pattern": "^[/|%2F]{1,}([/a-zA-Z0-9-_%,\\.])*\\??[a-zA-Z0-9-_=~%,\\.]*&?[a-zA-Z0-9-_=~%,\\.]*$"
          },
          "to": {
            "type": "string",
            "pattern": "^[/|%2F]{1,}([/a-zA-Z0-9-_%,\\.])*\\??[a-zA-Z0-9-_=~%,\\.]*&?[a-zA-Z0-9-_=~%,\\.]*$"
          }
        }
      }
    },
    "redirect_to": {
      "$id": "#/properties/redirect_to",
      "type": "string",
      "pattern": "^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$"
    },
    "promotable": {
      "$id": "#/properties/promotable",
      "type": "boolean"
    },
    "code": {
      "$id": "#/properties/code",
      "type": "integer",
      "enum": [
        301,
        302,
        304,
        305,
        307,
        308
      ]
    },
    "description": {
      "$id": "#/properties/description",
      "type": "string"
    }
  },
  "additionalProperties": false
}`
