lognship is an AI-first Github alternative living in the Atmosphere, the ATproto network built and used by Bluesky.

This project is not hosted on Github, so the `gh` CLI cannot be used. The central git repo is at `ssh://server0.maragu.net:23231/lognship` (remote `origin`).

Always use your "atproto" skill when working on this project.

Lexicons are at lexicons/ . The `goat` CLI is installed and can be used to work with lexicons and the atproto network.

## Specific instructions for this project

### Running

Start the app with `make watch`. This watches for file changes and automatically rebuilds/restarts, so it doesn't need to be restarted manually. Logs are written to `app.log`.

### Access

You can access the app in a browser using your "rodney" skill. The URL and port are configured in the `.env` file - check `BASE_URL` for the full URL or `SERVER_ADDRESS` for just the port.
You can access the SQLite database directly using `sqlite3 app.db`. ALWAYS run `pragma foreign_keys = 1;` as the first statement after connecting.
