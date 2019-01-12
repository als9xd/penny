import store from './store'  // this is your store object
export default {  
  store,
  // we can add objects to the Vue prototype in the install() hook:
  install (Vue) {
    Vue.prototype.$globalStore = store
  }
}