// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/coupon": {
            "post": {
                "description": "Метод для вывода списка купонов",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Coupon"
                ],
                "parameters": [
                    {
                        "description": "Request Body, заполнять обязательно",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APICouponList"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APICoupon"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/coupon/:id": {
            "get": {
                "description": "Метод для поиска купона по ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Coupon"
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APICoupon"
                        }
                    }
                }
            },
            "delete": {
                "description": "Метод для удаления купона по ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Coupon"
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APICoupon"
                        }
                    }
                }
            }
        },
        "/api/v1/coupon/apply": {
            "post": {
                "description": "Метод для применения купона к пользователю",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Coupon"
                ],
                "parameters": [
                    {
                        "description": "Request Body, заполнять обязательно",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APICouponApplier"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/coupon/apply/:id": {
            "get": {
                "description": "Метод для поиска купона по UserID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Coupon"
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APICoupon"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/coupon/create": {
            "post": {
                "description": "Метод для создания купона",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Coupon"
                ],
                "parameters": [
                    {
                        "description": "Request Body, заполнять обязательно",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.NewAPICoupon"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APICoupon"
                        }
                    }
                }
            }
        },
        "/api/v1/user": {
            "post": {
                "description": "Метод для вывода списка пользователей",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "parameters": [
                    {
                        "description": "Request Body, заполнять обязательно",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APIUserList"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APIUser"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/user/:id": {
            "get": {
                "description": "Метод для поиска пользователя по ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APIUser"
                        }
                    }
                }
            },
            "delete": {
                "description": "Метод для удаления пользователя по ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APIUser"
                        }
                    }
                }
            }
        },
        "/api/v1/user/create": {
            "post": {
                "description": "Метод для создания пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "parameters": [
                    {
                        "description": "Request Body, заполнять обязательно",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.NewAPIUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/github_com_xamust_couponApp_internal_adapter_api_v1_models.APIUser"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_xamust_couponApp_internal_adapter_api_v1_models.APICoupon": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "maxRedemptions": {
                    "type": "integer"
                },
                "metadata": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "redeemBy": {
                    "type": "string"
                },
                "reward": {
                    "type": "string"
                },
                "timesRedeemed": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "github_com_xamust_couponApp_internal_adapter_api_v1_models.APICouponApplier": {
            "type": "object",
            "required": [
                "coupon_id",
                "user_id"
            ],
            "properties": {
                "coupon_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "github_com_xamust_couponApp_internal_adapter_api_v1_models.APICouponList": {
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                }
            }
        },
        "github_com_xamust_couponApp_internal_adapter_api_v1_models.APIUser": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_active": {
                    "type": "boolean"
                },
                "metadata": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "github_com_xamust_couponApp_internal_adapter_api_v1_models.APIUserList": {
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                }
            }
        },
        "github_com_xamust_couponApp_internal_adapter_api_v1_models.NewAPICoupon": {
            "type": "object",
            "properties": {
                "maxRedemptions": {
                    "type": "integer"
                },
                "metadata": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "redeemBy": {
                    "type": "string"
                },
                "reward": {
                    "type": "string"
                }
            }
        },
        "github_com_xamust_couponApp_internal_adapter_api_v1_models.NewAPIUser": {
            "type": "object",
            "properties": {
                "metadata": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
                    "example": {
                        "{\"key\"": "\"value\"}"
                    }
                },
                "name": {
                    "type": "string",
                    "example": "John Doe"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}