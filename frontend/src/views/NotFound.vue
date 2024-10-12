<template>
  <div class="error-container">
    <h1 class="error-code" v-show="showErrorCode">404</h1>
    <p class="error-message" v-show="showErrorMessage">
      抱歉，您访问的页面不存在。
    </p>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showErrorCode: false,
      showErrorMessage: false,
      showErrorSuggestion: false,
    };
  },

  mounted() {
    setTimeout(() => {
      this.showErrorCode = true;
    }, 500);

    setTimeout(() => {
      this.showErrorMessage = true;
    }, 1000);

    setTimeout(() => {
      this.showErrorSuggestion = true;
    }, 1500);
  },

  methods: {
  },
};
</script>

<style>
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  color: #333;
}

.error-code {
  font-size: 10rem;
  font-weight: bold;
  opacity: 0;
  animation: errorCode 0.5s forwards;
}

.error-message {
  font-size: 2rem;
  opacity: 0;
  animation: errorMessage 0.5s forwards;
}

.error-suggestion {
  font-size: 1.5rem;
  opacity: 0;
  animation: errorSuggestion 0.5s forwards;
}

.error-suggestion a {
  color: #333;
  text-decoration: underline;
  cursor: pointer;
}

@keyframes errorCode {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes errorMessage {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes errorSuggestion {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>
