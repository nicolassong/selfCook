const { request } = require('../../utils/request');

Page({
  data: {
    addresses: [],
    editingId: null,
    form: {
      contactName: '演示用户',
      contactPhone: '13800001111',
      province: '广东省',
      city: '深圳市',
      district: '南山区',
      detailAddress: '',
      communityName: '科技园',
      isDefault: true
    }
  },

  onShow() {
    this.loadAddresses();
  },

  async loadAddresses() {
    try {
      const data = await request('/me/addresses');
      this.setData({ addresses: data.list || [] });
    } catch (error) {
      wx.showToast({ title: '加载地址失败', icon: 'none' });
    }
  },

  setField(event) {
    const { field } = event.currentTarget.dataset;
    this.setData({ [`form.${field}`]: event.detail.value });
  },

  editAddress(event) {
    const { id } = event.currentTarget.dataset;
    const current = this.data.addresses.find((item) => item.id === id);
    if (!current) return;
    this.setData({
      editingId: id,
      form: {
        contactName: current.contactName,
        contactPhone: current.contactPhone,
        province: current.province,
        city: current.city,
        district: current.district,
        detailAddress: current.detailAddress,
        communityName: current.communityName,
        isDefault: current.isDefault
      }
    });
  },

  resetForm() {
    this.setData({
      editingId: null,
      form: {
        contactName: '演示用户',
        contactPhone: '13800001111',
        province: '广东省',
        city: '深圳市',
        district: '南山区',
        detailAddress: '',
        communityName: '科技园',
        isDefault: true
      }
    });
  },

  async saveAddress() {
    try {
      if (this.data.editingId) {
        await request(`/me/addresses/${this.data.editingId}`, 'PUT', this.data.form);
      } else {
        await request('/me/addresses', 'POST', this.data.form);
      }
      wx.showToast({ title: '地址已保存', icon: 'success' });
      this.resetForm();
      this.loadAddresses();
    } catch (error) {
      wx.showToast({ title: '保存失败', icon: 'none' });
    }
  },

  async setDefault(event) {
    const { id } = event.currentTarget.dataset;
    try {
      await request(`/me/addresses/${id}/default`, 'POST', {});
      wx.showToast({ title: '已设默认', icon: 'success' });
      this.loadAddresses();
    } catch (error) {
      wx.showToast({ title: '设置失败', icon: 'none' });
    }
  },

  async removeAddress(event) {
    const { id } = event.currentTarget.dataset;
    try {
      await request(`/me/addresses/${id}`, 'DELETE');
      wx.showToast({ title: '已删除', icon: 'success' });
      this.loadAddresses();
    } catch (error) {
      wx.showToast({ title: '删除失败', icon: 'none' });
    }
  }
});
