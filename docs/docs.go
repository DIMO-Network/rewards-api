// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2022-11-30 19:11:55.525613 -0500 EST m=+1.065657460
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
        "/user": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "A summary of the user's rewards.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_controllers.UserResponse"
                        }
                    }
                }
            }
        },
        "/user/history": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "A summary of the user's rewards for past weeks.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_controllers.HistoryResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github.com_DIMO-Network_rewards-api_internal_controllers.HistoryResponse": {
            "type": "object",
            "properties": {
                "weeks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github.com_DIMO-Network_rewards-api_internal_controllers.HistoryResponseWeek"
                    }
                }
            }
        },
        "github.com_DIMO-Network_rewards-api_internal_controllers.HistoryResponseWeek": {
            "type": "object",
            "properties": {
                "end": {
                    "description": "End is the starting time of the issuance week after this one.",
                    "type": "string",
                    "example": "2022-04-18T05:00:00Z"
                },
                "points": {
                    "description": "Points is the number of points the user earned this week.",
                    "type": "integer",
                    "example": 4000
                },
                "start": {
                    "description": "Start is the starting time of the issuance week.",
                    "type": "string",
                    "example": "2022-04-11T05:00:00Z"
                },
                "tokens": {
                    "description": "Tokens is the number of tokens the user earned this week.",
                    "type": "number",
                    "example": 4000
                }
            }
        },
        "github.com_DIMO-Network_rewards-api_internal_controllers.UserResponse": {
            "type": "object",
            "properties": {
                "devices": {
                    "description": "Devices is a list of the user's devices, together with some information about their\nconnectivity.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github.com_DIMO-Network_rewards-api_internal_controllers.UserResponseDevice"
                    }
                },
                "points": {
                    "description": "Points is the user's total number of points, across all devices and issuance weeks.",
                    "type": "integer",
                    "example": 5000
                },
                "thisWeek": {
                    "description": "ThisWeek describes the current issuance week.",
                    "$ref": "#/definitions/github.com_DIMO-Network_rewards-api_internal_controllers.UserResponseThisWeek"
                },
                "tokens": {
                    "description": "Tokens is the number of tokens the user has earned, across all devices and issuance\nweeks.",
                    "type": "number",
                    "example": 1.105e+24
                }
            }
        },
        "github.com_DIMO-Network_rewards-api_internal_controllers.UserResponseDevice": {
            "type": "object",
            "properties": {
                "connectedThisWeek": {
                    "description": "ConnectedThisWeek is true if we've seen activity from the device during the current issuance\nweek.",
                    "type": "boolean",
                    "example": true
                },
                "connectionStreak": {
                    "description": "ConnectionStreak is what we consider the streak of the device to be. This may not literally\nbe the number of consecutive connected weeks, because the user may disconnect for a week\nwithout penalty, or have the connection streak reduced after three weeks of inactivity.",
                    "type": "integer",
                    "example": 4
                },
                "disconnectionStreak": {
                    "description": "DisconnectionStreak is the number of consecutive issuance weeks that the device has been\ndisconnected. This number resets to 0 as soon as a device earns rewards for a certain week.",
                    "type": "integer",
                    "example": 0
                },
                "id": {
                    "description": "ID is the user device ID used across all services.",
                    "type": "string",
                    "example": "27cv7gVTh9h4RJuTsmJHpBcr4I9"
                },
                "integrationsThisWeek": {
                    "description": "IntegrationsThisWeek details the integrations we've seen active this week.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github.com_DIMO-Network_rewards-api_internal_controllers.UserResponseIntegration"
                    }
                },
                "lastActive": {
                    "description": "LastActive is the last time we saw activity from the vehicle.",
                    "type": "string",
                    "example": "2022-04-12T09:23:01Z"
                },
                "level": {
                    "description": "Level is the level 1-4 of the device. This is fully determined by ConnectionStreak.",
                    "$ref": "#/definitions/github.com_DIMO-Network_rewards-api_internal_controllers.UserResponseLevel"
                },
                "points": {
                    "description": "Points is the total number of points that the device has earned across all weeks.",
                    "type": "integer",
                    "example": 5000
                },
                "tokens": {
                    "description": "Tokens is the total number of tokens that the device has earned across all weeks.",
                    "type": "number",
                    "example": 5000
                }
            }
        },
        "github.com_DIMO-Network_rewards-api_internal_controllers.UserResponseIntegration": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "ID is the integration ID.",
                    "type": "string",
                    "example": "27egBSLazAT7njT2VBjcISPIpiU"
                },
                "points": {
                    "description": "Points is the number of points a user earns for being connected with this integration\nfor a week.",
                    "type": "integer",
                    "example": 1000
                },
                "vendor": {
                    "description": "Vendor is the name of the integration vendor. At present, this uniquely determines the\nintegration.",
                    "type": "string",
                    "example": "SmartCar"
                }
            }
        },
        "github.com_DIMO-Network_rewards-api_internal_controllers.UserResponseLevel": {
            "type": "object",
            "properties": {
                "maxWeeks": {
                    "description": "MaxWeeks is the last streak week at this level. In the next week, we enter the next level.",
                    "type": "integer",
                    "example": 20
                },
                "minWeeks": {
                    "description": "MinWeeks is the minimum streak of weeks needed to enter this level.",
                    "type": "integer",
                    "example": 4
                },
                "number": {
                    "description": "Number is the level number 1-4",
                    "type": "integer",
                    "example": 2
                },
                "streakPoints": {
                    "description": "StreakPoints is the number of points you earn per week at this level.",
                    "type": "integer",
                    "example": 1000
                }
            }
        },
        "github.com_DIMO-Network_rewards-api_internal_controllers.UserResponseThisWeek": {
            "type": "object",
            "properties": {
                "end": {
                    "description": "End is the timestamp of the start of the next issuance week.",
                    "type": "string",
                    "example": "2022-04-18T05:00:00Z"
                },
                "start": {
                    "description": "Start is the timestamp of the start of the issuance week.",
                    "type": "string",
                    "example": "2022-04-11T05:00:00Z"
                }
            }
        },
        "internal_controllers.HistoryResponse": {
            "type": "object",
            "properties": {
                "weeks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal_controllers.HistoryResponseWeek"
                    }
                }
            }
        },
        "internal_controllers.HistoryResponseWeek": {
            "type": "object",
            "properties": {
                "end": {
                    "description": "End is the starting time of the issuance week after this one.",
                    "type": "string",
                    "example": "2022-04-18T05:00:00Z"
                },
                "points": {
                    "description": "Points is the number of points the user earned this week.",
                    "type": "integer",
                    "example": 4000
                },
                "start": {
                    "description": "Start is the starting time of the issuance week.",
                    "type": "string",
                    "example": "2022-04-11T05:00:00Z"
                },
                "tokens": {
                    "description": "Tokens is the number of tokens the user earned this week.",
                    "type": "number",
                    "example": 4000
                }
            }
        },
        "internal_controllers.UserResponse": {
            "type": "object",
            "properties": {
                "devices": {
                    "description": "Devices is a list of the user's devices, together with some information about their\nconnectivity.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal_controllers.UserResponseDevice"
                    }
                },
                "points": {
                    "description": "Points is the user's total number of points, across all devices and issuance weeks.",
                    "type": "integer",
                    "example": 5000
                },
                "thisWeek": {
                    "description": "ThisWeek describes the current issuance week.",
                    "$ref": "#/definitions/internal_controllers.UserResponseThisWeek"
                },
                "tokens": {
                    "description": "Tokens is the number of tokens the user has earned, across all devices and issuance\nweeks.",
                    "type": "number",
                    "example": 1.105e+24
                }
            }
        },
        "internal_controllers.UserResponseDevice": {
            "type": "object",
            "properties": {
                "connectedThisWeek": {
                    "description": "ConnectedThisWeek is true if we've seen activity from the device during the current issuance\nweek.",
                    "type": "boolean",
                    "example": true
                },
                "connectionStreak": {
                    "description": "ConnectionStreak is what we consider the streak of the device to be. This may not literally\nbe the number of consecutive connected weeks, because the user may disconnect for a week\nwithout penalty, or have the connection streak reduced after three weeks of inactivity.",
                    "type": "integer",
                    "example": 4
                },
                "disconnectionStreak": {
                    "description": "DisconnectionStreak is the number of consecutive issuance weeks that the device has been\ndisconnected. This number resets to 0 as soon as a device earns rewards for a certain week.",
                    "type": "integer",
                    "example": 0
                },
                "id": {
                    "description": "ID is the user device ID used across all services.",
                    "type": "string",
                    "example": "27cv7gVTh9h4RJuTsmJHpBcr4I9"
                },
                "integrationsThisWeek": {
                    "description": "IntegrationsThisWeek details the integrations we've seen active this week.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal_controllers.UserResponseIntegration"
                    }
                },
                "lastActive": {
                    "description": "LastActive is the last time we saw activity from the vehicle.",
                    "type": "string",
                    "example": "2022-04-12T09:23:01Z"
                },
                "level": {
                    "description": "Level is the level 1-4 of the device. This is fully determined by ConnectionStreak.",
                    "$ref": "#/definitions/internal_controllers.UserResponseLevel"
                },
                "points": {
                    "description": "Points is the total number of points that the device has earned across all weeks.",
                    "type": "integer",
                    "example": 5000
                },
                "tokens": {
                    "description": "Tokens is the total number of tokens that the device has earned across all weeks.",
                    "type": "number",
                    "example": 5000
                }
            }
        },
        "internal_controllers.UserResponseIntegration": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "ID is the integration ID.",
                    "type": "string",
                    "example": "27egBSLazAT7njT2VBjcISPIpiU"
                },
                "points": {
                    "description": "Points is the number of points a user earns for being connected with this integration\nfor a week.",
                    "type": "integer",
                    "example": 1000
                },
                "vendor": {
                    "description": "Vendor is the name of the integration vendor. At present, this uniquely determines the\nintegration.",
                    "type": "string",
                    "example": "SmartCar"
                }
            }
        },
        "internal_controllers.UserResponseLevel": {
            "type": "object",
            "properties": {
                "maxWeeks": {
                    "description": "MaxWeeks is the last streak week at this level. In the next week, we enter the next level.",
                    "type": "integer",
                    "example": 20
                },
                "minWeeks": {
                    "description": "MinWeeks is the minimum streak of weeks needed to enter this level.",
                    "type": "integer",
                    "example": 4
                },
                "number": {
                    "description": "Number is the level number 1-4",
                    "type": "integer",
                    "example": 2
                },
                "streakPoints": {
                    "description": "StreakPoints is the number of points you earn per week at this level.",
                    "type": "integer",
                    "example": 1000
                }
            }
        },
        "internal_controllers.UserResponseThisWeek": {
            "type": "object",
            "properties": {
                "end": {
                    "description": "End is the timestamp of the start of the next issuance week.",
                    "type": "string",
                    "example": "2022-04-18T05:00:00Z"
                },
                "start": {
                    "description": "Start is the timestamp of the start of the issuance week.",
                    "type": "string",
                    "example": "2022-04-11T05:00:00Z"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "DIMO Rewards API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
