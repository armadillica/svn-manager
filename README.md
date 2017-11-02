# SVN Manager


## System requirements

Requirements are:

- Apache 2
- RabbitMQ 3

The SVNManager needs to be able to gracefully restart Apache after configuration files have been
created. This is done by invoking `sudo apache2ctl`, and requires that this command can be performed
without having to provide a password. Add the following to `/etc/sudoers` to set this up:

    Cmnd_Alias GRACEFUL = /usr/sbin/apache2ctl configtest, /usr/sbin/apache2ctl graceful
    www-data ALL = NOPASSWD: GRACEFUL

Replace `www-data` with the username of the SVNManager, and `/usr/sbin/apache2ctl` with the absolute
path of that executable.

SVNManager tests that the command `sudo --non-interactive apache2ctl configtest` can be run
successfully at startup. Failing to do so is considered a fatal error and will prevent SVNManager
from starting.


## Internal Structure

The HTTP interface is implemented in the `httphandler` subpackage. This package is responsible for
interpreting the HTTP request by parsing variables from the URL and JSON sent in the body.
Furthermore, it is responsible for validating these data.

The actual work managing on-disk files and directories is implemented in the `svnman` subpackage.
This package assumes the data is vetted as correct by the `httphandler` subpackage.

Apache life cycles, at this moment consisting of graceful restarts, is handled by the `apache`
subpackage.
