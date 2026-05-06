const config = require('./config');

App({
  globalData: {
    baseUrl: config.baseUrl,
    assetBaseUrl: config.assetBaseUrl
  }
});
