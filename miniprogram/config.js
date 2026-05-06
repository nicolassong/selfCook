const defaults = {
  baseUrl: '',
  assetBaseUrl: ''
};

let local = {};
try {
  local = require('./config.local');
} catch (error) {
  local = {};
}

module.exports = {
  ...defaults,
  ...local
};
