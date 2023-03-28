// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2023-03-28 09:07:04.575621 -0600 MDT m=+34.147678501
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
                            "$ref": "#/definitions/github.com_DIMO-Network_rewards-api_internal_controllers.UserResponse"
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
                            "$ref": "#/definitions/github.com_DIMO-Network_rewards-api_internal_controllers.HistoryResponse"
                        }
                    }
                }
            }
        },
        "/user/history/balance": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "A summary of the user's DIMO balance across all chains.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github.com_DIMO-Network_rewards-api_internal_controllers.BalanceHistory"
                        }
                    }
                }
            }
        },
        "/user/history/transactions": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "A summary of the user's DIMO transaction history, all time.",
                "parameters": [
                    {
                        "enum": [
                            "Baseline",
                            "Referrals",
                            "Marketplace",
                            "Other"
                        ],
                        "type": "string",
                        "description": "A label for a transaction type.",
                        "name": "type",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_controllers.TransactionHistory"
                        }
                    }
                }
            }
        },
        "/user/referrals": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "A summary of the user's referrals.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github.com_DIMO-Network_rewards-api_internal_controllers.UserResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github.com_DIMO-Network_rewards-api_internal_controllers.APITransaction": {
            "type": "object",
            "properties": {
                "chainId": {
                    "description": "ChainID is the chain id of the chain on which the transaction took place. Important\nvalues are 137 for Polygon, 1 for Ethereum.",
                    "type": "integer",
                    "example": 137
                },
                "description": {
                    "description": "Description is a short elaboration of the Type or a generic, e.g., \"Incoming transfer\" message.",
                    "type": "string"
                },
                "from": {
                    "description": "From is the address of the source of the value, in 0x-prefixed hex.",
                    "type": "string",
                    "example": "0xf316832fbfe49f90df09eee019c2ece87fad3fac"
                },
                "time": {
                    "description": "Time is the timestamp of the block in which the transaction took place, in RFC-3999 format.",
                    "type": "string",
                    "example": "2023-01-22T09:00:12Z"
                },
                "to": {
                    "description": "To is the address of the recipient of the value, in 0x-prefixed hex.",
                    "type": "string",
                    "example": "0xc66d80f5063677425270013136ef9fa2bf1f9f1a"
                },
                "type": {
                    "description": "Type is a transaction type.",
                    "type": "string",
                    "enum": [
                        "Baseline",
                        "Referrals",
                        "Marketplace"
                    ]
                },
                "value": {
                    "description": "Value is the amount of token being transferred. Divide by 10^18 to get what people\nnormally consider $DIMO.",
                    "type": "number",
                    "example": 10000000000000000
                }
            }
        },
        "github.com_DIMO-Network_rewards-api_internal_controllers.Balance": {
            "type": "object",
            "properties": {
                "balance": {
                    "description": "Balance is the total amount of $DIMO held at this time, across all chains.",
                    "type": "number",
                    "example": 237277217092548850000
                },
                "time": {
                    "description": "Time is the block timestamp of this balance update.",
                    "type": "string",
                    "example": "2023-03-06T09:11:00Z"
                }
            }
        },
        "github.com_DIMO-Network_rewards-api_internal_controllers.BalanceHistory": {
            "type": "object",
            "properties": {
                "balanceHistory": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github.com_DIMO-Network_rewards-api_internal_controllers.Balance"
                    }
                }
            }
        },
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
        "github.com_DIMO-Network_rewards-api_internal_controllers.TransactionHistory": {
            "type": "object",
            "properties": {
                "transactions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github.com_DIMO-Network_rewards-api_internal_controllers.APITransaction"
                    }
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
                },
                "walletBalance": {
                    "description": "WalletBalance is the number of tokens held in the users's wallet, if he has a wallet\nattached to the present account.",
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
                "minted": {
                    "description": "Minted is true if the device has been minted on-chain.",
                    "type": "boolean"
                },
                "optedIn": {
                    "description": "OptedIn is true if the user has agreed to the terms of service.",
                    "type": "boolean"
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
                "dataThisWeek": {
                    "type": "boolean"
                },
                "id": {
                    "description": "ID is the integration ID.",
                    "type": "string",
                    "example": "27egBSLazAT7njT2VBjcISPIpiU"
                },
                "onChainPairingStatus": {
                    "description": "OnChainPairingStatus is the on-chain pairing status of the integration.",
                    "type": "string",
                    "enum": [
                        "Paired",
                        "Unpaired",
                        "NotApplicable"
                    ],
                    "example": "Paired"
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
        "internal_controllers.APITransaction": {
            "type": "object",
            "properties": {
                "chainId": {
                    "description": "ChainID is the chain id of the chain on which the transaction took place. Important\nvalues are 137 for Polygon, 1 for Ethereum.",
                    "type": "integer",
                    "example": 137
                },
                "description": {
                    "description": "Description is a short elaboration of the Type or a generic, e.g., \"Incoming transfer\" message.",
                    "type": "string"
                },
                "from": {
                    "description": "From is the address of the source of the value, in 0x-prefixed hex.",
                    "type": "string",
                    "example": "0xf316832fbfe49f90df09eee019c2ece87fad3fac"
                },
                "time": {
                    "description": "Time is the timestamp of the block in which the transaction took place, in RFC-3999 format.",
                    "type": "string",
                    "example": "2023-01-22T09:00:12Z"
                },
                "to": {
                    "description": "To is the address of the recipient of the value, in 0x-prefixed hex.",
                    "type": "string",
                    "example": "0xc66d80f5063677425270013136ef9fa2bf1f9f1a"
                },
                "type": {
                    "description": "Type is a transaction type.",
                    "type": "string",
                    "enum": [
                        "Baseline",
                        "Referrals",
                        "Marketplace"
                    ]
                },
                "value": {
                    "description": "Value is the amount of token being transferred. Divide by 10^18 to get what people\nnormally consider $DIMO.",
                    "type": "number",
                    "example": 10000000000000000
                }
            }
        },
        "internal_controllers.Balance": {
            "type": "object",
            "properties": {
                "balance": {
                    "description": "Balance is the total amount of $DIMO held at this time, across all chains.",
                    "type": "number",
                    "example": 237277217092548850000
                },
                "time": {
                    "description": "Time is the block timestamp of this balance update.",
                    "type": "string",
                    "example": "2023-03-06T09:11:00Z"
                }
            }
        },
        "internal_controllers.BalanceHistory": {
            "type": "object",
            "properties": {
                "balanceHistory": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal_controllers.Balance"
                    }
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
        "internal_controllers.TransactionHistory": {
            "type": "object",
            "properties": {
                "transactions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal_controllers.APITransaction"
                    }
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
                },
                "walletBalance": {
                    "description": "WalletBalance is the number of tokens held in the users's wallet, if he has a wallet\nattached to the present account.",
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
                "minted": {
                    "description": "Minted is true if the device has been minted on-chain.",
                    "type": "boolean"
                },
                "optedIn": {
                    "description": "OptedIn is true if the user has agreed to the terms of service.",
                    "type": "boolean"
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
                "dataThisWeek": {
                    "type": "boolean"
                },
                "id": {
                    "description": "ID is the integration ID.",
                    "type": "string",
                    "example": "27egBSLazAT7njT2VBjcISPIpiU"
                },
                "onChainPairingStatus": {
                    "description": "OnChainPairingStatus is the on-chain pairing status of the integration.",
                    "type": "string",
                    "enum": [
                        "Paired",
                        "Unpaired",
                        "NotApplicable"
                    ],
                    "example": "Paired"
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
