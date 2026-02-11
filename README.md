# Synopsis

This project is about tracking the user’s energy levels.

There are 3 types of energy :

1. Physical : is the user physically able to accomplish their tasks ? Are they physically tired ?
2. Mental : is the user mentally capable to focus on their tasks?
3. Emotional : is the user willing to engage in their tasks?

Everyday and for each type of energy, the user will give a grade from 0 to 10, 0 being the lowest level and 10 the highest.

The app shows :

- how the energy levels vary
- a list of marking events every day and which type of energy and how it has been affected
- a summary of how the user spent their time week after week

## Run Front + Back In One Terminal

1. Create backend env file:
`cp go/.env.example go/.env`
2. Add Firebase service-account credentials (one of):
- Create `go/firebase-credentials.json` (not committed) and keep `FIREBASE_CREDENTIALS_FILE=/app/firebase-credentials.json` in `go/.env`.
- Or set `FIREBASE_CREDENTIALS` in `go/.env` to a base64-encoded service-account JSON.
3. Fill required values in `go/.env` (`GCP_PROJECT_ID`, `FIREBASE_API_KEY`).
4. Start all services:
`docker compose up --build`

Services:
- Frontend: `http://localhost:8080`
- Backend: `http://localhost:8888`
- Mailpit UI: `http://localhost:8025`
