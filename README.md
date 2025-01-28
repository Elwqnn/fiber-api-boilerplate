# Fiber API Boilerplate [![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL-red.svg)](https://www.gnu.org/licenses/agpl-3.0.en.html)

A production-ready RESTful API boilerplate using [Fiber](https://gofiber.io/), a fast and lightweight web framework for Go inspired by Express.js.

## Implemented Features

✅ Clean Architecture

- Organized into layers: handlers, services, repositories, and models
- Clear separation of concerns and dependencies
- Modular and maintainable code structure

✅ Authentication & Authorization

- JWT-based authentication
- Session management with Redis
- Role-based access control (user/admin)
- OAuth2 integration (Google, GitHub, Discord)
- Multiple authentication methods (credentials/OAuth)

✅ Database Integration

- PostgreSQL with GORM ORM
- Auto migrations
- UUID primary keys
- Soft delete support
- Relationship handling

✅ Security

- Password hashing with bcrypt
- JWT token validation
- Session management
- CORS configuration
- HTTP-only cookies

✅ Request Validation

- Request payload validation using validator/v10
- Custom validation middleware
- Standardized error responses

✅ Error Handling

- Custom error middleware
- Standardized error responses
- Logging for internal server errors

✅ Configuration Management

- Environment-based configuration
- Easy environment variable management
- Secure secrets handling

✅ API Features

- RESTful endpoints
- Standardized response format
- User management
- Profile updates
- Session management

## Getting Started

1. Clone this repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Run the Docker compose:

```sh
docker-compose up -d
```

## API Documentation

Coming soon...

## Upcoming features

- Unit tests coverage
- Rate limiting
- File upload handling
- Swagger documentation

## Contributing

This project is under active development. Contributions are welcome!

## License

This project is under the AGPL-3.0 license. The AGPL-3.0 is a strong copyleft license that requires any modified or derived work to be distributed under the same license terms. Failure to comply with this license's terms may result in:

- Legal action for copyright infringement
- Requirement to release proprietary source code
- Statutory damages
- Termination of rights to use the software

For full license terms, see [LICENSE](LICENSE) file or visit [AGPL-3.0](https://www.gnu.org/licenses/agpl-3.0.en.html).
