# Crowdsec-helper-service

This service provides access over a REST API to a CrowdSec feature, that can normally only be called from the command line.
Originally the REST API of CrowdSec only allows to query decisions, but not to create or delete decisions. But because in case you want, for example, to trigger the deletion of a CAPTCHA decision from a frontend, you need the ability to call an REST API endpoint from the frontend.

This services requires, to run on the same system, where CrowdSec is installed. It requires from the caller the knowledge of all arguments that should be passed to the CrowdSec commandline. Also this service is only listening on `localhost`. This should prevent, that any not authorized user can call the API.

