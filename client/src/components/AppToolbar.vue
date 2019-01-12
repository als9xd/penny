<template>
  <div>
    <!-- Toolbar (navbar) -->
    <v-toolbar color="primary" app>
      <v-toolbar-side-icon @click="navigationDrawer = !navigationDrawer" dark></v-toolbar-side-icon>
      <v-toolbar-items class="hidden-sm-and-down">
        <v-btn flat dark @click="$router.push({path:'/'})">
          <img class="penny-logo" alt="Penny logo" src="../assets/penny_new.png">
          <v-toolbar-title id="toolbar-title">Penny</v-toolbar-title>
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

    <!-- Navigation Drawer (sidebar) -->
    <v-navigation-drawer v-model="navigationDrawer" app>
      <v-list>
        <v-list-tile @click="$router.push({name:'home'})" avatar>
          <v-list-tile-avatar>
            <v-icon :color="$route.name==='home'?'primary':null">home</v-icon>
          </v-list-tile-avatar>
          <v-list-tile-content>
            <v-list-tile-title>Home</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile @click="$router.push({name:'trending'})" avatar>
          <v-list-tile-avatar>
            <v-icon :color="$route.name==='trending'?'primary':null">trending_up</v-icon>
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
          <v-list-tile
            :key="`profile-${index}`"
            @click="$router.push({path:`/profile/${profile.id}`});"
            avatar
          >
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
          <v-list-tile
            :key="`thread-${index}`"
            @click="$router.push({path:`/thread/${thread.id}`});"
            avatar
          >
            <v-list-tile-avatar size="32">
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
  </div>
</template>

<script>
import store from "../store";

import CreateThreadDialog from "../components/CreateThreadDialog";
import CreatePostDialog from "../components/CreatePostDialog";

import LoginDialog from "../components/LoginDialog";
import SignupDialog from "../components/SignupDialog";

export default {
  name: "AppToolbar",
  computed: {
    profile() {
      return store.state.profile;
    },
    drawerOpen() {
      return this.open;
    },
    subscriptions() {
      return store.state.subscriptions;
    },
    threads() {
      return store.state.threads;
    }
  },
  data() {
    return {
      searchInput: "",
      navigationDrawer: false
    };
  },
  components: {
    CreateThreadDialog,
    CreatePostDialog,
    LoginDialog,
    SignupDialog
  },
  methods: {
    submitSearch: function(event) {
      this.$router.push({ name: "search", query: { s: event.target.value } });
    }
  }
};
</script>

<style scoped>
</style>
