This is backend

sebelum jalanin, go mod tidy


Update Log
(+) New
(-) Delete
(++) Update add minor
(--) Update delete minor

v0.1 (init?)
- (+) Main.go
- (+) go module :  agatra (Arcade GAme TRAcker)
    - (+) fiber v2
    - (+) godotenv

v0.2 (coffeeboosted afternoon)
- (+) Folder api, handlers, repos, utils [all empty]
- (+) Folder db > (+) postgres.go
        - (+) struct config, dbstruct
        - (+) func connect, reset, automigrate [empty]
- (+) Folder middleware > (+) auth.go > (+) func auth
- (+) Folder models 
    - (+) jwt.go > (+) struct claims
    - (+) response.go
        - (+) struct Error + Success response
        - (+) func new error + success response
