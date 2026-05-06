const { request } = require('../../utils/request');
// Voice input is temporarily disabled.
// To enable later, add the WechatSI plugin in app.json and restore:
// const plugin = requirePlugin('WechatSI');
// const recognitionManager = plugin.getRecordRecognitionManager();

Page({
  data: {
    inputValue: '',
    loading: false,
    recording: false,
    voiceEnabled: false,
    scrollIntoView: '',
    messages: [
      {
        id: 1,
        role: 'assistant',
        content: '你好，我是 selfCook 菜单助手。告诉我人数、预算、口味偏好，我帮你从团餐场景里搭配一份好下单的菜单。'
      }
    ],
    quickPrompts: [
      { label: '三人晚餐', text: '三个人晚餐，预算100以内，少辣，想吃得清爽一点。' },
      { label: '家庭套餐', text: '一家四口，有老人和小孩，想要营养均衡，不要太油。' },
      { label: '工作餐', text: '办公室六个人午餐，预算180以内，希望出餐快、方便分餐。' }
    ]
  },

  // initVoiceRecognition() {
  //   recognitionManager.onStart = () => {
  //     this.setData({ recording: true });
  //   };
  //   recognitionManager.onRecognize = (res) => {
  //     const text = (res.result || '').trim();
  //     if (text) {
  //       this.setData({ inputValue: text });
  //     }
  //   };
  //   recognitionManager.onStop = (res) => {
  //     const text = (res.result || '').trim();
  //     this.setData({ recording: false });
  //     if (!text) {
  //       wx.showToast({ title: '没有识别到语音', icon: 'none' });
  //       return;
  //     }
  //     this.setData({ inputValue: text });
  //   };
  //   recognitionManager.onError = () => {
  //     this.setData({ recording: false });
  //     wx.showToast({ title: '语音识别失败', icon: 'none' });
  //   };
  // },

  onInput(event) {
    this.setData({ inputValue: event.detail.value });
  },

  useQuickPrompt(event) {
    this.setData({ inputValue: event.currentTarget.dataset.text || '' });
  },

  startVoiceInput() {
    wx.showToast({ title: '语音输入暂未开启', icon: 'none' });
    // if (this.data.loading || this.data.recording) return;
    // wx.authorize({
    //   scope: 'scope.record',
    //   success: () => {
    //     recognitionManager.start({
    //       lang: 'zh_CN'
    //     });
    //   },
    //   fail: () => {
    //     wx.showModal({
    //       title: '需要麦克风权限',
    //       content: '开启麦克风权限后，可以用语音输入菜单需求。',
    //       confirmText: '去设置',
    //       success: (res) => {
    //         if (res.confirm) wx.openSetting();
    //       }
    //     });
    //   }
    // });
  },

  stopVoiceInput() {
    // if (!this.data.recording) return;
    // recognitionManager.stop();
  },

  async sendMessage() {
    const content = (this.data.inputValue || '').trim();
    if (!content || this.data.loading) return;

    const userMessage = {
      id: Date.now(),
      role: 'user',
      content
    };
    const messages = this.data.messages.concat(userMessage);
    this.setData({
      messages,
      inputValue: '',
      loading: true,
      scrollIntoView: `msg-${userMessage.id}`
    });

    try {
      const result = await request('/ai/recommend', 'POST', {
        message: content,
        history: messages.map((item) => ({
          role: item.role,
          content: item.content
        }))
      });
      const assistantMessage = {
        id: Date.now() + 1,
        role: 'assistant',
        content: result.reply || '暂时没有生成结果，请稍后再试。'
      };
      this.setData({
        messages: this.data.messages.concat(assistantMessage),
        scrollIntoView: `msg-${assistantMessage.id}`
      });
    } catch (error) {
      wx.showToast({ title: error.message || '发送失败', icon: 'none' });
    } finally {
      this.setData({ loading: false });
    }
  }
});
