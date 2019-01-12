<template>
  <v-app id="app">
    <div v-if="this.$route.name !== 'page_not_found'">

      <!-- Toolbar (navbar) -->
      <AppToolbar />

      <!-- Snackbar -->
      <v-snackbar v-model="snackBar" :bottom="true" :right="true">
        {{ snackBarText }}
        <v-btn color="red" flat @click="snackBar = false">Close</v-btn>
      </v-snackbar>
    </div>
    <!-- Router -->
    <v-content>
      <router-view class="view"></router-view>
    </v-content>
  </v-app>
</template>

<script>
import VueCookies from "vue-cookies";

import store from "./store";

import AppToolbar from "./components/AppToolbar";

export default {
  name: "app",
  computed: {
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
    AppToolbar
  },
  watch: {
    $route: {
      deep: true,
      handler: function() {
        this.activeRoutePath = this.$route.path;
      }
    }
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
