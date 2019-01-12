<template>
  <v-dialog v-model="signupDialog" width="500" class="vert-dialogue">
    <v-btn slot="activator" class="toolbar-btn" flat outline dark small>Sign Up</v-btn>

    <v-card>
      <v-card-title class="headline grey lighten-2" primary-title>
        <img class="penny-logo mr-2" alt="Penny logo" src="../assets/penny_new.png">Sign Up
      </v-card-title>

      <v-card-text>
        <v-form v-model="signupValid">
          <v-text-field
            label="Username"
            v-model="signupUsername"
            required
            :rules="usernameRules"
            outline
          />
          <v-text-field label="Email" v-model="signupEmail" outline :rules="emailRules" required/>
          <v-text-field
            label="Password"
            v-model="signupPassword"
            outline
            type="password"
            :rules="passwordRules"
            required
          />
          <v-text-field
            outline
            label="Profile Avatar"
            @click="pickFile"
            v-model="imageName"
            prepend-icon="attach_file"
          ></v-text-field>
          <input
            type="file"
            style="display: none"
            ref="image"
            accept="image/*"
            @change="onFilePicked"
          >
          <img :src="imageUrl" height="150" v-if="imageUrl">
        </v-form>
      </v-card-text>

      <v-divider></v-divider>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" dark @click="signUp()">Submit</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import VueCookies from "vue-cookies";

import store from "../store";

export default {
  name: "CreatePostDialog",
  computed: {
    profile() {
      return store.state.profile;
    }
  },
  data() {
    return {
      signupDialog: false,

      signupValid: false,
      signupUsername: null,
      signupPassword: null,
      signupEmail: null,

      usernameRules: [v => !!v || "Username is required"],
      passwordRules: [v => !!v || "Password is required"],
      emailRules: [
        v => !!v || "E-mail is required",
        v =>
          /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(.\w{2,3})+$/.test(v) ||
          "E-mail must be valid"
      ],

      imageName: "",
      imageFile: "",
      imageUrl: "",
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
    signUp() {
      const fd = new FormData();
      fd.append('username',this.signupUsername);
      fd.append('password',this.signupPassword);
      fd.append('email',this.signupEmail);
      fd.append('avatar',this.imageFile);

      this.axios
        .post("http://localhost:3000/api/v1/u",fd,{
            headers: {
            'Content-Type': 'multipart/form-data'
            },
        })
        .then(() => {
          this.axios
            .post("http://localhost:3000/api/v1/login", {
              username: this.signupUsername,
              password: this.signupPassword,
              email: this.signupEmail
            })
            .then(response => {
              store.commit("setJWT", response.data.token);
              VueCookies.set("loginToken", response.data.token);
              this.signupDialog = false;
              this.axios
                .get("http://localhost:3000/api/v1/protected/u", {
                  headers: {
                    Authorization: `Bearer ${response.data.token}`,
                    "Content-Type": "application/json"
                  }
                })
                .then(response => {
                  store.commit("setProfile", response.data.data.profile);
                  this.snackbarText = `Welcome ${this.profile.username}`;
                  this.snackbar = true;
                });
            })
            .catch(err => {
              this.snackbarText = err.response.data.error.message;
              this.snackbar = true;
            });
        })
        .catch(err => {
          this.snackbarText = err.response.data.error.message;
          this.snackbar = true;
        });
    }
  }
};
</script>

<style scoped>
</style>
