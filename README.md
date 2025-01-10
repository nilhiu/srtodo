# Speedrunning roadmap.sh Todo List API

The following code is a solution to the [roadmap.sh](https://roadmap.sh)'s
[Todo List API](https://roadmap.sh/projects/todo-list-api) project, written in Go.

As a heads-up, I wrote this to see how quickly I could learn to build an API
in Go, having not worked with front- or back-end tech in more than a year, and
being somewhat new to Go.

## Building and Running

To build the server, just run

```bash
go build
```

And to run it

```bash
./srtodo
```

## The API

The API is divided into six routes:

- `POST /register` - Registers the user;
- `POST /login` - Authenticates the user;
- `GET /todos?page=1&limit=10` - Gets the user's todos (1st page, 10 todos);
- `POST /todos` - Creates a todo for the authenticated user;
- `PUT /todos/{id}` - Updates a todo for the authenticated user;
- `DELETE /todos/{id}` - Deletes a todo.

### Registration Route

The registration route expects the following JSON:

```json
{
  "name": "John Doe",
  "email": "some@email.com",
  "password": "password"
}
```

For a successful registration, `email` has to be unique. Afterwards the route
returns

```json
{
  "token": "JWT"
}
```

The returned token will be a valid JWT token, and should be used in the
request's header as `Authorization`.

### Login Route

The login route expects the following JSON:

```json
{
  "email": "some@email.com",
  "password": "password"
}
```

And after a successful authentication, returns the JWT token like the
registration route.

### Fetching Todos

To fetch the current users todos, there's the `GET /todos` route. It supports
pagination via the `page` and `limit` query parameters. The returned JSON looks
like the following:

```json
{
  "data": [
    {
      "id": 1,
      "title": "todo title",
      "description": "todo description"
    },
    {
      "id": 2,
      "title": "todo title",
      "description": "todo description"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 2
}
```

The `page` and `limit` parameters, by default, are equal to `1` and `10`.

### Adding a Todo

To add a todo for the current user, use the `POST /todos` route. It expects
the following request JSON:

```json
{
  "title": "todo title",
  "description": "todo description"
}
```

And after a successful addition, the response JSON should be:

```json
{
  "id": 1,
  "title": "todo title",
  "description": "todo description"
}
```

The `id` can, of course, be different.

### Updating a Todo

To update a todo, there's the `PUT /todos/{id}` route. The `id` is the ID of the
todo you want to change. The request and response JSONs are the same as above.

### Deleting a Todo

To delete a todo, use the `DELETE /todos/{id}` route, where `id` is the ID of
the todo you want to delete. The API should return a `204` status code upon
success.
