import Vue from 'vue'

import Vuex from 'vuex'
Vue.use(Vuex)

Vue.use(require('vue-cookies'))

import VueRouter from 'vue-router'
Vue.use(VueRouter)

Vue.use(require('vue-moment'));

import HomePage from './components/HomePage'
import ProfilePage from './components/ProfilePage'
import SettingsPage from './components/SettingsPage'
import HelpPage from './components/HelpPage'
import TrendingPage from './components/TrendingPage'
import PostPage from './components/PostPage'
import SearchPage from "./components/SearchPage";
import ThreadPage from "./components/ThreadPage";

import PageNotFound from "./components/PageNotFound";


const router = new VueRouter({
  mode: 'history',
  base: __dirname,
  routes: [
    { name: 'home', path: '/', component: HomePage, meta: { title: 'Home' } },
    { name: 'profile', path: '/profile/:id(\\d+)', component: ProfilePage, meta: { title: 'Profile' }, props: (route) => ({ profileId: parseInt(route.params.id) }) },
    { name: 'settings', path: '/settings', component: SettingsPage, meta: { title: 'Settings' } },
    { name: 'help', path: '/help', component: HelpPage, meta: { title: 'Help' } },
    { name: 'trending', path: '/trending', component: TrendingPage, meta: { title: 'Trending' } },
    { name: 'post', path: '/post/:id(\\d+)', component: PostPage, meta: { title: 'Post' }, props: (route) => ({ postId: parseInt(route.params.id) }) },
    { name: 'thread', path: '/thread/:id(\\d+)', component: ThreadPage, meta: { title: 'Thread' }, props: (route) => ({ threadId: parseInt(route.params.id) }) },
    { name: 'search', path: '/search', component: SearchPage, meta: { title: 'Search' } },
    { name: 'page_not_found', path: "*", component: PageNotFound, meta: {title: 'Page Not Found'} },
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
Vue.use(Vuetify, {
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
