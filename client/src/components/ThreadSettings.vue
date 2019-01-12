<template>
  <v-dialog v-model="threadSettingsDialog" width="500" class="vert-dialogue">
    <v-list-tile slot="activator" avatar>
      <v-list-tile-avatar>
        <v-icon class="justify-end" large>settings</v-icon>
      </v-list-tile-avatar>
    </v-list-tile>
    <v-card>
      <v-card-title class="headline grey lighten-2" primary-title>
        <v-list-tile-avatar class="justify-end pr-4" size="64">
          <img v-if="thread.avatar" :src="`http://localhost:3000/uploads/${thread.avatar}`">
          <v-icon v-else size="64">people</v-icon>
        </v-list-tile-avatar>Thread Settings
      </v-card-title>

      <v-card-text>
        <v-form v-model="threadSettingsValid">
          <v-text-field label="Name" v-model="threadName" required outline/>
          <v-text-field
            outline
            label="Avatar"
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
        <v-btn color="primary" dark @click="submiThreadUpdate()">Edit</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>


<script>
export default {
  name: "ThreadSettings",
  data() {
    return {
      threadSettingsDialog: false,

      threadSettingsValid: false,
      threadName: null,
      threadDescription: null,

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
      imageUrl: ""
    };
  },
      props: {
        thread: Object
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
    submiThreadUpdate() {
      const fd = new FormData();
      fd.append("name", this.threadName);
      fd.append("description", this.threadDescription);
      fd.append("avatar", this.imageFile);

      this.axios
        .post("http://localhost:3000/api/v1/t", fd, {
          headers: {
            "Content-Type": "multipart/form-data"
          }
        })
        .then(() => {
          this.snackbarText = `Updated thread ${this.threadName}`;
          this.snackbar = true;
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
