<template>
  <div>
    <v-tabs slot="extension" v-model="model" color="secondary" slider-color="primary">
      <v-tab href="#tab-profiles">Profiles</v-tab>
      <v-tab href="#tab-threads">Threads</v-tab>
      <v-tab href="#tab-posts">Posts</v-tab>
    </v-tabs>

    <v-tabs-items v-model="model">
      <v-tab-item value="tab-profiles">
        <v-container>
          <v-layout row wrap>
            <v-flex xs12>
              <v-list two-line>
                <template v-for="(result,index) in profiles">
                  <v-list-tile
                    :key="result.username"
                    avatar
                    v-on:click.self="profile && result.id===profile.id?changeRoute('/profile'):changeRoute(`/profile/${result.id}`)"
                  >
                    <v-list-tile-avatar v-if="profile" @click="submitProfileSubscribe(result)">
                      <v-icon
                        v-on:mouseover="activeAddIcon=result.id"
                        v-on:mouseleave="activeAddIcon=null"
                        :color="activeAddIcon===result.id||subscriptions.profiles.map(s=>s.id).filter(id=>id==result.id).length?'primary':null"
                      >person_add</v-icon>
                    </v-list-tile-avatar>
                    <v-list-tile-avatar
                      @click="profile && result.id===profile.id?changeRoute('/profile'):changeRoute(`/profile/${result.id}`)"
                    >
                      <img
                        v-if="result.avatar"
                        :src="`http://localhost:3000/uploads/${result.avatar}`"
                      >
                      <v-icon v-else x-large>account_circle</v-icon>
                    </v-list-tile-avatar>

                    <v-list-tile-content
                      @click="profile && result.id===profile.id?changeRoute('/profile'):changeRoute(`/profile/${result.id}`)"
                    >
                      <v-list-tile-title v-html="result.username"></v-list-tile-title>
                      <!-- <v-list-tile-sub-title v-html="post.subtitle"></v-list-tile-sub-title> -->
                    </v-list-tile-content>
                  </v-list-tile>
                  <v-divider :key="`divider-${result.username}`" v-if="index!==profiles.length-1"/>
                </template>
              </v-list>
            </v-flex>
          </v-layout>
        </v-container>
      </v-tab-item>
      <v-tab-item value="tab-threads">
        <v-container>
          <v-layout row wrap>
            <v-flex xs12>
              <v-list two-line>
                <template v-for="(result,index) in threads">
                  <v-list-tile
                    :key="result.username"
                    avatar
                    v-on:click.self="changeRoute(`/thread/${result.id}`)"
                  >
                    <v-list-tile-avatar v-if="profile" @click="submitThreadSubscribe(result)">
                      <v-icon
                        v-on:mouseover="activeAddIcon=result.id"
                        v-on:mouseleave="activeAddIcon=null"
                        :color="activeAddIcon===result.id||subscriptions.threads.map(s=>s.id).filter(id=>id==result.id).length?'primary':null"
                      >group_add</v-icon>
                    </v-list-tile-avatar>
                    <v-list-tile-avatar size="32" @click="changeRoute(`/thread/${result.id}`)">
                      <img
                        v-if="result.avatar"
                        :src="`http://localhost:3000/uploads/${result.avatar}`"
                      >
                      <v-icon v-else dark style="background:gray">people</v-icon>
                    </v-list-tile-avatar>

                    <v-list-tile-content @click="changeRoute(`/thread/${result.id}`)">
                      <v-list-tile-title v-html="result.name"></v-list-tile-title>
                      <!-- <v-list-tile-sub-title v-html="post.subtitle"></v-list-tile-sub-title> -->
                    </v-list-tile-content>
                  </v-list-tile>
                  <v-divider :key="`divider-${result.id}`" v-if="index!==threads.length-1"/>
                </template>
              </v-list>
            </v-flex>
          </v-layout>
        </v-container>
      </v-tab-item>
      <v-tab-item value="tab-posts">
        <v-container>
          <ThreadPosts :posts="posts" />
        </v-container>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script>
