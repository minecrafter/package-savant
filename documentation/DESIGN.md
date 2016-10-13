# Design

## Exposures

For most average users, the only thing that is relevant is the "exposures" that
Sage provides. It looks up metadata and returns packages from storage. Currently,
Sage only supports a Maven exposure.

## Package storage

Backing all exposures is a generic storage backend that allows package information
to be retrieved on demand.

Sage decouples metadata from package content to allow flexibility in how you can
serve your packages. For instance, you might look up metadata in MySQL but serve
packages over Amazon S3.
