# SVN Manager


## Internal Structure

The HTTP interface is implemented in the `httphandler` subpackage. This package is responsible for
interpreting the HTTP request by parsing variables from the URL and JSON sent in the body.
Furthermore, it is responsible for validating these data.

The actual work managing on-disk files and directories is implemented in the `svnman` subpackage.
This package assumes the data is vetted as correct by the `httphandler` subpackage.

Apache life cycles, at this moment consisting of graceful restarts, is handled by the `apache`
subpackage. At this moment this isn't implemented at all.
