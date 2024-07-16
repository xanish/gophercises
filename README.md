# Gophercises

Gophercises consists of a series of mini-exercises that progressively introduce different aspects of the Go programming language.
Material can be found [here](https://courses.calhoun.io/courses/cor_gophercises).

| #  | Topic                             | Solution                                    |
|----|-----------------------------------|---------------------------------------------|
| 01 | Quiz Game                         | [here](./quiz_game/README.md)               |
| 02 | URL Shortener                     | [here](./url_shortener/README.md)           |
| 03 | Choose Your Own Adventure         | [here](./choose_your_adventure/README.md)   |
| 04 | HTML Link Parser                  | [here](./html_link_parser/README.md)        |
| 05 | Sitemap Builder                   | [here](./sitemap_builder/README.md)         |
| 06 | Hacker Rank Problem               | [here](./strings_and_bytes/README.md)       |
| 07 | CLI Task Manager                  | [here](task_manager/README.md)        |

## Checking Test Coverage

- To generate report run `go test -coverprofile=".\PACKAGE_NAME\coverage.out" .\PACKAGE_NAME\`
- To view coverage in pretty html format run `go tool cover -html=".\PACKAGE_NAME\coverage.out"`
