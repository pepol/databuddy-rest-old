// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate_swagger = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Peter Polacik"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/kv": {
            "get": {
                "description": "List all keys.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "kv"
                ],
                "summary": "List all keys",
                "parameters": [
                    {
                        "type": "string",
                        "default": "",
                        "description": "Key prefix",
                        "name": "prefix",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "default",
                        "description": "Namespace",
                        "name": "ns",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/kv/{key}": {
            "get": {
                "description": "Get key value.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "kv"
                ],
                "summary": "Get key",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "default",
                        "description": "Namespace",
                        "name": "ns",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "default": false,
                        "description": "Return only value",
                        "name": "raw",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.KVItem"
                        }
                    },
                    "400": {
                        "description": "Returned when 'raw' parameter is not parseable as boolean",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.RequestError"
                        }
                    },
                    "404": {
                        "description": "Returned when either key or namespace doesn't exist",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.RequestError"
                        }
                    }
                }
            },
            "put": {
                "description": "Store the provided value under key.\nIf namespace doesn't exist, it gets created.",
                "consumes": [
                    "text/plain",
                    "application/octet-stream"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "kv"
                ],
                "summary": "Put key",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "default",
                        "description": "Namespace",
                        "name": "ns",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "User-defined metadata",
                        "name": "flags",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Time-To-Live (in seconds), 0 means the item won't expire",
                        "name": "ttl",
                        "in": "query"
                    },
                    {
                        "description": "Value to store",
                        "name": "value",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "boolean"
                        }
                    },
                    "400": {
                        "description": "Returned when no value is provided",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.RequestError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete provided key.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "kv"
                ],
                "summary": "Delete key",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "default",
                        "description": "Namespace",
                        "name": "ns",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.KVItem"
                        }
                    },
                    "404": {
                        "description": "Returned when either key or namespace doesn't exist",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.RequestError"
                        }
                    }
                }
            }
        },
        "/namespace": {
            "get": {
                "description": "List all namespaces.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "namespace"
                ],
                "summary": "List all namespaces",
                "parameters": [
                    {
                        "type": "string",
                        "default": "",
                        "description": "Namespace name prefix",
                        "name": "prefix",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/namespace/{name}": {
            "get": {
                "description": "Get namespace by name.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "namespace"
                ],
                "summary": "Get namespace",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Namespace name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.NamespaceStatus"
                        }
                    },
                    "404": {
                        "description": "Returned when namespace doesn't exist",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.RequestError"
                        }
                    }
                }
            },
            "put": {
                "description": "Create the namespace with given name and spec.\nUpdate fields of given namespace based on body if it already\nexists.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "namespace"
                ],
                "summary": "Create/update namespace",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Namespace name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Namespace fields to update",
                        "name": "spec",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.NamespaceSpec"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.NamespaceStatus"
                        }
                    },
                    "400": {
                        "description": "Returned when 'spec' doesn't conform to model",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.RequestError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Mark given namespace as deleted.\nAll the objects stored within the namespace are scheduled for\ndeletion asynchronously. While the namespace is in the process\nof being deleted, GET-ing it will return the object with status\nattribute \"DeleteIndex\" set to index of the delete operation.\nOnce all the contents of the namespace are deleted, GET on\nthe namespace will return HTTP 404.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "namespace"
                ],
                "summary": "Delete namespace",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Namespace name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.NamespaceStatus"
                        }
                    },
                    "404": {
                        "description": "Returned when namespace doesn't exist",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.RequestError"
                        }
                    }
                }
            }
        },
        "/status/cluster/leader": {
            "get": {
                "description": "Get cluster's Raft leader information.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Get Raft leader",
                "parameters": [
                    {
                        "type": "string",
                        "default": "local",
                        "description": "Datacenter to query, defaults to local datacenter",
                        "name": "dc",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Requested datacenter doesn't exist",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.RequestError"
                        }
                    }
                }
            }
        },
        "/status/cluster/peers": {
            "get": {
                "description": "Get the peers for given cluster.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Get cluster peers",
                "parameters": [
                    {
                        "type": "string",
                        "default": "local",
                        "description": "Datacenter to query, defaults to local datacenter",
                        "name": "dc",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Requested datacenter doesn't exist",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.RequestError"
                        }
                    }
                }
            }
        },
        "/status/node": {
            "get": {
                "description": "Get current node information and status.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Get node status",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1alpha3.NodeInfo"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "v1alpha3.KVItem": {
            "type": "object",
            "properties": {
                "createIndex": {
                    "type": "string"
                },
                "expiresAt": {
                    "type": "integer"
                },
                "flags": {
                    "type": "integer"
                },
                "key": {
                    "type": "string"
                },
                "updateIndex": {
                    "type": "string"
                },
                "value": {
                    "description": "Value is encoded using base64.",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "v1alpha3.NamespaceSpec": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "labels": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "v1alpha3.NamespaceStatus": {
            "type": "object",
            "properties": {
                "createIndex": {
                    "type": "string"
                },
                "deleteIndex": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "labels": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "updateIndex": {
                    "type": "string"
                }
            }
        },
        "v1alpha3.NodeCPUInfo": {
            "type": "object",
            "properties": {
                "cacheSize": {
                    "type": "integer"
                },
                "cores": {
                    "type": "integer"
                },
                "cpuindex": {
                    "type": "integer"
                },
                "family": {
                    "type": "string"
                },
                "flags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "mhz": {
                    "type": "number"
                },
                "microcode": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "modelName": {
                    "type": "string"
                },
                "vendorID": {
                    "type": "string"
                }
            }
        },
        "v1alpha3.NodeDiskInfo": {
            "type": "object",
            "properties": {
                "free": {
                    "type": "integer"
                },
                "fstype": {
                    "type": "string"
                },
                "inodesFree": {
                    "type": "integer"
                },
                "inodesTotal": {
                    "type": "integer"
                },
                "inodesUsed": {
                    "type": "integer"
                },
                "inodesUsedPercent": {
                    "type": "number"
                },
                "path": {
                    "type": "string"
                },
                "total": {
                    "type": "integer"
                },
                "used": {
                    "type": "integer"
                },
                "usedPercent": {
                    "type": "number"
                }
            }
        },
        "v1alpha3.NodeInfo": {
            "type": "object",
            "properties": {
                "annotations": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "cluster": {
                    "type": "string"
                },
                "cpu": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1alpha3.NodeCPUInfo"
                    }
                },
                "disk": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1alpha3.NodeDiskInfo"
                    }
                },
                "hostname": {
                    "type": "string"
                },
                "labels": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "memory": {
                    "$ref": "#/definitions/v1alpha3.NodeMemoryInfo"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "v1alpha3.NodeMemoryInfo": {
            "type": "object",
            "properties": {
                "available": {
                    "type": "integer"
                },
                "free": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                },
                "used": {
                    "type": "integer"
                },
                "usedPercent": {
                    "type": "number"
                }
            }
        },
        "v1alpha3.RequestError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo_swagger holds exported Swagger Info so clients can modify it
var SwaggerInfo_swagger = &swag.Spec{
	Version:          "1.0.0-alpha3",
	Host:             "localhost:8080",
	BasePath:         "/v1alpha3",
	Schemes:          []string{},
	Title:            "DataBuddy",
	Description:      "API to use DataBuddy data storage system",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate_swagger,
}

func init() {
	swag.Register(SwaggerInfo_swagger.InstanceName(), SwaggerInfo_swagger)
}
