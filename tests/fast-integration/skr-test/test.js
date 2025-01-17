const {
  gatherOptions,
  oidcE2ETest,
} = require('./index');
const {getOrProvisionSKR} = require('./provision/provision-skr');
const {deprovisionAndUnregisterSKR} = require('./provision/deprovision-skr');

const skipProvisioning = process.env.SKIP_PROVISIONING === 'true';
const provisioningTimeout = 1000 * 60 * 30; // 30m
const deprovisioningTimeout = 1000 * 60 * 95; // 95m
let globalTimeout = 1000 * 60 * 70; // 70m
const slowTime = 5000;

describe('SKR test', function() {
  if (!skipProvisioning) {
    globalTimeout += provisioningTimeout + deprovisioningTimeout;
  }
  this.timeout(globalTimeout);
  this.slow(slowTime);

  let options = gatherOptions(); // with default values
  let skr;
  const getShootInfoFunc = function() {
    return skr.shoot;
  };

  before('Ensure SKR is provisioned', async function() {
    this.timeout(provisioningTimeout);
    skr = await getOrProvisionSKR(options, skipProvisioning, provisioningTimeout);
    options = skr.options;
  });

  // Run the OIDC tests
  oidcE2ETest(options, getShootInfoFunc);

  after('Cleanup the resources', async function() {
    this.timeout(deprovisioningTimeout);
    await deprovisionAndUnregisterSKR(options, deprovisioningTimeout, skipProvisioning, false);
  });
});
