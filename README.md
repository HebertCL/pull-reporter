# pull-reporter
A simple program to generate Pull requests reports

## TODO

There's a bunch of improvements that can be done to this (very basic) codebase:

- Write tests.
- Extend functionality to:
    - Tag PRs which are older than X days and no activity with `stale` tag.
    - Close PRs which have been hanging more than X days with `stale` tag.
    - Delete branches which are not associated to open PRs.
    - Add additional notification mechanisms besides email (i.e. Slack, Discord).
- (Optionally) enable it to work as webserver/cli.
- (Optionally) support GitHub and GitLab.
