<template>
  <div id="app">
    <form v-on:submit.prevent="connect">
      <p>
        <label for="user">User: </label>
        <input
          id="user"
          type="text"
          v-model="user"
          placeholder="Your name is here"
          maxlength="50"
          :disabled="isConnnected"
        />
        &nbsp;<label for="server">Chat Server: </label>
        <input
          id="server"
          type="text"
          v-model="server"
          :disabled="isConnnected"
        />
        &nbsp;<input type="submit" value="Connect" v-show="!isConnnected" />
        &nbsp;<input
          type="button"
          value="Disconnect"
          v-show="isConnnected"
          v-on:click="disconnect"
        />
        &nbsp;<span id="error">{{ error }}</span>
      </p>
      <hr />
    </form>
    <div id="messages" v-show="isConnnected">
      <div v-for="(message, i) in messages" :key="i">
        [{{ message.user }} - {{ message.time }}] {{ message.message }}
      </div>
    </div>
    <div v-show="isConnnected">
      <hr />
      <form v-on:submit.prevent="send">
        <p>
          <label for="message">Message</label>:
          <input
            id="message"
            type="text"
            v-model="message"
            placeholder="Say something here"
            maxlength="255"
          />
          &nbsp;<button :disabled="!readyToSend()">Send</button>
        </p>
      </form>
    </div>
  </div>
</template>

<script>
import axios from "axios";
// import http from "http";
// import https from "https";
export default {
  name: "App",
  data() {
    return {
      isConnnected: false,
      server: "",
      user: "",
      message: "",
      messages: [],
      error: "",
      source: null,
    };
  },

  created() {
    this.loadRecent();
    this.reset();
    this.clearError();
  },
  methods: {
    loadRecent() {
      this.server = localStorage.server;
      this.user = localStorage.user;
      if (this.server === undefined) {
        this.server = "https://chat-room-be.herokuapp.com";
      }
      if (this.user === undefined) {
        this.user = this.genUsername();
      }
    },
    clearError() {
      setInterval(() => this.error = "", 3000);
    },
    genUsername() {
      const characters = [
        "Moon",
        "Cony",
        "Brown",
        "James",
        "Sally",
        "Jessica",
        "Leonard",
        "Edward",
        "Boss",
        "Choco",
        "Pangyo",
      ];
      return (
        characters[Math.floor(Math.random() * characters.length)] +
        new Date().getTime()
      );
    },
    readyToSend() {
      return this.user !== "" && this.message != "";
    },
    send() {
      this.sendMessageToServer(this.user, this.message);
    },
    reset() {
      this.isConnnected = false;
      this.source = null;
      this.message = "";
    },
    connect() {
      localStorage.server = this.server;
      localStorage.user = this.user;
      this.isConnnected = true;
      this.$nextTick(() => document.getElementById("message").focus());
      this.source = new EventSource(this.server + "/message");
      this.source.onerror = (event) => {
        console.log(event);
        this.isConnnected = false;
        this.error = "Failed to Connect";
        this.disconnect();
      };
      this.source.onmessage = (event) => {
        console.log(event);
        this.messages.push(JSON.parse(event.data));
        this.$nextTick(() => {
          const messagesEl = document.getElementById("messages");
          messagesEl.scrollTop = messagesEl.scrollHeight;
        });
      };
      this.sendMessageToServer(this.user, "Connnected to " + this.server);
    },
    sendMessageToServer(user, message) {
      const data = {
        user: user,
        message: message,
      };
      axios
        .post(this.server + "/message", data)
        .then(() => this.message = "")
        .catch((error) => {
          console.log(error);
          this.error = error;
        });
    },
    disconnect() {
      this.source.close();
      this.reset();
    },
  },
};
</script>

<style>
* {
  font-family: Verdana;
}
html,
body,
#app {
  height: 100%;
}
#app {
  display: flex;
  flex-direction: column;
}
#error {
  color: #f00;
}
#server {
  width: 300px;
}
#user {
  width: 150px;
}
#message {
  width: 300px;
}
#messages {
  padding: 10px;
  margin-bottom: 10px;
  flex: 1;
  background-color: #eee;
  overflow: scroll;
}
</style>
