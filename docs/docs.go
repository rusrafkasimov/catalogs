// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Ruslan Kasimov"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/catalog": {
            "get": {
                "security": [
                    {
                        "TokenJWT": []
                    }
                ],
                "description": "Get JSON CatalogsRequest, return JSON GetCatalogsResponse",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Catalog"
                ],
                "summary": "Get catalogs",
                "parameters": [
                    {
                        "description": "Catalogs",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CatalogsRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetCatalogsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "TokenJWT": []
                    }
                ],
                "description": "Get JSON CatalogRequest, return JSON UpdateCatalogResponse",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Catalog"
                ],
                "summary": "Update catalog",
                "parameters": [
                    {
                        "description": "Catalog",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CatalogRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/UpdateCatalogResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "TokenJWT": []
                    }
                ],
                "description": "Get JSON CatalogRequest, return JSON CreateCatalogResponse",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Catalog"
                ],
                "summary": "Create catalog",
                "parameters": [
                    {
                        "description": "Catalog",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CatalogRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CatalogRequest"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            }
        },
        "/catalog/:id": {
            "get": {
                "security": [
                    {
                        "TokenJWT": []
                    }
                ],
                "description": "Get id in path, return JSON GetCatalogResponse",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Catalog"
                ],
                "summary": "Get catalog by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Catalog ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetCatalogResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "TokenJWT": []
                    }
                ],
                "description": "Get id from path, return JSON DeleteCatalogResponse",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Catalog"
                ],
                "summary": "Delete catalog",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Catalog ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/DeleteCatalogResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            }
        },
        "/categories": {
            "get": {
                "security": [
                    {
                        "TokenJWT": []
                    }
                ],
                "description": "Return JSON GetCategoriesResponse",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Catalog"
                ],
                "summary": "Get categories",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetCategoriesResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Return answer from server for checking what server is stay alive",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "system"
                ],
                "summary": "Health check",
                "operationId": "health-check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Just ping-pong endpoint, can be used as health indicator",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "testing"
                ],
                "summary": "Ping",
                "operationId": "ping",
                "responses": {
                    "200": {
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "CatalogRequest": {
            "type": "object",
            "properties": {
                "active": {
                    "type": "boolean"
                },
                "category": {
                    "type": "string"
                },
                "desc": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "CatalogResponse": {
            "type": "object",
            "properties": {
                "active": {
                    "type": "boolean"
                },
                "category": {
                    "type": "string"
                },
                "desc": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "CatalogsRequest": {
            "type": "object",
            "required": [
                "category"
            ],
            "properties": {
                "category": {
                    "type": "string"
                },
                "query": {
                    "type": "string"
                },
                "sorted": {
                    "type": "boolean",
                    "default": true
                }
            }
        },
        "DeleteCatalogResponse": {
            "type": "object",
            "properties": {
                "meta": {
                    "$ref": "#/definitions/ResponseMeta"
                },
                "payload": {
                    "type": "object",
                    "properties": {
                        "active": {
                            "type": "boolean"
                        }
                    }
                }
            }
        },
        "Error": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "GetCatalogResponse": {
            "type": "object",
            "properties": {
                "meta": {
                    "$ref": "#/definitions/ResponseMeta"
                },
                "payload": {
                    "type": "object",
                    "properties": {
                        "active": {
                            "type": "boolean"
                        },
                        "category": {
                            "type": "string"
                        },
                        "desc": {
                            "type": "string"
                        },
                        "id": {
                            "type": "string"
                        },
                        "name": {
                            "type": "string"
                        },
                        "value": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "GetCatalogsResponse": {
            "type": "object",
            "properties": {
                "meta": {
                    "$ref": "#/definitions/ResponseMetaList"
                },
                "payload": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/CatalogResponse"
                    }
                }
            }
        },
        "GetCategoriesResponse": {
            "type": "object",
            "properties": {
                "meta": {
                    "$ref": "#/definitions/ResponseMetaList"
                },
                "payload": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "ResponseMeta": {
            "type": "object"
        },
        "ResponseMetaList": {
            "type": "object",
            "properties": {
                "current_page": {
                    "type": "integer"
                },
                "node_id": {
                    "type": "string"
                },
                "num_of_pages": {
                    "type": "integer"
                },
                "num_of_results": {
                    "type": "integer"
                },
                "page_size": {
                    "type": "integer"
                }
            }
        },
        "UpdateCatalogResponse": {
            "type": "object",
            "properties": {
                "meta": {
                    "$ref": "#/definitions/ResponseMeta"
                },
                "payload": {
                    "type": "object",
                    "properties": {
                        "active": {
                            "type": "boolean"
                        },
                        "category": {
                            "type": "string"
                        },
                        "desc": {
                            "type": "string"
                        },
                        "id": {
                            "type": "string"
                        },
                        "name": {
                            "type": "string"
                        },
                        "value": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "TokenJWT": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1",
	Host:        "127.0.0.1:8091",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Swagger Catalogs Service",
	Description: "Catalogs Microservice (Golang)",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}