import store from "../store";
import ThreadPosts from "../components/ThreadPosts";

export default {
  name: "Home",
  computed: {
    profile() {
      return store.state.profile;
    },
    subscriptions() {
      return store.state.subscriptions;
    },
    jwt() {
      return store.state.jwt;
    }
  },
  methods: {
    submitProfileSubscribe: function(profile) {
      this.axios
        .post(
          `http://localhost:3000/api/v1/protected/subscribe/u/${profile.id}`,
          {},
          {
            headers: {
              Authorization: `Bearer ${this.jwt}`,
              "Content-Type": "application/json"
            }
          }
        )
        .then(response => {
          const subscriptions = Object.assign({}, this.subscriptions);
          if (response.data.data.subscribed) {
            store.commit(
              "setSnackBarText",
              `Subscribed to ${profile.username}`
            );
            subscriptions.profiles = [profile, ...subscriptions.profiles];
          } else {
            store.commit(
              "setSnackBarText",
              `Unsubscribed from ${profile.username}`
            );
            subscriptions.profiles = subscriptions.profiles.filter(
              p => p.id !== response.data.data.subscription.to_profile_id
            );
          }
          store.commit("setSubscriptions", subscriptions);
        })
        .catch(err => {
          if (err.response && err.response.data.error.message) {
            store.commit("setSnackBarText", err.response.data.error.message);
            return;
          }
          throw err;
        });
    },
    submitThreadSubscribe: function(thread) {
      this.axios
        .post(
          `http://localhost:3000/api/v1/protected/subscribe/t/${thread.id}`,
          {},
          {
            headers: {
              Authorization: `Bearer ${this.jwt}`,
              "Content-Type": "application/json"
            }
          }
        )
        .then(response => {
          const subscriptions = Object.assign({}, this.subscriptions);
          if (response.data.data.subscribed) {
            store.commit("setSnackBarText", `Subscribed to ${thread.name}`);
            subscriptions.threads = [thread, ...subscriptions.threads];
          } else {
            store.commit("setSnackBarText", `Unsubscribed from ${thread.name}`);
            subscriptions.threads = subscriptions.threads.filter(
              t => t.id !== response.data.data.subscription.thread_id
            );
          }
          store.commit("setSubscriptions", subscriptions);
        })
        .catch(err => {
          if (err.response && err.response.data.error.message) {
            store.commit("setSnackBarText", err.response.data.error.message);
            return;
          }
          throw err;
        });
    },
    changeRoute: function(path) {
      this.$router.push({ path });
    },
    search: function(newSearchInput) {
      this.axios
        .get(`http://localhost:3000/api/v1/search`, {
          params: { search: newSearchInput }
        })
        .then(response => {
          this.profiles = response.data.data.profiles;
          this.threads = response.data.data.threads;
          this.posts = response.data.data.posts;
        });
    }
  },
  mounted() {
    this.search(this.$route.query.s);
  },
  data() {
    return {
      activeAddIcon: null,
      profiles: [],
      threads: [],
      posts: [],
      model: "tab-profiles"
    };
  },
  components:{
      ThreadPosts,
  },
};
</script>

<style scoped>
.post-list-card {
  padding: 1rem;
}

.post-title {
  font-size: 18px;
  font-weight: 400;
  max-height: 4.8rem;
  overflow: hidden;
}

.profile-avatar {
  height: 32px;
  width: 32px;
  border-radius: 32px;
  cursor: pointer;
}

.post-avatar {
  height: 138px;
}
.profile-title {
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
}

.post-subtext {
  font-size: 13px;
  font-weight: 400;
  color: #606060;
}

.post-subtext .profile-username{
    display: inline;
}

.post-subtext .post-time {
    display: inline;
}

.post-subtext .post-time::before {
    content: '•';
    margin:  0 0.25rem;
}

.post-description {
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical; 
  overflow: hidden;
}

</style>
