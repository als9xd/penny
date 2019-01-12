<template>
  <div v-if="profile">
    <v-layout row wrap id="profile-header">
      <v-container>
        <v-layout row wrap>
          <v-flex xs2>
            <v-list-tile-avatar class="justify-end pr-4" size="64">
              <img v-if="profile.avatar" :src="`http://localhost:3000/uploads/${profile.avatar}`">
              <v-icon v-else size="64">account_box</v-icon>
            </v-list-tile-avatar>
          </v-flex>
          <v-flex xs10>
            <div id="profile-username">{{profile.username}}</div>
            {{profile.email}} • Joined {{ profile.created | moment("from") }}
          </v-flex>
        </v-layout>
      </v-container>
    </v-layout>
    <v-container></v-container>
  </div>
</template>

<script>
import store from "../store";

export default {
  name: "Profile",
  mounted: function() {
    const otherProfileId = this.$route.params.id;
    if (otherProfileId) {
      this.axios
        .get(`http://localhost:3000/api/v1/u/${otherProfileId}`)
        .then(response => {
          this.otherProfile = response.data.data.profile;
        });
    } else {
      if (!this.profile) {
        this.$router.push({ path: "/" });
      }
    }
  },
  data() {
    return {
      otherProfile: null
    };
  },
  computed: {
    profile() {
      return this.$route.params.id ? this.otherProfile : store.state.profile;
    }
  }
};
</script>

<style scoped>
#profile-header {
  background-color: #eee;
}
#profile-username {
  font-family: "Roboto", sans-serif;
  font-weight: 400;
  font-size: 2.4rem;
  line-height: 3rem;
}
</style>
