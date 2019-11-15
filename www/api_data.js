define({ "api": [
  {
    "type": "delete",
    "url": "/v1/composition/:compositionId",
    "title": "",
    "name": "Delete",
    "group": "Composition",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "compositionId",
            "description": "<p>Composition ID</p>"
          }
        ]
      }
    },
    "description": "<p>Deletes an existing Composition by ID. To be deleted it doesn't have to be used as dependency in another Composition.</p>",
    "success": {
      "examples": [
        {
          "title": "Response",
          "content": "HTTP/1.1 200 OK\n\n{\n\t\"status\": \"DELETED\"\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./infrastructure/composition/rest.go",
    "groupTitle": "Composition"
  },
  {
    "type": "get",
    "url": "/v1/comosition/:compositionId",
    "title": "",
    "name": "GetByID",
    "group": "Composition",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "compositionId",
            "description": "<p>Composition ID</p>"
          }
        ]
      }
    },
    "description": "<p>Gets a composition by ID</p>",
    "success": {
      "examples": [
        {
          "title": "Response",
          "content": "HTTP/1.1 200 OK\n{\n  \"composition\": {\n    \"id\": \"9dc9c429b9aa2a3c82801007\",\n    \"name\": \"Comp 7\",\n    \"cost\": 475.75,\n    \"unit\": {\n      \"quantity\": 3,\n      \"unit\": \"u\"\n    },\n    \"stock\": {\n      \"quantity\": 1,\n      \"unit\": \"u\"\n    },\n    \"dependencies\": [\n      {\n        \"of\": \"9dc9c429b9aa2a3c82801005\",\n        \"quantity\": {\n          \"quantity\": 2,\n          \"unit\": \"u\"\n        },\n        \"subvalue\": 82\n      },\n      {\n        \"of\": \"9dc9c429b9aa2a3c82801006\",\n        \"quantity\": {\n          \"quantity\": 1.5,\n          \"unit\": \"u\"\n        },\n        \"subvalue\": 393.75\n      }\n    ],\n    \"autoupdateCost\": true,\n    \"enabled\": true,\n    \"validated\": true,\n    \"usesUpdatedSinceLastChange\": true,\n    \"createdAt\": \"2019-11-11T22:15:59.301Z\",\n    \"updatedAt\": \"2019-11-15T01:35:19.024Z\"\n }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./infrastructure/composition/rest.go",
    "groupTitle": "Composition"
  },
  {
    "type": "post",
    "url": "/v1/composition",
    "title": "",
    "name": "Post",
    "group": "Composition",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "name",
            "description": "<p>Name</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "cost",
            "defaultValue": "0",
            "description": "<p>Initial cost</p>"
          },
          {
            "group": "Parameter",
            "type": "Quantity",
            "optional": false,
            "field": "unit",
            "description": "<p>Composition base unit</p>"
          },
          {
            "group": "Parameter",
            "type": "Quantity",
            "optional": true,
            "field": "stock",
            "description": "<p>Stock quantity. Same units as &quot;unit&quot;.</p>"
          },
          {
            "group": "Parameter",
            "type": "[]Dependency",
            "optional": true,
            "field": "dependencies",
            "description": "<p>Dependencies: foreign key &quot;of&quot;, and &quot;quantity&quot;.</p>"
          },
          {
            "group": "Parameter",
            "type": "Boolean",
            "optional": true,
            "field": "autoupdateCost",
            "defaultValue": "true",
            "description": "<p>Auto update cost based on dependencies.</p>"
          }
        ]
      }
    },
    "description": "<p>Creates a new Composition. &quot;id&quot; is optional but it can be specified. In case &quot;id&quot; was not specified, a new ObjectID would be assigned. The only required field in body is &quot;unit&quot;. &quot;cost&quot; can be set to any value greater than 0 (default: 0). If &quot;dependencies&quot; are added, &quot;cost&quot; will be calculated automatically based on these, unless &quot;autoupdateCost&quot; is set to &quot;false&quot;.</p>",
    "examples": [
      {
        "title": "Body",
        "content": "{\n  \"composition\": {\n    \"id\": \"9dc9c429b9aa2a3c82801007\",\n    \"name\": \"Comp 7\",\n    \"cost\": 0,\n    \"unit\": {\n      \"quantity\": 3,\n      \"unit\": \"u\"\n    },\n    \"stock\": {\n      \"quantity\": 1,\n      \"unit\": \"u\"\n    },\n    \"dependencies\": [\n      {\n        \"of\": \"9dc9c429b9aa2a3c82801005\",\n        \"quantity\": {\n          \"quantity\": 2,\n          \"unit\": \"u\"\n        },\n        \"subvalue\": 82\n      },\n      {\n        \"of\": \"9dc9c429b9aa2a3c82801006\",\n        \"quantity\": {\n          \"quantity\": 1.5,\n          \"unit\": \"u\"\n        },\n        \"subvalue\": 393.75\n      }\n    ],\n    \"autoupdateCost\": true,\n  }\n}",
        "type": "json"
      }
    ],
    "success": {
      "examples": [
        {
          "title": "Response",
          "content": "HTTP/1.1 200 OK\n\n{\n  \"composition\": {\n    \"id\": \"9dc9c429b9aa2a3c82801007\",\n    \"name\": \"Comp 7\",\n    \"cost\": 475.75,\n    \"unit\": {\n      \"quantity\": 3,\n      \"unit\": \"u\"\n    },\n    \"stock\": {\n      \"quantity\": 1,\n      \"unit\": \"u\"\n    },\n    \"dependencies\": [\n      {\n        \"of\": \"9dc9c429b9aa2a3c82801005\",\n        \"quantity\": {\n          \"quantity\": 2,\n          \"unit\": \"u\"\n        },\n        \"subvalue\": 82\n      },\n      {\n        \"of\": \"9dc9c429b9aa2a3c82801006\",\n        \"quantity\": {\n          \"quantity\": 1.5,\n          \"unit\": \"u\"\n        },\n        \"subvalue\": 393.75\n      }\n    ],\n    \"autoupdateCost\": true,\n    \"enabled\": true,\n    \"validated\": false,\n    \"usesUpdatedSinceLastChange\": true,\n    \"createdAt\": \"2019-11-11T22:15:59.301Z\",\n    \"updatedAt\": \"2019-11-15T01:35:19.024Z\"\n  },\n\t\"status\": \"CREATED\"\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./infrastructure/composition/rest.go",
    "groupTitle": "Composition"
  },
  {
    "type": "put",
    "url": "/v1/composition/:compositionId",
    "title": "",
    "name": "Put",
    "group": "Composition",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "compositionId",
            "description": "<p>Composition ID</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "name",
            "description": "<p>Name</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "cost",
            "description": "<p>Initial cost</p>"
          },
          {
            "group": "Parameter",
            "type": "Quantity",
            "optional": true,
            "field": "unit",
            "description": "<p>Composition base unit. Cannot be changed.</p>"
          },
          {
            "group": "Parameter",
            "type": "Quantity",
            "optional": true,
            "field": "stock",
            "description": "<p>Stock quantity. Same units as &quot;unit&quot;.</p>"
          },
          {
            "group": "Parameter",
            "type": "[]Dependency",
            "optional": true,
            "field": "dependencies",
            "description": "<p>Dependencies: foreign key &quot;of&quot;, and &quot;quantity&quot;.</p>"
          },
          {
            "group": "Parameter",
            "type": "Boolean",
            "optional": true,
            "field": "autoupdateCost",
            "description": "<p>Auto update cost based on dependencies.</p>"
          }
        ]
      }
    },
    "description": "<p>Updates an existing Composition based on its ID. &quot;cost&quot; can be set to any value greater than 0 (default: 0). If &quot;dependencies&quot; are added, &quot;cost&quot; will be calculated automatically based on these, unless &quot;autoupdateCost&quot; is set to &quot;false&quot;. All fields are optional. &quot;unit&quot; unit cannot be changed of type.</p>",
    "examples": [
      {
        "title": "Body",
        "content": "{\n  \"composition\": {\n    \"name\": \"Comp 8\",\n    \"cost\": 125,\n    \"unit\": {\n      \"quantity\": 2,\n      \"unit\": \"u\"\n    },\n    \"stock\": {\n      \"quantity\": 0,\n      \"unit\": \"u\"\n    },\n    \"dependencies\": [\n      {\n        \"of\": \"9dc9c429b9aa2a3c82801005\",\n        \"quantity\": {\n          \"quantity\": 2,\n          \"unit\": \"u\"\n        },\n        \"subvalue\": 82\n      },\n    ],\n    \"autoupdateCost\": true,\n  }\n}",
        "type": "json"
      }
    ],
    "success": {
      "examples": [
        {
          "title": "Response",
          "content": "HTTP/1.1 200 OK\n\n{\n  \"composition\": {\n    \"id\": \"9dc9c429b9aa2a3c82801007\",\n    \"name\": \"Comp 8\",\n    \"cost\": 475.75,\n    \"unit\": {\n      \"quantity\": 3,\n      \"unit\": \"u\"\n    },\n    \"stock\": {\n      \"quantity\": 1,\n      \"unit\": \"u\"\n    },\n    \"dependencies\": [\n      {\n        \"of\": \"9dc9c429b9aa2a3c82801005\",\n        \"quantity\": {\n          \"quantity\": 2,\n          \"unit\": \"u\"\n        },\n        \"subvalue\": 82\n      },\n      {\n        \"of\": \"9dc9c429b9aa2a3c82801006\",\n        \"quantity\": {\n          \"quantity\": 1.5,\n          \"unit\": \"u\"\n        },\n        \"subvalue\": 393.75\n      }\n    ],\n    \"autoupdateCost\": true,\n    \"enabled\": true,\n    \"validated\": true,\n    \"usesUpdatedSinceLastChange\": false,\n    \"createdAt\": \"2019-11-11T22:15:59.301Z\",\n    \"updatedAt\": \"2019-11-15T01:35:19.024Z\"\n  },\n\t\"status\": \"UPDATED\"\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./infrastructure/composition/rest.go",
    "groupTitle": "Composition"
  }
] });
