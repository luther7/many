
# Many

[![Build Status](https://travis-ci.org/rubberydub/many.svg?branch=master)](https://travis-ci.org/rubberydub/many)

# To Do

- [X] Idea in read me
- [ ] Test harness
- [ ] Basic working git-backed example
- [ ] Docker
- [ ] Read me
- [ ] 80% test coverage

# Introduction

Many is a tool for versioning collections of services and applications.

Consider the follow scenario: A product consists of a backend API service and a
frontend web application. While the backend and frontend are in separate git 
repositories and deployed as separate services, they are tested and released 
together as part of the overall product. The versions of the overall product
and it's constituent services might hypothetically be like so:

| Overall Product | Backend API     | Frontend App    |
|-----------------|-----------------|-----------------|
| v1.0.0          | ae9d8d5         | ca039c2         |
|                 | b76c124         | 3e35198         |
| v1.0.1          | a2bbb7a         | c4e3b93         |
|                 | 04db8d4         | 13ffd3e         |
|                 | dcb6d2d         | 030be43         |
| v1.1.0          | e16e3d2         | 48eee8d         |

**Note:** The overall product's version might not be strictly correct semantic
versioning. This is just an example.

Many aims to provide tooling to simplify the management of these overall
versions. It provides a CLI tool to manage versions as a TOML file, and aims to
be CI friendly. The TOML file is stored in a git repository providing a single
source of truth in a familiar environment.

# Usage

**TODO**









