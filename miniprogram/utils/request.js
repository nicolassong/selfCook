const app = getApp();

function request(path, method = 'GET', data = undefined) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${app.globalData.baseUrl}${path}`,
      method,
      data,
      success(res) {
        const body = res.data || {};
        if (res.statusCode >= 200 && res.statusCode < 300 && body.code === 0) {
          resolve(body.data);
          return;
        }
        reject(new Error(body.message || `request failed: ${res.statusCode}`));
      },
      fail(err) {
        reject(err);
      }
    });
  });
}

module.exports = { request };
