<template>
  <v-dialog v-model="loginDialog" width="500" class="vert-dialogue">
    <v-btn slot="activator" class="toolbar-btn mx-2" flat outline dark small>Log In</v-btn>

    <v-form v-model="valid" @submit.prevent="logIn()">
      <v-card>
        <v-card-title class="headline grey lighten-2" primary-title>
          <img class="penny-logo mr-2" alt="Penny logo" src="../assets/penny_new.png">Log In
        </v-card-title>

        <v-card-text>
          <v-text-field
            label="Username"
            v-model="loginUsername"
            required
            :rules="usernameRules"
            outline
          />
          <v-text-field
            label="Password"
            v-model="loginPassword"
            outline
            type="password"
            :rules="passwordRules"
            required
          />
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn type="submit" color="primary" dark>Submit</v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>

<script>
import VueCookies from "vue-cookies";

import store from "../store";

export default {
  name: "CreatePostDialog",
  data() {
    return {
      valid: false,
      loginUsername: null,
      loginPassword: null,
      loginDialog: false,

      usernameRules: [v => !!v || "Username is required"],
      passwordRules: [v => !!v || "Password is required"]
    };
  },
  methods: {
    pickFile() {
      this.$refs.image.click();
    },
    onFilePicked(e) {
      const files = e.target.files;
      if (files[0] !== undefined) {
        this.imageName = files[0].name;
        if (this.imageName.lastIndexOf(".") <= 0) {
          return;
        }
        const fr = new FileReader();
        fr.readAsDataURL(files[0]);
        fr.addEventListener("load", () => {
          this.imageUrl = fr.result;
          this.imageFile = files[0]; // this is an image file that can be sent to server...
        });
      } else {
        this.imageName = "";
        this.imageFile = "";
        this.imageUrl = "";
      }
    },
    logIn() {
      this.axios
        .post("http://localhost:3000/api/v1/login", {
          username: this.loginUsername,
          password: this.loginPassword,
          email: this.loginEmail
        })
        .then(response => {
          store.commit("setJWT", response.data.token);
          VueCookies.set("loginToken", response.data.token);
          this.loginDialog = false;
          this.axios
            .get("http://localhost:3000/api/v1/protected/u", {
              headers: {
                Authorization: `Bearer ${response.data.token}`,
                "Content-Type": "application/json"
              }
            })
            .then(response => {
              store.commit("setProfile", response.data.data.profile);
              store.commit(
                "setSubscriptions",
                response.data.data.subscriptions
              );
              store.commit(
                "setSnackBarText",
                `Welcome ${store.state.profile.username}`
              );
            });
        })
        .catch(err => {
          if (err.response.status === 401) {
            store.commit("setSnackBarText", "Bad username or password");
            return;
          } else {
            store.commit("setSnackBarText", err.response.data.error.message);
          }
        });
    }
  }
};
</script>

<style scoped>
</style>
