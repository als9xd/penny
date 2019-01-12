import Vue from 'vue'
import Vuex from 'vuex'
import VueCookies from "vue-cookies";

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    jwt: null,
    profile: null,
    activePostId: null,
    snackBar: false,
    snackBarText: "",
    subscriptions: {profiles:[],threads:[]},
    threads: [],
  },
  mutations: {
    setJWT (state,jwt) {
      state.jwt = jwt
    },
    setProfile (state,profile) {
        state.profile = profile
    },
    setSnackBarText (state,snackBarText) {
        state.snackBar = true;
        state.snackBarText = snackBarText;
    },
    setSnackBar (state,enabled) {
        state.snackBar = enabled;
    },
    setSubscriptions (state,subscriptions){
        state.subscriptions = subscriptions;
    },
    setThreads (state,threads) {
        state.threads = threads;
    },
    logOut (state){
        state.jwt = null;
        state.profile = null;
        state.subscriptions = {profiles:[],threads:[]}
        VueCookies.remove("loginToken");
    }
  }
})
