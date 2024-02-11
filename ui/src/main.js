import Vue from 'vue'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import koLang from 'element-ui/lib/locale/lang/ko'
import App from './App.vue'
import router from './router'
import store from './store'

Vue.config.productionTip = false

Vue.use(ElementUI, {
  locale: koLang
})

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
