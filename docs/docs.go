// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Johnson Lee"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/blogs": {
            "post": {
                "description": "Create a new blog",
                "tags": [
                    "blogs"
                ],
                "summary": "Create Blog",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Blog Title",
                        "name": "title",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Blog Content",
                        "name": "content",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Blog Image",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.createBlogResponse"
                        }
                    }
                }
            }
        },
        "/invite": {
            "post": {
                "description": "Invite new user to create a pair",
                "tags": [
                    "invite"
                ],
                "summary": "Invite",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Invite User",
                        "name": "invite_info",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.InviteRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.InviteResponse"
                        }
                    }
                }
            }
        },
        "/users/invitee_signup": {
            "post": {
                "description": "for invitee to sign up",
                "tags": [
                    "users"
                ],
                "summary": "Invitee SignUp",
                "parameters": [
                    {
                        "description": "Create User",
                        "name": "signup_info",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.invitedUserSignUpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.invitedUserSignUpResponse"
                        }
                    }
                }
            }
        },
        "/users/login": {
            "post": {
                "description": "Login to an user account",
                "tags": [
                    "users"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "Login User",
                        "name": "login_info",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.loginUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.loginUserResponse"
                        }
                    }
                }
            }
        },
        "/users/signup": {
            "post": {
                "description": "Create a new user account",
                "tags": [
                    "users"
                ],
                "summary": "SignUp",
                "parameters": [
                    {
                        "description": "Create User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.createUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.userResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/verify/verify_email": {
            "get": {
                "description": "Verify the email of created account.",
                "tags": [
                    "verify"
                ],
                "summary": "Verify email",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "email_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "secret_code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.VerifyEmailResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.InviteRequest": {
            "type": "object",
            "properties": {
                "invitee_email": {
                    "type": "string"
                },
                "inviter_id": {
                    "type": "string"
                }
            }
        },
        "api.InviteResponse": {
            "type": "object",
            "properties": {
                "create_time": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "invitation_token": {
                    "type": "string"
                },
                "invitee_email": {
                    "type": "string"
                },
                "inviter_id": {
                    "type": "string"
                },
                "is_accepted": {
                    "type": "boolean"
                }
            }
        },
        "api.VerifyEmailResponse": {
            "type": "object",
            "properties": {
                "is_email_verified": {
                    "type": "boolean"
                }
            }
        },
        "api.createBlogResponse": {
            "type": "object",
            "properties": {
                "blog_id": {
                    "type": "string"
                },
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "pair_id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "api.createUserRequest": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 2
                },
                "password": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 6
                }
            }
        },
        "api.invitedUserSignUpRequest": {
            "type": "object",
            "properties": {
                "invitation_id": {
                    "type": "integer"
                },
                "invitation_token": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "api.invitedUserSignUpResponse": {
            "type": "object",
            "properties": {
                "create_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "api.loginUserRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 6
                }
            }
        },
        "api.loginUserResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "access_token_expires_at": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "refresh_token_expires_at": {
                    "type": "string"
                },
                "session_id": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/api.userResponse"
                }
            }
        },
        "api.userResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password_changed_at": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Couple Website API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
