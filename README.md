# Gophercises

Gophercises consists of a series of mini-exercises that progressively introduce different aspects of the Go programming language.
Material can be found [here](https://courses.calhoun.io/courses/cor_gophercises).

| #  | Topic                          | Solution                                    | Tests |
|----|--------------------------------|---------------------------------------------|-------|
| 01 | Quiz Game                      | [here](./quiz_game/README.md)               | Done  |
| 02 | URL Shortener                  | [here](./url_shortener/README.md)           | Done  |
| 03 | Choose Your Own Adventure      | [here](./choose_your_adventure/README.md)   | Done  |
| 04 | HTML Link Parser               | [here](./html_link_parser/README.md)        | Done  |
| 05 | Sitemap Builder                | [here](./sitemap_builder/README.md)         | Done  |
| 06 | Hacker Rank Problem            | [here](./strings_and_bytes/README.md)       | Done  |
| 07 | CLI Task Manager               | [here](./cli_task_manager/README.md)        | Done  |
| 08 | Phone Number Normalizer        | [here](./phone_number_normalizer/README.md) | Done  |
| 09 | Deck of Cards                  | [here](./deck_of_cards/README.md)           | Done  |
| 10 | Blackjack                      | [here](./blackjack/README.md)               | Todo  |
| 11 | Blackjack AI                   | [here](./blackjack_ai/README.md)            | Todo  |
| 12 | File Renaming Tool             | [here](./file_renaming_tool/README.md)      | Todo  |
| 13 | Quiet Hacker News              | [here](./hacker_news/README.md)             | Done  |
| 14 | Recover Middleware             | [here](./recover_middleware/README.md)      | Todo  |
| 15 | Development Recover Middleware | [here](./dev_recover_middleware/README.md)  | Todo  |

## Checking Test Coverage

- To generate report run `go test -coverprofile=".\PACKAGE_NAME\coverage.out" .\PACKAGE_NAME\`
- To view coverage in pretty html format run `go tool cover -html=".\PACKAGE_NAME\coverage.out"`
