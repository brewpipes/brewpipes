# BrewPipes "Identity" Service Agent Guide

This guide is for agentic coding tools working in this particular service in the BrewPipes repo.
It captures commands and conventions observed in the current codebase.
When making changes to this service, its data models, its logic, or any other aspect, this document MUST be kept-up-to-date do that it accurately reflects the actual implementation.

## Service Domain

The Identity service models user accounts and token issuance for BrewPipes.

## Overview

The big picture: the system stores local users with credentials and issues access/refresh tokens for API access.

User
- A user is a local account with a username and a stored password hash.
- Users are identified by UUIDs for token subjects and API lookups.

Credential
- Credentials are currently username/password pairs stored as password hashes.
- Passwords are verified via bcrypt comparisons.

Access Token
- An access token is a JWT used to authorize API requests.
- It includes issuer, subject (user UUID), and role claims.

Refresh Token
- A refresh token is a longer-lived JWT used to mint new access tokens.
- It includes issuer and subject claims.

## User Journey: Brewery Admin

Here is a simple identity story from setup to daily use.

You create a local user account for a brewer. The user signs in with a username and password and receives an access token for API calls plus a refresh token for longer sessions. When the access token expires, the client exchanges the refresh token for a new token pair and continues working without re-entering credentials.

In short:
- users define who can sign in,
- credentials validate identity,
- access tokens authorize requests,
- refresh tokens extend sessions.

## Acceptance Criteria

- A brewery admin can create, list, fetch, update, and deactivate users.
- Login accepts a username and password and returns access and refresh tokens.
- Refresh accepts a refresh token and returns a new token pair.
- Passwords are stored as hashes and never returned in API responses.
- Tokens include issuer and subject claims and are signed with the service secret.
- The identity service remains self-contained with no cross-service foreign keys.
