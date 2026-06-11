# chirpy

Boot.dev's guided HTTP server project

Chirpy is a Go backend that supports 
1. authentication with JWT and refresh tokens
2. creating and filtering "chirps" or posts 
3. a proof of concept webhook for future payment integration

## Development

Go and PostgreSQL are required. Additional dependencies are managed by [Mise](https://github.com/jdx/mise).

### Database

```bash
sudo -u postgres psql
# create database chirpy;
```

### Environment Variables

A sample file is provided in `.env.example`

Note that `POLKA_API_KEY` is not a real key and is simply a mock service for Boot.dev's testing. You can review its actual usage in `internal/routes/webhooks`

### Running Chirpy 

```bash
mise install # installs Air for hot reloading, as well as other tools like Goose
mise dev     # start the development server with Air
```

