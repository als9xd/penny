import Vue from 'vue'

import Vuex from 'vuex'
Vue.use(Vuex)

Vue.use(require('vue-cookies'))

import VueRouter from 'vue-router'
Vue.use(VueRouter)

Vue.use(require('vue-moment'));

import Home from './components/Home.vue'
import Profile from './components/Profile.vue'
import Settings from './components/Settings.vue'
import Help from './components/Help.vue'
import Trending from './components/Trending.vue'
import Post from './components/Post.vue'
import Search from "./components/Search";
import Thread from "./components/Thread";



const router = new VueRouter({
  mode: 'history',
  base: __dirname,
  routes: [
    { path: '/', component: Home, meta:{ title: 'Home'} },
    { path: '/profile', component: Profile, meta:{ title: 'Profile'} },
    { path: '/profile/:id(\\d+)', component: Profile, meta:{ title: 'Profile'} },
    { path: '/settings', component: Settings, meta:{ title: 'Settings'} },
    { path: '/help', component: Help, meta:{ title: 'Help'} },
    { path: '/trending', component: Trending, meta:{ title: 'Trending'} },
    { path: '/post/:id(\\d+)', component: Post, meta:{ title: 'Post'},props:(route)=>({postId:parseInt(route.params.id)}) },
    { name:'thread', path: '/thread/:id(\\d+)', component: Thread, meta:{ title: 'Thread'},props:(route)=>({threadId:parseInt(route.params.id)}) },
    { name:'search', path: '/search', component: Search, meta:{ title: 'Search'} },
  ]
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title
  next()
})

import storePlugin from './storePlugin'
Vue.use(storePlugin) 

import axios from 'axios'
import VueAxios from 'vue-axios'
 
Vue.use(VueAxios, axios)

import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
Vue.use(Vuetify,{
  theme: {
    primary: '#051046',
    secondary: '#eee'
  }
})

import App from './App.vue'

Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
