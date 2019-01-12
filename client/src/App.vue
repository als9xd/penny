<template>
  <v-app id="app">
    <v-navigation-drawer v-model="sidebar" app>
      <v-list>
        <v-list-tile @click="changeRoute('/');" avatar>
          <v-list-tile-avatar>
            <v-icon :color="activeRoutePath==='/'?'primary':null">home</v-icon>
          </v-list-tile-avatar>
          <v-list-tile-content>
            <v-list-tile-title>Home</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile @click="changeRoute('/trending');" avatar>
          <v-list-tile-avatar>
            <v-icon :color="activeRoutePath==='/trending'?'primary':null">trending_up</v-icon>
          </v-list-tile-avatar>
          <v-list-tile-content>
            <v-list-tile-title>Trending</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-divider/>
        <v-list-tile>
          <h5>Profiles</h5>
        </v-list-tile>
        <template v-for="(profile, index) in subscriptions.profiles">
          <v-list-tile :key="`profile-${index}`" @click="changeRoute(`/profile/${profile.id}`);" avatar>
            <v-list-tile-avatar>
              <img v-if="profile.avatar" :src="`http://localhost:3000/uploads/${profile.avatar}`">
              <v-icon v-else large>account_box</v-icon>
            </v-list-tile-avatar>
            <v-list-tile-content>
              <v-list-tile-title>{{ profile.username }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </template>
        <v-divider/>
        <v-list-tile>
          <h5>Threads</h5>
        </v-list-tile>
        <template v-for="(thread, index) in subscriptions.threads">
          <v-list-tile :key="`thread-${index}`" @click="changeRoute(`/thread/${thread.id}`);" avatar>
            <v-list-tile-avatar  size="32">
              <img v-if="thread.avatar" :src="`http://localhost:3000/uploads/${thread.avatar}`">
              <v-icon v-else dark style="background:gray">people</v-icon>
            </v-list-tile-avatar>
            <v-list-tile-content>
              <v-list-tile-title>{{ thread.name }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </template>
      </v-list>
    </v-navigation-drawer>

    <v-toolbar color="primary" app>
      <v-toolbar-side-icon @click="sidebar = !sidebar" dark></v-toolbar-side-icon>
      <v-toolbar-items class="hidden-sm-and-down">
        <v-btn flat dark @click="()=>changeRoute('/')">
          <img class="penny-logo" alt="Penny logo" src="./assets/penny_new.png">
          <v-toolbar-title id="toolbar-title">{{appTitle.toUpperCase()}}</v-toolbar-title>
        </v-btn>
      </v-toolbar-items>
      <v-spacer></v-spacer>
      <v-toolbar-items>
        <div class="mr-2">
          <v-text-field
            v-model="searchInput"
            @keyup.enter.native="submitSearch"
            dark
            color="white"
            label="Search"
            clearable
          />
        </div>
        <v-btn flat dark class="toolbar-icon-btn">
          <v-icon>email</v-icon>
        </v-btn>
        <v-btn flat dark class="toolbar-icon-btn">
          <v-icon>notifications</v-icon>
        </v-btn>
        <v-menu offset-y left v-if="profile">
          <v-btn class="toolbar-icon-btn" slot="activator" flat dark>
            <v-icon>add</v-icon>
            <v-icon>arrow_drop_down</v-icon>
          </v-btn>
          <v-list>
            <CreateThreadDialog/>
            <CreatePostDialog v-if="$route.name==='thread'"/>
          </v-list>
        </v-menu>
        <v-menu offset-y left v-if="profile">
          <v-btn class="toolbar-icon-btn" slot="activator" flat dark>
            <v-icon>account_box</v-icon>
            <v-icon>arrow_drop_down</v-icon>
          </v-btn>
          <v-list>
            <v-list-tile @click="changeRoute('/profile')" avatar>
              <v-list-tile-title>My Profile</v-list-tile-title>
            </v-list-tile>
            <v-list-tile @click="changeRoute('/settings')" avatar>
              <v-list-tile-title>Settings</v-list-tile-title>
            </v-list-tile>
            <v-list-tile @click="changeRoute('/help')" avatar>
              <v-list-tile-title>Help</v-list-tile-title>
            </v-list-tile>
            <v-divider/>
            <v-list-tile @click="store.commit('logOut')" avatar>
              <v-list-tile-title>Log Out</v-list-tile-title>
            </v-list-tile>
          </v-list>
        </v-menu>
        <template v-else>
          <LoginDialog/>
          <SignupDialog/>
        </template>
      </v-toolbar-items>
    </v-toolbar>

    <v-snackbar v-model="snackBar" :bottom="true" :right="true">
      {{ snackBarText }}
      <v-btn color="red" flat @click="snackBar = false">Close</v-btn>
    </v-snackbar>

    <v-content>
      <router-view class="view"></router-view>
    </v-content>
  </v-app>
</template>

<script>
import VueCookies from "vue-cookies";

import store from "./store";

import CreateThreadDialog from "./components/CreateThreadDialog";
import CreatePostDialog from "./components/CreatePostDialog";

import LoginDialog from "./components/LoginDialog";
import SignupDialog from "./components/SignupDialog";

export default {
  name: "app",
  computed: {
    store (){
      return store;
    },
    subscriptions() {
      return store.state.subscriptions;
    },
    threads() {
      return store.state.threads;
    },
    profile() {
      return store.state.profile;
    },
    snackBar: {
      set(enabled) {
        store.commit("setSnackBar", enabled);
      },
      get() {
        return store.state.snackBar;
      }
    },
    snackBarText() {
      return store.state.snackBarText;
    }
  },
  components: {
    CreateThreadDialog,
    CreatePostDialog,
    LoginDialog,
    SignupDialog
  },
  watch: {
    $route: {
      deep: true,
      handler: function() {
        this.activeRoutePath = this.$route.path;
      }
    }
  },
  methods: {
    submitSearch: function(event) {
      store.commit("commitSearchInput", event.target.value);
      this.$router.push({ name: "search", query: { s: event.target.value } });
    },
    changeRoute: function(path) {
      this.$router.push({ path });
    },
  },
  mounted: function() {
    const jwt = VueCookies.get("loginToken");
    store.commit("setJWT", jwt);

    if (jwt) {
      this.axios
        .get("http://localhost:3000/api/v1/protected/u", {
          headers: {
            Authorization: `Bearer ${jwt}`,
            "Content-Type": "application/json"
          }
        })
        .then(response => {
          store.commit("setProfile", response.data.data.profile);
          store.commit("setSubscriptions", response.data.data.subscriptions);
          store.commit("setThreads", response.data.data.threads);
        })
        .catch(err => {
          if (err.response && err.response.status === 401) {
            store.commit("logOut");
            return;
          }
          throw err;
        });
    }
  },
  data() {
    return {
      appTitle: "Penny",
      searchInput: "",
      imageName: "",
      imageUrl: "",
      imageFile: "",
      activeRoutePath: this.$route.path,
      sidebar: false,
    };
  }
};
</script>

<style>
html {
  overflow: auto;
}

.penny-logo {
  height: 2.5rem;
}

#toolbar-title {
  margin-left: 0.5rem;
  font-weight: 900;
}

.vert-dialogue {
  align-items: center;
  display: flex !important;
}

.toolbar-btn {
  height: unset !important;
  padding: 0.5rem 0;
}

.toolbar-icon-btn {
  padding: 0 16px;
  min-width: 0;
}

#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
